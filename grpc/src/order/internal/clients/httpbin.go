package clients

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type publicIP struct {
	Origin string
}

func GetPublicIP() string {

	log.Printf("[Order] Invoking httpbin public IP service...")

	httpsInsecure, err := strconv.ParseBool(os.Getenv("HTTPBIN_INSECURE"))

	if err != nil {
		httpsInsecure = false
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: httpsInsecure},
	}

	client := &http.Client{Transport: tr}

	httpbinProtocol := os.Getenv("HTTPBIN_PROTOCOL")
	httpbinHost := os.Getenv("HTTPBIN_HOST")
	httpbinPort := os.Getenv("HTTPBIN_PORT")

	if httpbinProtocol == "" {
		httpbinProtocol = "https"
	}

	if httpbinHost == "" {
		httpbinHost = "httpbin.org"
	}

	if httpbinPort == "" {
		httpbinPort = "443"
	}

	queryURL := httpbinProtocol + "://" + httpbinHost + ":" + httpbinPort + "/ip"

	log.Printf("[Order] httpbin query URL: %s", queryURL)

	response, err := client.Get(queryURL)

	var result string

	if err != nil {
		log.Fatalf("[Order] httpbin public IP service failed with error: %s", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		var response publicIP
		json.Unmarshal([]byte(data), &response)

		result = response.Origin

		log.Printf("[Order] httpbin returned public IP: %s", result)
	}

	return result
}
