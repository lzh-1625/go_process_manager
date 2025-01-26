package cgroup_test

import (
	"testing"
	"time"

	_ "github.com/lzh-1625/go_process_manager/boot"
	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/internal/app/service"

	"github.com/containerd/cgroups/v3/cgroup1"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func TestCgroup(t *testing.T) {
	period := uint64(100000) // 100ms = 100000微秒
	// 设置 CPU 配额为 20% (20ms)
	quota := int64(20000 * 8) // 20ms = 20000微秒
	control, err := cgroup1.New(cgroup1.StaticPath("/test"), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Period: &period,
			Quota:  &quota,
		},
		Memory: &specs.LinuxMemory{},
	})
	if err != nil {
		panic(err)
	}
	defer control.Delete()
	p, err := service.ProcessCtlService.RunNewProcess(model.Process{
		Name:     "test",
		Cmd:      "bash",
		Cwd:      `/root`,
		TermType: constants.TERMINAL_PTY,
	})
	if err != nil {
		t.FailNow()
	}
	control.AddProc(uint64(p.Pid))
	time.Sleep(time.Second * 20)
	p.Kill()
	control.Delete()
}
