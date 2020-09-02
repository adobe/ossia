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

	"github.com/kataras/iris/v12"
)

// Default Handlers

func defaultHandler(c iris.Context) {
	c.JSON(iris.Map{"message": "ok"})

}

func notFoundHandler(c iris.Context) {
	c.JSON(iris.Map{
		"message": "Path not found.",
	})

}

// OpenStack Resource View handlers implemenation

// deploymentsHandler represents registered OpenStack deployments view
// swagger:operation GET /deployments resources listDeployments
//
// Registered OpenStack Deployments
//
// Returns all registered OpenStack Deployments
//
// ---
// responses:
//   '200':
//     description: Available Deployments
//     schema:
//       type: object
//       properties:
//         deployments:
//           description: List of registered deployments
//           type: array
//           items:
//             type: string
//   '501':
//      description: Application Issue
func deploymentsHandler(c iris.Context) {
	c.JSON(iris.Map{
		"deployments": listDeployments(),
	})
}

// projectsHandler represents OpenStack projects view
// swagger:operation GET /deployment/{deployment}/projects resources listProjects
//
// OpenStack Projects
//
// Returns all projects for the deployment
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
// responses:
//   '200':
//     description: "List of OpenStack Projects"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         projects:
//           description: List of registered deployments
//           type: array
//           items:
//             $ref: '#/definitions/Project'
//   '404':
//     description: "Returns 404 Code if there is no deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func projectsHandler(c iris.Context) {

	deployment := c.Params().Get("deployment")

	c.StatusCode(iris.StatusNotFound)
	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}

	if deploymentRegistered(deployment) {
		projects := listProjects(deployment)
		c.StatusCode(iris.StatusOK)
		response = iris.Map{
			"deployment": deployment,
			"projects":   projects,
		}
	}

	c.JSON(response)
}

// instancesHandler represents OpenStack instances view
// swagger:operation GET /deployment/{deployment}/instances resources listInstances
//
// OpenStack Instances
//
// Returns all instances for the deployment
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
// responses:
//   '200':
//     description: "List of OpenStack Instances"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         instances:
//           description: List of instances
//           type: array
//           items:
//             $ref: '#/definitions/Instance'
//   '404':
//     description: "Returns 404 Code if there is no deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func instancesHandler(c iris.Context) {

	deployment := c.Params().Get("deployment")

	c.StatusCode(iris.StatusNotFound)
	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}

	if deploymentRegistered(deployment) {
		instances := listInstances(deployment, "")
		c.StatusCode(iris.StatusOK)
		response = iris.Map{
			"deployment": deployment,
			"instances":  instances,
		}
	}
	c.JSON(response)
}

// projectInstancesHandler represents OpenStack instances view for
// a particular project
// swagger:operation GET /deployment/{deployment}/project/{project}/instances resources filterInstancesByProject
//
// OpenStack Instances by Project
//
// Returns all instances filtered by Project
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
//  - name: project
//    in: path
//    description: OpenStack Project Name
//    type: string
//    required: true
//    example: rtb
// responses:
//   '200':
//     description: "List of OpenStack Instances"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         "project:instances":
//           description: List of instances by Project
//           type: array
//           items:
//             $ref: '#/definitions/Instance'
//   '404':
//     description: "Returns 404 Code if there is no deployment or project"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func projectInstancesHandler(c iris.Context) {

	deployment := c.Params().Get("deployment")
	project := c.Params().Get("project")

	c.StatusCode(iris.StatusNotFound)
	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}

	if deploymentRegistered(deployment) {
		if projectExists(deployment, project) {
			instances := filterInstancesByProject(deployment, project)
			c.StatusCode(iris.StatusOK)
			response = iris.Map{
				"deployment":                         deployment,
				fmt.Sprintf("%s:instances", project): instances,
			}
		} else {
			c.StatusCode(iris.StatusNotFound)
			response = iris.Map{
				"message": fmt.Sprintf("Project %s not found", project),
			}

		}

	}
	c.JSON(response)
}

