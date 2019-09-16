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

// OpenStack Simple Inventory API
//
// Ossia [osˈsiːa] is a musical term for an alternative passage which may be played instead of the original passage.
// Ossia passages are very common in opera and solo-piano works. They are usually an easier version of the preferred form of passage.
//
//     Schemes: http
//     BasePath: /v1
//     Version: 1.0.2
//     License: Apache License 2.0 http://www.apache.org/licenses/LICENSE-2.0
//     Contact: Mykola Mogylenko <mykola@adobe.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"ossia/application"
	"ossia/scheduler"

	log "github.com/Sirupsen/logrus"
)

// Reload is used for the call-back from
// configuration file changes
// var Reload bool

func main() {
	defer scheduler.Stop()
	defer application.CloseDB()

	// Init Application
	app := application.NewApp()

	// // Init Service
	service := application.NewService(app)
	defer service.Run()

	for deployment := range app.Config.Deployments {

		log.WithFields(log.Fields{
			"deployment": deployment,
		}).Info("Registered new deployment")

		go application.UpdateInventory(
			deployment,
			app.Config.PollInterval,
		)

	}

	go application.DataStoreMetrics()

	scheduler.Run()

}
