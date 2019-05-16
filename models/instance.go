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
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

// Instance represents the OpenStack Instance
//
// swagger:model
type Instance struct {
	// the id for the instance
	//
	// required: true
	ID string `storm:"id"`
	// the name for the instance
	//
	// required: true
	Name string `storm:"index"`
	// the Image id for the instance
	//
	// required: true
	ImageID string
	// the status of the instance
	//
	// required: true
	Status string
	// the HostID for the instance
	//
	// required: true
	HostID string `storm:"index"`
	// the hypervisor for the instance
	//
	// required: false
	Hypervisor string
	// the projectID for the instance
	//
	// required: true
	ProjectID string `storm:"index"`
	// the flavorID for the instance
	//
	// required: true
	Flavor string
	// the Fixed IPv4 for the instance
	//
	// required: true
	FixedIPv4 string
	// the Floating IPv4 for the instance
	//
	// required: true
	FloatingIPv4 string
	// the Fixed IPv6 for the instance
	//
	// required: true
	FixedIPv6 string
	// the Floating IPv6 for the instance
	//
	// required: true
	FloatingIPv6 string
	// the metadata for the instance
	//
	// required: true
	Metadata map[string]string
	// the time of the instance creation
	//
	// required: true
	Created time.Time
	// the security groups
	//
	// required: true
	SecurityGroups []map[string]interface{}
	// the time of the instance modification
	//
	// required: true
	Updated time.Time
	// OSSIA update time
	//
	// required: true
	PollTime time.Time
}

// InstanceNIC is a structured representation of a Gophercloud servers.Server
// virtual NIC.
// swagger:ignore
type InstanceNIC struct {
	FixedIPv4    string // Instance Private IPv4
	FixedIPv6    string // Instance "Private" IPv6
	FloatingIPv4 string // Instance Floating IPv4
	FloatingIPv6 string // Instance Floating IPv6
	MAC          string // MAC Address of the primary Interface
}

// InstanceAddresses is a collection of InstanceNICs, grouped by the
// network name. An instance/server could have multiple NICs on the same
// network.
// swagger:ignore
type InstanceAddresses struct {
	NetworkName  string
	InstanceNICs []InstanceNIC
}

// Exists method checking if instance is in API Response Slice
func (i *Instance) Exists(instances []servers.Server) bool {
	for _, v := range instances {
		if v.ID == i.ID {
			return true
		}
	}
	return false

}

// GetInstanceAddresses metod extracs IPAddresses from the API Instance Object
func GetInstanceAddresses(addresses map[string]interface{}) []InstanceAddresses {
	var allInstanceAddresses []InstanceAddresses

	for networkName, v := range addresses {
		instanceAddresses := InstanceAddresses{
			NetworkName: networkName,
		}

		for _, v := range v.([]interface{}) {
			instanceNIC := InstanceNIC{}
			var exists bool

			v := v.(map[string]interface{})
			if v, ok := v["OS-EXT-IPS-MAC:mac_addr"].(string); ok {
				instanceNIC.MAC = v
			}

			if v["OS-EXT-IPS:type"] == "fixed" {
				switch v["version"].(float64) {
				case 6:
					instanceNIC.FixedIPv6 = fmt.Sprintf("[%s]", v["addr"].(string))
				default:
					instanceNIC.FixedIPv4 = v["addr"].(string)
				}
			} else {
				switch v["version"].(float64) {
				case 6:
					instanceNIC.FloatingIPv6 = fmt.Sprintf("[%s]", v["addr"].(string))
				default:
					instanceNIC.FloatingIPv4 = v["addr"].(string)
				}
			}

			// To associate IPv4 and IPv6 on the right NIC,
			// key on the mac address and fill in the blanks.
			for i, v := range instanceAddresses.InstanceNICs {
				if v.MAC == instanceNIC.MAC {
					exists = true
					if instanceNIC.FixedIPv6 != "" {
						instanceAddresses.InstanceNICs[i].FixedIPv6 = instanceNIC.FixedIPv6
					}
					if instanceNIC.FixedIPv4 != "" {
						instanceAddresses.InstanceNICs[i].FixedIPv4 = instanceNIC.FixedIPv4
					}
					if instanceNIC.FloatingIPv6 != "" {
						instanceAddresses.InstanceNICs[i].FloatingIPv6 = instanceNIC.FloatingIPv6
					}
					if instanceNIC.FloatingIPv4 != "" {
						instanceAddresses.InstanceNICs[i].FloatingIPv4 = instanceNIC.FloatingIPv4
					}
				}
			}

			if !exists {
				instanceAddresses.InstanceNICs = append(instanceAddresses.InstanceNICs, instanceNIC)
			}
		}

		allInstanceAddresses = append(allInstanceAddresses, instanceAddresses)
	}

	return allInstanceAddresses
}
