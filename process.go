package dingtalk

import (
	"encoding/json"
	
	"github.com/tidwall/gjson"
)

type ProcessClient struct {
	client *Client
}

type InstanceIdsRequest struct {
	ProcessCode string  `json:"process_code"`
	StartTime   int64   `json:"start_time"` // 毫秒
	EndTime     int64   `json:"end_time"`   // 毫秒
	Size        int     `json:"size"`
	Cursor      *int64  `json:"cursor"`
	UseridList  *string `json:"userid_list"`
}

type InstanceIdsReply struct {
	List       []string `json:"list"`
	NextCursor int64    `json:"next_cursor"`
}

// InstanceIds 查询审批流实例列表
// https://open.dingtalk.com/document/orgapp-server/operation-to-retrieve-a-list-of
func (p *ProcessClient) InstanceIds(params *InstanceIdsRequest) (reply *InstanceIdsReply, err error) {
	
	req := p.client.newRestyRequest()
	req.SetBody(params)
	
	err = p.client.wrapAccessToken(req)
	if err != nil {
		return
	}
	
	resp, err := req.Post(p.client.url("/topapi/processinstance/listids"))
	if err != nil {
		return
	}
	data, err := extractResult(resp)
	if err != nil {
		return
	}
	
	reply = new(InstanceIdsReply)
	err = json.Unmarshal(data, reply)
	if err != nil {
		return
	}
	
	return
}

// Instance 获取单个审批实例详情
func (p *ProcessClient) Instance(id string) (reply string, err error) {
	req := p.client.newRestyRequest()
	req.SetBody(map[string]string{
		"process_instance_id": id,
	})
	err = p.client.wrapAccessToken(req)
	if err != nil {
		return
	}
	resp, err := req.Post(p.client.url("/topapi/processinstance/get"))
	if err != nil {
		return
	}
	
	reply = gjson.Parse(resp.String()).Get("process_instance").Raw
	
	return
}
