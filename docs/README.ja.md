# Aliyun マルチアカウント資産管理とクエリシステム

## プロジェクト紹介

大規模な企業やチームでは、複数の Aliyun アカウントが異なる部門やプロジェクトチームに分かれて存在することがよくあります。日常的にクラウドリソース（ECS、RDS、SLB、PolarDB など）の管理とクエリが必要です。特にサーバーの数が膨大で、複数のリージョンに分散している場合、手動での維持は煩雑で、エラーが発生しやすいです。

このプロジェクトの目的は、企業内部の複数の Aliyun アカウントリソースをローカルデータベース（例えば SQLite）に統合し、スクリプトや今後のセルフサービスクエリページを通じて、ユーザーがどのサーバーの詳細情報、所属する Aliyun アカウントやリージョンを簡単にクエリできるようにすることです。

## 背景と動機

- **複数アカウント、分散管理**：会社やチームが複数の Aliyun アカウントを所有しており、従来の方法では各 Aliyun コンソールにログインして確認する必要があり、効率が悪く、見逃しが発生しやすいです。
- **手動維持が不便**：資産情報を Excel や他のオープンソースの CMDB ツールに記録する場合、手動で更新する必要があり、忘れたり更新が遅れたりすると、情報が古くなったり、不正確になったりします。
- **新人に優しくない**：新しい同僚がサーバー、データベースインスタンス、ロードバランサーなどのリソースがどのアカウントやリージョンに属しているかを迅速に把握する必要があります。全体的なビューがない新人には非常に辛いです。
- **将来的な拡張**：入力したサーバー名、IP などのキーワードで資産の詳細情報、所属アカウント、将来的にユーザーとパスワード管理を提供できるように、自分自身でクエリできるシステムを構築する予定です。

これらを基に、このプロジェクトは生まれました。Go プロジェクトを通じて、手動操作とコミュニケーションコストを削減し、チームメンバーが常に最新のクラウドリソース情報を確認できるようにします。

## 機能

1. **複数アカウント管理**
    - 複数の Aliyun アカウント情報（AccessKey、SecretKey など）を設定ファイルで定義することができ、一度の実行で全アカウント下のリソースを同期できます。
2. **複数リージョン対応**
    - ECS、RDS、SLB、PolarDB など、複数リージョンでリソースをクエリし、自動的にページネーションして全てのインスタンスデータを取得します。これにより部分的なリソースの取得を防ぎます。
3. **弾性公网 IP の収集**
    - ECS による自動公共 IP と弾性公共 IP を区別します。ECS が複数のネットワークインターフェースと複数の EIP をバインドしている場合でも、全てを収集し、データベースに書き込みます。
4. **ローカルデータベース保存**
    - SQLite を使用して、全てのクラウドリソース情報をローカルに保存します。追加のデータベース展開が不要で、すぐに始めることができます。
    - プロジェクト構造内でデータベース操作が封装されており、MySQL や PostgreSQL などへの変更が容易です。
5. **拡張可能なクエリインターフェース**
    - プロジェクトには API エンドポイントの例（Gin などのフレームワーク使用）が含まれており、将来的には完全な RESTful クエリサービスとして拡張できます。

## プロジェクト構造

```
.
├── cmd
│   └── main.go                   // プロジェクトのエントリーポイント、Aliyun 資産の同期とデータベース初期化
├── config
│   ├── config.go                 // 設定ファイルの解析
│   └── config.yaml               // Aliyun アカウントとデータベース設定情報
├── internal
│   └── services                  // 同期ロジックと Aliyun API 呼び出し
│       ├── ecs.go                // ECS 同期
│       ├── rds.go                // RDS 同期
│       ├── slb.go                // SLB 同期
│       └── polar.go              // PolarDB 同期
├── pkg
│   ├── config                    // 統一的な設定管理（Viper や環境変数などで拡張可能）
│   ├── database                  // データベース操作のラッピング
│   │   └── database.go
│   └── logger                    // ログ管理（logrus や zap を使用）
│       └── logger.go
├── go.mod
├── go.sum
└── README.md                     // プロジェクトドキュメント
```

## クイックスタート