// instancesHandler represents OpenStack instances view
// swagger:operation GET /deployment/{deployment}/instances/filter/{name} resources filterInstances
//
// OpenStack Instances by Name
//
// Returns all instances fitlered by instance name
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
//  - name: name
//    in: path
//    description: Instance Name Contains
//    type: string
//    required: true
//    example: instance
// responses:
//   '200':
//     description: "List of OpenStack Instances"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         instances:
//           description: List of instances
//           type: array
//           items:
//             $ref: '#/definitions/Instance'
//   '404':
//     description: "Returns 404 Code if there is no deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func filterInstancesHandler(c iris.Context) {

	deployment := c.Params().Get("deployment")
	name := c.Params().Get("name")

	c.StatusCode(iris.StatusNotFound)
	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}
	if deploymentRegistered(deployment) {
		instances := listInstances(deployment, name)

		response = iris.Map{
			"deployment": deployment,
			"instances":  instances,
		}
		c.StatusCode(iris.StatusOK)
	}

	c.JSON(response)
}

// clustersHandler represents OpenStack instances view
// filtered by metadata [cluster] key
// swagger:operation GET /deployment/{deployment}/instances/clusters resources listClusters
//
// OpenStack Instances Count by Metadata Cluster key
//
// Returns instances count by Metadata Cluster key
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
// responses:
//   '200':
//     description: "OpenStack Instances Count by Metadata Cluster key"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         clusters:
//           description: Cluster Names
//           type: object
//           properties:
//             cluster_name:
//               description: Instances count by cluster name
//               type: integer
//   '404':
//     description: "Returns 404 Code if there is no deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func clustersHandler(c iris.Context) {

	deployment := c.Params().Get("deployment")

	c.StatusCode(iris.StatusNotFound)
	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}

	if deploymentRegistered(deployment) {

		instances := listInstances(deployment, "")
		clusters := make(map[string]int, len(instances))
		for _, i := range instances {
			_, ok := clusters[i.Metadata["cluster"]]
			if ok {
				clusters[i.Metadata["cluster"]]++
			} else {
				clusters[i.Metadata["cluster"]] = 1
			}
		}

		c.StatusCode(iris.StatusOK)
		response = iris.Map{
			"deployment": deployment,
			"clusters":   clusters,
		}
	}

	c.JSON(response)
}

// imagesHandler represents OpenStack images view
// swagger:operation GET /deployment/{deployment}/images resources listImages
//
// OpenStack Images
//
// Returns all images for the deployment
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
// responses:
//   '200':
//     description: "List of OpenStack Images"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         images:
//           description: list of images
//           type: array
//           items:
//             $ref: '#/definitions/Image'
//   '404':
//     description: "Returns 404 Code if there is no deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func imagesHandler(c iris.Context) {
	deployment := c.Params().Get("deployment")

	c.StatusCode(iris.StatusNotFound)
	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}

	if deploymentRegistered(deployment) {
		images := listImages(deployment)

		response = iris.Map{
			"deployment": deployment,
			"images":     images,
		}
		c.StatusCode(iris.StatusOK)
	}
	c.JSON(response)
}

// hypervisorsHandler represents OpenStack hypervisors view
// swagger:operation GET /deployment/{deployment}/hypervisors resources listHypervisors
//
// OpenStack Hypervisors
//
// Returns all hypervisors for the deployment
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
// responses:
//   '200':
//     description: "List of OpenStack Hypervisors"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         hypervisors:
//           description: list of hypervisors
//           type: array
//           items:
//             $ref: '#/definitions/Hypervisor'
//   '404':
//     description: "Returns 404 Code if there is no deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func hypervisorsHandler(c iris.Context) {
	deployment := c.Params().Get("deployment")

	c.StatusCode(iris.StatusNotFound)
	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}

	if deploymentRegistered(deployment) {
		hypervisors := listHypervisors(deployment)

		response = iris.Map{
			"deployment":  deployment,
			"hypervisors": hypervisors,
		}
		c.StatusCode(iris.StatusOK)
	}
	c.JSON(response)
}

// hypervisorsEmptyHandler returns Empty OpenStack Hypervisors
// swagger:operation GET /deployment/{deployment}/hypervisors/empty resources listEmptyHypervisors
//
// OpenStack Empty Hypervisors
//
// Returns OpenStack Empty Hypervisors
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
// responses:
//   '200':
//     description: "List of OpenStack Hypervisors"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         hypervisors:
//           description: list of hypervisors
//           type: array
//           items:
//             $ref: '#/definitions/EmptyHypervisor'
//   '404':
//     description: "Returns 404 Code if there is no deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func hypervisorsEmptyHandler(c iris.Context) {

	deployment := c.Params().Get("deployment")

	c.StatusCode(iris.StatusNotFound)
	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}

	if deploymentRegistered(deployment) {
		//var emptyHypervisors []models.EmptyHypervisor
		var emptyHypervisors []string
		hypervisors := listEmptyHypervisors(deployment)

		for _, h := range hypervisors {
			//emptyHypervisors = append(emptyHypervisors, models.EmptyHypervisor{Hostname: h.Hostname, VCPUs: h.VCPUs, TotalRAMMB: h.TotalRAMMB})
			emptyHypervisors = append(emptyHypervisors, h.Hostname)
		}

		response = iris.Map{
			"deployment":        deployment,
			"total_empty":       len(emptyHypervisors),
			"empty_hypervisors": emptyHypervisors,
		}
		c.StatusCode(iris.StatusOK)
	}

	c.JSON(response)

}

