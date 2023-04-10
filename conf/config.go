// Package conf parse config to configuration
package conf

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config contains configuration items for sdk
type Config struct {
	IsHTTP           bool
	ChainID          int64
	CAFile           string
	TLSCAContext     []byte
	Key              string
	TLSKeyContext    []byte
	Cert             string
	TLSCertContext   []byte
	IsSMCrypto       bool
	PrivateKey       []byte
	GroupID          int
	NodeURL          string
	MslUrl           string
	MslUsername      string
	MslPasswd        string
	MslName          string
	MslProtocol      string
	LogDBUrl         string
	LogDBUsername    string
	LogDBPasswd      string
	LogDBName        string
	LogDBProtocol    string
	CanalIP          string
	CanalPort        int
	CanalUsername    string
	CanalPassword    string
	CanalDestination string
	CanalConnectedDB string
}

// ParseConfigFile parses the configuration from toml config file
func ParseConfigFile(cfgFile string) ([]Config, error) {
	file, err := os.Open(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("open file failed, err: %v", err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			logrus.Fatalf("close file failed, err: %v", err)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("file is not found, err: %v", err)
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("read file failed, err: %v", err)
	}
	return ParseConfig(buffer)
}

// ParseConfig parses the configuration from []byte
func ParseConfig(buffer []byte) ([]Config, error) {
	viper.SetConfigType("toml")
	viper.SetDefault("SMCrypto", false)
	err := viper.ReadConfig(bytes.NewBuffer(buffer))
	if err != nil {
		return nil, fmt.Errorf("viper .ReadConfig failed, err: %v", err)
	}
	config := new(Config)
	var configs []Config
	if viper.IsSet("Mysql") {
		config.MslUrl = viper.GetString("Mysql.MslUrl")
		config.MslUsername = viper.GetString("Mysql.MslUsername")
		config.MslPasswd = viper.GetString("Mysql.MslPasswd")
		config.MslName = viper.GetString("Mysql.MslName")
		config.MslProtocol = viper.GetString("Mysql.MslProtocol")
	}
	if viper.IsSet("LogDB") {
		config.LogDBUrl = viper.GetString("LogDB.LogDBUrl")
		config.LogDBUsername = viper.GetString("LogDB.LogDBUsername")
		config.LogDBPasswd = viper.GetString("LogDB.LogDBPasswd")
		config.LogDBName = viper.GetString("LogDB.LogDBName")
		config.LogDBProtocol = viper.GetString("LogDB.LogDBProtocol")
	}
	if viper.IsSet("Canal") {
		config.CanalIP = viper.GetString("Canal.CanalIP")
		config.CanalPort = viper.GetInt("Canal.CanalPort")
		config.CanalUsername = viper.GetString("Canal.CanalUsername")
		config.CanalPassword = viper.GetString("Canal.CanalPassword")
		config.CanalDestination = viper.GetString("Canal.CanalDestination")
		config.CanalConnectedDB = viper.GetString("Canal.CanalConnectedDB")
	}
	configs = append(configs, *config)
	return configs, nil
}

// ParseConfigOptions parses from arguments
func ParseConfigOptions(caFile string, key string, cert, keyFile string, groupId int, ipPort string, isHttp bool, chainId int64, isSMCrypto bool) (*Config, error) {
	config := Config{
		IsHTTP:     isHttp,
		ChainID:    chainId,
		CAFile:     caFile,
		Key:        key,
		Cert:       cert,
		IsSMCrypto: isSMCrypto,
		GroupID:    groupId,
		NodeURL:    ipPort,
	}
	keyBytes, curve, err := LoadECPrivateKeyFromPEM(keyFile)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed, err: %v", err)
	}
	if config.IsSMCrypto && curve != sm2p256v1 {
		return nil, fmt.Errorf("smcrypto must use sm2p256v1 private key, but found %s", curve)
	}
	if !config.IsSMCrypto && curve != secp256k1 {
		return nil, fmt.Errorf("must use secp256k1 private key, but found %s", curve)
	}
	config.PrivateKey = keyBytes
	return &config, nil
}
