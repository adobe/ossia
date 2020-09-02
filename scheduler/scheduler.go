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

package scheduler

import (
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

var scheduler *cron.Cron

type funcName func()

func init() {
	scheduler = cron.New()
}

// func AddTask(schedule string, task funcName, message string) error {
// 	log.WithFields(log.Fields{
// 		"schedule": schedule,
// 		"message":  message,
// 	}).Warn("Scheduled the new task")
// 	return scheduler.AddFunc(schedule, task)
// }

// AddTask registers a new task
func AddTask(schedule string, task funcName, log_message log.Fields) error {
	log_message["schedule"] = schedule
	log.WithFields(log_message).Info("Scheduled a new task")
	return scheduler.AddFunc(schedule, task)
}

// Run starts scheduler
func Run() {
	//defer utils.TimeTrack(time.Now(), "Scheduler")
	scheduler.Start()
	//select {}
}

// Entries implements view for the scheduled tasks
func Entries() []*cron.Entry {
	return scheduler.Entries()
}

// Stop implements scheduler termination
func Stop() {
	//defer utils.TimeTrack(time.Now(), "Scheduler")
	log.Info("Stopping Scheduler...")
	scheduler.Stop()
	//select {}
}

// Clean implements Tasks cleanup
func Clean() {
	//defer utils.TimeTrack(time.Now(), "Scheduler")
	log.Info("Stopping Scheduler...")
	scheduler.Stop()
	//select {}
}

// func (s *scheduler) Entries() []*cron.Entry {
// 	return s.entrySnapshot()
// }
