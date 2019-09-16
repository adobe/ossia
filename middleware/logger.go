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
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kataras/iris/context"
)

func init() {
	Register(func(ctx context.Context) {
		start := time.Now()
		ctx.Next()
		log.WithFields(log.Fields{
			"status_code": strconv.Itoa(ctx.GetStatusCode()),
			"remote_addr": ctx.RemoteAddr(),
			"method":      ctx.Method(),
			"path":        ctx.Path(),
			"duration":    time.Since(start),
		}).Info("External API Call")
	})
}
