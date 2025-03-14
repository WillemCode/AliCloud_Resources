## 阿里雲多帳戶資產管理與查詢系統

## 專案簡介

在大型企業或團隊中，往往會存在多個阿里雲帳戶，這些帳戶分散在不同的部門或專案組中。日常需要對雲資源（ECS、RDS、SLB、PolarDB 等）進行管理和查詢，尤其當伺服器數量龐大且分布在多個區域時，手動維護繁瑣、容易出錯。

本專案的目標是將企業內部的多個阿里雲帳戶資源統一整合到本地資料庫（如 SQLite）中，通過腳本或未來的自助查詢頁面，讓使用者一鍵查詢任何伺服器的詳細資訊，以及它所在的具體阿里雲帳戶及區域。

## 背景與動機

- **多帳戶、分散管理**：公司/團隊擁有多個阿里雲帳戶，傳統方式需要逐個登錄阿里雲控制台查看，效率低、容易漏查。
- **手動維護不便**：將資產資訊記錄在 Excel 或其他開源 CMDB 工具時，往往需要手動更新，一旦忘記或延遲更新，就會導致資訊過期、不準確。
- **新人不友好**：新同事需要快速了解伺服器、資料庫實例、負載均衡等資源所在的帳戶、區域以及用途，這對於沒有整體視圖的新人而言非常痛苦。
- **未來擴展**：計畫進一步開發成自助查詢系統，輸入任意查詢條件（如伺服器名稱、IP 或其他關鍵字），就能快速定位該資產的詳細資訊、所屬帳戶，以及未來可能的用戶和密碼管理。

基於此，本專案應運而生，通過統一的 Go 工程實現自動化同步，減少人工操作和溝通成本，讓團隊成員可隨時查看最新的雲資源資訊。

## 功能特性

1. **多帳戶管理**
    - 支援在配置檔案中定義多個阿里雲帳戶資訊（AccessKey、SecretKey 等），一次執行即可同步全部帳戶下的資源。
2. **多區域支援**
    - 支援 ECS、RDS、SLB、PolarDB 等多種資源的多區域查詢，自動分頁拉取全部實例資料，避免只獲取部分資源。
3. **彈性公網 IP 收集**
    - 區分 ECS 自帶公網 IP 與彈性公網 IP；當 ECS 綁定了多張網卡和多個 EIP 時，也能全部收集並寫入資料庫。
4. **本地資料庫儲存**
    - 采用 SQLite 將所有雲資源資訊儲存於本地，無需額外部署資料庫，便於快速上手。
    - 專案結構中已經封裝了資料庫操作，可輕鬆擴展或更換資料庫類型（如 MySQL、PostgreSQL）。
5. **可擴展查詢介面**
    - 專案留有 API 端點範例（使用 Gin 等框架），後續可輕鬆擴展成完整的 RESTful 查詢服務，為新同事或其他系統提供自助查詢功能。

## 專案結構

```
.
├── cmd
│   └── main.go                   // 專案入口，執行阿里雲資產同步並初始化資料庫等
├── config
│   ├── config.go                 // 解析配置檔案
│   └── config.yaml               // 阿里雲帳戶及資料庫等配置信息
├── internal
│   └── services                  // 同步邏輯及阿里雲 API 呼叫
│       ├── ecs.go                // ECS 同步
│       ├── rds.go                // RDS 同步
│       ├── slb.go                // SLB 同步
│       └── polar.go              // PolarDB 同步
├── pkg
│   ├── config                    // 統一配置管理（可擴展 Viper、環境變數等）
│   ├── database                  // 資料庫操作封裝
│   │   └── database.go
│   └── logger                    // 日誌管理，使用 logrus 或 zap
│       └── logger.go
├── go.mod
├── go.sum
└── README.md                     // 專案說明文檔
```

## 快速開始

1. **克隆專案並進入目錄**
    
```bash
git clone https://github.com/WillemCode/AliCloud_Resources.git
cd AliCloud_Resources
```

> 注意⚠️: go version go1.23+
    
2. **修改配置**

* 打開 `config/config.yaml` ，填入你的阿里雲帳戶資訊，比如：

```yaml
aliyun_accounts:
  - name: "AccountA"
    access_key: "YourAccessKey"
    access_secret: "YourSecret"
    ecs_region_ids:
      - "cn-hangzhou"
      - "cn-beijing"
    rds_region_id: "cn-beijing"
    slb_region_id: "cn-hangzhou"
    polardb_region_id: "cn-beijing"
database:
  path: "./sqlite.db"   # SQLite 資料庫檔案路徑
log_level: "info"     # 預設日誌級別
```

