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
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/asdine/storm"
	bolt "github.com/coreos/bbolt"
)

var (
	// DB is Storm Database Object
	DB               *storm.DB
	dbReady          bool
	datastoreMetrics bolt.Stats
)

// CloseDB will shutdown DB gracefuly
func CloseDB() {
	log.Info("Closing DB Connection...")
	DB.Close()
}

// InitDB will initialize BoltDB and Storm Toolkit
func InitDB(stumb string) {
	var err error

	DB, err = storm.Open(stumb)
	if err != nil {
		log.WithFields(log.Fields{
			"DB":    stumb,
			"error": err,
		}).Fatal("Failed to initialize DB")
	} else {
		dbReady = true
		log.WithFields(log.Fields{
			"DB":      stumb,
			"dbReady": dbReady,
		}).Info("Database initialized")
	}

}

func save(i interface{}, db storm.Node) {
	//defer utils.TimeTrack(time.Now(), save)
	err := db.Save(i)
	if err != nil {
		log.Error(err)
	}
}

func updateOrSave(key string, value interface{}, item interface{}, db storm.Node) error {
	//defer utils.TimeTrack(time.Now(), updateOrSave)
	obj := item

	err := db.One(key, value, obj)
	if err != nil {
		if err == storm.ErrNotFound {
			log.WithFields(log.Fields{
				"key":   key,
				"value": value,
			}).Debug("Object not found, adding it to the inventory")
			db.Save(item)
			return nil
		}
		log.Error(err)
		return err
	}

	err = db.Update(item)
	if err != nil {
		log.Error(err)
	}
	return nil
}

// DataStoreMetrics provides BOLT DB Statistics
func DataStoreMetrics() {
	// Grab the initial stats.
	prev := DB.Bolt.Stats()

	for {
		// Wait for 10s.
		time.Sleep(5 * time.Second)

		// Grab the current stats and diff them.
		stats := DB.Bolt.Stats()
		datastoreMetrics = stats.Sub(&prev)

		// Encode stats to JSON and print to STDERR.
		//json.NewEncoder(os.Stderr).Encode(diff)

		// Save stats for the next loop.
		prev = stats
	}
}