1. **プロジェクトをクローンし、ディレクトリに移動**

```bash
git clone https://github.com/WillemCode/AliCloud_Resources.git
cd AliCloud_Resources
```

> 注意⚠️: go version go1.23+

2. **設定を変更**

* `config/config.yaml` を開き、Aliyun アカウント情報を入力します。

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
  path: "./sqlite.db"   # SQLite データベースファイルパス
log_level: "info"     # デフォルトのログレベル
```

3. **依存関係をインストール**

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

4. **同期を実行**

```bash
go run cmd/main.go      // 実行
```

5. **結果を見る**

* 任意の SQLite クライアント（またはコード内）を使用して、収集した ECS、RDS、SLB などのリソース情報をクエリできます：

```bash
sqlite3 sqlite.db
sqlite> .tables
ecs   rds   slb   polardb
sqlite> select * from ecs limit 5;
```

6. **プロジェクトをビルド**

Linux/Mac システム：
```
cd cmd/
go build -o AliCloud_Resource
```

ARM チップ：
```
GOOS=darwin GOARCH=arm64 go build -o AliCloud_Resource
```

AMD チップ：
```
GOOS=darwin GOARCH=amd64 go build -o AliCloud_Resource
```

Windows：
```
cd cmd/
GOOS=windows GOARCH=amd64 go build -o AliCloud_Resource.exe
```

---

## 使用シナリオ

1. **企業/チーム内の複数アカウント管理**
    - 異なるプロジェクトチーム、部門の Aliyun アカウントリソースを統合し、すべての資産を迅速に確認。
2. **新人がすぐに把握できる**
    - 統一されたクエリエンドポイントやページを提供し、新しい同僚がすぐ

にサーバーの所有者、IP アドレス、地域などの重要情報を把握できるようにする。
3. **資産監査**
    - 手動または定期的に同期スクリプトを実行し、最新のリソース情報を取得し、コスト評価やセキュリティ監査に使用。
4. **将来の拡張**
    - RESTful API を通じて、企業内部の情報ポータルや運用システムに統合し、メンテナンスプロセスを簡素化する。

## 今後の計画

1. **セルフサービスクエリシステム**
    - Web フロントページを作成し、ユーザーがサーバー名や IP などのキーワードを入力すると、資産の完全な詳細情報とそのアカウントおよび地域を表示します。
2. **アクセス制御とパスワード管理**
    - ECS 同期情報に責任者、ユーザー、パスワードなどのフィールドを追加し、CMDB をより完全にし、カスタマイズ可能な権限とアクセス制御を形成。
3. **データベース変更**
    - チームが拡大し、より高い並列処理や信頼性のあるストレージが必要な場合、MySQL や PostgreSQL などに切り替え、プロジェクト構造の大きな変更なしで対応。
4. **自動アラート**
    - 定期的なタスク（crontab/クラウド関数）を使い、同期後に新規/異常インスタンスに対してメールや SMS アラートを送信し、管理者が迅速に対応できるようにする。

## 貢献とサポート

- **Issues**：使用中に問題が発生したり、新しい要求があったりした場合は、リポジトリの Issue セクションに提出してください。
- **Pull Requests**：新機能、最適化、またはドキュメントの修正への貢献を歓迎します。

---

## ライセンス

このプロジェクトは [GNU General Public License (GPL)](./LICENSE) の下で公開されています。

これにより：

- このプロジェクトのソースコードを自由にコピー、修正、配布できますが、変更後のプロジェクトも GPL または互換性のあるライセンスで公開する必要があります；
- 配布または公開時には、元の著作権表示と GPL 契約書を含め、完全なソースコードへのアクセス方法を提供する必要があります。

詳細な条項については [LICENSE](./LICENSE) ファイルをご覧ください。GPL の使用と遵守について質問がある場合は、[GNU のウェブサイト](https://www.gnu.org/licenses/)を参照するか、専門家に相談してください。

---

## Star 歴史

[![Star History Chart](https://api.star-history.com/svg?repos=WillemCode/AliCloud_Resources&type=Date)](https://www.star-history.com/#WillemCode/AliCloud_Resources&Date)
