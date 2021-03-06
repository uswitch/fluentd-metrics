package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	buflen    = 10240
	namespace = "fluentd."
)

var (
	statsdAddr  = kingpin.Flag("statsd", "Host:Port of Datadog Statsd agent").Required().String()
	clusterName = kingpin.Flag("cluster", "Name of kubernetes cluster").Required().String()
	fluentURL   = kingpin.Flag("fluent", "Fluentd HTTP API endpoint").Default("http://127.0.0.1:24220").URL()
	interval    = kingpin.Flag("interval", "Gap between metric probes").Default("10s").Duration()
)

type fluentStats struct {
	Plugins []struct {
		PluginID             string  `json:"plugin_id"`
		Type                 string  `json:"type"`
		OutputPlugin         bool    `json:"output_plugin"`
		BufferQueueLen       float64 `json:"buffer_queue_length"`
		BufferTotalQueueSize float64 `json:"buffer_total_queued_size"`
		RetryCount           float64 `json:"retry_count"`
	} `json:"plugins"`
}

func getMetrics(URL string) (*fluentStats, error) {
	var fs fluentStats
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(body, &fs)
	if err != nil {
		return nil, err
	}
	return &fs, nil
}

type config struct {
	fluentURL   *url.URL
	hostname    string
	clusterName string
}

func publishMetrics(c *config, statsd *statsd.Client) error {
	fs, err := getMetrics(fmt.Sprintf("%s/api/plugins.json", c.fluentURL))
	if err != nil {
		return err
	}

	for _, i := range fs.Plugins {
		tags := []string{
			fmt.Sprintf("nodename:%s", c.hostname),
			fmt.Sprintf("kube_cluster:%s", c.clusterName),
			fmt.Sprintf("plugin_id:%s", i.PluginID),
			fmt.Sprintf("plugin_type:%s", i.Type),
		}
		if i.OutputPlugin && i.Type != "null" {
			statsd.Gauge("buffer_queue_len", i.BufferQueueLen, tags, 1)
			statsd.Gauge("buffer_total_queued_size", i.BufferTotalQueueSize, tags, 1)
			statsd.Gauge("retry_count", i.RetryCount, tags, 1)
		}
	}

	return nil
}

func main() {
	kingpin.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error getting hostname:%s", err)
	}

	// Statds Client
	log.Infof("Starting a buffered datadog statsd client at: %s", *statsdAddr)
	statsdClient, err := statsd.NewBuffered(*statsdAddr, buflen)
	if err != nil {
		log.Fatalf("Error starting statsd client: %s", err)
	}
	statsdClient.Namespace = namespace
	defer statsdClient.Close()

	cfg := &config{
		fluentURL:   *fluentURL,
		hostname:    hostname,
		clusterName: *clusterName,
	}

	// Ticker & Main Loop
	for {
		err = publishMetrics(cfg, statsdClient)
		if err != nil {
			log.Errorf("error forwarding metrics to statsd: %s", err)
		}
		time.Sleep(*interval)
	}
}
