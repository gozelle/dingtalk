package dingtalk_test

import (
	"testing"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/gozelle/dingtalk"
	"github.com/gozelle/pointer"
	"github.com/gozelle/spew"
	"github.com/stretchr/testify/require"
)

func TestProcessIds(t *testing.T) {

	client := NewTestClient()

	r, err := client.ProcessClient().InstanceIds(&dingtalk.InstanceIdsRequest{
		ProcessCode: "PROC-6F6C034B-0C79-4D56-A290-94D678FB011C",
		StartTime:   1675008000000,
		EndTime:     1675180800000,
		Size:        20,
		Cursor:      0,
		UseridList:  pointer.ToString("016961205832722717"),
	})
	require.NoError(t, err)

	spew.Json(r)
}

func TestProcessInstance(t *testing.T) {
	client := NewTestClient()

	r, err := client.ProcessClient().Instance("XJ8cI0CaR5CrCq3ZQbFL0A02101675174858")
	require.NoError(t, err)
	spew.Json(r)
}

func TestProcessIds2(t *testing.T) {

	client := NewTestClient()

	start, err := time.Parse(carbon.DateTimeLayout, "2023-03-24 09:00:00")
	require.NoError(t, err)
	t.Log(start.UnixMilli())
	r, err := client.ProcessClient().InstanceIds(&dingtalk.InstanceIdsRequest{
		ProcessCode: "PROC-6F6C034B-0C79-4D56-A290-94D678FB011C",
		StartTime:   start.UnixMilli(),
		EndTime:     start.Add(24 * time.Hour).UnixMilli(),
		Size:        20,
		Cursor:      0,
		UseridList:  pointer.ToString("3242435424239462"),
	})
	require.NoError(t, err)

	spew.Json(r)
}

func TestProcessInstance2(t *testing.T) {
	client := NewTestClient()
	r, err := client.ProcessClient().Instance("baJPkl8-RgquXxmcqGoxdA02101678844941")
	require.NoError(t, err)
	t.Log(r)
}
