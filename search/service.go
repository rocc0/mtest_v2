package search

import (
	"context"
	"database/sql"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"
)

type Idx struct {
	Mid        string `json:"mid"`
	Name       string `json:"name"`
	Region     int    `json:"region"`
	Govern     int    `json:"govern"`
	Business   int    `json:"business"`
	MathResult int    `json:"math_result"`
	CorrResult int    `json:"corr_result"`
	Author     string `json:"author"`
	RegAct     string `json:"reg_act"`
	Synonyms   string `json:"synonyms"`
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
			"math_result":{
				"type":"long"
			},
			"corr_result":{
				"type":"long"
			},
			"author":{
				"type":"text"
			},
			"reg_act":{
				"type":"text",
				"analyzer": "ukrainian"
			},
			"synonyms":{
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
		mathResult, corrResult   int
		name, author             string
		regAct                   string
		trk                      Idx
	)

	res, err := s.Query(`
		SELECT 
		mt.mid, mt.name, mt.region, mt.govern, mt.business, mt.author, 
		       COALESCE(rg.doc_text, ''), mt.corr_result, mt.math_result 
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
		if err := res.Scan(&mid, &name, &region, &govern, &business, &author, &regAct, &corrResult, &mathResult); err != nil {
			logrus.Error(err.Error(), " | ", mid, "\n")
			return err
		}

		trk = Idx{Mid: mid.String(), Name: name, Region: region, Govern: govern, Business: business, Author: author, RegAct: regAct}
		syns, err := s.getSynonyms(mid.String())
		if err != nil {
			return err
		}
		trk.Synonyms = syns
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
		mathResult, corrResult   int
		name, author, regAct     string
	)

	ind, err := s.Query(`
        SELECT 
		mtests.mid, mtests.name, mtests.region, mtests.govern, mtests.business, 
               mtests.author, reg_acts.doc_text, mtests.corr_result, mtests.math_result  
		FROM mtests
		JOIN reg_acts 
		ON mtests.mid=reg_acts.mid WHERE mtests.mid=?;`, id)
	if err != nil {
		return err
	}

	for ind.Next() {
		if err = ind.Scan(&mid, &name, &region, &govern, &business, &author, &regAct, &corrResult, &mathResult); err != nil {
			logrus.Error(err)
			return err
		}
	}

	syns, err := s.getSynonyms(mid.String())
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	idx := Idx{
		Mid:        mid.String(),
		Name:       name,
		Region:     region,
		Govern:     govern,
		Business:   business,
		Author:     author,
		RegAct:     regAct,
		MathResult: mathResult,
		CorrResult: corrResult,
		Synonyms:   syns,
	}

	if _, err = s.Index().Index("mtests").Id(mid.String()).BodyJson(idx).Do(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateIndexWithFile(id string, text string) error {
	var (
		mid                      uuid.UUID
		region, govern, business int
		mathResult, corrResult   int
		name, author             string
	)

	ind, err := s.Query("SELECT mid, name, region, govern, business, author, corr_result, math_result FROM mtests WHERE mid=?;", id)
	if err != nil {
		return err
	}

	for ind.Next() {
		if err = ind.Scan(&mid, &name, &region, &govern, &business, &author, &mathResult, &corrResult); err != nil {
			logrus.Error(err)
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	syns, err := s.getSynonyms(mid.String())
	if err != nil {
		return err
	}

	idx := Idx{
		Mid:        mid.String(),
		Name:       name,
		Region:     region,
		Govern:     govern,
		Business:   business,
		Author:     author,
		RegAct:     text,
		MathResult: mathResult,
		CorrResult: corrResult,
		Synonyms:   syns,
	}
	if _, err = s.Index().Index("mtests").Id(mid.String()).BodyJson(idx).Do(ctx); err != nil {
		return err
	}

	return nil
}

const listSynonymsQuery = `SELECT synonym FROM synonyms WHERE mtest_id=?;`

func (s *Service) getSynonyms(mtestID string) (string, error) {
	var synonyms []string
	res, err := s.Query(listSynonymsQuery, mtestID)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := res.Close(); err != nil {
			log.Error(err)
		}
	}()
	for res.Next() {
		var s string
		if err := res.Scan(&s); err != nil {
			return "", err
		}

		synonyms = append(synonyms, s)
	}

	return strings.Join(synonyms, " "), nil

}
