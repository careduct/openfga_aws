package main

import (
	"context"
	"fmt"
	"log"
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

func waitForOpenfga() {
	// Define the URL you want to make a request to
	url := "http://localhost:4000"

	// Create a context with a timeout of 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // It's important to call cancel to avoid leaking resources

	// Create an HTTP request with the context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("The OpenFGA didn't started within 10s on the cold start: %v\n", err)
		return
	}

	// Create an HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("The OpenFGA didn't started within 10s on the cold start: %v\n", err)
		return
	}
	defer resp.Body.Close() // Don't forget to close the response body
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		origin, err := url.Parse("http://localhost:4000")
		if err != nil {
			panic(err)
		}

		//create teh reverse çroxy
		proxy := httputil.NewSingleHostReverseProxy(origin)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/openfga")

		if err != nil {
			log.Fatal(err)
		}
		//wait until OpenFGA is ready on a cold start
		waitForOpenfga()
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
