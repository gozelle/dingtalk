package dingtalk

import (
	"encoding/json"
	"fmt"
)

type MessageClient struct {
	client *Client
}

type SendMsgByRobotReq struct {
	RobotCode string   `json:"robotCode"`
	UserIDs   []string `json:"userIds"`
	MsgKey    string   `json:"msgKey"`
	MsgParam  string   `json:"msgParam"`
}

type SendMsgByRobotResp struct {
	Code                      string   `json:"code,omitempty"`
	ReqID                     string   `json:"requestid,omitempty"`
	Message                   string   `json:"message,omitempty"`
	ProcessQueryKey           string   `json:"processQueryKey,omitempty"`           // 消息id
	InvalidStaffIdList        []string `json:"invalidStaffIdList,omitempty"`        // 无效的用户userid列表。
	FlowControlledStaffIdList []string `json:"flowControlledStaffIdList,omitempty"` // 被限流的userid列表。
}

type MsgContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// SendRobotMessage 机器人发送消息
// https://open.dingtalk.com/document/group/chatbots-send-one-on-one-chat-messages-in-batches
func (d *MessageClient) SendRobotMessage(title, content string, to []string) (err error) {
	req := d.client.newRestyRequest()

	msg, err := json.Marshal(&MsgContent{Title: title, Text: content})
	if err != nil {
		return fmt.Errorf("生成消息失败: %v", err)
	}

	robotCode := d.client.appKey
	reqObj := &SendMsgByRobotReq{
		RobotCode: robotCode,
		UserIDs:   to,
		//MsgKey:    "officialMarkdownMsg",
		MsgKey:   "sampleMarkdown",
		MsgParam: string(msg),
	}
	req.SetBody(reqObj)

	_, err = req.Post("https://api.dingtalk.com/v1.0/robot/oToMessages/batchSend")
	if err != nil {
		return
	}

	return
}

// SendUserMessage 给用户发送消息
func (d *MessageClient) SendUserMessage(agentId int64, userId string, msg string) (err error) {

	req := d.client.newRestyRequest()
	req.SetBody(map[string]interface{}{
		"agent_id":    agentId,
		"userid_list": userId,
		"msg": map[string]interface{}{
			"msgtype": "text",
			"text": map[string]string{
				"content": msg,
			},
		},
	})
	err = d.client.wrapAccessToken(req)
	if err != nil {
		return
	}

	_, err = req.Post(d.client.url("/topapi/message/corpconversation/asyncsend_v2"))
	if err != nil {
		return
	}
	return
}
