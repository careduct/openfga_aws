package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"

	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

type Model struct {
	Id       string `json:"id"`
	FieldOne string `json:"fieldOne"`
	FieldTwo string `json:"fieldTwo"`
}

func waitForOpenfga(targetUrl string) {

	client := http.Client{
		Timeout: 10 * time.Second, // Adjust the timeout as needed
	}

	maxRetries := 10
	baseRetryInterval := 500 * time.Millisecond // Base time to wait before retries

	for i := 0; i < maxRetries; i++ {
		resp, err := client.Get(targetUrl)
		if err != nil {
			//fmt.Printf("Attempt #%d: error making GET request: %v\n", i+1, err)
			sleepDuration := time.Duration(math.Pow(2, float64(i))) * baseRetryInterval
			//fmt.Printf("Waiting %v before retrying...\n", sleepDuration)
			time.Sleep(sleepDuration)
			continue
		}

		defer resp.Body.Close()
		//fmt.Printf("Attempt #%d: Server responded with status code: %d\n", i+1, resp.StatusCode)
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
			//fmt.Println("Server is up and responded successfully")
			break
		} else {
			sleepDuration := time.Duration(math.Pow(2, float64(i))) * baseRetryInterval
			//fmt.Printf("Server is not ready, waiting %v before retrying...\n", sleepDuration)
			time.Sleep(sleepDuration)
		}
	}
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		origin, err := url.Parse("http://localhost:4000")
		if err != nil {
			panic(err)
		}
		//wait until the cold start ends if we are in one
		/*if waitForHTTPServer("http://localhost:4000", 10*time.Second) {
			fmt.Println("The server is up and running, and responded with 200/404 OK!")
		} else {
			fmt.Println("The server did not respond with 200 OK within the specified timeout.")
			panic("Error instantiating the backend openfga on the lambda extension")
		}
		*/
		waitForOpenfga(origin.String())
		//create teh reverse çroxy
		proxy := httputil.NewSingleHostReverseProxy(origin)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/openfga")

		if err != nil {
			log.Fatal(err)
		}
		//wait until OpenFGA is ready on a cold start

		proxy.ServeHTTP(w, r)
		fmt.Printf("Go version: %s\n", runtime.Version())
	})

	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)

}

func init() {
	//client = &http.Client{}
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
