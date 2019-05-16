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
	"fmt"
	"ossia/models"
	"ossia/scheduler"
	"ossia/utils"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/asdine/storm"
)

// UpdateInventory executes initial API calls to
// OpenStack API to prepopulate the Inventory Database
// In addition, it also schedules periodic tasks based on
// application cofiguration file
func UpdateInventory(deployment string, pollinterval models.PollInterval) {
	defer utils.TimeTrack(time.Now(), UpdateInventory)

	updateProjects(deployment)
	updateAggregates(deployment)
	updateHypervisors(deployment)
	updateImages(deployment)
	updateFlavors(deployment)
	updateInstances(deployment)
	usageSnapshot(deployment)
	dbCleanup()

	// Need to re-work this awful code
	scheduler.AddTask(
		fmt.Sprintf("@every %s", pollinterval.Projects),
		func() { updateProjects(deployment) },
		log.Fields{"task": "projects", "deployment": deployment},
	)
	scheduler.AddTask(
		fmt.Sprintf("@every %s", pollinterval.Aggregates),
		func() { updateAggregates(deployment) },
		log.Fields{"task": "aggregates", "deployment": deployment},
	)
	scheduler.AddTask(
		fmt.Sprintf("@every %s", pollinterval.Hypervisors),
		func() { updateHypervisors(deployment) },
		log.Fields{"task": "hypervisors", "deployment": deployment},
	)
	scheduler.AddTask(
		fmt.Sprintf("@every %s", pollinterval.Images),
		func() { updateImages(deployment) },
		log.Fields{"task": "images", "deployment": deployment},
	)
	scheduler.AddTask(
		fmt.Sprintf("@every %s", pollinterval.Flavors),
		func() { updateFlavors(deployment) },
		log.Fields{"task": "flavors", "deployment": deployment},
	)
	scheduler.AddTask(
		fmt.Sprintf("@every %s", pollinterval.Instances),
		func() { updateInstances(deployment) },
		log.Fields{"task": "instances", "deployment": deployment},
	)
	scheduler.AddTask(
		fmt.Sprintf("@every 24h"),
		func() { usageSnapshot(deployment) },
		log.Fields{"task": "usageSnapshot", "deployment": deployment},
	)
	scheduler.AddTask(
		"@every 24h",
		func() { dbCleanup() },
		log.Fields{"task": "dbCleanup"},
	)

}

// updateDeployment executes initial API calls to
// OpenStack API to fetch information about resources
func updateDeployment(deployment string) {
	defer utils.TimeTrack(time.Now(), updateDeployment)

	updateProjects(deployment)
	updateHypervisors(deployment)
	updateImages(deployment)
	updateFlavors(deployment)
	updateInstances(deployment)
	updateAggregates(deployment)
	dbCleanup()
}

// OpenStack Resources view methods implemenation
// Returns list of resources for the deployment

// listDeployments method returns list of
// registered OpenStack deployments
func listDeployments() []string {
	defer utils.TimeTrack(time.Now(), listDeployments)

	var deployments []string
	log.Info("Fetching registered deployments")

	buckets := DB.PrefixScan("")
	for _, i := range buckets {
		deployment := i.Bucket()[0]
		if !strings.Contains(string(deployment), "storm") {
			deployments = append(deployments, deployment)
		}

	}
	return deployments
}

// listProjects method returns list of projects
// for the deployment
func listProjects(deployment string) []models.Project {
	defer utils.TimeTrack(time.Now(), listProjects)

	var projects []models.Project
	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching inventory Projects for the deployment")
	err := bucket.All(&projects)
	if err != nil {
		log.Error(err)
	}
	return projects
}

// listInstances method returns list of instances
// for the deployment
func listInstances(deployment string, hostname string) []models.Instance {
	defer utils.TimeTrack(time.Now(), listInstances)

	var instances []models.Instance

	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching inventory Instances for the deployment")

	err := bucket.All(&instances)
	if err != nil {
		log.Error(err)
	}

	// Filter Instances
	if len(hostname) > 0 {
		for i := 0; i < len(instances); i++ {
			if !strings.Contains(instances[i].Name, hostname) {
				instances = append(instances[:i], instances[i+1:]...)
				i--
			}
		}
	}

	return instances
}

// filterInstancesByProject method returns list of instances
// filteted by Project ID
func filterInstancesByProject(deployment string, projectName string) []models.Instance {
	defer utils.TimeTrack(time.Now(), listInstances)

	var instances []models.Instance
	var project models.Project
	var err error

	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
		"project":    projectName,
	}).Info("Fetching inventory Instances for the deployment")

	if len(projectName) > 0 {
		err = bucket.One("Name", projectName, &project)
		if err != nil {
			log.WithFields(log.Fields{
				"deployment": deployment,
				"project":    projectName,
			}).Error(err)
		} else {
			err = bucket.All(&instances)
			if err != nil {
				log.Error(err)
			}
			for i := 0; i < len(instances); i++ {
				if instances[i].ProjectID != project.ID {
					instances = append(instances[:i], instances[i+1:]...)
					i--
				}
			}
			return instances
		}

	}
	return instances
}

// listImages method returns list of images
// for the deployment
func listImages(deployment string) []models.Image {
	defer utils.TimeTrack(time.Now(), listImages)

	var images []models.Image
	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching inventory Images for the deployment")
	err := bucket.All(&images)
	if err != nil {
		log.Error(err)
	}
	return images
}

