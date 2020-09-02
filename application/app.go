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
	"ossia/models"
	"ossia/openstack"

	"github.com/gophercloud/gophercloud"
	log "github.com/sirupsen/logrus"
)

// App is application object
type App struct {
	Config *models.Configuration
}

// APIConnection reprsents Identity and Compute
// API Connections. Token is the same in terms of
// one session
type APIConnection struct {
	Nova     *gophercloud.ServiceClient
	Keystone *gophercloud.ServiceClient
	Token    string
}

var (
	//Connection is map of Initialized APIs
	Connection = make(map[string]APIConnection)
)

// Nova establishes Nova API Connection
func Nova(deploymentName string) *gophercloud.ServiceClient {
	deployment := Cfg.Deployments[deploymentName]
	value, ok := Connection[deploymentName]
	var (
		err error
		tmp = Connection[deploymentName]
	)

	// Some CNX for deploymentName initialized
	if ok {
		if value.Nova != nil {
			err = value.Nova.Reauthenticate(value.Token)
			if err != nil {
				log.Error(err)
				return nil
			}
			tmp.Token = Connection[deploymentName].Nova.TokenID
			tmp.Nova = value.Nova
			tmp.Keystone = value.Keystone
			Connection[deploymentName] = tmp
			return Connection[deploymentName].Nova
		}
		cnx := openstack.ComputeConnection(deployment.OsAuthURL, deployment.OsUsername, deployment.OsPassword, deployment.OsProjectName)
		if cnx == nil {
			return nil
		}
		Connection[deploymentName] = APIConnection{
			Nova:     cnx,
			Keystone: nil,
			Token:    "",
		}
		tmp.Nova = Connection[deploymentName].Nova
		tmp.Keystone = value.Keystone
		tmp.Token = Connection[deploymentName].Nova.TokenID
		Connection[deploymentName] = tmp
		return Connection[deploymentName].Nova
	}

	// No CNX for deploymentName. Establish both connections
	log.WithFields(log.Fields{
		"deployment": deploymentName,
		"AuthUrl":    deployment.OsAuthURL,
	}).Debug("Making Compute Connection")

	cnx := openstack.ComputeConnection(deployment.OsAuthURL, deployment.OsUsername, deployment.OsPassword, deployment.OsProjectName)
	if cnx == nil {
		return nil
	}
	Connection[deploymentName] = APIConnection{
		Nova:     cnx,
		Keystone: nil,
		Token:    "",
	}
	tmp.Token = Connection[deploymentName].Nova.TokenID
	tmp.Nova = Connection[deploymentName].Nova
	Connection[deploymentName] = tmp
	return Connection[deploymentName].Nova

}

// Keystone establishes Nova API Connection
func Keystone(deploymentName string) *gophercloud.ServiceClient {
	deployment := Cfg.Deployments[deploymentName]
	value, ok := Connection[deploymentName]
	var (
		err error
		tmp = Connection[deploymentName]
	)

	// Some API's for deploymentName are initialized
	if ok {
		if value.Keystone != nil {
			err = value.Keystone.Reauthenticate(value.Token)
			if err != nil {
				log.Error(err)
				return nil
			}
			tmp.Token = Connection[deploymentName].Keystone.TokenID
			tmp.Nova = value.Nova
			tmp.Keystone = value.Keystone
			Connection[deploymentName] = tmp
			return Connection[deploymentName].Keystone
		}
		cnx := openstack.IdentityConnection(deployment.OsAuthURL, deployment.OsUsername, deployment.OsPassword, deployment.OsProjectName)
		if cnx == nil {
			return nil
		}
		Connection[deploymentName] = APIConnection{
			Keystone: cnx,
			Nova:     nil,
			Token:    "",
		}
		tmp.Nova = value.Nova
		tmp.Keystone = Connection[deploymentName].Keystone
		tmp.Token = Connection[deploymentName].Keystone.TokenID
		Connection[deploymentName] = tmp
		return Connection[deploymentName].Keystone
	}

	// No CNX for deploymentName. Establish both connections
	log.WithFields(log.Fields{
		"deployment": deploymentName,
		"AuthUrl":    deployment.OsAuthURL,
	}).Debug("Making Keystone Connection")

	cnx := openstack.IdentityConnection(deployment.OsAuthURL, deployment.OsUsername, deployment.OsPassword, deployment.OsProjectName)
	if cnx == nil {
		return nil
	}
	Connection[deploymentName] = APIConnection{
		Keystone: cnx,
		Nova:     nil,
		Token:    "",
	}
	tmp.Token = Connection[deploymentName].Keystone.TokenID
	tmp.Keystone = Connection[deploymentName].Keystone
	Connection[deploymentName] = tmp
	return Connection[deploymentName].Keystone

}

// func Keystone(deploymentName string) *gophercloud.ServiceClient {
// 	deployment := Cfg.Deployments[deploymentName]
// 	value, ok := Connection[deploymentName]
// 	var (
// 		err error
// 		tmp = Connection[deploymentName]
// 	)
// 	// API for deploymentName initialized
// 	if ok {
// 		err = value.Keystone.Reauthenticate(value.Token)
// 		if err != nil {
// 			log.Error(err)
// 			return nil
// 		}
// 		tmp.Token = Connection[deploymentName].Nova.TokenID
// 		tmp.Nova = value.Nova
// 		tmp.Keystone = value.Keystone
// 		Connection[deploymentName] = tmp
// 		return Connection[deploymentName].Keystone
// 	}
// 	// No CNX for deploymentName. Establish both connections
// 	log.WithFields(log.Fields{
// 		"deployment": deploymentName,
// 		"AuthUrl":    deployment.OsAuthURL,
// 	}).Debug("Making Identity Connection")

// 	Connection[deploymentName] = ApiConnection{
// 		Nova:     nil,
// 		Keystone: openstack.IdentityConnection(deployment.OsAuthURL, deployment.OsUsername, deployment.OsPassword, deployment.OsProjectName),
// 		Token:    "",
// 	}
// 	tmp.Token = Connection[deploymentName].Keystone.TokenID
// 	tmp.Keystone = Connection[deploymentName].Keystone
// 	Connection[deploymentName] = tmp
// 	return Connection[deploymentName].Keystone

// }

// NewApp is responsible for the core
// functions
func NewApp() *App {

	config := GetConfig()
	setupLogger()
	InitDB(Cfg.Database)

	app := &App{
		Config: config,
	}

	return app
}
