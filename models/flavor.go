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

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

// Flavor represents the OpenStack Flavor
//
// swagger:model
type Flavor struct {
	// the id for the flavor
	//
	// required: true
	ID string `storm:"id"`
	// the name for the flavor
	//
	// required: true
	Name string `storm:"index"`
	// memory for the flavor
	//
	// required: true
	RAM int
	// vpcus for the flavor
	//
	// required: true
	VCPUs int
	// disk for the flavor
	//
	// required: true
	Disk int
	// swap for the flavor
	//
	// required: true
	Swap int
	// the rxtxfactor for the flavor
	//
	// required: true
	RxTxFactor float64
	// visibility for the flavor
	//
	// required: true
	IsPublic bool
	// ephemeral for the flavor
	//
	// required: true
	Ephemeral int
	// OSSIA update time
	//
	// required: true
	PollTime time.Time
}

// Exists method checking if flavor is in API Response Slice
func (f *Flavor) Exists(flavors []flavors.Flavor) bool {
	for _, v := range flavors {
		if v.ID == f.ID {
			return true
		}
	}
	return false

}
