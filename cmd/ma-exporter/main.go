package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	minioanalytics "github.com/codeandship/minio-analytics"
	"github.com/peterbourgon/ff"
)

var (
	getRequestsGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "minio",
		Subsystem: "analytics",
		Name:      "get_requests_count",
		Help:      "Counts number of http get requests on a given object",
	}, []string{"file"})
	headRequestsGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "minio",
		Subsystem: "analytics",
		Name:      "head_requests_count",
		Help:      "Counts number of http head requests on a given object",
	}, []string{"file"})
	userAgentsGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "minio",
		Subsystem: "analytics",
		Name:      "user_agents_count",
		Help:      "Counts number of reported user agents",
	}, []string{"file", "user_agent"})
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	fs := flag.NewFlagSet("minio-analytics-exporter", flag.ExitOnError)
	var (
		endpoint     string
		authUser     string
		authPassword string
		interval     string
	)
	fs.StringVar(&endpoint, "endpoint", "http://localhost:8080/analytics", "set endpoint to export from")
	fs.StringVar(&authUser, "auth-user", "", "set auth user")
	fs.StringVar(&authPassword, "auth-passwd", "", "set auth password")
	fs.StringVar(&interval, "interval", "10s", "interval for endpoint polling")

	ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
		ff.WithEnvVarPrefix("MINIO_ANALYTICS_EXPORTER"),
	)

	errs := []string{}
	d, err := time.ParseDuration(interval)
	if err != nil {
		errs = append(errs, err.Error())
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		log.Println(len(errs), "errors:")
		for i, s := range errs {
			fmt.Println("\t", i+1, s)
		}
		os.Exit(1)
	}

	go func() {
		prometheus.MustRegister(getRequestsGauge)
		prometheus.MustRegister(headRequestsGauge)
		prometheus.MustRegister(userAgentsGauge)
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9601", nil); err != nil {
			log.Fatal(err.Error())
		}
	}()

	c := http.Client{
		Timeout: d,
	}

	log.Println("checking every", d.String())
	var analytics []minioanalytics.Analytics
	for {
		time.Sleep(d)
		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if authUser != "" && authPassword != "" {
			req.SetBasicAuth(authUser, authPassword)
		}
		res, err := c.Do(req)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if err = json.Unmarshal(data, &analytics); err != nil {
			log.Println(err.Error())
			continue
		}
		log.Printf("%+v", analytics)
		setPromMetrics(analytics)
	}

}

func setPromMetrics(al []minioanalytics.Analytics) {

	aggr := minioanalytics.MapAnayltics(al)

	for _, a := range aggr {
		getRequestsGauge.With(prometheus.Labels{"file": a.Filename}).Set(float64(a.GetRequestCount))
		headRequestsGauge.With(prometheus.Labels{"file": a.Filename}).Set(float64(a.HeadRequestCount))
		for ua, count := range a.UserAgentCount {
			userAgentsGauge.With(prometheus.Labels{"file": a.Filename, "user_agent": ua}).Set(float64(count))
		}
	}
}
