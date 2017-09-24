package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const DatadogBaseURL = "https://app.datadoghq.com/api/v1/series?api_key="

type Heartbeat struct {
	datadogURL     string
	appIndex       int
	deploymentName string
	cellIP         string
	tags           []string
	client         *http.Client
}

func NewHeartbeat(appIndex int, datadogAPIKey, deploymentName, cellIP string, includeCellIPTag bool) *Heartbeat {
	tags := []string{
		"deployment:" + deploymentName,
		fmt.Sprintf("diego-canary-app:%d", appIndex),
	}

	if includeCellIPTag {
		tags = append(tags, "cell-ip:"+cellIP)
	}

	return &Heartbeat{
		datadogURL:     DatadogBaseURL + datadogAPIKey,
		appIndex:       appIndex,
		deploymentName: deploymentName,
		cellIP:         cellIP,
		tags:           tags,
		client:         http.DefaultClient,
	}
}

type DatadogSeries struct {
	Series []DatadogMetric `json:"series"`
}

type DatadogMetric struct {
	Metric string   `json:"metric"`
	Points []Point  `json:"points"`
	Tags   []string `json:"tags"`
}

type Point [2]int

func (h *Heartbeat) Post() {
	fmt.Printf("instance '%d' in deployment '%s' emitting from host IP '%s'\n", h.appIndex, h.deploymentName, h.cellIP)

	series := DatadogSeries{
		Series: []DatadogMetric{{
			Metric: "diego.canary.app.instance",
			Tags:   h.tags,
			Points: []Point{{
				int(time.Now().Unix()),
				1,
			}},
		}},
	}
	jsonPayload, err := json.Marshal(series)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	req, err := http.NewRequest("POST", h.datadogURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	req.Header.Set("Content-type", "application/json")
	resp, err := h.client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
	resp.Body.Close()

	fmt.Println("Datadog status: " + resp.Status)
}
