package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite 驱动
)

// 全局数据库连接对象（sqlite3）
var db *sql.DB

// 初始化数据库，建立连接并创建表（如不存在）
func Init(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("无法打开数据库: %w", err)
	}
	// 测试数据库连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}
	// 创建所需的数据表（如果不存在）
	// ECS 实例信息表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ecs (
        instance_id TEXT PRIMARY KEY,
		cloud_name TEXT,
        instance_name TEXT,
        status TEXT,
        region_id TEXT,
        os_name TEXT,
        instance_type TEXT,
        cpu INTEGER,
        memory INTEGER,
        public_ip TEXT,
        private_ip TEXT,
		remarks TEXT,
		login_user TEXT,
		login_passwd TEXT
    );`)
	if err != nil {
		return fmt.Errorf("创建 ecs 表失败: %w", err)
	}
	// RDS 实例信息表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS rds (
        instance_id TEXT PRIMARY KEY,
		cloud_name TEXT,
        engine TEXT,
        region_id TEXT,
        status TEXT,
        memory INTEGER,
        instance_description TEXT,
        connection_string TEXT,
		remarks TEXT
    );`)
	if err != nil {
		return fmt.Errorf("创建 rds 表失败: %w", err)
	}
	// Tair Redis 实例信息表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS redis (
		instance_id TEXT PRIMARY KEY,
		cloud_name TEXT,
		instance_name TEXT,
		port INTEGER,
		region_id TEXT,
		capacity INTEGER,
		instance_class TEXT,
		qps INTEGER,
		band_width INTEGER,
		connections INTEGER,
		instance_type TEXT,
		connection_string TEXT,
		ip_address TEXT,
		remarks TEXT
	);`)
	if err != nil {
		return fmt.Errorf("创建 Tair 表失败: %w", err)
	}
	// SLB 实例信息表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS slb (
        lb_id TEXT PRIMARY KEY,
		cloud_name TEXT,
        lb_name TEXT,
        ip_address TEXT,
        band_width INTEGER,
        network_type TEXT,
        region_id TEXT,
        lb_status TEXT,
		remarks TEXT
    );`)
	if err != nil {
		return fmt.Errorf("创建 slb 表失败: %w", err)
	}
	// PolarDB 实例信息表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS polardb (
        dbcluster_id TEXT PRIMARY KEY,
		cloud_name TEXT,
        engine TEXT,
        region_id TEXT,
		db_cluster_status TEXT,
        dbnode_number INTEGER,
        dbcluster_description TEXT,
        memory_size INTEGER,
        connection_string TEXT,
		remarks TEXT
    );`)
	if err != nil {
		return fmt.Errorf("创建 polardb 表失败: %w", err)
	}

	return nil
}

// 关闭数据库连接
func Close() {
	if db != nil {
		_ = db.Close()
	}
}

// ------ 以下是各类资源的 CRUD 操作封装 ------

// ECSRecord 定义 ECS 记录的本地结构，用于数据库读写
type ECSRecord struct {
	InstanceID   string // 实例ID
	CloudName    string // 账户名称
	InstanceName string // 实例名称
	Status       string // 实例状态
	RegionID     string // 区域ID
	OSName       string // 操作系统名称
	InstanceType string // 实例规格
	CPU          int64  // CPU核数
	Memory       int64  // 内存大小
	PublicIP     string // 公网IP地址(逗号分隔)
	PrivateIP    string // 内网IP地址
}

// SaveECSRecords 将一组 ECS 实例记录保存到数据库
func SaveECSRecords(records []ECSRecord) error {
	for _, rec := range records {
		_, err := db.Exec(
			`INSERT OR REPLACE INTO ecs 
             (instance_id, cloud_name, instance_name, status, region_id, os_name, instance_type, cpu, memory, public_ip, private_ip) 
             VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			rec.InstanceID, rec.CloudName, rec.InstanceName, rec.Status, rec.RegionID, rec.OSName,
			rec.InstanceType, rec.CPU, rec.Memory, rec.PublicIP, rec.PrivateIP,
		)
		if err != nil {
			// 返回封装了上下文的错误，包含出错的实例ID
			return fmt.Errorf("插入 ECS 记录失败 (InstanceID=%s): %w", rec.InstanceID, err)
		}
	}
	return nil
}

