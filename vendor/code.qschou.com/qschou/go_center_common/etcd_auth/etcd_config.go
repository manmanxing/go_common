package etcd_auth

import (
	"fmt"
	"os"
	"strings"
)

func GetEtcdConfig() (userName, userPass string) {
	v, ok := os.LookupEnv("QSCHOU_CONFIG_CENTER")
	if ok {
		env := strings.SplitN(v, ":", 2)
		if len(env) == 2 {
			fmt.Println("level", "INFO", "etcd_config", "QSCHOU_CONFIG_CENTER", "config_env", env)
			userName = env[0]
			userPass = env[1]
		} else {
			fmt.Println("level", "WARN", "etcd_config", "QSCHOU_CONFIG_CENTER", "config_env", "fail")
		}
	} else {
		fmt.Println("level", "WARN", "etcd_config", "QSCHOU_CONFIG_CENTER", "config_env", "empty")
	}
	return
}
