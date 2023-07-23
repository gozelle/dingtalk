package process

import (
	"sort"
	"testing"

	"github.com/gozelle/dingtalk"
	"github.com/gozelle/dingtalk/test"
	user "github.com/gozelle/dingtalk/user"
	"github.com/gozelle/spew"
	"github.com/gozelle/testify/require"
)

func TestParseJson(t *testing.T) {

	var err error
	item := new(Instance)
	client := test.NewTestClient()

	um := user.NewUsers(client)

	for _, v := range item.Tasks {
		if v.TaskStatus != "COMPLETED" {
			continue
		}
		var u *dingtalk.DepartmentUser
		u, err = um.GetUserOrNil(v.Userid)
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
		u, err = um.GetUserOrNil(v.Userid)
		require.NoError(t, err)
		v.Username = u.Name
	}

	spew.Json(item)

}
