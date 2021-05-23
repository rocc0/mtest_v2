package search

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"

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

type Service struct {
	*elastic.Client
	*sql.DB
}

func NewService(db *sql.DB) (*Service, error) {
	return &Service{}, nil
}

func (s *Service) Connect(address string) error {
	if address == "" {
		address = "http://localhost:9200"
	}

	client, err := elastic.NewClient(
		elastic.SetURL(address, address),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "changeme"),
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping(address).Do(ctx)
	if err != nil {
		return err
	}

	logrus.Infof("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esVersion, err := client.ElasticsearchVersion(address)
	if err != nil {
		return err
	}

	logrus.Infof("Elasticsearch version %s\n", esVersion)
	s.Client = client

	return nil
}

func (s *Service) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// Use the IndexExists service to check if a specified index exists.
	if _, err := s.DeleteIndex("mtests").Do(ctx); err != nil {
		return err
	}
	exists, err := s.IndexExists("mtests").Do(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	createIndex, err := s.CreateIndex("mtests").BodyString(mapping).Do(ctx)
	if err != nil {
		return err
	}

	if !createIndex.Acknowledged {
		// Not acknowledged
	}

	return nil
}

func (s *Service) ElasticIndex() error {
	var (
		mid            uuid.UUID
		region, govern int
		name, author   string
		trk            Idx
	)

	res, err := s.Query("SELECT mid, name, region, govern, author FROM mtests")
	if err != nil {
		return err
	}

	logrus.Info("Indexing started")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for res.Next() {
		if err := res.Scan(&mid, &name, &region, &govern, &author); err != nil {
			logrus.Error(err.Error(), " | ", mid, "\n")
			return err
		}
		trk = Idx{Mid: mid.String(), Name: name, Region: region, Govern: govern, Author: author}
		if _, err = s.Index().
			Index("mtests").Type("mtest").Id(mid.String()).BodyJson(trk).Do(ctx); err != nil {
			return err
		}
	}
	logrus.Info("Indexing completed!")
	return nil
}

func (s *Service) UpdateIndex(id int64) error {
	var (
		mid            uuid.UUID
		region, govern int
		name, author   string
	)

	ind, err := s.Query("SELECT mid, name, region, govern, author FROM mtests WHERE id=?;", id)
	if err != nil {
		return err
	}

	for ind.Next() {
		if err = ind.Scan(&mid, &name, &region, &govern, &author); err != nil {
			logrus.Error(err)
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	idx := Idx{Mid: mid.String(), Name: name, Region: region, Govern: govern, Author: author}

	if _, err = s.Index().Index("mtests").Type("mtest").Id(mid.String()).BodyJson(idx).Do(ctx); err != nil {
		return err
	}

	return nil
}