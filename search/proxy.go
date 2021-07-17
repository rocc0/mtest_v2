package search

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	datapkg "mtest.com.ua/db/dataprocessor"

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
	LoadHandler
	url   string
	cache []datapkg.GlobalSynonym
}

type LoadHandler func() ([]datapkg.GlobalSynonym, error)

func NewElasticProxy(url string, l LoadHandler) ElasticProxy {
	return ElasticProxy{
		LoadHandler: l,
		url:         url,
		cache:       []datapkg.GlobalSynonym{},
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

func (e *ElasticProxy) LoadCache() error {
	c, err := e.LoadHandler()
	if err != nil {
		return err
	}
	e.cache = c
	return nil
}

func (e *ElasticProxy) findSynonyms(phrase request) ([]byte, error) {
	if phrase.Query.Bool.Filter == nil || phrase.Query.Bool.Filter.Bool == nil || phrase.Query.Bool.Filter.Bool.Must == nil {
		phrase.Query.Bool.Filter = nil
	}

	words := strings.Split(phrase.Query.Bool.Should.MultiMatch.Query, " ")

	for _, word := range words {
		for _, w := range e.cache {
			if w.Word == word {
				for _, s := range w.Synonyms {
					if !contains(s, words) {
						words = append(words, s)
					}
				}
			}
		}
	}
	q := strings.Join(words, " ")
	phrase.Query.Bool.Should.MultiMatch.Query = q
	b, err := json.Marshal(phrase)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func contains(word string, arr []string) bool {
	for _, s := range arr {
		if s == word {
			return true
		}
	}
	return false
}
