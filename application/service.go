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
	stdContext "context"
	"fmt"
	"os"
	"os/signal"
	"ossia/middleware"
	"strings"
	"syscall"
	"time"

	"github.com/iris-contrib/middleware/cors"
	log "github.com/sirupsen/logrus"

	//prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
	"github.com/kataras/iris/v12"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Service is an API service Object
type Service struct {
	app *App
}

// NewService initiliazes service instance
func NewService(app *App) *Service {
	return &Service{
		app: app,
	}
}

// Run start Web Server for the API
func (s *Service) Run() {
	configuration := iris.WithConfiguration(iris.Configuration{
		DisableStartupLog: true,
	})

	irisAddress := iris.Addr(s.app.Config.ListenOn)

	if s.app.Config.AutoTLS.Enabled {
		irisAddress = iris.AutoTLS(
			":"+strings.Split(s.app.Config.ListenOn, ":")[1],
			s.app.Config.AutoTLS.Domain,
			s.app.Config.AutoTLS.AdminEmail,
		)
	}

	// iris.RegisterOnInterrupt(func() {
	// 	log.Info("Stopping API Web Server...")
	// 	timeout := 1 * time.Second
	// 	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
	// 	defer cancel()
	// 	//close all hosts
	// 	s.Engine().Shutdown(ctx)
	// })

	go func() {
		ch := make(chan os.Signal, 1)
		//signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		signal.Notify(ch,
			// kill -SIGINT XXXX or Ctrl+c
			os.Interrupt,
			syscall.SIGINT, // register that too, it should be ok
			// os.Kill  is equivalent with the syscall.Kill
			os.Kill,
			syscall.SIGKILL, // register that too, it should be ok
			// kill -SIGTERM XXXX
			syscall.SIGTERM,
		)
		select {
		case <-ch:
			log.Info("Stopping API Web Server...")
			ctx, cancel := stdContext.WithTimeout(stdContext.Background(), 1*time.Second) //, timeout)
			defer cancel()
			s.Engine().Shutdown(ctx)
			//engine.Shutdown(ctx)
		}
	}()

	if err := s.Engine().Run(
		irisAddress,
		iris.WithOptimizations,
		iris.WithoutServerError(iris.ErrServerClosed),
		// iris.WithoutInterruptHandler,
		configuration,
	); err != nil {
		fmt.Println("start error", err)
	}

}

// Engine is reponsible for routing and middleware
func (s *Service) Engine() *iris.Application {

	engine := iris.New()

	engine.Configure()

	for _, ware := range middleware.Provider {
		engine.UseGlobal(ware)
	}

	//metrics := prometheusMiddleware.New(AppName, 300, 1200, 5000)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
	})

	engine.Use(crs)
	//engine.Use(metrics.ServeHTTP)

	engine.HandleDir("/", AssetFile())
	v1 := engine.Party("/v1")

	// Application Handlers
	engine.Get("/", apiReference)
	//engine.Get("/metrics", iris.FromStd(promhttp.Handler()))
	engine.OnErrorCode(iris.StatusNotFound, notFoundHandler)
	v1.Get("/status", statusHandler)

	// OpenStack Resource View (all resources per deployment)
	v1.Get("/deployments", deploymentsHandler)
	v1.Get("/deployment/{deployment:string}/projects", projectsHandler)
	v1.Get("/deployment/{deployment:string}/images", imagesHandler)
	v1.Get("/deployment/{deployment:string}/flavors", flavorsHandler)
	v1.Get("/deployment/{deployment:string}/aggregates", aggregatesHandler)
	v1.Get("/deployment/{deployment:string}/hypervisors", hypervisorsHandler)
	v1.Get("/deployment/{deployment:string}/hypervisors/empty", hypervisorsEmptyHandler)
	v1.Get("/deployment/{deployment:string}/instances", instancesHandler)
	v1.Get("/deployment/{deployment:string}/project/{project:string}/instances", projectInstancesHandler)
	v1.Get("/deployment/{deployment:string}/instances/clusters", clustersHandler)

	// OpenStack Resource (by resource name)
	v1.Get("/deployment/{deployment:string}", deploymentHandler)
	v1.Get("/deployment/{deployment:string}/snapshots", deploymentSnapshotHandler)
	v1.Get("/deployment/{deployment:string}/project/{project:string}", projectHandler)
	v1.Get("/deployment/{deployment:string}/image/{image:string}", imageHandler)
	v1.Get("/deployment/{deployment:string}/flavor/{flavor:string}", flavorHandler)
	v1.Get("/deployment/{deployment:string}/aggregate/{aggregate:string}", aggregateHandler)
	v1.Get("/deployment/{deployment:string}/instance/{instance:string}", instanceHandler)
	v1.Get("/deployment/{deployment:string}/instances/cluster/{cluster:string}", clusterHandler)
	v1.Get("/deployment/{deployment:string}/instances/filter/{name:string}", filterInstancesHandler)
	v1.Get("/deployment/{deployment:string}/hypervisor/{hostname:string}", hypervisorHandler)

	// Resource Update (POST)
	v1.Post("/deployment/{deployment:string}/update", deploymentUpdateHandler)

	return engine

}