// flavorsHandler represents OpenStack flavors view
// swagger:operation GET /deployment/{deployment}/flavors resources listFlavors
//
// OpenStack Flavors
//
// Returns all flavors for the deployment
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
// responses:
//   '200':
//     description: "List of OpenStack Flavors"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         flavors:
//           description: list of flavors
//           type: array
//           items:
//             $ref: '#/definitions/Flavor'
//   '404':
//     description: "Returns 404 Code if there is no deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func flavorsHandler(c iris.Context) {
	deployment := c.Params().Get("deployment")

	c.StatusCode(iris.StatusNotFound)
	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}

	if deploymentRegistered(deployment) {
		flavors := listFlavors(deployment)

		response = iris.Map{
			"deployment": deployment,
			"flavors":    flavors,
		}
		c.StatusCode(iris.StatusOK)
	}

	c.JSON(response)
}

// aggregatesHandler represents OpenStack aggregates view
// swagger:operation GET /deployment/{deployment}/aggregates resources listAggregates
//
// OpenStack Aggregates
//
// Returns all aggregates for the deployment
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
// responses:
//   '200':
//     description: "List of OpenStack Aggregates"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         aggregates:
//           description: list of aggregates
//           type: array
//           items:
//             $ref: '#/definitions/Aggregate'
//   '404':
//     description: "Returns 404 Code if there is no deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func aggregatesHandler(c iris.Context) {
	deployment := c.Params().Get("deployment")

	c.StatusCode(iris.StatusNotFound)
	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}

	if deploymentRegistered(deployment) {
		aggregates := listAggregates(deployment)

		response = iris.Map{
			"deployment": deployment,
			"aggregates": aggregates,
		}
		c.StatusCode(iris.StatusOK)
	}

	c.JSON(response)
}

// OpenStack Resource handlers implemenation (by resource name)

// deploymentHandler returns statistics about particular deployment
// Rework this in order to not iterate over hashes
// swagger:operation GET /deployment/{deployment} resources getDeployment
//
// OpenStack Deployment Statistics
//
// Returns OpenStack Deployment Statistics
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
// responses:
//   '200':
//     description: "Returns deployment statistics"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         projects:
//           description: Number of Projects for the deployment
//           type: integer
//         hypervisors:
//           description: Number of Hypervisors for the deployment
//           type: integer
//         instances:
//           description: Number of Instances for the deployment
//           type: integer
//         flavors:
//           description: Number of Flavors for the deployment
//           type: integer
//         images:
//           description: Number of Flavors for the deployment
//           type: integer
//   '404':
//     description: "Returns 404 Code if there is no deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func deploymentHandler(c iris.Context) {

	deployment := c.Params().Get("deployment")

	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}
	c.StatusCode(iris.StatusNotFound)
	for _deployment := range Cfg.Deployments {

		if _deployment == deployment {
			images := listImages(deployment)
			flavors := listFlavors(deployment)
			projects := listProjects(deployment)
			instances := listInstances(deployment, "")
			hypervisors := listHypervisors(deployment)

			response = iris.Map{
				"images":      len(images),
				"flavors":     len(flavors),
				"projects":    len(projects),
				"instances":   len(instances),
				"deployment":  deployment,
				"hypervisors": len(hypervisors),
			}
			c.StatusCode(iris.StatusOK)
		}

	}

	c.JSON(response)

}

