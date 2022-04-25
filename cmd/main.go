package main

import (
	"enterpret/config"
	"enterpret/dataaccess"
	"enterpret/server"
	sources2 "enterpret/sources"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello World")

	sources := []config.SourceConfig{
		{
			Name: "discourse",
		},
		{
			Name: "template",
		},
	}

	_ = sources

	ds := dataaccess.NewDataStore()
	sourceProcessor := sources2.NewSourceProcessor(ds)

	for _, t := range config.TC {
		for _, s := range t.Sources {
			mp, err := sourceProcessor.GetProcessor(s.Name)
			if err != nil {
				fmt.Println("Failed fetching message processor for source "+s.Name, err.Error())
				return
			}

			_, err = mp.FetchAndStoreFeedbacks(s.Params, t.Name)
			if err != nil {
				fmt.Println("Failed fetching and storing messages: ", err.Error())
				return
			}
		}
	}

	messages, _ := ds.FetchFeedbacks("zoom.us", "discourse", 1, 20)

	for _, m := range messages {
		log.Default().Println(m.Meta.ID, ": ", m.Data.Message)
	}

	err := server.NewServer(ds)
	if err != nil {
		log.Default().Println("Server ended with error ", err.Error())
	}
}
