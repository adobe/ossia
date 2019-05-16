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

import (
	"time"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

// Project represents the OpenStack Project (Tenant)
//
// swagger:model
type Project struct {
	// the id for the project
	//
	// required: true
	ID string `storm:"index"`
	// the name for the project
	//
	// required: true
	Name string
	// the status for the project
	//
	// required: true
	Enabled bool
	// the description for the project
	//
	// required: true
	Description string
	// OSSIA update time
	//
	// required: true
	PollTime time.Time
}

// Exists method checking if project is in API Response Slice
func (p *Project) Exists(projects []projects.Project) bool {
	for _, v := range projects {
		if v.ID == p.ID {
			return true
		}
	}
	return false

}
