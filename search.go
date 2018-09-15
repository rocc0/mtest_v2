package main

import (
	"context"
	"log"

	"gopkg.in/olivere/elastic.v5"
	"github.com/google/uuid"
)


type Idx struct {
	Mid string `json:"mid"`
	Name string `json:"name"`
	Region int `json:"region"`
	Govern int `json:"govern"`
	Author string `json:"author"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"mtest":{
			"properties":{
				"mid":{
					"type":"string"
				},
				"name":{
					"type":"text",
					"analyzer": "ukrainian"
				},
				"region":{
					"type":"integer"
				},
				"govern":{
					"type":"integer"
				},
				"author":{
					"type":"string"
				}
			}
		}
	}
}`



func elasticConnect() (context.Context, *elastic.Client, error){
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200", "http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "changeme"),
	)
	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://localhost:9200").Do(ctx)
	check(err)

	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://localhost:9200")
	check(err)

	log.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	client.DeleteIndex("mtests").Do(ctx)
	exists, err := client.IndexExists("mtests").Do(ctx)
	check(err)

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("mtests").BodyString(mapping).Do(ctx)
		check(err)

		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return ctx, client, nil
}


func elasticIndex(){
	var (
		mid uuid.UUID
		region, govern int
		name, author string
		trk Idx
	)
	ctx, client, err := elasticConnect()
	check(err)

	res, err := db.Query("SELECT mid, name, region, govern, author FROM mtests")
	check(err)

	log.Print("Indexing started")
	for res.Next(){
		err := res.Scan(&mid, &name, &region, &govern, &author)
		if err != nil {
			log.Print(err.Error(), " | " ,mid, "\n")
		}
		trk = Idx{mid.String(), name, region, govern,author}
		_, err = client.Index().
			Index("mtests").
			Type("mtest").
			Id(mid.String()).
			BodyJson(trk).
			Do(ctx)
		check(err)
	}
	log.Print("Indexing complited!")
}

func updateIndex(id int64) error {
	var (
		mid uuid.UUID
		region, govern int
		name, author string
	)
	ctx, client, err := elasticConnect()

	ind, err := db.Query("SELECT mid, name, region, govern, author FROM mtests WHERE id=?;", id)
	check(err)

	for ind.Next() {
		err = ind.Scan(&mid, &name, &region, &govern, &author)
		if err != nil {
			log.Print(err.Error())
			return err
		}
	}
	idx := Idx{mid.String(), name, region, govern,author}
	_, err = client.Index().
		Index("mtests").
		Type("mtest").
		Id(string(mid.String())).
		BodyJson(idx).
		Do(ctx)

		if err != nil {
			log.Print(err.Error())
			return err
		}

	return nil
}