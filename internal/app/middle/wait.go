package middle

import (
	"sync"
	"time"

	"github.com/lzh-1625/go_process_manager/config"

	"github.com/gin-gonic/gin"
)

type waitCond struct {
	cond    sync.Cond
	ts      int64
	timeMap sync.Map
	trigger chan struct{}
}

var (
	ProcessWaitCond = newWaitCond()
	TaskWaitCond    = newWaitCond()
)

var waitCondList []*waitCond

func InitCondTiming() {
	for _, v := range waitCondList {
		go v.timing()
	}
}

func newWaitCond() *waitCond {
	wc := &waitCond{
		cond:    *sync.NewCond(&sync.Mutex{}),
		ts:      time.Now().UnixMicro(),
		timeMap: sync.Map{},
		trigger: make(chan struct{}),
	}
	waitCondList = append(waitCondList, wc)
	return wc
}

func (p *waitCond) Trigger() {
	p.trigger <- struct{}{}
	p.ts = time.Now().UnixMicro()
}

func (p *waitCond) WaitGetMiddel(c *gin.Context) {
	reqUser := c.GetHeader("token")
	defer p.timeMap.Store(reqUser, p.ts)
	if ts, ok := p.timeMap.Load(reqUser); !ok || ts.(int64) > p.ts {
		c.Next()
		return
	}
	p.cond.L.Lock()
	defer p.cond.L.Unlock()
	p.cond.Wait()
	c.Next()
}

func (p *waitCond) WaitTriggerMiddel(c *gin.Context) {
	defer p.Trigger()
	c.Next()
}

func (p *waitCond) timing() { // 添加定时信号清理阻塞协程
	ticker := time.NewTicker(time.Second * time.Duration(config.CF.CondWaitTime))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
		case <-p.trigger:
		}
		ticker.Reset(time.Second * time.Duration(config.CF.CondWaitTime))
		p.cond.Broadcast()
	}
}