// deploymentSnapshotHandler returns usage snapshots for particular deployment
// swagger:operation GET /deployment/{deployment}/snapshots resources getSnapshots
//
// OpenStack Deployment Usage Snapshots
//
// Returns OpenStack Deployment Usage Snapshots
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
// responses:
//   '200':
//     description: "Returns usage snapshots"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         usage_snapshots:
//           description: Number of Projects for the deployment
//           type: array
//           items:
//             $ref: '#/definitions/Snapshot'
//   '404':
//     description: "Returns 404 Code if there is no snapshots or registered deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func deploymentSnapshotHandler(c iris.Context) {

	deployment := c.Params().Get("deployment")

	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}
	c.StatusCode(iris.StatusNotFound)

	for _deployment := range Cfg.Deployments {

		if _deployment == deployment {
			snapshots, err := getSnapshots(deployment)
			if err != nil {
				response = iris.Map{
					"message": fmt.Sprintf("No Usage Snapshots for the %s deployment", deployment),
				}
				c.StatusCode(iris.StatusNotFound)
			} else {
				response = iris.Map{
					"usage_snapshots": snapshots,
					"deployment":      deployment,
				}
			}
			c.StatusCode(iris.StatusOK)
		}

	}

	c.JSON(response)
}

// deploymentUpdateHandler triggers deployment Update
// swagger:operation POST /deployment/{deployment}/update resources updateDeployment
//
// Update OpenStack Deployment
//
// Returns 200 on succcess
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
// responses:
//   '200':
//     description: "Returns 200 on success"
//     schema:
//       type: object
//       properties:
//         message:
//           description: Success Message
//           type: string
//   '404':
//     description: "Returns 404 Code if there is no deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func deploymentUpdateHandler(c iris.Context) {

	c.StatusCode(iris.StatusNotImplemented)

	response := iris.Map{
		"message": fmt.Sprintf("Method %s not implemented", c.Method()),
	}
	if c.Method() == "POST" {
		deployment := c.Params().Get("deployment")

		response = iris.Map{
			"message": fmt.Sprintf("Deployment %s not found", deployment),
		}
		c.StatusCode(iris.StatusNotFound)
		for _deployment := range Cfg.Deployments {

			if _deployment == deployment {
				go updateDeployment(deployment)
				c.StatusCode(iris.StatusOK)
				response = iris.Map{
					"message": fmt.Sprintf("Triggered update for %s deployment", deployment),
				}

			}
		}

	}
	c.JSON(response)

}

// imageHandler returns OpenStack Image Object
// swagger:operation GET /deployment/{deployment}/image/{image} resources getImage
//
// OpenStack Image
//
// Returns OpenStack Image Object
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
//  - name: image
//    in: path
//    description: OpenStack Image Name
//    type: string
//    required: true
//    example: ubuntu-18.04-x86_64
// responses:
//   '200':
//     description: "Returns OpenStack Image"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         "image:image_name":
//           $ref: '#/definitions/Image'
//   '404':
//     description: "Returns 404 Code if there is no deployment or image"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func imageHandler(c iris.Context) {
	deployment := c.Params().Get("deployment")
	imageName := c.Params().Get("image")

	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}
	c.StatusCode(iris.StatusNotFound)

	if deploymentRegistered(deployment) {
		image, err := getImage(deployment, imageName)
		if err != nil {
			if err.Error() == "not found" {
				response = iris.Map{"message": fmt.Sprintf("Image %s not found", imageName)}
				c.StatusCode(iris.StatusNotFound)
			} else {
				c.StatusCode(iris.StatusInternalServerError)
				response = iris.Map{"message": err.Error()}
			}

		} else {
			response = iris.Map{
				"deployment":                        deployment,
				fmt.Sprintf("image:%s", image.Name): image,
			}
			c.StatusCode(iris.StatusOK)
		}
	}
	c.JSON(response)

}

// projectHandler returns OpenStack Project Object
// swagger:operation GET /deployment/{deployment}/project/{project} resources getProject
//
// OpenStack Project
//
// Returns OpenStack Project Object
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
//  - name: project
//    in: path
//    description: OpenStack Project Name
//    type: string
//    required: true
//    example: admin
// responses:
//   '200':
//     description: "Returns OpenStack Project"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         "project:project_name":
//           $ref: '#/definitions/Project'
//   '404':
//     description: "Returns 404 Code if there is no deployment or project"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func projectHandler(c iris.Context) {
	deployment := c.Params().Get("deployment")
	projectName := c.Params().Get("project")

	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}
	c.StatusCode(iris.StatusNotFound)

	if deploymentRegistered(deployment) {

		project, err := getProject(deployment, projectName)
		if err != nil {
			if err.Error() == "not found" {
				response = iris.Map{"message": fmt.Sprintf("Project %s not found", projectName)}
				c.StatusCode(iris.StatusNotFound)
			} else {
				c.StatusCode(iris.StatusInternalServerError)
				response = iris.Map{"message": err.Error()}
			}

		} else {
			c.StatusCode(iris.StatusOK)
			response = iris.Map{
				"deployment":                            deployment,
				fmt.Sprintf("project:%s", project.Name): project,
			}
		}

	}
	c.JSON(response)
}

