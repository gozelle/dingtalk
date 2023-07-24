package processes

import (
	"sort"
	"testing"
	"time"
	
	"github.com/gozelle/dingtalk"
	"github.com/gozelle/dingtalk/test"
	user "github.com/gozelle/dingtalk/user"
	"github.com/gozelle/spew"
	"github.com/gozelle/testify/require"
)

func TestPorcessInstances(t *testing.T) {
	client := test.NewTestClient()
	
	ids, err := client.ProcessClient().InstanceIds(&dingtalk.InstanceIdsRequest{
		ProcessCode: "PROC-755AFE69-7834-407C-9060-84C673BC9876",
		StartTime:   time.Now().Add(-30 * 24 * time.Hour).UnixMilli(),
		EndTime:     time.Now().UnixMilli(),
		Size:        1,
	})
	
	require.NoError(t, err)
	spew.Json(ids)
}

func TestProcessInstanceDetail(t *testing.T) {
	
	var err error
	client := test.NewTestClient()
	
	um := user.NewUserManger(client)
	
	pm := NewProcessManger("PROC-755AFE69-7834-407C-9060-84C673BC9876", um, client)
	
	item, err := pm.Instance("3r43pRm1TDunVy_yeotMhQ02101689066368")
	require.NoError(t, err)
	
	spew.Json(item)
	
	for _, v := range item.Tasks {
		if v.TaskStatus != "COMPLETED" {
			continue
		}
		var u *dingtalk.DepartmentUser
		u, err = um.GetUser(v.Userid)
		require.NoError(t, err)
		v.Username = u.Name
		
		for _, vv := range item.OperationRecords {
			if vv.Date > v.CreateTime && vv.Date <= v.FinishTime {
				v.OperationRecords = append(v.OperationRecords, vv)
			}
		}
		sort.Sort(v.OperationRecords)
	}
	
	sort.Sort(item.Tasks)
	
	for _, v := range item.OperationRecords {
		var u *dingtalk.DepartmentUser
		u, err = um.GetUser(v.Userid)
		require.NoError(t, err)
		if u != nil {
			v.Username = u.Name
		}
	}
	
	spew.Json(item)
	
}
