# 阿里云多账户资产管理与查询系统

## 项目简介

在大型企业或团队中，往往会存在多个阿里云账户，这些账户分散在不同的部门或项目组中。日常需要对云资源（ECS、RDS、SLB、PolarDB 等）进行管理和查询，尤其当服务器数量庞大且分布在多个区域时，手动维护繁琐、容易出错，而且对于新同事极不友好。

本项目的目标是将企业内部的多个阿里云账户资源统一整合到本地数据库（如 SQLite）中，通过脚本或未来的自助查询页面，让用户一键查询任何服务器的详细信息，以及它所在的具体阿里云账户及区域。

## 背景与动机

- **多账户、分散管理**：公司/团队拥有多个阿里云账户，传统方式需要逐个登录阿里云控制台查看，效率低、容易漏查。
- **手动维护不便**：将资产信息记录在 Excel 或其他开源 CMDB 工具时，往往需要手动更新，一旦忘记或延迟更新，就会导致信息过期、不准确。
- **新人不友好**：新同事需要快速了解服务器、数据库实例、负载均衡等资源所在的账户、区域以及用途，这对于没有整体视图的新人而言非常痛苦。
- **未来扩展**：计划进一步开发成自助查询系统，输入任意查询条件（如服务器名称、IP 或其他关键词），就能快速定位该资产的详细信息、所属账户，以及未来可能的用户和密码管理。

基于此，本项目应运而生，通过统一的 Go 工程实现自动化同步，减少人工操作和沟通成本，让团队成员可随时查看最新的云资源信息。

## 功能特性

1. **多账户管理**
    - 支持在配置文件中定义多个阿里云账户信息（AccessKey、SecretKey 等），一次运行即可同步全部账户下的资源。
2. **多区域支持**
    - 支持 ECS、RDS、SLB、PolarDB 等多种资源的多区域查询，自动分页拉取全部实例数据，避免只获取部分资源。
3. **弹性公网 IP 收集**
    - 区分 ECS 自带公网 IP 与弹性公网 IP；当 ECS 绑定了多张网卡和多个 EIP 时，也能全部收集并写入数据库。
4. **本地数据库存储**
    - 采用 SQLite 将所有云资源信息存储在本地，无需额外部署数据库，便于快速上手。
    - 项目结构中已经封装了数据库操作，可轻松扩展或更换数据库类型（如 MySQL、PostgreSQL）。
5. **可扩展查询接口**
    - 项目留有 API 端点示例（使用 Gin 等框架），后续可轻松扩展成完整的 RESTful 查询服务，为新同事或其他系统提供自助查询功能。

## 项目结构

```
.
├── cmd
│   └── main.go                   // 项目入口，执行阿里云资产同步并初始化数据库等
├── config
│   ├── config.go                 // 解析配置文件
│   └── config.yaml               // 阿里云账户及数据库等配置信息
├── internal
│   └── services                  // 同步逻辑及阿里云 API 调用
│       ├── ecs.go                // ECS 同步
│       ├── rds.go                // RDS 同步
│       ├── slb.go                // SLB 同步
│       └── polar.go              // PolarDB 同步
├── pkg
│   ├── config                    // 统一配置管理（可扩展 Viper、环境变量等）
│   ├── database                  // 数据库操作封装
│   │   └── database.go
│   └── logger                    // 日志管理，使用 logrus 或 zap
│       └── logger.go
├── go.mod
├── go.sum
└── README.md                     // 项目说明文档
```

## 快速开始

1. **克隆项目并进入目录**
    
```bash
git clone https://github.com/WillemCode/AliCloud_Resources.git
cd AliCloud_Resources
```

    
2. **修改配置**


* 打开 `config/config.yaml` ，填入你的阿里云账户信息，比如：

        
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
  path: "./cmdb.db"   # SQLite 数据库文件路径
log_level: "info"     # 默认日志级别
```


* `ecs_region_ids` 为数组，可同时拉取多个区域的 ECS 资源。


3. **安装依赖**


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
    
4. **运行同步**
    
```bash
go run cmd/main.go      // 运行
```
    
    - 程序会根据 `config.yaml` 加载配置、初始化日志和数据库，然后依次同步 ECS、RDS、SLB、PolarDB 等资源，并将数据存储到本地数据库。
    - 默认在 `./cmdb.db` 中生成 SQLite 数据库文件（可在配置中修改路径）。
5. **查看结果**
    
    - 你可以通过任意 SQLite 客户端（或在代码中）查询采集到的 ECS、RDS、SLB 等资源信息：
        
```bash
sqlite3 cmdb.db
sqlite> .tables
ecs   rds   slb   polardb
sqlite> select * from ecs limit 5;
```
        
## 使用场景

1. **企业/团队内多账户管理**
    - 想要整合不同项目组、部门的阿里云账号资源，快速查看全部资产。
2. **新人快速接手**
    - 提供统一查询端点或页面，使新同事可立即了解服务器归属、IP 地址、地区等关键信息。
3. **资产审计**
    - 每次手动或定时执行同步脚本，都能获取最新资源信息，用于成本评估或安全审计。
4. **后续扩展**
    - 可集成到企业内部的信息门户或运维系统，通过 RESTful API 或直接查询数据库，简化维护流程。

## 未来规划

1. **自助查询系统**
    - 计划打造一个 Web 前端页面，用户输入关键字（服务器名、IP 等），能返回资产完整详情并指明其所属账户及区域。
2. **访问控制与密码管理**
    - 在 ECS 同步信息中增加负责人、用户、密码等字段，形成更完善的 CMDB，自定义权限与访问控制。
3. **更换数据库**
    - 如果团队规模扩大，需要更高并发或更可靠的存储，可切换至 MySQL、PostgreSQL 等，无需大改项目结构。
4. **自动化告警**
    - 定期任务（crontab/云函数）执行同步后，对新增/异常实例进行邮件或短信告警，方便管理员及时处理。

## 贡献与支持
- **Issues**：如在使用过程中遇到问题或有新需求，欢迎在仓库的 Issue 中提出。
- **Pull Requests**：欢迎参与贡献新的功能、优化或文档修复。
- **License**：本项目采用 MIT License / 其他协议，具体请查看 [LICENSE](https://chatgpt.com/c/LICENSE) 文件。

---

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=WillemCode/AliCloud_Resources&type=Date)](https://www.star-history.com/#WillemCode/AliCloud_Resources&Date)
