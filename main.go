package main

import (
	"bahamut/core/bahamut"
	"bahamut/core/component"
	"bahamut/core/types"
	"fmt"
	"time"
)

func main() {
	netw := types.Network{Host: "localhost", Port: 8080}
	fmt.Printf("Network : %v", netw.Address())

	opts := types.Options{
		Image:         "metabase/metabase:latest",
		RestartPolicy: types.UnlessStopped,
	}

	metabase := component.New("metabase", types.Reporting, opts)

	postgres := component.New("postgres", types.Database, types.Options{
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=bahamut",
			"POSTGRES_PASSWORD=really--secret",
		},
	})

	setting := types.BahamutSettings{types.Database: 2, types.Reporting: 1}
	fmt.Printf("Setting len %v", setting)

	bahamut := bahamut.New(setting)
	go bahamut.Start()
	fmt.Printf("Runners %v ", bahamut.Runners)
	bahamut.Schedule(*metabase)
	bahamut.Schedule(*postgres)

	time.Sleep(50 * time.Second)

	fmt.Printf("Scheluded a component to stop...")
	metabase.State = component.Completed
	postgres.State = component.Completed

	bahamut.Schedule(*metabase)
	bahamut.Schedule(*postgres)

	time.Sleep(2000 * time.Second)

	// Create a bahamut instance to schedule components
	// bahamut := bahamut.New(types.BahaumtSettings{types.Database: 4, types.Baremetal: 5},)

	// r := runner.New("BI")

	// go r.Start()

	// // runner should enqueue components not run directly.
	// _ = r.Schedule(*metabase)
	// _ = r.Schedule(*postgres)

	// fmt.Printf("Sleeping for 50 seconds...")
	// time.Sleep(50 * time.Second)

	// fmt.Printf("Scheluded a component to stop...")
	// metabase.State = component.Completed

	// _ = r.Schedule(*metabase)

	// time.Sleep(100 * time.Second)

}
