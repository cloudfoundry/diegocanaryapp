package main

import (
	"bytes"
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
	go postHeartbeat(appIndex, datadogApiKey, deploymentName)

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

func postHeartbeat(appIndex int, datadogApiKey string, deploymentName string) {
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
			`"tags":["deployment:%s", "diego-canary-app:%d"]`+
			`}]`+
			`}`,
			time.Now().Unix(),
			deploymentName,
			appIndex,
		))
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
		resp.Body.Close()

		println("datadog: " + resp.Status)
	}
}
