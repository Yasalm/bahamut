package main

import (
	"bahamut/core/bahamut"
	"bahamut/core/component"
	"bahamut/core/types"
	"time"
)

func main() {
	setting := types.BahamutSettings{
		types.Database:  2,
		types.General:   1,
		types.Reporting: 1,
	}
	bahamut := bahamut.New(setting)
	go bahamut.Start()

	postgres := component.New("ps:database", types.Database, types.Options{Image: "postgres", Env: []string{"POSTGRES_PASSWORD=root", "POSTGRES_USER=root", "POSTGRES_DB=db"}})
	metabase := component.New("metabase", types.Reporting, types.Options{Image: "metabase/metabase:latest"})

	bahamut.Schedule(*postgres)
	bahamut.Schedule(*metabase)
	time.Sleep(60 * time.Second)
	postgres.State = component.Completed
	bahamut.Schedule(*postgres)

	time.Sleep(60 * time.Second)

}
