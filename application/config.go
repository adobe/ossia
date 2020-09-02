/*
Copyright 2019 Adobe. All rights reserved.
This file is licensed to you under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License. You may obtain a copy
of the License at http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
OF ANY KIND, either express or implied. See the License for the specific language
governing permissions and limitations under the License.
*/

package application

import (
	"fmt"

	"ossia/models"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Cfg is global config object variable
var Cfg *models.Configuration

// GetConfig is returning config object
func GetConfig() *models.Configuration {
	viper.SetConfigType("yaml")
	viper.SetConfigName("conf")

	viper.AddConfigPath("./etc")
	viper.AddConfigPath("/etc/ossia/")
	viper.AddConfigPath("/opt/ossia/etc/")

	err := viper.ReadInConfig()

	viper.WatchConfig()

	if err != nil {
		log.Fatal("%v", err)
	}

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.WithFields(log.Fields{
			"file": e.Name,
		}).Warn("Config file changed. Please do a manual restart for now.")
	})

	return Cfg
}
