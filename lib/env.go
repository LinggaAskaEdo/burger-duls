package lib

import (
	"log"

	"github.com/spf13/viper"
)

// Env has environment stored
type Env struct {
	ServerPort      string `mapstructure:"SERVER_PORT"`
	Environment     string `mapstructure:"ENV"`
	LogOutput       string `mapstructure:"LOG_OUTPUT"`
	DBMaxOpenConn   int    `mapstructure:"DB_MAX_OPEN_CONN"`
	DBMaxIdleConn   int    `mapstructure:"DB_MAX_IDLE_CONN"`
	DBMaxLifeTime   int    `mapstructure:"DB_MAX_LIFE_TIME"`
	DBMaxIdleTime   int    `mapstructure:"DB_MAX_IDLE_TIME"`
	URLCloudStorage string `mapstructure:"URL_CLOUD_STORAGE"`
	ProjectId       string `mapstructure:"PROJECT_ID"`
	BucketName      string `mapstructure:"BUCKET_NAME"`
}

// NewEnv creates a new environment
func NewEnv() Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("☠️ cannot read configuration")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("☠️ environment can't be loaded: ", err)
	}

	log.Printf("%+v \n", env)
	return env
}
