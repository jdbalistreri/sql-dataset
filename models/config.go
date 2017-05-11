package models

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	MySQLDriver    = "mysql"
	PostgresDriver = "postgres"
	SQLiteDriver   = "sqlite3"
)

var supportedDrivers = []string{MySQLDriver, PostgresDriver, SQLiteDriver}
var errParseConfigFile = "Error occurred parsing the config: %s"

type Config struct {
	GeckoboardAPIKey string          `yaml:"geckoboard_api_key"`
	DatabaseConfig   *DatabaseConfig `yaml:"database_config"`
	RefreshTimeSec   uint16          `yaml:"refresh_time_sec"`
	Datasets         []Dataset       `yaml:"datasets"`
}

// DatabaseConfig holds the db type, url
// and other custom options such as tls config
type DatabaseConfig struct {
	Driver    string            `yaml:"driver"`
	URL       string            `yaml:"-"`
	Host      string            `yaml:"host"`
	Port      string            `yaml:"port"`
	Protocol  string            `yaml:"protocol"`
	Database  string            `yaml:"database"`
	Username  string            `yaml:"username"`
	Password  string            `yaml:"password"`
	TLSConfig *TLSConfig        `yaml:"tls_config"`
	Params    map[string]string `yaml:"params"`
}

type TLSConfig struct {
	KeyFile  string `yaml:"key_file"`
	CertFile string `yaml:"cert_file"`
	CAFile   string `yaml:"ca_file"`
	SSLMode  string `yaml:"ssl_mode"`
}

func LoadConfig(filepath string) (config *Config, err error) {
	var b []byte

	if filepath == "" {
		return nil, errors.New("File path is required to load config")
	}

	if b, err = ioutil.ReadFile(filepath); err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(b, &config); err != nil {
		return nil, fmt.Errorf(errParseConfigFile, err)
	}

	return config, nil
}

func (c Config) Validate() (errors []string) {
	if c.GeckoboardAPIKey == "" {
		errors = append(errors, "Geckoboard api key is required")
	}

	if c.DatabaseConfig == nil {
		errors = append(errors, "Database config is required")
	} else {
		errors = append(errors, c.DatabaseConfig.Validate()...)
	}

	for _, ds := range c.Datasets {
		errors = append(errors, ds.Validate()...)
	}

	return errors
}

func (dc DatabaseConfig) Validate() (errors []string) {
	if dc.Driver == "" {
		errors = append(errors, "Database driver is required")
	} else {
		var matched bool

		for _, d := range supportedDrivers {
			if d == dc.Driver {
				matched = true
				break
			}
		}

		if !matched {
			errors = append(errors, fmt.Sprintf("Unsupported driver '%s' only %s are supported", dc.Driver, supportedDrivers))
		}
	}

	return errors
}
