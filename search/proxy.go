package search

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type request struct {
	From  int `json:"from"`
	Size  int `json:"size"`
	Query struct {
		Bool struct {
			Should struct {
				MultiMatch struct {
					Query  string   `json:"query"`
					Fields []string `json:"fields"`
				} `json:"multi_match"`
			} `json:"should"`
			Filter *filter `json:"filter,omitempty"`
		} `json:"bool"`
	} `json:"query"`
}

type filter struct {
	Bool *struct {
		Must []struct {
			Term struct {
				Govern   int `json:"govern,omitempty"`
				Region   int `json:"region,omitempty"`
				Business int `json:"business,omitempty"`
			} `json:"term,omitempty"`
		} `json:"must,omitempty"`
	} `json:"bool,omitempty"`
}

type ElasticProxy struct {
	db  *sql.DB
	url string
}

func NewElasticProxy(url string, db *sql.DB) ElasticProxy {
	return ElasticProxy{
		db:  db,
		url: url,
	}
}

func (e *ElasticProxy) ElasticProxy(c *gin.Context) {
	remote, err := url.Parse(e.url)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
		req.Method = http.MethodPost
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Error(err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		var r request
		if err := json.Unmarshal(body, &r); err != nil {
			log.Error(err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		b, err := e.findSynonyms(r)
		if err != nil {
			log.Error(err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(b))
		req.ContentLength = int64(len(b))
	}

	proxy.ModifyResponse = func(response *http.Response) error {
		return nil
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func (e *ElasticProxy) findSynonyms(phrase request) ([]byte, error) {
	//Query.Bool.Should.MultiMatch.Query
	if phrase.Query.Bool.Filter == nil || phrase.Query.Bool.Filter.Bool == nil || phrase.Query.Bool.Filter.Bool.Must == nil {
		phrase.Query.Bool.Filter = nil
	}
	b, err := json.Marshal(phrase)
	if err != nil {
		return nil, err
	}

	log.Error(string(b))
	return b, nil
}
