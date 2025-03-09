package logic

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/middle"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
)

type clusterLogic struct {
	data        sync.Map
	ClusterMode constants.ClusterMode
	UpdateChan  chan struct{}
	clusterConn *websocket.Conn
}

var ClusterLogic = new(clusterLogic)

func (c *clusterLogic) Init() {
	middle.ProcessWaitCond.RegisterFunc(func() {
		c.UpdateChan <- struct{}{}
	})

}

func (c *clusterLogic) UploadData() {
	for {
		<-c.UpdateChan
		ProcessCtlLogic.GetProcessList()

	}
}

func (c *clusterLogic) ReciveData(nodeName string, data []model.ProcessInfo) {
	c.data.Store(nodeName, data)
}