// flavorHandler returns OpenStack Flavor Object
// swagger:operation GET /deployment/{deployment}/flavor/{flavor} resources getFlavor
//
// OpenStack Flavor
//
// Returns OpenStack Flavor Object
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
//  - name: flavor
//    in: path
//    description: OpenStack Flavor Name
//    type: string
//    required: true
//    example: m1.small
// responses:
//   '200':
//     description: "Returns OpenStack Flavor"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         "flavor:flavor_name":
//           $ref: '#/definitions/Flavor'
//   '404':
//     description: "Returns 404 Code if there is no deployment or flavor"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func flavorHandler(c iris.Context) {
	deployment := c.Params().Get("deployment")
	flavorName := c.Params().Get("flavor")

	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}
	c.StatusCode(iris.StatusNotFound)

	if deploymentRegistered(deployment) {
		flavor, err := getFlavor(deployment, flavorName)

		if err != nil {
			if err.Error() == "not found" {
				response = iris.Map{"message": fmt.Sprintf("Flavor %s not found", flavorName)}
				c.StatusCode(iris.StatusNotFound)
			} else {
				c.StatusCode(iris.StatusInternalServerError)
				response = iris.Map{"message": err.Error()}
			}
		} else {
			c.StatusCode(iris.StatusOK)
			response = iris.Map{
				"deployment":                          deployment,
				fmt.Sprintf("flavor:%s", flavor.Name): flavor,
			}
		}
	}
	c.JSON(response)
}

// aggregateHandler returns OpenStack Aggregate Object
// swagger:operation GET /deployment/{deployment}/aggregate/{aggregate} resources getAggregate
//
// OpenStack Aggregate
//
// Returns OpenStack Aggregate Object
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
//  - name: aggregate
//    in: path
//    description: OpenStack Aggregate Name
//    type: string
//    required: true
//    example: staging
// responses:
//   '200':
//     description: "Returns OpenStack Aggregate"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         "aggregate:aggregate_name":
//           $ref: '#/definitions/Aggregate'
//   '404':
//     description: "Returns 404 Code if there is no deployment or aggregate"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func aggregateHandler(c iris.Context) {
	deployment := c.Params().Get("deployment")
	aggregateName := c.Params().Get("aggregate")

	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}
	c.StatusCode(iris.StatusNotFound)

	if deploymentRegistered(deployment) {

		aggregate, err := getAggregate(deployment, aggregateName)

		if err != nil {
			if err.Error() == "not found" {
				response = iris.Map{"message": fmt.Sprintf("Aggregate %s not found", aggregateName)}
				c.StatusCode(iris.StatusNotFound)
			} else {
				c.StatusCode(iris.StatusInternalServerError)
				response = iris.Map{"message": err.Error()}
			}
		} else {
			c.StatusCode(iris.StatusOK)
			response = iris.Map{
				"deployment": deployment,
				fmt.Sprintf("aggregate:%s", aggregate.Name): aggregate,
			}
		}
	}
	c.JSON(response)
}

// hypervisorHandler returns OpenStack Hypervisor Object
// swagger:operation GET /deployment/{deployment}/hypervisor/{hypervisor} resources getHypervisor
//
// OpenStack Hypervisor
//
// Returns OpenStack Hypervisor Object
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
//  - name: hypervisor
//    in: path
//    description: OpenStack Hypervisor Name
//    type: string
//    required: true
//    example: cn07-1a
// responses:
//   '200':
//     description: "Returns OpenStack Hypervisor"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         "hypervisor:hypervisor_name":
//           $ref: '#/definitions/Hypervisor'
//   '404':
//     description: "Returns 404 Code if there is no deployment or hypervisor"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func hypervisorHandler(c iris.Context) {
	deployment := c.Params().Get("deployment")
	hostname := c.Params().Get("hostname")

	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}
	c.StatusCode(iris.StatusNotFound)

	if deploymentRegistered(deployment) {

		hypervisor, err := getHypervisor(deployment, hostname)

		if err != nil {
			if err.Error() == "not found" {
				response = iris.Map{"message": fmt.Sprintf("Hypervisor %s not found", hostname)}
				c.StatusCode(iris.StatusNotFound)
			} else {
				c.StatusCode(iris.StatusInternalServerError)
				response = iris.Map{"message": err.Error()}
			}
		} else {
			c.StatusCode(iris.StatusOK)
			response = iris.Map{
				"deployment":                           deployment,
				fmt.Sprintf("hypervisor:%s", hostname): hypervisor,
			}
		}
	}
	c.JSON(response)

}

