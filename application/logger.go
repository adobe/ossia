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
	"os"

	log "github.com/Sirupsen/logrus"
)

func getLogLevel() log.Level {
	if Cfg.Debug {
		return log.DebugLevel
	}
	return log.InfoLevel
}

func setupLogger() {
	// Log as JSON instead of the default ASCII formatter.
	customFormatter := new(log.TextFormatter)
	//customFormatter.FullTimestamp = true
	customFormatter.TimestampFormat = "15:04:05"

	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(customFormatter)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	if Cfg.LogFile != "" {
		f, err := os.OpenFile(Cfg.LogFile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			fmt.Println("Can't write to log: ", err)
			fmt.Println("Writing to Stdout/Stderr")
		} else {
			log.SetOutput(f)
		}
	} else {
		log.SetOutput(os.Stdout)
	}

	// Only log the warning severity or above.
	log.SetLevel(getLogLevel())

}
