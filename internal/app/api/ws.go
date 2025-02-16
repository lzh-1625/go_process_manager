package api

import (
	"context"
	"sync"
	"time"

	"github.com/lzh-1625/go_process_manager/config"
	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/logic"
	"github.com/lzh-1625/go_process_manager/internal/app/middle"
	"github.com/lzh-1625/go_process_manager/log"
	"github.com/lzh-1625/go_process_manager/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type wsApi struct{}

var WsApi = new(wsApi)

type WsConnetInstance struct {
	WsConnect  *websocket.Conn
	wsLock     sync.Mutex
	CancelFunc context.CancelFunc
}

func (w *WsConnetInstance) Write(b []byte) {
	w.wsLock.Lock()
	defer w.wsLock.Unlock()
	w.WsConnect.WriteMessage(websocket.BinaryMessage, b)
}

func (w *WsConnetInstance) WriteString(s string) {
	w.wsLock.Lock()
	defer w.wsLock.Unlock()
	w.WsConnect.WriteMessage(websocket.TextMessage, []byte(s))
}
func (w *WsConnetInstance) Cancel() {
	w.CancelFunc()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (w *wsApi) WebsocketHandle(ctx *gin.Context) {
	reqUser := getUserName(ctx)
	uuid := getQueryInt(ctx, "uuid")
	proc, err := logic.ProcessCtlLogic.GetProcess(uuid)
	errCheck(ctx, err != nil, "Operation failed!")
	errCheck(ctx, proc.HasWsConn(reqUser), "A connection already exists; unable to establish a new one!")
	errCheck(ctx, proc.State.State != 1, "The process is currently running.")
	errCheck(ctx, proc.Control.Controller != reqUser && !proc.VerifyControl(), "Insufficient permissions; please check your access rights!")
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	errCheck(ctx, err != nil, "WebSocket connection upgrade failed!")

	log.Logger.AddAdditionalInfo("processName", proc.Name)
	log.Logger.AddAdditionalInfo("userName", reqUser)
	defer log.Logger.DeleteAdditionalInfo(2)

	log.Logger.Infow("ws连接成功")

	proc.SetTerminalSize(utils.GetIntByString(ctx.Query("cols")), utils.GetIntByString(ctx.Query("rows")))
	wsCtx, cancel := context.WithCancel(context.Background())
	wci := &WsConnetInstance{
		WsConnect:  conn,
		CancelFunc: cancel,
		wsLock:     sync.Mutex{},
	}
	proc.ReadCache(wci)
	w.startWsConnect(wci, cancel, proc, hasOprPermission(ctx, uuid, constants.OPERATION_TERMINAL_WRITE))
	proc.AddConn(reqUser, wci)
	defer middle.ProcessWaitCond.Trigger()
	defer proc.DeleteConn(reqUser)
	conn.SetCloseHandler(func(_ int, _ string) error {
		middle.ProcessWaitCond.Trigger()
		cancel()
		return nil
	})
	middle.ProcessWaitCond.Trigger()
	select {
	case <-proc.StopChan:
		log.Logger.Infow("ws连接断开", "操作类型", "进程已停止，强制断开ws连接")
	case <-time.After(time.Minute * time.Duration(config.CF.TerminalConnectTimeout)):
		log.Logger.Infow("ws连接断开", "操作类型", "连接时间超过最大时长限制")
	case <-wsCtx.Done():
		log.Logger.Infow("ws连接断开", "操作类型", "tcp连接建立已被关闭")
	}
	conn.Close()
}

func (w *wsApi) startWsConnect(wci *WsConnetInstance, cancel context.CancelFunc, proc logic.Process, write bool) {
	log.Logger.Debugw("ws读取线程已启动")
	go func() {
		for {
			_, b, err := wci.WsConnect.ReadMessage()
			if err != nil {
				log.Logger.Debugw("ws读取线程已退出", "info", err)
				return
			}
			if write {
				proc.WriteBytes(b)
			}
		}
	}()

	// proactive health check
	pongChan := make(chan struct{})
	wci.WsConnect.SetPongHandler(func(appData string) error {
		pongChan <- struct{}{}
		return nil
	})
	timer := time.NewTicker(time.Second)
	go func() {
		defer timer.Stop()
		for {
			wci.wsLock.Lock()
			wci.WsConnect.WriteMessage(websocket.PingMessage, nil)
			wci.wsLock.Unlock()
			select {
			case <-timer.C:
				cancel()
				return
			case <-pongChan:
			}
			time.Sleep(time.Second * time.Duration(config.CF.WsHealthCheckInterval))
			timer.Reset(time.Second)
		}
	}()

}
