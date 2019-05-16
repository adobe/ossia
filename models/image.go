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

	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
)

// Image represents the OpenStack Glance Image
//
// swagger:model
type Image struct {
	// the id for the image
	//
	// required: true
	ID string `storm:"id"`
	// the name for the image
	//
	// required: true
	Name string `storm:"index"`
	// the status of the image
	//
	// required: true
	Status string
	// the time of the image creation
	//
	// required: true
	Created string
	// the time of the image modification
	//
	// required: true
	Updated string
	// the metadata of the image
	//
	// required: true
	Metadata map[string]interface{}
	// the image usage (VMs)
	//
	// required: false
	UsedBy []string
	// OSSIA update time
	//
	// required: true
	PollTime time.Time
}

// Exists method checking if image is in API Response Slice
func (i *Image) Exists(images []images.Image) bool {
	for _, v := range images {
		if v.ID == i.ID {
			return true
		}
	}
	return false

}
