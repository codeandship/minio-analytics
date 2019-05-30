package bbolt_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	store "github.com/iwittkau/minio-analytics/bbolt"

	minioanalytics "github.com/iwittkau/minio-analytics"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *store.Storage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := store.New(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Open(t *testing.T) {

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &store.Storage{}
			if err := s.Open(); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &store.Storage{}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_CreateRecord(t *testing.T) {

	type args struct {
		r minioanalytics.AWSRecord
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create-1", args: args{
				r: minioanalytics.AWSRecord{
					EventName: minioanalytics.S3ObjectAccessedGet,
					EventTime: time.Now(),
					Source: minioanalytics.Source{
						UserAgent: "agent-1"},
					S3: minioanalytics.S3{
						Object: minioanalytics.S3Object{ETag: "1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "create-2", args: args{
				r: minioanalytics.AWSRecord{
					EventName: minioanalytics.S3ObjectAccessedHead,
					EventTime: time.Now(),
					Source: minioanalytics.Source{
						UserAgent: "agent-2"},
					S3: minioanalytics.S3{
						Object: minioanalytics.S3Object{
							ETag: "1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "create-3", args: args{
				r: minioanalytics.AWSRecord{
					EventName: minioanalytics.S3ObjectAccessedGet,
					EventTime: time.Now(),
					Source: minioanalytics.Source{
						UserAgent: "agent-2"},
					S3: minioanalytics.S3{
						Object: minioanalytics.S3Object{
							ETag: "1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "create-4", args: args{
				r: minioanalytics.AWSRecord{
					EventName: minioanalytics.S3ObjectAccessedHead,
					EventTime: time.Now(),
					Source: minioanalytics.Source{
						UserAgent: "agent-2"},
					S3: minioanalytics.S3{
						Object: minioanalytics.S3Object{
							ETag: "1"},
					},
				},
			},
			wantErr: false,
		},
	}
	s := store.New("test-data-create")
	err := s.Open()
	require.NoError(t, err)
	require.NoError(t, os.RemoveAll("test-data-create"))
	defer func() {
		require.NoError(t, s.Close())
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.CreateRecord(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Storage.CreateRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_ListRecords(t *testing.T) {

	tests := []struct {
		name    string
		want    []minioanalytics.AWSRecord
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "list-1",
			want:    nil,
			wantErr: true,
		},
	}
	s := store.New("test-data-listr")
	err := s.Open()
	require.NoError(t, err)
	require.NoError(t, os.RemoveAll("test-data-listr"))
	defer func() {
		require.NoError(t, s.Close())
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ListRecords()
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ListRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.ListRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_ListAnalytics(t *testing.T) {
	recs := []minioanalytics.AWSRecord{
		{
			EventName: minioanalytics.S3ObjectAccessedGet,
			EventTime: time.Now(),
			Source: minioanalytics.Source{
				UserAgent: "agent-1"},
			S3: minioanalytics.S3{
				Object: minioanalytics.S3Object{
					ETag: "1",
					Key:  "file-1",
				},
			},
		},
		{
			EventName: minioanalytics.S3ObjectAccessedGet,
			EventTime: time.Now(),
			Source: minioanalytics.Source{
				UserAgent: "agent-2"},
			S3: minioanalytics.S3{
				Object: minioanalytics.S3Object{
					ETag: "1",
					Key:  "file-1",
				},
			},
		},
		{
			EventName: minioanalytics.S3ObjectAccessedHead,
			EventTime: time.Now(),
			Source: minioanalytics.Source{
				UserAgent: "agent-1"},
			S3: minioanalytics.S3{
				Object: minioanalytics.S3Object{
					ETag: "1",
					Key:  "file-1",
				},
			},
		},
		{
			EventName: minioanalytics.S3ObjectAccessedHead,
			EventTime: time.Now(),
			Source: minioanalytics.Source{
				UserAgent: "agent-1"},
			S3: minioanalytics.S3{
				Object: minioanalytics.S3Object{
					ETag: "2",
					Key:  "file-2",
				},
			},
		},
	}

	tests := []struct {
		name    string
		wantRes []minioanalytics.Analytics
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "list-1",
			wantRes: []minioanalytics.Analytics{
				{
					Filename:         "file-1",
					GetRequestCount:  2,
					HeadRequestCount: 1,
					UserAgentCount: map[string]int{
						"agent-1": 2,
						"agent-2": 1,
					},
				},
				{
					Filename:         "file-2",
					GetRequestCount:  0,
					HeadRequestCount: 1,
					UserAgentCount: map[string]int{
						"agent-1": 1,
					},
				},
			},
			wantErr: false,
		},
	}
	s := store.New("test-data-lista")
	err := s.Open()
	require.NoError(t, err)
	require.NoError(t, os.RemoveAll("test-data-lista"))
	defer func() {
		require.NoError(t, s.Close())
	}()

	for _, r := range recs {
		require.NoError(t, s.CreateRecord(r))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := s.ListAnalytics()
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ListAnalytics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Storage.ListAnalytics() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
