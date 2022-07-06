package cfg

import (
	"aws-helper/log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	AWS        AWS        `mapstructure:"aws"`
	Cloudflare Cloudflare `mapstructure:"cloudflare"`
}

type AWS struct {
	AccessKeyID     string `mapstructure:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
	DnsName         string `mapstructure:"dnsName"`
	NodeName        string `mapstructure:"nodeName"`
	Region          string `mapstructure:"region"`
}

type Cloudflare struct {
	ApiKey    string `mapstructure:"apiKey"`
	ZoneId    string `mapstructure:"zoneId"`
	KvId      string `mapstructure:"kvId"`
	AccountId string `mapstructure:"accountId"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")

	var cc Config
	// optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("create config file")
			viper.SafeWriteConfig()
		} else {
			// Config file was found but another error was produced
		}

		return nil, errors.Wrap(err, "load config")
	}

	if err := viper.Unmarshal(&cc); err != nil {
		return nil, errors.Wrap(err, "unmarshal file")
	}
	return &cc, err
}
