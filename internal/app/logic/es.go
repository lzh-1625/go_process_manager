package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/lzh-1625/go_process_manager/config"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/log"

	"github.com/olivere/elastic/v7"
)

type esLogic struct {
	esClient *elastic.Client
}

var (
	EsLogic = new(esLogic)
)

func (e *esLogic) InitEs() bool {
	if !config.CF.EsEnable {
		log.Logger.Debug("不使用es")
		return false
	}

	var err error
	EsLogic.esClient, err = elastic.NewClient(
		elastic.SetURL(config.CF.EsUrl),
		elastic.SetBasicAuth(config.CF.EsUsername, config.CF.EsPassword),
		elastic.SetSniff(false),
	)
	if err != nil {
		config.CF.EsEnable = false
		log.Logger.Warnw("Failed to connect to es", "err", err)
		return false
	}
	EsLogic.CreateIndexIfNotExists(config.CF.EsIndex)
	return true
}

func (e *esLogic) Insert(logContent string, processName string, using string, ts int64) {
	data := model.ProcessLog{
		Log:   logContent,
		Name:  processName,
		Using: using,
		Time:  ts,
	}
	_, err := e.esClient.Index().Index(config.CF.EsIndex).BodyJson(data).Do(context.TODO())
	if err != nil {
		log.Logger.Errorw("es数据插入失败", "err", err)
	}
}

func (e *esLogic) CreateIndexIfNotExists(index string) error {

	ctx := context.Background()
	exists, err := e.esClient.IndexExists(index).Do(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	info, err := e.esClient.CreateIndex(index).BodyString(e.structToJSON()).Do(ctx)
	if err != nil {
		return err
	}
	if !info.Acknowledged {
		return fmt.Errorf("ES 创建索引 [%s] 失败", index)
	}
	return nil
}

func (e *esLogic) Search(req model.GetLogReq, filterProcessName ...string) model.LogResp {
	// 检查 req 是否为 nil
	if req.Page.From < 0 || req.Page.Size <= 0 {
		log.Logger.Error("无效的分页请求参数")
		return model.LogResp{Total: 0, Data: []model.ProcessLog{}}
	}

	search := e.esClient.Search(config.CF.EsIndex).From(req.Page.From).Size(req.Page.Size).TrackScores(true)
	if req.Sort == "asc" {
		search.Sort("time", true)
	}
	if req.Sort == "desc" {
		search.Sort("time", false)
	}

	queryList := []elastic.Query{}
	timeRangeQuery := elastic.NewRangeQuery("time")
	if req.TimeRange.StartTime != 0 {
		queryList = append(queryList, timeRangeQuery.Gte(req.TimeRange.StartTime))
	}
	if req.TimeRange.EndTime != 0 {
		queryList = append(queryList, timeRangeQuery.Lte(req.TimeRange.EndTime))
	}
	if req.Match.Log != "" {
		queryList = append(queryList, elastic.NewMatchQuery("log", req.Match.Log))
	}
	if req.Match.Name != "" {
		queryList = append(queryList, elastic.NewMatchQuery("name", req.Match.Name))
	}
	if req.Match.Using != "" {
		queryList = append(queryList, elastic.NewMatchQuery("using", req.Match.Using))
	}

	if len(filterProcessName) != 0 { // 过滤进程名
		shouldQueryList := []elastic.Query{}
		for _, fpn := range filterProcessName {
			shouldQueryList = append(shouldQueryList, elastic.NewMatchQuery("name", fpn))
		}
		if len(shouldQueryList) > 0 {
			shouldQuery := elastic.NewBoolQuery().Should(shouldQueryList...)
			queryList = append(queryList, shouldQuery)
		}
	}

	result := model.LogResp{}
	resp, err := search.Query(elastic.NewBoolQuery().Must(queryList...)).Do(context.TODO())
	if err != nil {
		log.Logger.Errorw("es查询失败", "err", err, "reason", resp.Error.Reason)
		return result
	}

	// 遍历响应结果
	for _, v := range resp.Hits.Hits {
		if v.Source != nil {
			var data model.ProcessLog
			if err := json.Unmarshal(v.Source, &data); err == nil {
				result.Data = append(result.Data, data)
			} else {
				log.Logger.Errorw("JSON 解码失败", "err", err)
			}
		}
	}

	result.Total = resp.TotalHits()
	return result
}

// 通过反射得到mapping
func (e *esLogic) structToJSON() string {
	typ := reflect.TypeOf(model.ProcessLog{})
	properties := make(map[string]map[string]string)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldTag := field.Tag.Get("type")
		if fieldTag != "" {
			properties[field.Tag.Get("json")] = map[string]string{
				"type": fieldTag,
			}
		}
	}
	result := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": properties,
		},
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
