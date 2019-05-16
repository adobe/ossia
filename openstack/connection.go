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

package openstack

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

func initOpenStackProvider(IdentityEndpoint string, Username string, Password string, TenantName string) *gophercloud.ProviderClient {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: IdentityEndpoint,
		Username:         Username,
		Password:         Password,
		TenantName:       TenantName,
	}

	opts.DomainName = "default"
	opts.AllowReauth = true

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.WithFields(log.Fields{
			//"URL":   IdentityEndpoint,
			"Error": err,
		}).Error("Could not create OpenStack Provider.")
	}
	return provider
}

// ComputeConnection initializes Nova API Connection
func ComputeConnection(IdentityEndpoint string, Username string, Password string, TenantName string) *gophercloud.ServiceClient {

	provider := initOpenStackProvider(IdentityEndpoint, Username, Password, TenantName)
	if provider == nil {
		//need to fetch complete error at provider initialisation
		return nil
	}
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to decode into config struct.")
	}
	return client

}

// IdentityConnection initializes Keystone API Connection
func IdentityConnection(IdentityEndpoint string, Username string, Password string, TenantName string) *gophercloud.ServiceClient {

	provider := initOpenStackProvider(IdentityEndpoint, Username, Password, TenantName)
	if provider == nil {
		//need to fetch complete error at provider initialisation
		return nil
	}
	client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to decode into config struct.")
	}
	return client

}
