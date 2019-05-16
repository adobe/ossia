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

package models

// Configuration is the Global Object for the config file
type Configuration struct {
	ListenOn     string                `mapstructure:"listen_on"`
	Debug        bool                  `mapstructure:"debug"`
	Database     string                `mapstructure:"database"`
	LogFile      string                `mapstructure:"logfile"`
	PollInterval PollInterval          `mapstructure:"poll_interval"`
	AutoTLS      AutoTLS               `mapstructure:"auto_tls"`
	Deployments  map[string]Deployment `mapstructure:"deployments"`
}

// Deployment stanza representation (OpenStack Credentials)
type Deployment struct {
	OsAuthURL     string `mapstructure:"os_auth_url"`
	OsProjectName string `mapstructure:"os_project_name"`
	OsUsername    string `mapstructure:"os_username"`
	OsPassword    string `mapstructure:"os_password"`
}

// PollInterval intervals for the API calls.
// Used for the periodic tasks
type PollInterval struct {
	Images      string `mapstructure:"images"`
	Flavors     string `mapstructure:"flavors"`
	Projects    string `maptstructure:"projects"`
	Instances   string `maptstructure:"instances"`
	Aggregates  string `maptstructure:"aggregates"`
	Hypervisors string `maptstructure:"hypervisors"`
}

// AutoTLS is used for Let's Encrypt integration
type AutoTLS struct {
	Enabled    bool   `mapstructure:"enabled"`
	Domain     string `mapstructure:"domain"`
	AdminEmail string `mapstructure:"admin_email"`
}