* `ecs_region_ids` 為陣列，可同時拉取多個區域的 ECS 資源。

3. **安裝依賴**

```bash
go get github.com/mattn/go-sqlite3
go get github.com/sirupsen/logrus
go get github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests
go get github.com/aliyun/alibaba-cloud-sdk-go/services/ecs
go get github.com/aliyun/alibaba-cloud-sdk-go/services/polardb
go get github.com/aliyun/alibaba-cloud-sdk-go/services/rds
go get github.com/aliyun/alibaba-cloud-sdk-go/services/slb
go get github.com/spf13/viper
```
    
4. **執行同步**
    
```bash
go run cmd/main.go      // 執行
```

* 程式會根據 `config.yaml` 加載配置、初始化日誌和資料庫，然後依次同步 ECS、RDS、SLB、PolarDB 等資源，並將資料存儲到本地資料庫。
* 預設在 `./sqlite.db` 中生成 SQLite 資料庫檔案（可在配置中修改路徑）。

5. **查看結果**

* 你可以透過任意 SQLite 客戶端（或在程式碼中）查詢採集到的 ECS、RDS、SLB 等資源資訊：

```bash
sqlite3 sqlite.db
sqlite> .tables
ecs   rds   slb   polardb
sqlite> select * from ecs limit 5;
```

6. **編譯構建**

Linux/Mac 系統:
```
cd cmd/
go build -o AliCloud_Resource
```

ARM 芯片：
```
GOOS=darwin GOARCH=arm64 go build -o AliCloud_Resource
```

AMD 芯片：
```
GOOS=darwin GOARCH=amd64 go build -o AliCloud_Resource
```

Windows：
```
cd cmd/
GOOS=windows GOARCH=amd64 go build -o AliCloud_Resource.exe
```

## 使用場景

1. **企業/團隊內多帳戶管理**
    - 想要整合不同專案組、部門的阿里雲帳戶資源，快速查看全部資產。
2. **新人快速接手**
    - 提供統一查詢端點或頁面，使新同事可立即了解伺服器歸屬、IP 地址、地區等關鍵資訊。
3. **資產審計**
    - 每次手動或定時執行同步腳本，都能獲取最新資源資訊，用於成本評估或安全審計。
4. **後續擴展**
    - 可整合到企業內部的信息門戶或運維系統，通過 RESTful API 或直接查詢資料庫，簡化維護流程。

## 未來規劃

1. **自助查詢系統**
    - 計劃打造一個 Web 前端頁面，使用者輸入關鍵字（伺服器名、IP 等），能返回資產完整詳情並指明其所屬帳戶及區域。
2. **訪問控制與密碼管理**
    - 在 ECS 同步資訊中增加負責人、用戶、密碼等欄位，形成更完善的 CMDB，自定義權限與訪問控制。
3. **更換資料庫**
    - 如果團隊規模擴大，需要更高並發或更可靠的存儲，可切換至 MySQL、PostgreSQL 等，無需大改專案結構。
4. **自動化告警**
    - 定期任務（crontab/雲函數）執行同步後，對新增/異常實例進行郵件或簡訊告警，方便管理員及時處理。

## 貢獻與支持
- **Issues**：如在使用過程中遇到問題或有新需求，歡迎在倉庫的 Issue 中提出。
- **Pull Requests**：歡迎參與貢獻新的功能、優化或文檔修復。

---

## 授權說明

本專案採用 [GNU General Public License (GPL)](./LICENSE) 進行開源發布。  
這意味著：

- 你可以自由複製、修改和分發本專案的源代碼，但修改後的專案也必須繼續以 GPL 或兼容的許可證進行發布；
- 分發或發布時，需包含本專案的原始版權聲明與 GPL 協議文本，並提供完整的源代碼獲取方式。

請參閱 [LICENSE](./LICENSE) 文件獲取詳細條款。若你對 GPL 的使用及合規性有任何疑問，請查閱 [GNU 官網](https://www.gnu.org/licenses/) 或諮詢相關專業人士。

---

## Star 歷史

[![Star History Chart](https://api.star-history.com/svg?repos=WillemCode/AliCloud_Resources&type=Date)](https://www.star-history.com/#WillemCode/AliCloud_Resources&Date)
