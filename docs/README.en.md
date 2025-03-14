# Aliyun Multi-Account Asset Management and Query System

## Project Introduction

In large enterprises or teams, multiple Aliyun accounts are often spread across different departments or project teams. Managing and querying cloud resources (ECS, RDS, SLB, PolarDB, etc.) is a daily requirement. When dealing with a large number of servers distributed across multiple regions, manual maintenance is cumbersome and prone to errors.

The goal of this project is to consolidate multiple Aliyun account resources within an enterprise into a local database (e.g., SQLite). Through scripts or future self-service query pages, users can instantly retrieve detailed information about any server, including the corresponding Aliyun account and region.

## Background and Motivation

- **Multiple Accounts, Decentralized Management**  
  Companies/teams possess multiple Aliyun accounts, requiring users to log into the Aliyun console for each account manually—an inefficient and error-prone process.  
- **Inconvenience of Manual Maintenance**  
  Storing asset information in Excel or other open-source CMDB tools often necessitates manual updates. If updates are missed or delayed, information quickly becomes outdated and inaccurate.  
- **Unfriendly to Newcomers**  
  New employees need to quickly understand which accounts and regions host various servers, database instances, and load balancers. Without a comprehensive overview, this process can be highly inefficient.  
- **Future Expansion**  
  Plans are underway to develop this into a self-service query system where users can enter search criteria (e.g., server name, IP address) and quickly retrieve asset details, associated accounts, and potential user credentials.  

This project aims to automate asset management through a unified Go-based solution, reducing manual operations and communication costs while ensuring that team members can access the latest cloud resource data at any time.

## Features

1. **Multi-Account Management**  
   - Supports defining multiple Aliyun account credentials (AccessKey, SecretKey, etc.) in a configuration file. Running the script once synchronizes resources across all accounts.  
2. **Multi-Region Support**  
   - Enables querying ECS, RDS, SLB, and PolarDB resources across multiple regions. The system automatically paginates through all instances, preventing partial data retrieval.  
3. **Elastic Public IP Collection**  
   - Differentiates between built-in ECS public IPs and Elastic IPs (EIPs). If ECS instances have multiple network interfaces and EIPs, all details are collected and stored in the database.  
4. **Local Database Storage**  
   - Uses SQLite to store cloud asset information locally, eliminating the need for an external database, making setup easier.  
   - The database operations are encapsulated, allowing seamless migration to MySQL or PostgreSQL if required.  
5. **Expandable Query Interface**  
   - Provides API endpoint examples (using Gin and other frameworks), enabling future development into a full RESTful query service for self-service asset retrieval.  

## Project Structure

```
.
├── cmd
│   └── main.go                   // Project entry point, synchronizes Aliyun assets and initializes the database
├── config
│   ├── config.go                 // Parses configuration files
│   └── config.yaml               // Stores Aliyun account credentials and database settings
├── internal
│   └── services                  // Sync logic and Aliyun API integration
│       ├── ecs.go                // ECS synchronization
│       ├── rds.go                // RDS synchronization
│       ├── slb.go                // SLB synchronization
│       └── polar.go              // PolarDB synchronization
├── pkg
│   ├── config                    // Centralized configuration management
│   ├── database                  // Database abstraction
│   │   └── database.go
│   └── logger                    // Logging (using logrus or zap)
│       └── logger.go
├── go.mod
├── go.sum
└── README.md                     // Project documentation
```

## Quick Start

1. **Clone the repository and enter the directory**

```bash
git clone https://github.com/WillemCode/AliCloud_Resources.git
cd AliCloud_Resources
```

> ⚠️ **Note**: Requires Go version 1.23+

2. **Modify the Configuration**

* Open `config/config.yaml` and enter your Aliyun account credentials:

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
  path: "./sqlite.db"   # SQLite database path
log_level: "info"       # Default log level
```

3. **Install Dependencies**

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

4. **Run the Synchronization Process**

```bash
go run cmd/main.go
```

* The program loads configurations from `config.yaml`, initializes logs and the database, and then synchronizes ECS, RDS, SLB, PolarDB, and other resources, storing them locally.
* By default, the SQLite database is saved in `./sqlite.db` (modifiable in the configuration).

5. **Query Results**

* Use any SQLite client to query the collected ECS, RDS, and SLB data:

```bash
sqlite3 sqlite.db
sqlite> .tables
ecs   rds   slb   polardb
sqlite> select * from ecs limit 5;
```

6. **Build the Binary**

Linux/Mac:

```
cd cmd/
go build -o AliCloud_Resource
```

ARM Chip:

```
GOOS=darwin GOARCH=arm64 go build -o AliCloud_Resource
```

AMD Chip:

```
GOOS=darwin GOARCH=amd64 go build -o AliCloud_Resource
```

Windows:

```
cd cmd/
GOOS=windows GOARCH=amd64 go build -o AliCloud_Resource.exe
```

---

## License

This project is released under the [GNU General Public License (GPL)](./LICENSE).

This means:

- You are free to copy, modify, and distribute this project’s source code, but any modifications must also be released under the GPL or a compatible license.
- When distributing or publishing, you must include the original copyright notice, GPL license text, and provide access to the full source code.

For details, refer to the [LICENSE](./LICENSE) file. If you have any questions about GPL compliance, please visit [GNU's official website](https://www.gnu.org/licenses/) or consult a legal professional.

---

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=WillemCode/AliCloud_Resources&type=Date)](https://www.star-history.com/#WillemCode/AliCloud_Resources&Date)
