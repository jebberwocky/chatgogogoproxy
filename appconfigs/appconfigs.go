package appconfigs

import (
	"chatproxy/models"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
)

var Appconfigs models.Apps

func init() {
	//init app config
	content, err := os.ReadFile("./apps.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		panic(err)
	}
	var val models.Apps
	if err := json.Unmarshal(content, &val); err != nil {
		log.Fatal("Error when parsing appconfig file: ", err)
		panic(err)
	}
	Appconfigs = val
}

func GetContextByHost(host string) models.AppContext {
	for i := range Appconfigs.Apps {
		appc := Appconfigs.Apps[i]
		if appc.Hostname == host {
			return appc
		}
	}
	return Appconfigs.Apps[0] //return default
}
