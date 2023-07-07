package es

import (
	"errors"
	"github.com/star-table/common/core/config"
	"github.com/olivere/elastic/v7"
	"net/http"
	"sync"
	"time"
)

const (
	ESVersion7 = "7"
	ESVersion6 = "6"
)

var (
	_lock    sync.Mutex
	esClient = &ElasticSearchClient{}
)

type ElasticSearchClient struct {
	Client  *elastic.Client
	Version string
}

func initClient() (*ElasticSearchClient, error) {
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(config.GetElasticSearch().ServerUrls...),
		elastic.SetSniff(config.GetElasticSearch().Sniff),
	}

	if auth := config.GetElasticSearch().Auth; auth != nil {
		options = append(options, elastic.SetBasicAuth(auth.UserName, auth.Password))
	}

	timeout := config.GetElasticSearch().Timeout

	if timeout <= 0 {
		timeout = 1000
	}
	options = append(options, elastic.SetHttpClient(&http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}))

	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}

	esClient := &ElasticSearchClient{
		Client: client,
	}

	esClient, err = getEsVersion(esClient)
	if err != nil {
		return nil, err
	}
	return esClient, nil
}

func GetESClient() (*ElasticSearchClient, error) {
	if esClient != nil && esClient.Version != "" {
		return esClient, nil
	}

	_lock.Lock()
	defer _lock.Unlock()
	if esClient != nil && esClient.Version != "" {
		return esClient, nil
	}

	esClient, err := initClient()
	if err != nil {
		return nil, err
	}
	return esClient, nil
}

/**
获取es版本
*/
func getEsVersion(client *ElasticSearchClient) (*ElasticSearchClient, error) {
	url := config.GetElasticSearch().ServerUrls[0]
	v, err := client.Client.ElasticsearchVersion(url)
	if err != nil {
		return nil, err
	}
	if v == "" {
		return nil, errors.New(" get elasticsearch version fail. url: " + url)
	}
	vs := v[0:1]
	if vs == ESVersion7 || vs == ESVersion6 {
		client.Version = vs
		return client, nil
	}
	return nil, errors.New(" elasticsearch version not support. version: " + v)
}
