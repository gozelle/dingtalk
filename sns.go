package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

type SNSClient struct {
	client *Client
}

type SNSUserInfo struct {
	Nick                 string `json:"nick"`
	UnionId              string `json:"unionid"`
	OpenId               string `json:"openid"`
	MainOrgAuthHighLevel bool   `json:"main_org_auth_high_level"`
}

func EncodeSHA256(message, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := h.Sum(nil)
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

// GetUserInfoByCode 根据 sns 临时授权码获取用户信息
// https://open.dingtalk.com/document/orgapp-server/obtain-the-user-information-based-on-the-sns-temporary-authorization
func (s *SNSClient) GetUserInfoByCode(tmpAuthCode string) (userinfo *SNSUserInfo, err error) {

	req := s.client.newRestyRequest()

	req.SetBody(map[string]string{
		"tmp_auth_code": tmpAuthCode,
	})

	// 准备签名
	ts := time.Now().UnixMilli()
	s1 := EncodeSHA256(fmt.Sprintf("%d", ts), s.client.appSecret)

	// 注入参数
	req.QueryParam.Add("accessKey", s.client.appKey)
	req.QueryParam.Add("signature", s1)
	req.QueryParam.Add("timestamp", fmt.Sprintf("%d", ts))

	resp, err := req.Post(s.client.url("/sns/getuserinfo_bycode"))
	if err != nil {
		return
	}

	data, err := extractResult(resp)
	if err != nil {
		return
	}

	userinfo = new(SNSUserInfo)
	err = json.Unmarshal([]byte(gjson.Parse(string(data)).Get("user_info").Raw), userinfo)
	if err != nil {
		return
	}

	return
}
