package api

import (
	"context"
	"msm/consts/ctxflag"
	"msm/log"
	"msm/service/process"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type wsApi struct{}

var WsApi = new(wsApi)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (w *wsApi) WebsocketHandle(ctx *gin.Context) {
	reqUser := ctx.GetString(ctxflag.USER_NAME)
	uuid, err := strconv.Atoi(ctx.Query("uuid"))
	errCheck(ctx, err != nil, "参数有误")
	proc, err := process.ProcessCtlService.GetProcess(uuid)
	errCheck(ctx, err != nil, "进程获取失败")
	errCheck(ctx, proc.GetStateState() != 1, "进程未运行")
	errCheck(ctx, proc.GetControlController() != reqUser && !proc.VerifyControl(), "进程权限不足")
	errCheck(ctx, !proc.TryLock(), "进程已被占用")
	proc.SetWhoUsing(reqUser)
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	errCheck(ctx, err != nil, "ws升级失败")
	log.Logger.Infow("ws连接成功", "进程名称", proc.GetName(), "连接者", proc.GetWhoUsing())
	proc.SetControlController("")
	wsCtx, cancel := context.WithCancel(context.Background())
	w.startWsConnect(conn, proc, cancel)
	proc.SetWsConn(conn)
	proc.SetIsUsing(true)
	close := func(err string) {
		proc.SetWhoUsing("")
		proc.SetIsUsing(false)
		proc.SetWsConn(nil)
		conn.Close()
		proc.Unlock()
		log.Logger.Infow("ws连接断开", "操作类型", err, "进程名称", proc.GetName())
	}
	conn.SetCloseHandler(func(_ int, _ string) error {
		proc.ChangControlChan() <- 1
		close("ws连接被断开")
		return nil
	})
	select {
	case signal := <-proc.ChangControlChan():
		{
			if signal == 0 {
				close("强制断开ws连接")
			}
		}
	case <-proc.StopChan():
		{
			close("进程已停止，强制断开ws连接")
		}
	case <-time.After(time.Minute * 10):
		{
			close("连接时间超过最大时长限制")
		}
	case <-wsCtx.Done():
		{
			close("tcp连接建立已被关闭")
		}
	}

}

func (w *wsApi) startWsConnect(conn *websocket.Conn, proc process.Process, cancel context.CancelFunc) {
	proc.ReadCache(conn)
	log.Logger.Debugw("ws读取线程已启动")
	go func() {
		for {
			_, b, err := conn.ReadMessage()
			if err != nil {
				log.Logger.Debugw("ws读取线程已退出", "info", err)
				cancel()
				return
			}
			proc.WriteBytes(b)
		}
	}()
}