// listHypervisors method returns list of hypervisors
// for the deployment
func listHypervisors(deployment string) []models.Hypervisor {
	defer utils.TimeTrack(time.Now(), listHypervisors)
	var hypervisors []models.Hypervisor
	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching inventory Hypervisors for the deployment")
	err := bucket.All(&hypervisors)
	if err != nil {
		log.Error(err)
	}
	for e, p := range hypervisors {
		newHypervisor := &NewHypervisor{p}
		for _, i := range newHypervisor.Instances(deployment) {
			hypervisors[e].VMs = append(hypervisors[e].VMs, i.Name)
		}
		//hypervisors[e].Instances = newHypervisor.Instances(deployment)
	}
	return hypervisors
}

// listEmptyHypervisors method returns list of empty hypervisors
// for the deployment
func listEmptyHypervisors(deployment string) []models.Hypervisor {
	defer utils.TimeTrack(time.Now(), listEmptyHypervisors)
	var hypervisors, emptyHypervisors []models.Hypervisor
	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching inventory Hypervisors for the deployment")
	err := bucket.All(&hypervisors)
	if err != nil {
		log.Error(err)
	}
	for e, p := range hypervisors {
		newHypervisor := &NewHypervisor{p}
		if len(newHypervisor.Instances(deployment)) == 0 {
			emptyHypervisors = append(emptyHypervisors, hypervisors[e])
		}
	}
	return emptyHypervisors
}

// listFlavors method returns list of flavors
// for the deployment
func listFlavors(deployment string) []models.Flavor {
	defer utils.TimeTrack(time.Now(), listFlavors)

	var flavors []models.Flavor
	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching inventory Flavors for the deployment")
	err := bucket.All(&flavors)
	if err != nil {
		log.Error(err)
	}
	return flavors
}

// listAggregates method returns list of aggregates
// for the deployment
func listAggregates(deployment string) []models.Aggregate {
	defer utils.TimeTrack(time.Now(), listAggregates)

	var aggregates []models.Aggregate
	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching inventory Aggregates for the deployment")
	err := bucket.All(&aggregates)
	if err != nil {
		log.Error(err)
	}
	return aggregates
}

// OpenStack Resource methods implemenation (by resource name)
// Returns resource by its name (or 404)

// getProject method returns Inventory Project Object
func getProject(deployment string, name string) (models.Project, error) {
	var project models.Project
	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching Inventory Project for the deployment")
	err := bucket.One("Name", name, &project)
	if err != nil {
		if err == storm.ErrNotFound {
			return project, err
		}
		log.Error(err)
	}
	return project, nil
}

// getInstance method returns Inventory Instance Object
func getInstance(deployment string, hostname string) (models.Instance, error) {
	var instance models.Instance
	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
		"instance":   hostname,
	}).Info("Fetching inventory Instance for the deployment")

	err := bucket.One("Name", hostname, &instance)
	if err != nil {
		if err == storm.ErrNotFound {
			return instance, err
		}
		log.Error(err)
	}
	return instance, nil
}

// getHypervisor method returns Inventory Hypervisor Object
func getHypervisor(deployment string, hostname string) (models.Hypervisor, error) {
	var hypervisor models.Hypervisor

	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching inventory Hypervisor for the deployment")

	err := bucket.One("Hostname", hostname, &hypervisor)
	newHypervisor := &NewHypervisor{hypervisor}
	if err != nil {
		if err == storm.ErrNotFound {
			return hypervisor, err
		}
		log.Error(err)
	}
	for _, i := range newHypervisor.Instances(deployment) {
		hypervisor.VMs = append(hypervisor.VMs, i.Name)
	}
	return hypervisor, nil
}

// getImage method returns Inventory Image Object
func getImage(deployment string, name string) (models.Image, error) {
	var image models.Image
	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching Inventory Image for the deployment")
	err := bucket.One("Name", name, &image)
	if err != nil {
		if err == storm.ErrNotFound {
			return image, err
		}
		log.Error(err)
	}
	return image, nil
}

// getFlavor method returns Inventory Flavor Object
func getFlavor(deployment string, name string) (models.Flavor, error) {
	var flavor models.Flavor
	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching Inventory Flavor for the deployment")
	err := bucket.One("Name", name, &flavor)
	if err != nil {
		if err == storm.ErrNotFound {
			return flavor, err
		}
		log.Error(err)
	}
	return flavor, nil
}

// getAggregate method returns Inventory Aggregate Object
func getAggregate(deployment string, name string) (models.Aggregate, error) {
	var aggregate models.Aggregate
	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching Inventory Aggregate for the deployment")
	err := bucket.One("Name", name, &aggregate)
	if err != nil {
		if err == storm.ErrNotFound {
			return aggregate, err
		}
		log.Error(err)
	}
	return aggregate, nil
}

// getSnapshot method returns OpenStack Usage Snapshots per Deployment
//func getSnapshots(deployment string) ([]models.Snapshot, error) {
func getSnapshots(deployment string) (map[string]interface{}, error) {
	var (
		snapshots []models.Snapshot
	)

	bucket := DB.From(deployment)
	log.WithFields(log.Fields{
		"deployment": deployment,
	}).Info("Fetching Inventory Usage Snapshot for the deployment")
	err := bucket.All(&snapshots)
	if err != nil {
		log.Error(err)
	}

	x := make(map[string]interface{})

	for _, s := range snapshots {
		x[s.ID] = s.Public()
	}
	return x, nil
}