// 查询所有 ECS 记录（用于 API 层示例）
func ListECSRecords() ([]ECSRecord, error) {
	rows, err := db.Query(
		"SELECT instance_id, cloud_name, instance_name, status, region_id, os_name, instance_type, cpu, memory, public_ip, private_ip FROM ecs",
	)
	if err != nil {
		return nil, fmt.Errorf("查询 ECS 表失败: %w", err)
	}
	defer rows.Close()

	var results []ECSRecord
	for rows.Next() {
		var rec ECSRecord
		// 将查询结果的每一行扫描到 ECSRecord 结构体
		err := rows.Scan(&rec.InstanceID, &rec.CloudName, &rec.InstanceName, &rec.Status, &rec.RegionID,
			&rec.OSName, &rec.InstanceType, &rec.CPU, &rec.Memory, &rec.PublicIP, &rec.PrivateIP)
		if err != nil {
			return nil, fmt.Errorf("读取 ECS 行数据失败: %w", err)
		}
		results = append(results, rec)
	}
	return results, nil
}

// （类似地，我们为 RDS、SLB、PolarDB 定义各自的 Record 结构和保存函数）

// RDS 数据结构和保存
type RDSRecord struct {
	InstanceID       string
	CloudName        string
	Engine           string
	RegionID         string
	Status           string
	Memory           int64
	Description      string
	ConnectionString string
}

func SaveRDSRecords(records []RDSRecord) error {
	for _, rec := range records {
		_, err := db.Exec(
			`INSERT OR REPLACE INTO rds 
             (instance_id, cloud_name, engine, region_id, status, memory, instance_description, connection_string)
             VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			rec.InstanceID, rec.CloudName, rec.Engine, rec.RegionID, rec.Status, rec.Memory, rec.Description, rec.ConnectionString,
		)
		if err != nil {
			return fmt.Errorf("插入 RDS 记录失败 (InstanceID=%s): %w", rec.InstanceID, err)
		}
	}
	return nil
}

// 查询所有 RDS 记录（用于 API 层示例）
func ListRDSRecords() ([]RDSRecord, error) {
	rows, err := db.Query(
		"SELECT instance_id, cloud_name, engine, region_id, status, memory, instance_description, connection_string FROM rds",
	)
	if err != nil {
		return nil, fmt.Errorf("查询 RDS 表失败: %w", err)
	}
	defer rows.Close()

	var results []RDSRecord
	for rows.Next() {
		var rec RDSRecord
		// 将查询结果的每一行扫描到 RDSRecord 结构体
		err := rows.Scan(&rec.InstanceID, &rec.CloudName, &rec.Engine, &rec.RegionID,
			&rec.Status, &rec.Memory, &rec.Description, &rec.ConnectionString)
		if err != nil {
			return nil, fmt.Errorf("读取 RDS 行数据失败: %w", err)
		}
		results = append(results, rec)
	}
	return results, nil
}

// SLB 数据结构和保存
type SLBRecord struct {
	InstanceID       string
	CloudName        string
	LoadBalancerName string
	IPAddress        string
	Bandwidth        int64
	NetworkType      string
	RegionID         string
	Status           string
}

func SaveSLBRecords(records []SLBRecord) error {
	for _, rec := range records {
		_, err := db.Exec(
			`INSERT OR REPLACE INTO slb 
             (lb_id, cloud_name, lb_name, ip_address, band_width, network_type, region_id, lb_status)
             VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			rec.InstanceID, rec.CloudName, rec.LoadBalancerName, rec.IPAddress, rec.Bandwidth, rec.NetworkType, rec.RegionID, rec.Status,
		)
		if err != nil {
			return fmt.Errorf("插入 SLB 记录失败 (LoadBalancerID=%s): %w", rec.InstanceID, err)
		}
	}
	return nil
}

// 查询所有 SLB 记录（用于 API 层示例）
func ListSLBRecords() ([]SLBRecord, error) {
	rows, err := db.Query(
		"SELECT lb_id, cloud_name, lb_name, ip_address, band_width, network_type, region_id, lb_status FROM slb",
	)
	if err != nil {
		return nil, fmt.Errorf("查询 SLB 表失败: %w", err)
	}
	defer rows.Close()

	var results []SLBRecord
	for rows.Next() {
		var rec SLBRecord
		// 将查询结果的每一行扫描到 SLBRecord 结构体
		err := rows.Scan(&rec.InstanceID, &rec.CloudName, &rec.LoadBalancerName, &rec.IPAddress,
			&rec.Bandwidth, &rec.NetworkType, &rec.RegionID, &rec.Status)
		if err != nil {
			return nil, fmt.Errorf("读取 SLB 行数据失败: %w", err)
		}
		results = append(results, rec)
	}
	return results, nil
}

