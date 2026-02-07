package config

import (
	"os"
	"path/filepath"

	"kefu-server/utils/logger"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Admin AdminConfig `yaml:"admin"`
}

type AdminConfig struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
}

var AppConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	configPath = filepath.Clean(configPath)

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		logger.Errorf("failed to read config file: %v", err)
		return nil, err
	}

	// 解析配置
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		logger.Errorf("failed to parse config file: %v", err)
		return nil, err
	}

	AppConfig = &config
	logger.Infof("config loaded successfully: %+v", config)
	return &config, nil
}

// GetConfig 获取配置
func GetConfig() *Config {
	return AppConfig
}
