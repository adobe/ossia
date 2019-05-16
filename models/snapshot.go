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

// Snapshot represents OpenStack Usage Snapshot
//
// swagger:model
type Snapshot struct {
	// the id for the snapshot
	//
	// required: true
	ID string `storm:"id", json:"-"`
	// Amount of Flavors
	//
	// required: true
	Flavors int
	// Amount of Hypervisors
	//
	// required: true
	Hypervisors int
	// Amount of Images
	//
	// required: true
	Images int
	// Amount of Instances
	//
	// required: true
	Instances int
	// Amount of Projects
	//
	// required: true
	Projects int
	// Total VCPUs
	//
	// required: true
	VCPUs int
	// vCPU Usage
	//
	// required: true
	VCPUsUsed int
	// Total Memory
	//
	// required: true
	MemoryMB int
	// Memory Usage
	//
	// required: true
	MemoryUsedMB int
}

// Public method checking if project is in API Response Slice
//func (s *Snapshot) Public() map[string]int {
func (s *Snapshot) Public() interface{} {
	return map[string]int{
		"Flavors":      s.Flavors,
		"Hypervisors":  s.Hypervisors,
		"Images":       s.Images,
		"Instances":    s.Instances,
		"Projects":     s.Projects,
		"VCPUs":        s.VCPUs,
		"VCPUsUsed":    s.VCPUsUsed,
		"MemoryMB":     s.MemoryMB,
		"MemoryUsedMB": s.MemoryUsedMB,
	}

}
