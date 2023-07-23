package dingtalk_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	client := NewTestClient()
	err := client.MessageClient().SendUserMessage(1401579096, "016961205832722717", "Hello!")
	require.NoError(t, err)
}
