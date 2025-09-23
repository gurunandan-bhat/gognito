package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	defaultConfigFileName = ".gognito.json"
)

type Config struct {
	InProduction bool   `mapstructure:"inProduction"`
	AppRoot      string `mapstructure:"appRoot"`
	AppPort      string `mapstructure:"appPort"`
	Db           struct {
		User                 string `mapstructure:"user"`
		Passwd               string `mapstructure:"passwd"`
		Net                  string `mapstructure:"net"`
		Addr                 string `mapstructure:"addr"`
		DBName               string `mapstructure:"dbName"`
		ParseTime            bool   `mapstructure:"parseTime"`
		Loc                  string `mapstructure:"loc"`
		AllowNativePasswords bool   `mapstructure:"allowNativePasswords"`
	} `mapstructure:"db"`
	Security struct {
		CSRFKey string `mapstructure:"csrfKey"`
	} `mapstructure:"security"`
	Session struct {
		Name              string `mapstructure:"name"`
		Path              string `mapstructure:"path"`
		Domain            string `mapstructure:"domain"`
		MaxAgeHours       int    `mapstructure:"maxAgeHours"`
		AuthenticationKey string `mapstructure:"authenticationKey"`
		EncryptionKey     string `mapstructure:"encryptionKey"`
	} `mapstructure:"session"`
	AWS struct {
		ClientID     string `mapstructure:"clientID"`
		ClientSecret string `mapstructure:"clientSecret"`
		AppDomain    string `mapstructure:"appDomain"`
		IssuerURL    string `mapstructure:"issuerURL"`
		RedirectURL  string `mapstructure:"redirectURL"`
		KeySetURL    string `mapstructure:"keySetURL"`
		LogoutURL    string `mapstructure:"logoutURL"`
		State        string `mapstructure:"state"`
	} `mapstructure:"aws"`
}

var c = Config{}

func Configuration(configFileName ...string) (*Config, error) {

	if (c == Config{}) {

		var cfName string
		switch len(configFileName) {
		case 0:
			dirname, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			cfName = fmt.Sprintf("%s/%s", dirname, defaultConfigFileName)
		case 1:
			cfName = configFileName[0]
		default:
			return nil, fmt.Errorf("incorrect arguments for configuration file name")
		}

		viper.SetConfigFile(cfName)
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}

		if err := viper.Unmarshal(&c); err != nil {
			return nil, err
		}
	}

	return &c, nil
}
