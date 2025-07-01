package config

import (
	"AppointmentSummary_Assignment_Code/model"

	"log"
	"path/filepath"
	"runtime"

	"os"

	"github.com/spf13/viper"
)

var AppConfig *model.Configurations = GetConfigurations()

func GetConfigurations() *model.Configurations {

	var configPath string

	// get current directory path
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(b)
	)

	// where to fetch config file from
	if os.Getenv("GO_ENV") != "local" {
		configPath = "/etc/config"
	} else {
		configPath = basepath + "/data/"
	}

	// add path where to search for config file
	viper.AddConfigPath(configPath)

	// set name of the file to read
	viper.SetConfigName("config")

	// set the file type
	viper.SetConfigType("json")

	// read config from the specified path
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config file error: ", err)
	}

	var configJson model.Configurations

	// unmarshal read configurations into a struct
	err = viper.Unmarshal(&configJson)
	if err != nil {
		log.Fatal("Error unmarshalling config json: ", err)
	}

	return &configJson
}
