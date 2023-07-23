package process

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/gozelle/dingtalk"
	user "github.com/gozelle/dingtalk/user"
)

const (
	RUNNING   = "RUNNING"
	COMPLETED = "COMPLETED"
	CANCELED  = "CANCELED"
)

type IProcess interface {
	Iterate(start, end time.Time, handler Handler) (err error)
	Instance(instanceId string) (item *Instance, err error)
}

func NewProcess(processCode string, client *dingtalk.Client, users *user.Users) *Process {
	return &Process{processCode: processCode, client: client, users: users}
}

var _ IProcess = (*Process)(nil)

type Process struct {
	processCode string
	client      *dingtalk.Client
	users       *user.Users
}

type Handler func(instances []string) error

func (p Process) Iterate(start, end time.Time, handler Handler) (err error) {
	var cursor *int64
	for {
		var r *dingtalk.InstanceIdsReply
		r, err = p.client.ProcessClient().InstanceIds(&dingtalk.InstanceIdsRequest{
			ProcessCode: p.processCode,
			StartTime:   start.UnixMilli(),
			EndTime:     end.UnixMilli(),
			Size:        20,
			Cursor:      cursor,
			UseridList:  nil,
		})
		if err != nil {
			return
		}
		err = handler(r.List)
		if err != nil {
			return
		}
		if r.NextCursor == 0 {
			break
		}
		cursor = &r.NextCursor
	}
	return
}

func (p Process) Instance(instanceId string) (item *Instance, err error) {
	r, err := p.client.ProcessClient().Instance(instanceId)
	if err != nil {
		return
	}
	item = new(Instance)
	err = json.Unmarshal([]byte(r), item)
	if err != nil {
		return
	}
	for _, v := range item.Tasks {
		var u *dingtalk.DepartmentUser
		u, err = p.users.GetUserOrNil(v.Userid)
		if err != nil {
			return
		}
		if u != nil {
			v.Username = u.Name
		}
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
		u, err = p.users.GetUserOrNil(v.Userid)
		if err != nil {
			return
		}
		if u != nil {
			v.Username = u.Name
		}
	}

	return
}
