package boot

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"

	"github.com/lzh-1625/go_process_manager/config"
	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/middle"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"
	"github.com/lzh-1625/go_process_manager/internal/app/service"
	"github.com/lzh-1625/go_process_manager/internal/app/termui"
	logger "github.com/lzh-1625/go_process_manager/log"
	"github.com/lzh-1625/go_process_manager/utils"
)

func init() {
	initDb()
	initResetConfig()
	initConfiguration()
	initArgs()
	initLogHandle()
	initLog()
	initEs()
	initLogHanler()
	initCondTiming()
	initProcess()
	initJwtSecret()
	initTui()
	InitTask()
	initListenKillSignal()
}

func initDb() {
	repository.InitDb()
}

func initConfiguration() {
	defer func() {
		if err := recover(); err != nil {
			panic("config init fail")
		}
	}()
	typeElem := reflect.TypeOf(config.CF).Elem()
	valueElem := reflect.ValueOf(config.CF).Elem()
	for i := 0; i < typeElem.NumField(); i++ {
		typeField := typeElem.Field(i)
		valueField := valueElem.Field(i)
		value, err := repository.ConfigRepository.GetConfigValue(typeField.Name)
		if err != nil {
			value = typeField.Tag.Get("default")
		}
		if value == "-" {
			continue
		}
		switch typeField.Type.Kind() {
		case reflect.String:
			valueField.SetString(value)
		case reflect.Bool:
			valueField.SetBool(utils.UnwarpIgnore(strconv.ParseBool(value)))
		case reflect.Float64:
			valueField.SetFloat(utils.UnwarpIgnore(strconv.ParseFloat(value, 64)))
		case reflect.Int64, reflect.Int:
			valueField.SetInt(utils.UnwarpIgnore(strconv.ParseInt(value, 10, 64)))
		default:
			continue
		}
	}
}

func initArgs() {
	if len(os.Args) >= 2 && os.Args[1] == "tui" {
		config.CF.UserTui = true
	}
}

func initLog() {
	logger.InitLog()
}

func initEs() {
	service.EsService.InitEs()
}

func initProcess() {
	service.ProcessCtlService.ProcessInit()
}

func initJwtSecret() {
	if secret, err := repository.ConfigRepository.GetConfigValue(constants.SECRET_KEY); err == nil {
		utils.SetSecret([]byte(secret))
		return
	}
	secret := utils.RandString(32)
	repository.ConfigRepository.SetConfigValue(constants.SECRET_KEY, secret)
	utils.SetSecret([]byte(secret))
}

func initLogHanler() {
	service.InitLog()
}

func initTui() {
	go termui.Tui.TermuiInit()
}

func InitTask() {
	service.TaskService.InitTaskJob()
}

func initResetConfig() {
	if len(os.Args) >= 2 && os.Args[1] == "reset" {
		err := service.ConfigService.ResetSystemConfiguration()
		if err != nil {
			log.Panic(err)
		}
		log.Print("reset system config to deafult success!")
		os.Exit(0)
	}
}

func initListenKillSignal() {
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		service.ProcessCtlService.KillAllProcess()
		log.Print("已停止所有进程")
		os.Exit(0)
	}()
}

func initCondTiming() {
	middle.InitCondTiming()
}

func initLogHandle() {
	service.InitLogHandle()
}
