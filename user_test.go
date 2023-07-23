package dingtalk_test

import (
	"encoding/json"
	"testing"

	"github.com/gozelle/dingtalk"
	"github.com/gozelle/dingtalk/test"
	"github.com/stretchr/testify/require"
)

func TestUserClient(t *testing.T) {
	client := test.NewTestClient()
	list, err := client.UserClient().ListIDs(1)
	require.NoError(t, err)
	d, err := json.MarshalIndent(list, "", "\t")
	require.NoError(t, err)
	t.Log(string(d))

	for _, id := range list {
		var user *dingtalk.User
		user, err = client.UserClient().UserGet(dingtalk.UserGetParams{
			UserId: id,
		})
		require.NoError(t, err)
		var dd []byte
		dd, err = json.MarshalIndent(user, "", "\t")
		t.Log(string(dd))
	}
}
