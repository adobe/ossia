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

package middleware

import (
	"github.com/kataras/iris/context"
)

var serverName string

// UpdateHeader sets server header according
// to global name and global version
func UpdateServerHeader(header string) {
	serverName = header
}

//fmt.Sprintf("%s/%s", appName, appVersion)
func init() {
	Register(func(ctx context.Context) {
		ctx.Header("Server", serverName)
		ctx.Next()
	})
}
