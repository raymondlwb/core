package scheduler

import (
	"testing"

	"github.com/projecteru2/core/scheduler/complex"
	"github.com/projecteru2/core/types"
	"github.com/stretchr/testify/assert"
)

func TestSchedulerInvoke(t *testing.T) {
	coreCfg := types.Config{
		EtcdMachines:   []string{"http://127.0.0.1:2379"},
		EtcdLockPrefix: "/eru-core/_lock",
		Scheduler: types.SchedConfig{
			ShareBase: 10,
			MaxShare:  -1,
		},
	}
	scheduler, _ := complexscheduler.New(coreCfg)

	nodes := []types.NodeInfo{
		types.NodeInfo{
			CPUAndMem: types.CPUAndMem{
				CpuMap: types.CPUMap{
					"0": 10,
					"1": 10,
				},
				MemCap: 12400000,
			},
			Name: "node1",
		},
		types.NodeInfo{
			CPUAndMem: types.CPUAndMem{
				CpuMap: types.CPUMap{
					"0": 10,
					"1": 10,
				},
				MemCap: 12400000,
			},
			Name: "node2",
		},
	}
	_, _, err := scheduler.SelectCPUNodes(nodes, 1, 2)
	assert.NoError(t, err)
}
