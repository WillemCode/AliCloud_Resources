package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	// 注意：Viper 内部已支持 YAML，无需手动导入 yaml.v2/v3 包
)

// 阿里云账户配置结构体
type Account struct {
	Name             string   `yaml:"name" mapstructure:"name"`                             // 账户名称
	AccessKey        string   `yaml:"access_key" mapstructure:"access_key"`                 // 阿里云 AccessKey
	AccessSecret     string   `yaml:"access_secret" mapstructure:"access_secret"`           // 阿里云 AccessSecret
	ECSRegionIds     []string `yaml:"ecs_region_ids" mapstructure:"ecs_region_ids"`         // ECS 服务区域 ID
	RDSRegionIds     []string `yaml:"rds_region_ids" mapstructure:"rds_region_ids"`         // RDS 服务区域 ID
	SLBRegionIds     []string `yaml:"slb_region_ids" mapstructure:"slb_region_ids"`         // SLB 服务区域 ID
	RedisRegionIds   []string `yaml:"redis_region_ids" mapstructure:"redis_region_ids"`     // Tair Redis 服务区域 ID
	PolarDBRegionIds []string `yaml:"polardb_region_ids" mapstructure:"polardb_region_ids"` // PolarDB 服务区域 ID
}

// 数据库配置结构体
type DatabaseConfig struct {
	Path string `yaml:"path"` // 数据库文件路径
}

// 总配置结构体，包含所有配置项
type Config struct {
	AliyunAccounts []Account      `yaml:"aliyun_accounts" mapstructure:"aliyun_accounts"` // 阿里云账户列表
	Database       DatabaseConfig `yaml:"database" mapstructure:"database"`               // 数据库配置
	LogLevel       string         `yaml:"log_level" mapstructure:"log_level"`             // 日志级别
}

// LoadConfig 加载配置文件，并支持环境变量覆盖配置。
// 参数 configPath 可指定配置文件路径，如为空则使用默认路径或环境变量。
func LoadConfig(configPath string) (*Config, error) {
	// 如果未指定路径，检查环境变量 CONFIG_PATH，否则使用默认路径
	if configPath == "" {
		if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
			configPath = envPath
		} else {
			configPath = "./config.yaml" // 默认配置文件路径
		}
	}

	// 设置配置文件路径和类型，然后读取配置
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("无法读取配置文件(%s): %w", configPath, err)
	}

	// 绑定环境变量，以支持通过环境变量覆盖配置值
	// 例如：设置环境变量 DB_PATH 可覆盖配置文件中的 database.path
	_ = viper.BindEnv("database.path", "DB_PATH")
	_ = viper.BindEnv("log_level", "LOG_LEVEL")
	viper.AutomaticEnv() // 启用环境变量自动匹配

	// 反序列化配置到 Config 结构体
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("配置解析失败: %w", err)
	}

	return &cfg, nil
}
