package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const defaultEmissionInterval = 30 * time.Second

func main() {
	appIndex, err := strconv.Atoi(os.Getenv("CF_INSTANCE_INDEX"))
	if err != nil {
		panic(err.Error())
	}

	datadogAPIKey := os.Getenv("DATADOG_API_KEY")
	deploymentName := os.Getenv("DEPLOYMENT_NAME")
	cellIP := os.Getenv("CF_INSTANCE_IP")
	includeCellIPTag := (os.Getenv("INCLUDE_CELL_IP_TAG") == "true")

	heartbeat := NewHeartbeat(appIndex, datadogAPIKey, deploymentName, cellIP, includeCellIPTag)

	emissionInterval := constructEmissionInterval(os.Getenv("EMISSION_INTERVAL"))
	go postHeartbeat(heartbeat, emissionInterval)

	http.Handle("/", helloFromInstance(appIndex))

	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err.Error())
	}
}

func helloFromInstance(index int) http.Handler {
	instanceText := fmt.Sprintf("Diego canary app: tweet tweet from instance %d", index)
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-type", "text/plain")
		fmt.Fprintln(res, instanceText)
		fmt.Fprintln(res, "For more information about this app, please consult http://code.cloudfoundry.org/diegocanaryapp.")
	})
}

func postHeartbeat(heartbeat *Heartbeat, emissionInterval time.Duration) {
	for {
		heartbeat.Post()
		time.Sleep(emissionInterval)
	}
}

func constructEmissionInterval(intervalText string) time.Duration {
	if intervalText == "" {
		fmt.Printf("Using default emission interval of '%s'\n", defaultEmissionInterval.String())
		return defaultEmissionInterval
	}

	emissionInterval, err := time.ParseDuration(intervalText)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing EMISSION_INTERVAL value '%s': %s\n", intervalText, err.Error())
		fmt.Fprintf(os.Stderr, "Using default emission interval of '%s'\n", defaultEmissionInterval.String())
		return defaultEmissionInterval
	}

	return emissionInterval
}
