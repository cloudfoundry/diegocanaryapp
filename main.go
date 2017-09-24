package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	appIndex, err := strconv.Atoi(os.Getenv("CF_INSTANCE_INDEX"))
	if err != nil {
		panic(err.Error())
	}

	http.Handle("/", helloFromInstance(appIndex))

	datadogApiKey := os.Getenv("DATADOG_API_KEY")
	deploymentName := os.Getenv("DEPLOYMENT_NAME")
	cellIP := os.Getenv("CF_INSTANCE_IP")

	includeCellIPTag := (os.Getenv("INCLUDE_CELL_IP_TAG") == "true")

	go postHeartbeat(appIndex, datadogApiKey, deploymentName, cellIP, includeCellIPTag)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("listening and heartbeating...")
}

func helloFromInstance(index int) http.Handler {
	instanceText := fmt.Sprintf("Diego canary app: tweet tweet from instance %d", index)
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-type", "text/plain")
		fmt.Fprintln(res, instanceText)
		fmt.Fprintln(res, "For more information about this app, please consult http://code.cloudfoundry.org/diegocanaryapp.")
	})
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

func postHeartbeat(appIndex int, datadogApiKey, deploymentName, cellIP string, includeCellIPTag bool) {
	url := "https://app.datadoghq.com/api/v1/series?api_key=" + datadogApiKey

	client := http.DefaultClient
	tags := []string{
		"deployment:" + deploymentName,
		fmt.Sprintf("diego-canary-app:%d", appIndex),
	}
	if includeCellIPTag {
		tags = append(tags, "cell-ip:"+cellIP)
	}

	for {
		time.Sleep(5 * time.Second)

		fmt.Printf("instance '%d' in deployment '%s' emitting from host IP '%s'\n", appIndex, deploymentName, cellIP)

		series := DatadogSeries{
			Series: []DatadogMetric{{
				Metric: "diego.canary.app.instance",
				Tags:   tags,
				Points: []Point{{
					int(time.Now().Unix()),
					1,
				}},
			}},
		}
		jsonPayload, err := json.Marshal(series)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}

		req.Header.Set("Content-type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}
		resp.Body.Close()

		fmt.Println("Datadog status: " + resp.Status)
	}
}
