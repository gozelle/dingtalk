package dingtalk

import (
	"fmt"
	"sync"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gozelle/resty"
	"github.com/tidwall/gjson"
)

func NewClient(agentId int64, appKey string, appSecret string, proxyUrl string) *Client {
	return &Client{
		agentId:     agentId,
		appKey:      appKey,
		appSecret:   appSecret,
		restyClient: resty.New(),
		lock:        &sync.Mutex{},
		proxyUrl:    &proxyUrl,
	}
}

// Client 钉钉客户端封装
// 全局使用一个实例即可，以避免 AccessToken 过度请求而被限流。
type Client struct {
	agentId          int64
	appKey           string
	appSecret        string
	accessToken      string
	expiredAt        time.Time
	lock             *sync.Mutex // for accessToken and expiredAt
	restyClient      *resty.Client
	oAuth2Client     *oauth2_1_0.Client
	departmentClient *DepartmentClient
	userClient       *UserClient
	snsClient        *SNSClient
	messageClient    *MessageClient
	processClient    *ProcessClient
	proxyUrl         *string // 真实接口地址： https://oapi.dingtalk.com  使用 octopus 代理
}

func (c *Client) url(uri string) string {
	return fmt.Sprintf("%s%s", *c.proxyUrl, uri)
}

func (c *Client) SNSClient() *SNSClient {
	if c.snsClient == nil {
		c.snsClient = &SNSClient{
			client: c,
		}
	}
	return c.snsClient
}

func (c *Client) DepartmentClient() *DepartmentClient {
	if c.departmentClient == nil {
		c.departmentClient = &DepartmentClient{
			client: c,
		}
	}
	return c.departmentClient
}

func (c *Client) UserClient() *UserClient {
	if c.userClient == nil {
		c.userClient = &UserClient{
			client: c,
		}
	}
	return c.userClient
}

func (c *Client) MessageClient() *MessageClient {
	if c.messageClient == nil {
		c.messageClient = &MessageClient{
			client: c,
		}
	}
	return c.messageClient
}

func (c *Client) ProcessClient() *ProcessClient {
	if c.processClient == nil {
		c.processClient = &ProcessClient{
			client: c,
		}
	}
	return c.processClient
}

//func (c *Client) execute(method, url string, req *resty.Request, parser pavo.Parser) *pavo.APIExecutor {
//
//	resp, err := req.Post(url)
//	if err != nil {
//		return
//	}
//
//	return c.pavo.Request(method, "DingTalk", url, req).
//		WithBeforeHandler(func(req *resty.Request) error {
//			return c.wrapAccessToken(req)
//		}).
//		WithAfterHandler(c.extractResult).
//		WithParser(parser)
//}

func extractResult(resp *resty.Response) (data []byte, err error) {

	defer func() {
		if err != nil && resp.Request != nil {
			err = fmt.Errorf("[%s]%s error: %s", resp.Request.Method, resp.Request.URL, err)
		}
	}()

	if resp.IsError() {
		err = fmt.Errorf(resp.Status())
		return
	}

	gr := gjson.ParseBytes(resp.Body())
	if !gr.IsObject() {
		err = fmt.Errorf("illegal data")
		return
	}

	if gr.Get("errcode").Exists() {
		if code := gr.Get("errcode").Raw; code != "0" {
			err = fmt.Errorf("response error: %s(%s)", code, gr.Get("errmsg"))
			return
		}
	} else if gr.Get("processQueryKey").Exists() {
		data = []byte(gr.Raw)
		return
	}

	data = []byte(gr.Get("result").Raw)

	return
}

func (c *Client) wrapAccessToken(req *resty.Request) error {
	token, err := c.AccessToken()
	if err != nil {
		return err
	}
	req.QueryParam.Add("access_token", token)
	req.SetHeader("x-acs-dingtalk-access-token", token)
	return nil
}

func (c *Client) newRestyRequest() *resty.Request {
	return c.restyClient.R()
}

// AccessToken 获取
func (c *Client) AccessToken() (token string, err error) {

	c.lock.Lock()
	defer func() {
		c.lock.Unlock()
	}()

	if c.accessToken != "" && time.Now().Before(c.expiredAt) {
		token = c.accessToken
		return
	}

	c.oAuth2Client, err = createOAuth2Client()
	if err != nil {
		err = fmt.Errorf("new dingtalk oath2 client error: %s", err)
		return
	}

	var request = &oauth2_1_0.GetAccessTokenRequest{
		AppKey:    tea.String(c.appKey),
		AppSecret: tea.String(c.appSecret),
	}
	var resp *oauth2_1_0.GetAccessTokenResponse

	resp, err = c.oAuth2Client.GetAccessToken(request)
	if err != nil {
		err = fmt.Errorf("request error: %s", err)
		return
	}

	if resp.Body == nil {
		err = fmt.Errorf("resp.Body is nil")
		return
	}

	c.accessToken = tea.StringValue(resp.Body.AccessToken)
	c.expiredAt = time.Now().Add(time.Duration(tea.Int64Value(resp.Body.ExpireIn)) * time.Second)

	token = c.accessToken

	return
}

func createOAuth2Client() (c *oauth2_1_0.Client, err error) {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	return oauth2_1_0.NewClient(config)
}
