package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", hello)
	fmt.Println("listening and heartbeating ...")

	appIndex, err := extractAppIndex(os.Getenv("VCAP_APPLICATION"))
	if err != nil {
		panic(err.Error())
	}

	datadogApiKey := os.Getenv("DATADOG_API_KEY")
	go postHeartbeat(appIndex, datadogApiKey)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err.Error())
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "go, world")
}

func postHeartbeat(appIndex int, datadogApiKey string) {
	url := fmt.Sprintf(
		"https://app.datadoghq.com/api/v1/series?api_key=%s",
		datadogApiKey,
	)

	client := http.DefaultClient

	for {
		time.Sleep(5 * time.Second)

		jsonString := []byte(fmt.Sprintf(`{"series":`+
			`[{`+
			`"metric":"diego.canary.app.instance",`+
			`"points":[[%d, 1]],`+
			`"tags":["deployment:ketchup", "diego-canary-app-%d"]`+
			`}]`+
			`}`,
			time.Now().Unix(),
			appIndex,
		))
		println(jsonString)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
		if err != nil {
			println(err.Error())
		}

		req.Header.Set("Content-type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			println(err.Error())
			continue
		}
		println(resp.Status)
	}
}

func extractAppIndex(vcapApplicationJson string) (int, error) {
	type vcapApplication struct {
		InstanceIndex int `json:"instance_index"`
	}

	var v vcapApplication
	err := json.Unmarshal([]byte(vcapApplicationJson), &v)
	if err != nil {
		return 0, err
	}

	return v.InstanceIndex, nil
}
