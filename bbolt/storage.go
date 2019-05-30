package bbolt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	minioanalytics "github.com/codeandship/minio-analytics"
	bolt "go.etcd.io/bbolt"
)

var _ minioanalytics.Storage = &Storage{}

const (
	defaultBucketName = "records"
)

// Storage is a bbolt storage implementation
type Storage struct {
	path string
	db   *bolt.DB
}

// New returns a new bbolt storage
func New(path string) *Storage {
	return &Storage{
		path: path,
	}
}

// Open opens the Storage for r/w
func (s *Storage) Open() (err error) {
	s.db, err = bolt.Open(s.path, 0666, nil)
	if err != nil {
		return err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(defaultBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	return err
}

// Close closes the storage
func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// CreateRecord creates a new record in the storage
func (s *Storage) CreateRecord(r minioanalytics.AWSRecord) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucketName))
		data, err := json.Marshal(r)
		if err != nil {
			return err
		}
		if err = b.Put([]byte(fmt.Sprintf("record-%d", r.EventTime.UnixNano())), data); err != nil {
			return err
		}
		res := b.Get([]byte(fmt.Sprintf("object-%s", r.S3.Object.ETag)))
		var a minioanalytics.Analytics
		if res != nil {
			if err = json.Unmarshal(res, &a); err != nil {
				return err
			}
			if a.Filename != r.S3.Object.Key {
				return fmt.Errorf("analytics.filename (%q) != record.s3.obect.key (%q)", a.Filename, r.S3.Object.Key)
			}
			if r.EventName == minioanalytics.S3ObjectAccessedGet {
				a.GetRequestCount++
			} else if r.EventName == minioanalytics.S3ObjectAccessedHead {
				a.HeadRequestCount++
			}
			if v, ok := a.UserAgentCount[r.Source.UserAgent]; ok {
				a.UserAgentCount[r.Source.UserAgent] = v + 1
			} else {
				a.UserAgentCount[r.Source.UserAgent] = 1
			}

		} else {
			a.Filename = r.S3.Object.Key
			a.UserAgentCount = map[string]int{
				r.Source.UserAgent: 1,
			}

			if r.EventName == minioanalytics.S3ObjectAccessedGet {
				a.GetRequestCount = 1
			} else if r.EventName == minioanalytics.S3ObjectAccessedHead {
				a.HeadRequestCount = 1
			}
		}
		if data, err = json.Marshal(a); err != nil {
			return err
		}
		return b.Put([]byte(fmt.Sprintf("object-%s", r.S3.Object.ETag)), data)
	})
}

// ListRecords lists records from storage
func (s *Storage) ListRecords() ([]minioanalytics.AWSRecord, error) {
	return nil, errors.New("not implemented")
}

// ListAnalytics lists all Analytics for Objects
func (s *Storage) ListAnalytics() (res []minioanalytics.Analytics, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucketName))
		c := b.Cursor()
		prefix := []byte("object")
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var a minioanalytics.Analytics
			if err = json.Unmarshal(v, &a); err != nil {
				return err
			}
			res = append(res, a)
		}
		return nil
	})
	return
}
