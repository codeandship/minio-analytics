package minioanalytics_test

import (
	"reflect"
	"regexp"
	"testing"

	minioanalytics "github.com/codeandship/minio-analytics"
	"github.com/stretchr/testify/require"
)

func TestNewUserAgentMater(t *testing.T) {
	tests := []struct {
		name    string
		want    *minioanalytics.UserAgentMatcher
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := minioanalytics.NewUserAgentMater()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserAgentMater() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserAgentMater() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserAgentMatcher_Match(t *testing.T) {
	type fields struct {
		matches map[string]*regexp.Regexp
	}
	type args struct {
		ua string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uam, err := minioanalytics.NewUserAgentMater()
			require.NoError(t, err)
			if got := uam.Match(tt.args.ua); got != tt.want {
				t.Errorf("UserAgentMatcher.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserAgentMatcher_MatchMap(t *testing.T) {

	uasArgs := map[string]int{
		"AppleCoreMedia/1.0.0.16F203 (iPad; U; CPU OS 12_3_1 like Mac OS X; de_de)":   1,
		"AppleCoreMedia/1.0.0.13G36 (iPad; U; CPU OS 9_3_5 like Mac OS X; de_de)":     1,
		"AppleCoreMedia/1.0.0.13G36 (iPhone; U; CPU OS 9_3_5 like Mac OS X; en_gb)":   1,
		"AppleCoreMedia/1.0.0.16A404 (iPhone; U; CPU OS 12_0_1 like Mac OS X; de_de)": 1,
		"AppleCoreMedia/1.0.0.16B92 (iPhone; U; CPU OS 12_1 like Mac OS X; de_de)":    1,
		"AppleCoreMedia/1.0.0.16D57 (iPhone; U; CPU OS 12_1_4 like Mac OS X; de_de)":  1,
		"AppleCoreMedia/1.0.0.16E227 (iPhone; U; CPU OS 12_2 like Mac OS X; de_de)":   1,
		"AppleCoreMedia/1.0.0.16E227 (iPhone; U; CPU OS 12_2 like Mac OS X; en_us)":   1,
		"AppleCoreMedia/1.0.0.16F156 (iPhone; U; CPU OS 12_3 like Mac OS X; de_at)":   1,
		"AppleCoreMedia/1.0.0.16F203 (iPhone; U; CPU OS 12_3_1 like Mac OS X; de_at)": 1,
		"AppleCoreMedia/1.0.0.16F203 (iPhone; U; CPU OS 12_3_1 like Mac OS X; de_de)": 1,
		"AppleCoreMedia/1.0.0.16F203 (iPhone; U; CPU OS 12_3_1 like Mac OS X; en_us)": 1,
		"AppleCoreMedia/1.0.0.16F250 (iPhone; U; CPU OS 12_3_2 like Mac OS X; de_de)": 1,
	}
	uasWant := map[string]int{
		"iPad":   2,
		"iPhone": 11,
	}

	type args struct {
		uas map[string]int
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		// TODO: Add test cases.
		{name: "match-map-1", args: args{uas: uasArgs}, want: uasWant},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uam, err := minioanalytics.NewUserAgentMater()
			require.NoError(t, err)
			if got := uam.MatchMap(tt.args.uas); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserAgentMatcher.MatchMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
