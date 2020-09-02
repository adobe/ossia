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
	"ossia/utils"
	"strings"
	"time"

	"github.com/asdine/storm"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	log "github.com/sirupsen/logrus"
)

func updateInstances(deployment string) {
	defer utils.TimeTrack(time.Now(), updateInstances)

	var inventoryInstances []models.Instance
	var hash models.HypervisorHash

	bucket := DB.From(deployment)

	err := bucket.All(&inventoryInstances)
	if err != nil {
		log.Error(err)
	}

	cnx := Nova(deployment)
	if cnx == nil {
		log.WithFields(log.Fields{
			"deployment": deployment,
			"task":       "instances",
		}).Error("Nothing to update for the deployment. No OpenStack Connectivity")
	} else {
		log.WithFields(log.Fields{
			"deployment": deployment,
		}).Info("Updating Instances for the deployment")

		allPages, err := servers.List(cnx, servers.ListOpts{AllTenants: true}).AllPages()

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Unable to fetch Instances from OpenStack")
		} else {
			instances, err := servers.ExtractServers(allPages)
			if err != nil {
				log.Error(err)
			}

			for _, i := range instances {

				// networks := []map[string]interface{}{}
				// for _, instanceAddresses := range models.GetInstanceAddresses(i.Addresses) {
				// 	for _, instanceNIC := range instanceAddresses.InstanceNICs {
				// 		v := map[string]interface{}{
				// 			"name":           instanceAddresses.NetworkName,
				// 			"fixed_ip_v4":    instanceNIC.FixedIPv4,
				// 			"fixed_ip_v6":    instanceNIC.FixedIPv6,
				// 			"floating_ip_v4": instanceNIC.FloatingIPv4,
				// 			"floating_ip_v6": instanceNIC.FloatingIPv6,
				// 			"mac":            instanceNIC.MAC,
				// 		}
				// 		networks = append(networks, v)
				// 	}
				// }
				networks := models.GetInstanceAddresses(i.Addresses)

				_ = bucket.One("Hash", i.HostID, &hash)

				var FixedIPv4, FloatingIPv4, FixedIPv6, FloatingIPv6 string

				if len(networks) > 0 {
					FixedIPv4 = networks[0].InstanceNICs[0].FixedIPv4
					FloatingIPv4 = networks[0].InstanceNICs[0].FloatingIPv4
					FixedIPv6 = networks[0].InstanceNICs[0].FixedIPv6
					FloatingIPv6 = networks[0].InstanceNICs[0].FloatingIPv6
				}

				inst := &models.Instance{
					ID:     i.ID,
					Status: i.Status,
					//StatusMessage: i.Fault.Message,
					Name:           i.Name,
					HostID:         i.HostID,
					ProjectID:      i.TenantID,
					ImageID:        i.Image["id"].(string),
					Flavor:         i.Flavor["id"].(string),
					FixedIPv4:      FixedIPv4,
					FloatingIPv4:   FloatingIPv4,
					FixedIPv6:      FixedIPv6,
					FloatingIPv6:   FloatingIPv6,
					Hypervisor:     hash.Hostname,
					Metadata:       i.Metadata,
					Created:        i.Created,
					SecurityGroups: i.SecurityGroups,
					Updated:        i.Updated,
					PollTime:       time.Now(),
				}
				bucket.Save(inst)
				//updateOrSave("ID", inst.ID, inst, bucket)

			}

			// Instances DB Cleanup
			for _, i := range inventoryInstances {
				if !i.Exists(instances) {
					log.Debug("Deleting instance ", i.ID)
					err := bucket.DeleteStruct(&i)
					if err != nil {
						log.Error(err)
					}
				}
			}

		}

	}

}

