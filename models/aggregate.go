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

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
)

// Aggregate represents the OpenStack Aggregate
//
// swagger:model
type Aggregate struct {
	// the id for the aggregate
	//
	// required: true
	ID int `storm:"id"`
	// the name for the aggregate
	//
	// required: true
	Name string `storm:"index"`
	// the AvailabilityZone for the aggregate
	//
	// required: true
	AvailabilityZone string
	// the Hosts for the aggregate
	//
	// required: true
	Hosts []string
	// the metadata for the aggregate
	//
	// required: true
	Metadata map[string]string
	// the time of the aggregate creation
	//
	// required: true
	Created time.Time
	// the time of the aggregate modification
	//
	// required: true
	Updated time.Time
	// OSSIA update time
	//
	// required: true
	PollTime time.Time
}

// Exists method checking if flavor is in API Response Slice
func (a *Aggregate) Exists(aggregates []aggregates.Aggregate) bool {
	for _, v := range aggregates {
		if v.ID == a.ID {
			return true
		}
	}
	return false

}
