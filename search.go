package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"gopkg.in/olivere/elastic.v5"
)

type Idx struct {
	Mid    string `json:"mid"`
	Name   string `json:"name"`
	Region int    `json:"region"`
	Govern int    `json:"govern"`
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

func elasticConnect() (context.Context, *elastic.Client, error) {
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200", "http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "changeme"),
	)
	if err != nil {
		return nil, nil, err
	}
	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://localhost:9200").Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esVersion, err := client.ElasticsearchVersion("http://localhost:9200")
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Elasticsearch version %s\n", esVersion)

	// Use the IndexExists service to check if a specified index exists.

	if _, err := client.DeleteIndex("mtests").Do(ctx); err != nil {
		return nil, nil, err
	}
	exists, err := client.IndexExists("mtests").Do(ctx)
	if err != nil {
		return nil, nil, err
	}

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("mtests").BodyString(mapping).Do(ctx)
		if err != nil {
			return nil, nil, err
		}

		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return ctx, client, nil
}

func elasticIndex() error {
	var (
		mid            uuid.UUID
		region, govern int
		name, author   string
		trk            Idx
	)
	ctx, client, err := elasticConnect()
	if err != nil {
		return err
	}

	res, err := db.Query("SELECT mid, name, region, govern, author FROM mtests")
	if err != nil {
		return err
	}

	log.Info("Indexing started")
	for res.Next() {
		if err := res.Scan(&mid, &name, &region, &govern, &author); err != nil {
			log.Error(err.Error(), " | ", mid, "\n")
			return err
		}
		trk = Idx{Mid: mid.String(), Name: name, Region: region, Govern: govern, Author: author}
		if _, err = client.Index().
			Index("mtests").Type("mtest").Id(mid.String()).BodyJson(trk).Do(ctx); err != nil {
			return err
		}
	}
	log.Info("Indexing completed!")
	return nil
}

func updateIndex(id int64) error {
	var (
		mid            uuid.UUID
		region, govern int
		name, author   string
	)
	ctx, client, err := elasticConnect()
	ind, err := db.Query("SELECT mid, name, region, govern, author FROM mtests WHERE id=?;", id)
	if err != nil {
		return err
	}

	for ind.Next() {
		if err = ind.Scan(&mid, &name, &region, &govern, &author); err != nil {
			log.Error(err)
			return err
		}
	}
	idx := Idx{Mid: mid.String(), Name: name, Region: region, Govern: govern, Author: author}
	if _, err = client.Index().Index("mtests").Type("mtest").Id(mid.String()).BodyJson(idx).Do(ctx); err != nil {
		return err
	}

	return nil
}
