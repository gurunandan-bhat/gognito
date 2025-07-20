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
	InProduction bool   `json:"inProduction,omitempty"`
	AppRoot      string `json:"appRoot,omitempty"`
	AppPort      string `json:"appPort,omitempty"`
	Db           struct {
		User                 string `json:"user,omitempty"`
		Passwd               string `json:"passwd,omitempty"`
		Net                  string `json:"net,omitempty"`
		Addr                 string `json:"addr,omitempty"`
		DBName               string `json:"dbName,omitempty"`
		ParseTime            bool   `json:"parseTime,omitempty"`
		Loc                  string `json:"loc,omitempty"`
		AllowNativePasswords bool   `json:"allowNativePasswords,omitempty"`
	} `json:"db,omitzero"`
	Security struct {
		CSRFKey string `json:"csrfKey,omitempty"`
	} `json:"security,omitzero"`
	Session struct {
		Name              string `json:"name,omitempty"`
		Path              string `json:"path,omitempty"`
		Domain            string `json:"domain,omitempty"`
		MaxAgeHours       int    `json:"maxAgeHours,omitempty"`
		AuthenticationKey string `json:"authenticationKey,omitempty"`
		EncryptionKey     string `json:"encryptionKey,omitempty"`
	} `json:"session,omitzero"`
	AWS struct {
		ClientID     string `json:"clientID,omitempty"`
		ClientSecret string `json:"clientSecret,omitempty"`
		AppDomain    string `json:"appDomain,omitempty"`
		IssuerURL    string `json:"issuerURL,omitempty"`
		RedirectURL  string `json:"redirectURL,omitempty"`
		LogoutURL    string `json:"logoutURL,omitempty"`
		State        string `json:"state,omitempty"`
	} `json:"aws,omitzero"`
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
