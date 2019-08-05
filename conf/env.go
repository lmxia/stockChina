package conf

import (
	log "github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
)

//EnvConf describe env config
var EnvConf struct {
	Log struct {
		ToStdout bool   `envconfig:"default=false"`
		Level    string `envconfig:"default=INFO"`
		Size     int    `envconfig:"default=50"`
	}
}

//InitEnvConfig init envs
func InitEnvConfig() {
	if err := envconfig.Init(&EnvConf); err != nil {
		log.Warnf("warn=%s\n", err.Error())
	}
}
