package main

import (
	"errors"
	"flag"
	"log"
)

// Backdays to check for updates
var backdays int

// Collection to sync
var collection string

func init() {
	flag.IntVar(&backdays, "b", 0, "Specify backdays as an integer > 0")
	flag.StringVar(&collection, "c", "", "Specify what to sync - 'members', 'modules', 'resources' or 'all'")
}

func main() {

	err := flagCheck()
	if err != nil {
		log.Fatalf("flagCheck() err = %s", err)
	}
	log.Printf("Running syncr with backdays: %d on collection: %s", backdays, collection)

	err = sync()
	if err != nil {
		log.Fatalf("sync() err = %s", err)
	}
}

func flagCheck() error {
	flag.Parse()
	if backdays < 1 {
		return errors.New("Backdays (-b) required, -h for help")
	}
	if collection == "" {
		return errors.New("Sync target (-c) required, -h for help")
	}
	return nil
}

func sync() error {
	switch collection {
	case "member", "members":
		return syncMembers()
	case "module", "modules":
		return syncModules()
	case "resource", "resources":
		return syncResources()
	case "all":
		return syncAll()
	}
	return nil
}

func syncAll() error {
	err := syncMembers()
	if err != nil {
		return err
	}
	err = syncModules()
	if err != nil {
		return err
	}
	err = syncResources()
	if err != nil {
		return err
	}
	return nil
}

func syncMembers() error {
	log.Println("Syncing members")
	return nil
}

func syncModules() error {
	log.Println("Syncing modules")
	return nil
}

func syncResources() error {
	log.Println("Syncing resources")
	return nil
}
