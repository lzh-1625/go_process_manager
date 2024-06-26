package boot

import (
	"msm/config"
	"msm/dao"
	"msm/log"
	"msm/service/es"
	"msm/service/process"
	"msm/utils"
	"reflect"
	"strconv"
)

func Boot() {
	initConfiguration()
	initEs()
	initProcess()
}

func initConfiguration() {
	typeElem := reflect.TypeOf(config.CF).Elem()
	valueElem := reflect.ValueOf(config.CF).Elem()
	for i := 0; i < typeElem.NumField(); i++ {
		typeField := typeElem.Field(i)
		valueField := valueElem.Field(i)
		value, err := dao.ConfigDao.GetConfigValue(typeField.Name)
		if err != nil {
			value = typeField.Tag.Get("default")
		}
		switch typeField.Type.Kind() {
		case reflect.String:
			valueField.SetString(value)
		case reflect.Bool:
			valueField.SetBool(utils.Unwarp(strconv.ParseBool(value)))
		case reflect.Float64:
			valueField.SetFloat(utils.Unwarp(strconv.ParseFloat(value, 64)))
		case reflect.Int64, reflect.Int:
			valueField.SetInt(utils.Unwarp(strconv.ParseInt(value, 10, 64)))
		default:
			continue
		}
	}
	log.Logger.Debugw("获取配置信息完成", "Configuration", config.CF)
}

func initEs() {
	es.InitEs()
}

func initProcess() {
	process.ProcessCtlService.ProcessInit()
}