// Tair 数据结构和保存
type RedisRecord struct {
	InstanceID       string
	CloudName        string
	InstanceName     string
	Port             int64
	RegionId         string
	Capacity         int64
	InstanceClass    string
	QPS              int64
	Bandwidth        int64
	Connections      int64
	InstanceType     string
	ConnectionString string
	IPAddress        string
}

func SaveRedisRecords(records []RedisRecord) error {
	for _, rec := range records {
		_, err := db.Exec(
			`INSERT OR REPLACE INTO redis 
             (instance_id, cloud_name, instance_name, port, region_id, capacity, instance_class, qps, band_width, connections, instance_type, connection_string, ip_address)
             VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			rec.InstanceID, rec.CloudName, rec.InstanceName, rec.Port, rec.RegionId, rec.Capacity, rec.InstanceClass, rec.QPS,
			rec.Bandwidth, rec.Connections, rec.InstanceType, rec.ConnectionString, rec.IPAddress,
		)
		if err != nil {
			return fmt.Errorf("插入 Tair Redis 记录失败 (InstanceID=%s): %w", rec.InstanceID, err)
		}
	}
	return nil
}

// 查询所有 Tair Redis 记录（用于 API 层示例）
func ListRedisRecords() ([]RedisRecord, error) {
	rows, err := db.Query(
		"SELECT instance_id, cloud_name, instance_name, port, region_id, capacity, instance_class, qps, band_width, connections, instance_type, connection_string, ip_address FROM redis",
	)
	if err != nil {
		return nil, fmt.Errorf("查询 Tair Redis 表失败: %w", err)
	}
	defer rows.Close()

	var results []RedisRecord
	for rows.Next() {
		var rec RedisRecord
		// 将查询结果的每一行扫描到 RDSRecord 结构体
		err := rows.Scan(&rec.InstanceID, &rec.CloudName, &rec.InstanceName, &rec.Port, &rec.RegionId, &rec.Capacity, &rec.InstanceClass, &rec.QPS,
			&rec.Bandwidth, &rec.Connections, &rec.InstanceType, &rec.ConnectionString, &rec.IPAddress)
		if err != nil {
			return nil, fmt.Errorf("读取 Tair Redis 行数据失败: %w", err)
		}
		results = append(results, rec)
	}
	return results, nil
}

// PolarDB 数据结构和保存
type PolarDBRecord struct {
	InstanceID       string
	CloudName        string
	Engine           string
	RegionID         string
	Status           string
	DBNodeCount      int64
	Description      string
	MemorySize       int64
	ConnectionString string
}

func SavePolarDBRecords(records []PolarDBRecord) error {
	for _, rec := range records {
		_, err := db.Exec(
			`INSERT OR REPLACE INTO polardb 
             (dbcluster_id, cloud_name, engine, region_id, db_cluster_status, dbnode_number, dbcluster_description, memory_size, connection_string)
             VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			rec.InstanceID, rec.CloudName, rec.Engine, rec.RegionID, rec.Status, rec.DBNodeCount, rec.Description, rec.MemorySize, rec.ConnectionString,
		)
		if err != nil {
			return fmt.Errorf("插入 PolarDB 记录失败 (DBClusterID=%s): %w", rec.InstanceID, err)
		}
	}
	return nil
}

// 查询所有 Polardb 记录（用于 API 层示例）
func ListPolarDBRecords() ([]PolarDBRecord, error) {
	rows, err := db.Query(
		"SELECT dbcluster_id, cloud_name, engine, region_id, db_cluster_status, dbnode_number, dbcluster_description, memory_size, connection_string FROM polardb",
	)
	if err != nil {
		return nil, fmt.Errorf("查询 PolarDB 表失败: %w", err)
	}
	defer rows.Close()

	var results []PolarDBRecord
	for rows.Next() {
		var rec PolarDBRecord
		// 将查询结果的每一行扫描到 PolarDBRecord 结构体
		err := rows.Scan(&rec.InstanceID, &rec.CloudName, &rec.Engine, &rec.RegionID, &rec.Status,
			&rec.DBNodeCount, &rec.Description, &rec.MemorySize, &rec.ConnectionString)
		if err != nil {
			return nil, fmt.Errorf("读取 PolarDB 行数据失败: %w", err)
		}
		results = append(results, rec)
	}
	return results, nil
}
