package dingtalk_test

import "github.com/gozelle/dingtalk"

func NewTestClient() *dingtalk.Client {
	return dingtalk.NewClient(testAgentI, testKey, testSecret, testProxyUrl)
}