func updateImages(deployment string) {
	defer utils.TimeTrack(time.Now(), updateImages)

	var inventoryImages []models.Image
	var inventoryInstances []models.Instance
	var err error

	bucket := DB.From(deployment)

	//Get all Images from inventory
	err = bucket.All(&inventoryImages)
	if err != nil {
		log.Error(err)
	}

	//Get all Instances from inventory
	err = bucket.All(&inventoryInstances)
	if err != nil {
		log.Error(err)
	}

	cnx := Nova(deployment)
	if cnx == nil {
		log.WithFields(log.Fields{
			"deployment": deployment,
			"task":       "images",
		}).Error("Nothing to update for the deployment. No OpenStack Connectivity")
	} else {
		log.WithFields(log.Fields{
			"deployment": deployment,
		}).Info("Updating Images for the deployment")
		allPages, err := images.ListDetail(cnx, images.ListOpts{}).AllPages()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Unable to fetch Images from OpenStack")
		} else {
			images, err := images.ExtractImages(allPages)
			if err != nil {
				log.Error(err)
			}

			for _, i := range images {
				img := &models.Image{
					ID:       i.ID,
					Name:     i.Name,
					Status:   i.Status,
					Metadata: i.Metadata,
					Created:  i.Created,
					Updated:  i.Updated,
					PollTime: time.Now(),
				}

				for _, i := range inventoryInstances {
					if i.ImageID == img.ID {
						img.UsedBy = append(img.UsedBy, i.Name)
					}
				}

				bucket.Save(img)
				//updateOrSave("ID", img.ID, img, bucket)
			}

			// Images DB Cleanup
			for _, i := range inventoryImages {
				if !i.Exists(images) {
					log.Debug("Deleting image ", i.ID)
					err := bucket.DeleteStruct(&i)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}

}

func updateHypervisors(deployment string) {
	defer utils.TimeTrack(time.Now(), updateHypervisors)

	var inventoryProjects []models.Project
	var inventoryHypervisors []models.Hypervisor

	bucket := DB.From(deployment)

	err := bucket.All(&inventoryHypervisors)
	if err != nil {
		log.Error(err)
	}

	cnx := Nova(deployment)
	if cnx == nil {
		log.WithFields(log.Fields{
			"deployment": deployment,
			"task":       "hypervisors",
		}).Error("Nothing to update for the deployment. No OpenStack Connectivity")
	} else {
		log.WithFields(log.Fields{
			"deployment": deployment,
		}).Info("Updating Hypervisors for the deployment")

		allPages, err := hypervisors.List(cnx).AllPages()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Unable to fetch Hypervisors from OpenStack")
		} else {
			hypervisors, err := hypervisors.ExtractHypervisors(allPages)
			if err != nil {
				log.Error(err)
			}

			for _, h := range hypervisors {
				hyp := &models.Hypervisor{
					ID:          h.ID,
					Hostname:    strings.Split(h.HypervisorHostname, ".")[0],
					FQDN:        h.HypervisorHostname,
					Status:      h.Status,
					State:       h.State,
					HostIP:      h.HostIP,
					VCPUs:       h.VCPUs,
					VCPUsUsed:   h.VCPUsUsed,
					FreeDiskGB:  h.FreeDiskGB,
					TotalDiskGB: h.LocalGB,
					FreeRAMMB:   h.FreeRamMB,
					TotalRAMMB:  h.MemoryMB,
					RunningVMs:  h.RunningVMs,
					PollTime:    time.Now(),
				}
				//updateOrSave("ID", hyp.ID, hyp, bucket)
				bucket.Save(hyp)
			}

			// Hashes
			err = bucket.All(&inventoryProjects)
			if err != nil {
				log.Error(err)
			}

			for _, p := range inventoryProjects {
				for _, h := range hypervisors {
					hostname := strings.Split(h.HypervisorHostname, ".")[0]
					hh := &models.HypervisorHash{
						Hash:      utils.HashHypervisor(p.ID, hostname),
						Hostname:  hostname,
						ProjectID: p.ID,
					}

					err := bucket.Save(hh)
					if err != nil {
						log.Error(err)
					}

				}
			}

			// Hypervisors DB Cleanup
			for _, i := range inventoryHypervisors {
				if !i.Exists(hypervisors) {
					log.Debug("Deleting hypervisor ", i.ID)
					err := bucket.DeleteStruct(&i)
					if err != nil {
						log.Error(err)
					}
				}
			}

		}
	}

}

func updateFlavors(deployment string) {
	defer utils.TimeTrack(time.Now(), updateFlavors)

	var inventoryFlavors []models.Flavor

	bucket := DB.From(deployment)

	//Get all Flavors from inventory
	err := bucket.All(&inventoryFlavors)
	if err != nil {
		log.Error(err)
	}

	cnx := Nova(deployment)
	if cnx == nil {
		log.WithFields(log.Fields{
			"deployment": deployment,
			"task":       "flavors",
		}).Error("Nothing to update for the deployment. No OpenStack Connectivity")
	} else {
		log.WithFields(log.Fields{
			"deployment": deployment,
		}).Info("Updating Flavors for the deployment")
		allPages, err := flavors.ListDetail(cnx, flavors.ListOpts{}).AllPages()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Unable to fetch Flavors from OpenStack")
		} else {
			flavors, err := flavors.ExtractFlavors(allPages)
			if err != nil {
				log.Error(err)
			}

			for _, f := range flavors {
				flv := &models.Flavor{
					ID:         f.ID,
					Name:       f.Name,
					RAM:        f.RAM,
					VCPUs:      f.VCPUs,
					Disk:       f.Disk,
					Swap:       f.Swap,
					RxTxFactor: f.RxTxFactor,
					IsPublic:   f.IsPublic,
					Ephemeral:  f.Ephemeral,
					PollTime:   time.Now(),
				}

				bucket.Save(flv)
				//updateOrSave("ID", flv.ID, flv, bucket)
			}

			for _, i := range inventoryFlavors {
				if !i.Exists(flavors) {
					log.Debug("Deleting flavor ", i.ID)
					err := bucket.DeleteStruct(&i)
					if err != nil {
						log.Error(err)
					}
				}
			}

		}
	}

}

