package couchdbrus_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/codehakase/couchdbrus"
	"github.com/ryanjyoder/couchdb"
	"github.com/sirupsen/logrus"
)

var client *couchdb.Client
var db couchdb.DatabaseService

func TestHook(t *testing.T) {
	client, db = couchdbHelper(os.Getenv("COUCHDB_DBNAME"))
	log := logrus.New()
	couchdbHook, err := couchdbrus.NewHook(client, "mylogdatabaseName")
	if err != nil {
		// do proper error handling here...
		log.Fatalf("could not create couchdb hook: %v",
			err)
	}
	log.Hooks.Add(couchdbHook)

	log.WithFields(logrus.Fields{
		"hostname": "hakaselabs",
		"source":   "spacex",
		"tag":      "test-tag",
	}).Info("Hello Captain we can see the moon!")

	if len(couchdbHook.Levels()) < 1 {
		t.Error("Error setting level, level length less than 1")
	}
}

func ExampleHook() {
	client, db = couchdbHelper(os.Getenv("COUCHDB_DBNAME"))
	log := logrus.New()
	couchdbHook, err := couchdbrus.NewHook(client, "mylogdatabaseName")
	if err != nil {
		// do proper error handling here...
		log.Fatalf("could not create couchdb hook: %v",
			err)
	}
	log.Hooks.Add(couchdbHook)

	log.WithFields(logrus.Fields{
		"hostname": "hakaselabs",
		"source":   "spacex",
		"tag":      "test-tag",
	}).Info("Hello Captain we can see the moon!")
}

func couchdbHelper(dbname string) (*couchdb.Client, couchdb.DatabaseService) {
	u, err := url.Parse(os.Getenv("COUCHDB_SERVER"))
	if err != nil {
		logrus.Fatalf("could not parse cdb server url: %v", err)
	}

	// create couchdb client
	dbClient, err := couchdb.NewAuthClient(os.Getenv("COUCHDB_USER"), os.Getenv("COUCHDB_PASS"), u)
	if err != nil {
		logrus.Fatalf("could not create a cdb client: %v", err)
	}

	db, err := CreateOrUseDb(dbClient, dbname)
	if err != nil {
		logrus.Fatalf("could not create or use couchdb database: %v", err)
	}

	return dbClient, db
}

func CreateOrUseDb(client *couchdb.Client, db string) (couchdb.DatabaseService, error) {
	_, err := client.Create(db)
	if err != nil {
		dberror, ok := err.(*couchdb.Error)
		if !ok {
			return nil, dberror
		}
		// the error = file_exists, that's a good one
	}
	dbs := client.Use(db)
	return dbs, nil
}
