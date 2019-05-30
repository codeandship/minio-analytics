package minioanalytics_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	minioanalytics "github.com/iwittkau/minio-analytics"
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