// instanceHandler returns OpenStack Instance Object
// swagger:operation GET /deployment/{deployment}/instance/{instance} resources getInstance
//
// OpenStack Instance
//
// Returns OpenStack Instance Object
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
//  - name: instance
//    in: path
//    description: OpenStack Instance Name
//    type: string
//    required: true
//    example: instance01
// responses:
//   '200':
//     description: "Returns OpenStack Instance"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         "instance:instance_name":
//           $ref: '#/definitions/Instance'
//   '404':
//     description: "Returns 404 Code if there is no deployment or instance"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func instanceHandler(c iris.Context) {
	deployment := c.Params().Get("deployment")
	instanceName := c.Params().Get("instance")

	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}
	c.StatusCode(iris.StatusNotFound)

	if deploymentRegistered(deployment) {
		instance, err := getInstance(deployment, instanceName)

		if err != nil {
			if err.Error() == "not found" {
				response = iris.Map{"message": fmt.Sprintf("Instance %s not found", instanceName)}
				c.StatusCode(iris.StatusNotFound)
			} else {
				c.StatusCode(iris.StatusInternalServerError)
				response = iris.Map{"message": err.Error()}
			}
		} else {
			c.StatusCode(iris.StatusOK)
			response = iris.Map{
				"deployment":                             deployment,
				fmt.Sprintf("instance:%s", instanceName): instance,
			}
		}
	}
	c.JSON(response)

}

// clusterHandler returns OpenStack Instances
// filtered by Metadata Cluster Key
// swagger:operation GET /deployment/{deployment}/instances/cluster/{cluster} resources listClusters
//
// OpenStack Instances Count by Metadata Cluster key
//
// Returns OpenStack Instances Count by Metadata Cluster key
//
// ---
// parameters:
//  - name: deployment
//    in: path
//    description: OpenStack Deployment Name
//    type: string
//    required: true
//    example: tm-lab-1a
//  - name: cluster
//    in: path
//    description: Cluster Name
//    type: string
//    required: true
//    example: instance
// responses:
//   '200':
//     description: "Returns OpenStack Instances by cluster name"
//     schema:
//       type: object
//       properties:
//         deployment:
//           description: Name of the deployment
//           type: string
//         "cluster:cluster_name":
//           type: object
//           properties:
//             instance_name_with_ip:
//               description: "Instance name and IP"
//               type: string
//   '404':
//     description: "Returns 404 Code if there is no deployment"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func clusterHandler(c iris.Context) {

	deployment := c.Params().Get("deployment")
	cluster := c.Params().Get("cluster")

	response := iris.Map{
		"message": fmt.Sprintf("Deployment %s not found", deployment),
	}
	c.StatusCode(iris.StatusNotFound)

	if deploymentRegistered(deployment) {
		c.StatusCode(iris.StatusOK)
		instances := listInstances(deployment, "")
		members := make(map[string]string)

		for _, i := range instances {
			if i.Metadata["cluster"] == cluster {
				members[i.Name] = i.FixedIPv4
			}
		}
		response = iris.Map{
			"deployment":                       deployment,
			fmt.Sprintf("cluster:%s", cluster): members,
		}
	}
	c.JSON(response)
}

// Admin and Monitoring Handlers

// statusHandler returns application health status
// swagger:operation GET /status application getStatus
//
// OSSIA Operational Status
//
// Used for the monitoring purposes
//
// ---
// produces:
//  - application/json
// responses:
//   '200':
//     description: "Returns OSSIA Status"
//     schema:
//       type: object
//       properties:
//         status:
//           description: Current Status
//           type: string
//         datastore:
//           description: DB Statistics
//           type: object
//           properties:
//             metric:
//               type: integer
//   '404':
//     description: "Returns 404 Code if there is no path"
//     schema:
//       type: object
//       properties:
//         message:
//           type: string
//           description: Error Message
func statusHandler(c iris.Context) {
	c.JSON(iris.Map{
		"status":    "alive",
		"datastore": datastoreMetrics,
	})

}
