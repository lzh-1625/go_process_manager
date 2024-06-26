package api

import (
	"msm/config"
	"msm/dao"
	"msm/model"
	"msm/service/es"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

type configApi struct{}

var ConfigApi = new(configApi)

func (c *configApi) GetSystemConfiguration(ctx *gin.Context) {
	typeElem := reflect.TypeOf(config.CF).Elem()
	valueElem := reflect.ValueOf(config.CF).Elem()
	result := []model.SystemConfigurationResp{}
	for i := 0; i < typeElem.NumField(); i++ {
		typeField := typeElem.Field(i)
		valueField := valueElem.Field(i)
		var value any
		switch typeField.Type.Kind() {
		case reflect.Int64, reflect.Int:
			value = valueField.Int()
		case reflect.String:
			value = valueField.String()
		case reflect.Bool:
			value = valueField.Bool()
		case reflect.Float64:
			value = valueField.Float()
		default:
			continue
		}
		result = append(result, model.SystemConfigurationResp{
			Key:      typeField.Name,
			Value:    value,
			Default:  typeField.Tag.Get("default"),
			Describe: typeField.Tag.Get("describe"),
		})
	}
	rOk(ctx, "获取系统配置成功", result)
}

func (c *configApi) SetSystemConfiguration(ctx *gin.Context) {
	data := map[string]string{}
	errCheck(ctx, ctx.ShouldBindJSON(&data) != nil, "请求参数错误")
	typeElem := reflect.TypeOf(config.CF).Elem()
	valueElem := reflect.ValueOf(config.CF).Elem()
	for i := 0; i < typeElem.NumField(); i++ {
		typeField := typeElem.Field(i)
		valueField := valueElem.Field(i)
		for k, v := range data {
			if typeField.Name == k {
				var err error
				switch typeField.Type.Kind() {
				case reflect.String:
					valueField.SetString(v)
				case reflect.Bool:
					value, errV := strconv.ParseBool(v)
					err = errV
					if err == nil {
						valueField.SetBool(value)
					}
				case reflect.Float64:
					value, errV := strconv.ParseFloat(v, 64)
					err = errV
					if err == nil {
						valueField.SetFloat(value)
					}
				case reflect.Int64, reflect.Int:
					value, errV := strconv.ParseInt(v, 10, 64)
					err = errV
					if err == nil {
						valueField.SetInt(value)
					}
				default:
					continue
				}
				errCheck(ctx, err != nil, k+"类似错误")
				errCheck(ctx, dao.ConfigDao.SetConfigValue(k, v) != nil, "修改配置失败")
			}
		}
	}
	rOk(ctx, "修改配置成功", nil)
}

func (c *configApi) EsConfigReload(ctx *gin.Context) {
	errCheck(ctx, !es.InitEs(), "es连接失败，请检查是否启用es或账号密码是否存在错误")
	rOk(ctx, "已连接上es", nil)
}
