package search

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"
)

type Idx struct {
	Mid      string `json:"mid"`
	Name     string `json:"name"`
	Region   int    `json:"region"`
	Govern   int    `json:"govern"`
	Business int    `json:"business"`
	Author   string `json:"author"`
	RegAct   string `json:"reg_act"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"aliases": {
    	"mtest": {}
  	},
	"mappings":{
		"properties":{
			"mid":{
				"type":"text"
			},
			"name":{
				"type":"text",
				"analyzer": "ukrainian"
			},
			"region":{
				"type":"long"
			},
			"govern":{
				"type":"long"
			},
			"business":{
				"type":"long"
			},
			"author":{
				"type":"text"
			},
			"reg_act":{
				"type":"text",
				"analyzer": "ukrainian"
			}
		}
	}
}`

type Service struct {
	*elastic.Client
	*sql.DB
}

func NewService(db *sql.DB) (*Service, error) {
	return &Service{DB: db}, nil
}

func (s *Service) Connect(address string) error {
	if address == "" {
		address = "http://localhost:9200"
	}

	client, err := elastic.NewClient(
		elastic.SetURL(address),
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

	exists, err := s.IndexExists("mtests").Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		if _, err := s.DeleteIndex("mtests").Do(ctx); err != nil {
			return err
		}
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
		mid                      uuid.UUID
		region, govern, business int
		name, author             string
		regAct                   string
		trk                      Idx
	)

	res, err := s.Query(`
		SELECT 
		mt.mid, mt.name, mt.region, mt.govern, mt.business, mt.author, COALESCE(rg.doc_text, '') 
		FROM mtests mt
		LEFT JOIN reg_acts rg
		ON mt.mid = rg.mid GROUP BY mt.mid`)
	if err != nil {
		return err
	}

	logrus.Info("Indexing started")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for res.Next() {
		if err := res.Scan(&mid, &name, &region, &govern, &business, &author, &regAct); err != nil {
			logrus.Error(err.Error(), " | ", mid, "\n")
			return err
		}

		trk = Idx{Mid: mid.String(), Name: name, Region: region, Govern: govern, Business: business, Author: author, RegAct: regAct}
		if _, err = s.Index().Index("mtests").Id(mid.String()).BodyJson(trk).Do(ctx); err != nil {
			return err
		}
	}
	logrus.Info("Indexing completed!")
	return nil
}

func (s *Service) UpdateIndex(id string) error {
	var (
		mid                      uuid.UUID
		region, govern, business int
		name, author, regAct     string
	)

	ind, err := s.Query(`
        SELECT 
		mtests.mid, mtests.name, mtests.region, mtests.govern, mtest.business, mtests.author, reg_acts.doc_text 
		FROM mtests
		JOIN reg_acts 
		ON mtests.mid=reg_acts.mid WHERE mtests.mid=?;`, id)
	if err != nil {
		return err
	}

	for ind.Next() {
		if err = ind.Scan(&mid, &name, &region, &govern, &business, &author, &regAct); err != nil {
			logrus.Error(err)
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	idx := Idx{Mid: mid.String(), Name: name, Region: region, Govern: govern, Business: business, Author: author, RegAct: regAct}

	if _, err = s.Index().Index("mtests").Id(mid.String()).BodyJson(idx).Do(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateIndexWithFile(id string, text string) error {
	var (
		mid                      uuid.UUID
		region, govern, business int
		name, author             string
	)

	ind, err := s.Query("SELECT mid, name, region, govern, business, author FROM mtests WHERE mid=?;", id)
	if err != nil {
		return err
	}

	for ind.Next() {
		if err = ind.Scan(&mid, &name, &region, &govern, &business, &author); err != nil {
			logrus.Error(err)
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	idx := Idx{Mid: mid.String(), Name: name, Region: region, Govern: govern, Business: business, Author: author, RegAct: text}
	if _, err = s.Index().Index("mtests").Id(mid.String()).BodyJson(idx).Do(ctx); err != nil {
		return err
	}

	return nil
}
