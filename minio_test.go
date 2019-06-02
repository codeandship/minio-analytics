package minioanalytics_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	minioanalytics "github.com/codeandship/minio-analytics"
	"github.com/stretchr/testify/require"
)

func TestEvent_Unmarshal(t *testing.T) {
	data, err := ioutil.ReadFile("events.json")
	require.NoError(t, err)
	var res []minioanalytics.MinioEvent
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	t.Logf("%+v", res)
}

func Test_MapAnalytics(t *testing.T) {
	data, err := ioutil.ReadFile("analytics.json")
	require.NoError(t, err)
	var res []minioanalytics.Analytics
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	mp := minioanalytics.MapAnayltics(res)
	data, err = json.MarshalIndent(mp, "", "    ")
	require.NoError(t, err)
	t.Log(string(data))

}
