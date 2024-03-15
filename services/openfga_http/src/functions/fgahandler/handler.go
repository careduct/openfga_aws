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

func waitForHTTPServer(url string, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second) // Retry every second
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Timeout reached without receiving a 200 OK
			return false
		case <-ticker.C:
			// Attempt to make a request
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue // Try again on the next tick
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error making HTTP request:", err)
				continue // Try again on the next tick
			}
			resp.Body.Close() // Close the response body

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
				// The server responded with 200 OK or 404 NOT FOUND
				return true
			}
		}
	}
}

func waitForOpenfga() int {
	// Define the URL you want to make a request to
	url := "http://localhost:4000"

	// Create a context with a timeout of 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // It's important to call cancel to avoid leaking resources

	// Create an HTTP request with the context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("The OpenFGA didn't started within 10s on the cold start: %v\n", err)
		return 504
	}

	// Create an HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("The OpenFGA didn't started within 10s on the cold start: %v\n", err)
		return 504
	}
	resp.Body.Close() // Don't forget to close the response body

	//fmt.Printf("Response status code: %d\n", resp.StatusCode)
	return resp.StatusCode
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
