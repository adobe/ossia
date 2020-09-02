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
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	"time"
)

// Hypervisor represents OpenStack Hypervisor Resource
//
// swagger:model
type Hypervisor struct {
	// the id for the hypervisor
	//
	// required: true
	ID string `storm:"id"`
	// the hostname for the hypervisor
	//
	// required: true
	Hostname string `storm:"index"`
	// the FQDN of the hypervisor
	//
	// required: true
	FQDN string
	// the status of the hypervisor
	//
	// required: true
	Status string
	// the state of the hypervisor
	//
	// required: true
	State string
	// the hostip of the hypervisor
	//
	// required: true
	HostIP string
	// the vcpus of the hypervisor
	//
	// required: true
	VCPUs int
	// the vcpus_used of the hypervisor
	//
	// required: true
	VCPUsUsed int
	// the free_disk_gb for the hypervisor
	//
	// required: true
	FreeDiskGB int
	// the total_disk_gb for the hypervisor
	//
	// required: true
	TotalDiskGB int
	// the free_ram_mb for the hypervisor
	//
	// required: true
	FreeRAMMB int
	// the total_ram_mb for the hypervisor
	//
	// required: true
	TotalRAMMB int
	// the running_vms on the hypervisor
	//
	// required: false
	RunningVMs int
	// the vms on the hypervisor
	//
	// required: false
	VMs []string
	// OSSIA update time
	//
	// required: true
	PollTime time.Time
}

// HypervisorHash model represents HostID Instance attribute
// swagger:ignore
type HypervisorHash struct {
	Hash      string `storm:"id"`
	ProjectID string
	Hostname  string
}

// EmptyHypervisor model represents OpenStack Hypervisor
// with a limited attributes
type EmptyHypervisor struct {
	// the hostname for the hypervisor
	//
	// required: true
	Hostname string
	// the vcpus of the hypervisor
	//
	// required: true
	VCPUs int
	// the total_ram_mb for the hypervisor
	//
	// required: true
	TotalRAMMB int
}

// Exists method checking if hypervisor is in API Response Slice
func (h *Hypervisor) Exists(hypervisors []hypervisors.Hypervisor) bool {
	for _, v := range hypervisors {
		if v.ID == h.ID {
			return true
		}
	}
	return false

}