func updateProjects(deployment string) {
	defer utils.TimeTrack(time.Now(), updateProjects)

	var inventoryProjects []models.Project

	bucket := DB.From(deployment)

	//Get all Projects from inventory
	err := bucket.All(&inventoryProjects)
	if err != nil {
		log.Error(err)
	}

	cnx := Keystone(deployment)
	if cnx == nil {
		log.WithFields(log.Fields{
			"deployment": deployment,
			"task":       "projects",
		}).Error("Nothing to update for the deployment. No OpenStack Connectivity")
	} else {
		log.WithFields(log.Fields{
			"deployment": deployment,
		}).Info("Updating Projects for the deployment")
		allPages, err := projects.List(cnx, projects.ListOpts{}).AllPages()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Unable to fetch Projects from OpenStack")
		} else {
			projects, err := projects.ExtractProjects(allPages)
			if err != nil {
				log.Error(err)
			}

			for _, p := range projects {
				prj := &models.Project{
					ID:          p.ID,
					Name:        p.Name,
					Enabled:     p.Enabled,
					Description: p.Description,
					PollTime:    time.Now(),
				}
				bucket.Save(prj)
				//updateOrSave("ID", prj.ID, prj, bucket)
			}

			// Projects DB Cleanup
			for _, i := range inventoryProjects {
				if !i.Exists(projects) {
					log.Debug("Deleting project ", i.ID)
					err := bucket.DeleteStruct(&i)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}

}

func updateAggregates(deployment string) {
	defer utils.TimeTrack(time.Now(), updateAggregates)

	var inventoryAggregates []models.Aggregate

	bucket := DB.From(deployment)

	//Get all Aggregates from inventory
	err := bucket.All(&inventoryAggregates)
	if err != nil {
		log.Error(err)
	}

	cnx := Nova(deployment)
	if cnx == nil {
		log.WithFields(log.Fields{
			"deployment": deployment,
			"task":       "aggregates",
		}).Error("Nothing to update for the deployment. No OpenStack Connectivity")
	} else {
		log.WithFields(log.Fields{
			"deployment": deployment,
		}).Info("Updating Aggregates for the deployment")

		allPages, err := aggregates.List(cnx).AllPages()

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Unable to fetch Aggregates from OpenStack")
		} else {
			aggregates, err := aggregates.ExtractAggregates(allPages)
			if err != nil {
				log.Error(err)
			}
			for _, a := range aggregates {
				if !a.Deleted {
					agg := &models.Aggregate{
						ID:               a.ID,
						Name:             a.Name,
						AvailabilityZone: a.AvailabilityZone,
						Hosts:            a.Hosts,
						Metadata:         a.Metadata,
						Created:          a.CreatedAt,
						Updated:          a.UpdatedAt,
						PollTime:         time.Now(),
					}
					bucket.Save(agg)
				}
			}

			// Aggregates DB Cleanup
			for _, a := range inventoryAggregates {
				if !a.Exists(aggregates) {
					log.Debug("Deleting ", a.Name, " aggregate")
					err := bucket.DeleteStruct(&a)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}

	}

}

// deploymentRegistered - check if deployment from DB is defined
// in configuration file
func deploymentRegistered(deployment string) bool {
	for i := range Cfg.Deployments {
		if i == deployment {
			return true
		}
	}
	return false
}

// projectExists - check if OpenStack project exists
// in the database
func projectExists(deployment string, projectName string) bool {
	var project models.Project

	bucket := DB.From(deployment)

	err := bucket.One("Name", projectName, &project)
	if err != nil {
		log.WithFields(log.Fields{"function": "projectExists"}).Debug(err)
		return false
	}
	return true
}

// dbCleanup - purge not defined deployments
func dbCleanup() {
	defer utils.TimeTrack(time.Now(), dbCleanup)

	log.Debug("Analyzing deployments in the DB")
	deployments := DB.PrefixScan("")
	for _, i := range deployments {
		deployment := i.Bucket()[0]
		if !strings.Contains(string(deployment), "storm") {
			if !deploymentRegistered(deployment) {
				log.WithFields(log.Fields{
					"deployment": deployment,
				}).Warn("Deployment is not registered, deleting")
				err := DB.Drop(deployment)
				if err != nil {
					log.Error(err)
				}
			}
		}

	}

}

// NewInstance copy of Instance model
type NewInstance struct{ models.Instance }

// Hypervisor method adds Hypervisor for the Instance
// based on its HostID Hash
func (i *NewInstance) Hypervisor(deployment string) string {
	var hyphash models.HypervisorHash
	bucket := DB.From(deployment)
	err := bucket.One("Hash", i.HostID, &hyphash)
	if err != nil {
		if err == storm.ErrNotFound {
			log.WithFields(log.Fields{
				"deployment": deployment,
				"function":   "Hypervisor",
			}).Warn("Hash not found")
		} else {
			log.Error(err)
		}
	}
	return hyphash.Hostname

}

// NewHypervisor copy of Hypervisor model
type NewHypervisor struct{ models.Hypervisor }

// Instances method adds Instances for the Hypervisor
// based on its Hash
func (i *NewHypervisor) Instances(deployment string) []models.Instance {
	var hashes []models.HypervisorHash
	var instances []models.Instance

	bucket := DB.From(deployment)
	err := bucket.Find("Hostname", i.Hostname, &hashes)
	if err != nil {
		if err == storm.ErrNotFound {
			log.WithFields(log.Fields{
				"deployment": deployment,
				"function":   "Instances",
			}).Warn("Hash not found")
		} else {
			log.Error(err)
		}
	}
	for _, h := range hashes {
		err := bucket.Find("HostID", h.Hash, &instances)
		if err != nil {
			if err != storm.ErrNotFound {
				log.Error(err)
			}
		}

	}
	return instances

}

// usageSnapshot - Represent OpenStack Resources Utilization
func usageSnapshot(deployment string) {
	defer utils.TimeTrack(time.Now(), usageSnapshot)

	var (
		instances    []models.Instance
		hypervisors  []models.Hypervisor
		projects     []models.Project
		flavors      []models.Flavor
		images       []models.Image
		err          error
		VCPUs        int
		VCPUsUsed    int
		MemoryMB     int
		FreeMemoryMB int
	)

	bucket := DB.From(deployment)

	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Updating Usage Snapshot for the deployment")

	err = bucket.All(&instances)
	if err != nil {
		log.Error(err)
	}

	err = bucket.All(&hypervisors)
	if err != nil {
		log.Error(err)
	}

	err = bucket.All(&projects)
	if err != nil {
		log.Error(err)
	}

	err = bucket.All(&flavors)
	if err != nil {
		log.Error(err)
	}

	err = bucket.All(&images)
	if err != nil {
		log.Error(err)
	}

	for _, h := range hypervisors {
		// Let's make it more accurate than OS does
		if h.Status == "enabled" && h.State == "up" {
			VCPUs += h.VCPUs
			VCPUsUsed += h.VCPUsUsed
			MemoryMB += h.TotalRAMMB
			FreeMemoryMB += h.FreeRAMMB
		}
	}

	snapshot := &models.Snapshot{
		ID:           time.Now().Format("2006-01-02"),
		Flavors:      len(flavors),
		Hypervisors:  len(hypervisors),
		Images:       len(images),
		Instances:    len(instances),
		Projects:     len(projects),
		VCPUs:        VCPUs,
		VCPUsUsed:    VCPUsUsed,
		MemoryMB:     MemoryMB,
		MemoryUsedMB: MemoryMB - FreeMemoryMB,
	}
	bucket.Save(snapshot)

}
