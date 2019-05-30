package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	minio "github.com/minio/minio-go/v6"
	nats "github.com/nats-io/nats.go"
	"github.com/peterbourgon/ff"

	minioanalytics "github.com/codeandship/minio-analytics"
	"github.com/codeandship/minio-analytics/bbolt"
	"github.com/codeandship/minio-analytics/http"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	fs := flag.NewFlagSet("my-program", flag.ExitOnError)
	var (
		natsSubject    string
		natsAddr       string
		minioAddr      string
		minioSecretKey string
		minioAccessKey string
		minioSSL       bool
		minioBucket    string
		minioExt       string
		storagePath    string
		addr           string
	)

	fs.StringVar(&natsSubject, "nats-subject", "bucketevents", "set minio nats subject")
	fs.StringVar(&natsAddr, "nats-addr", "nats:4222", "set nats server address")
	fs.StringVar(&minioAddr, "minio-addr", "minio:9000", "set minio server address")
	fs.StringVar(&minioSecretKey, "minio-sec-key", "", "set minio secret key")
	fs.StringVar(&minioAccessKey, "minio-acc-key", "", "set minio access key")
	fs.BoolVar(&minioSSL, "minio-ssl", false, "enable minio via ssl")
	fs.StringVar(&minioBucket, "minio-bucket", "", "set minio bucket name")
	fs.StringVar(&minioExt, "minio-ext", ".mp3", "set minio extensions that should be monitored, comma separated value")
	fs.StringVar(&storagePath, "store-path", "/data/analytics-db", "set storage path")
	fs.StringVar(&addr, "http-addr", ":80", "set http listen address")

	ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
		ff.WithEnvVarPrefix("MINIO_ANALYTICS"),
	)

	errs := []string{}
	if minioSecretKey == "" {
		errs = append(errs, "configuration: minio secrect can not be empty")
	}
	if minioAccessKey == "" {
		errs = append(errs, "configuration: minio access key can not be empty")
	}
	if minioBucket == "" {
		errs = append(errs, "configuration: minio bucket can not be empty")
	}
	if minioExt == "" {
		errs = append(errs, "configuration: minio file extensions can not be empty")
	}

	if len(errs) > 0 {
		log.Println(len(errs), "errors:")
		for i, s := range errs {
			fmt.Println("\t", i+1, s)
		}
		os.Exit(1)
	}

	s := bbolt.New(storagePath)
	if err := s.Open(); err != nil {
		log.Println("storage:", err.Error())
		os.Exit(1)
	}

	api := http.NewAPI(addr, s)

	go func() {
		err := api.Open()
		if err != nil {
			s.Close()
			log.Println(err.Error())
			os.Exit(1)
		}
	}()

	wg := &sync.WaitGroup{}
	// Connect to nats server
	nc, err := nats.Connect(natsAddr)
	if err != nil {
		log.Println("nats connection:", err.Error())
		os.Exit(1)
	}

	wg.Add(1)

	// Simple Async Subscriber
	nc.Subscribe(natsSubject, func(m *nats.Msg) {
		var e minioanalytics.MinioEvent
		if err := json.Unmarshal(m.Data, &e); err != nil {
			log.Println(err.Error())
		}
		for _, r := range e.Records {
			if err := s.CreateRecord(r); err != nil {
				log.Println(err.Error())
			}
		}
	})

	// Initialize minio client object.
	minioClient, err := minio.New(minioAddr, minioAccessKey, minioSecretKey, minioSSL)
	if err != nil {
		log.Println("minio connection:", err)
		os.Exit(1)
	}

	// arn:partition:service:region:account-id:resource
	queueArn := minio.NewArn("minio", "sqs", "", "1", "nats")

	queueConfig := minio.NewNotificationConfig(queueArn)
	queueConfig.AddEvents(minio.ObjectAccessedGet, minio.ObjectAccessedHead)

	exts := strings.Split(minioExt, ",")

	for _, ext := range exts {
		queueConfig.AddFilterSuffix(ext)
	}

	bucketNotification := minio.BucketNotification{}
	bucketNotification.AddQueue(queueConfig)

	err = minioClient.SetBucketNotification(minioBucket, bucketNotification)
	if err != nil {
		log.Println("unable to set the bucket notification: ", err)
		os.Exit(1)
	}

	log.Println("minio-analytics: started, waiting for events")
	wg.Wait()

}
