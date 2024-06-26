package es

import (
	"bytes"
	"context"
	"encoding/json"
	"msm/config"
	"msm/log"
	"msm/model"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var esClient *elasticsearch.Client

type esService struct{}

var EsService = new(esService)

func InitEs() bool {
	if config.CF.EsEnable {
		cfg := elasticsearch.Config{
			Addresses: []string{
				config.CF.EsUrl,
			},
			Username: config.CF.EsUsername,
			Password: config.CF.EsPassword,
		}
		var err error
		esClient, err = elasticsearch.NewClient(cfg)
		if err != nil {
			log.Logger.Fatalln("Failed to connect to es")
		}
		_, err = esClient.Info()
		if err != nil {
			log.Logger.Error("es启动失败", err)
			config.CF.EsEnable = false
		} else {
			return true 
		}
	} else {
		log.Logger.Debug("不使用es")
	}
	return false
}

// idx 为空，默认随机唯一字符串
func (e *esService) Index(index, idx string, doc map[string]interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		log.Logger.Error(err, "Error encoding doc")
		return
	}
	res, err := esClient.Index(
		index,
		&buf,
		esClient.Index.WithDocumentID(idx),
		esClient.Index.WithRefresh("true"),
	)
	if err != nil {
		log.Logger.Error(err, "Error create response")
	}
	defer res.Body.Close()
}

func (e *esService) Insert(log string, processName string, using string, ts int64) {
	doc := map[string]interface{}{
		"log":   log,
		"name":  processName,
		"using": using,
		"time":  ts,
	}
	e.Index(config.CF.EsIndex, "", doc)
}

func (e *esService) Search(req model.GetLogReq) model.LogResp {
	query := []func(*esapi.SearchRequest){
		esClient.Search.WithIndex(config.CF.EsIndex),
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithPretty(),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithFrom(req.Page.From),
		esClient.Search.WithSize(req.Page.Size),
	}
	if req.Sort == "asc" {
		query = append(query, esClient.Search.WithSort("time:asc"))
	}
	if req.Sort == "desc" {
		query = append(query, esClient.Search.WithSort("time:desc"))
	}
	body := e.buildQueryBody(req)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		log.Logger.Error(err)
		return model.LogResp{}
	}
	query = append(query, esClient.Search.WithBody(&buf))
	res, err := esClient.Search(query...)
	if err != nil {
		log.Logger.Error(err)
		return model.LogResp{}
	}
	resp := model.EsResp{}
	json.NewDecoder(res.Body).Decode(&resp)
	res.Body.Close()
	result := model.LogResp{}
	for _, v := range resp.Hits.Hits {
		result.Data = append(result.Data, model.Eslog{
			Log:   v.Source.Log,
			Name:  v.Source.Name,
			Using: v.Source.Using,
			Time:  v.Source.Time,
			Id:    v.ID,
		})
	}
	result.Total = resp.Hits.Total.Value
	return result
}

func (e *esService) buildQueryBody(req model.GetLogReq) model.QueryBody {
	result := model.QueryBody{}
	if req.TimeRange.EndTime != 0 || req.TimeRange.StartTime != 0 {
		result.Query.Bool.Must = append(result.Query.Bool.Must, map[string]any{
			"range": map[string]any{
				"time": map[string]any{
					"gte": req.TimeRange.StartTime,
					"lte": req.TimeRange.EndTime,
				},
			},
		})
	}
	if req.Match.Log != "" {
		result.Query.Bool.Must = append(result.Query.Bool.Must, map[string]any{
			"match": map[string]any{
				"log": req.Match.Log,
			},
		})
	}
	if req.Match.Name != "" {
		result.Query.Bool.Must = append(result.Query.Bool.Must, map[string]any{
			"match": map[string]any{
				"name": req.Match.Name,
			},
		})
	}
	if req.Match.Using != "" {
		result.Query.Bool.Must = append(result.Query.Bool.Must, map[string]any{
			"match": map[string]any{
				"using": req.Match.Using,
			},
		})
	}
	return result
}
