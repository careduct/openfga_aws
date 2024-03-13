package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"

	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

type Model struct {
	Id       string `json:"id"`
	FieldOne string `json:"fieldOne"`
	FieldTwo string `json:"fieldTwo"`
}

func getModel(id string) (string, error) {
	request, _ := http.NewRequest("GET", "http://localhost:4000/1", nil)
	c := &http.Client{}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	logrus.WithFields(logrus.Fields{"ret_id": id, "url": fmt.Sprintf("http://localhost:4000/%s", id)}).Debug("Sending the reqquest")

	response, error := c.Do(request)
	if error != nil {
		logrus.Error("We got an error while invoking the http server")
		return "", error
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		logrus.WithFields(logrus.Fields{"ret_code": response.StatusCode}).Debug("Item not found by key")
		return "", nil
	}

	resBody, _ := ioutil.ReadAll(response.Body)
	log.Printf("The http server response: %s", resBody)
	return string(resBody), nil
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//log.Printf("Received request: %s %s %s %s", event.HTTPMethod, event.Path, event.QueryStringParameters, event.PathParameters)
	//log.Printf("Received request: %s", string(json.Marshal(event)))
	requestBody, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	// Log the JSON representation of the entire request
	log.Printf("Full Request:\n%s", string(requestBody))
	m, err := getModel(event.PathParameters["id"])
	status := 200
	var body string

	if err != nil {
		b, _ := json.Marshal(err)
		body = string(b)
		status = 404
	} else {
		b, _ := json.Marshal(m)
		body = string(b)
	}

	e := events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Method":      "OPTIONS,POST,GET,PUT,DELETE",
		},
	}

	return e, nil
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//logrus.Info("The Handler broking")
		origin, err := url.Parse("http://localhost:4000")
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(origin)
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/openfga")
		/*if pathsufix != "" {
			// Add a backslash if pathprefix is not empty
			pathsufix = "/" + pathsufix
		}*/
		//r.URL.Path = "stores" + pathsufix
		//logrus.Info("Serving the proxy:" + r.URL.Path)
		//reqDump, err := httputil.DumpRequestOut(r, true)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("REQUEST:\n%s", string(reqDump))
		proxy.ServeHTTP(w, r)
		//logrus.Info("After Serving the proxy, for")
		fmt.Printf("Go version: %s\n", runtime.Version())
	})

	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
	//lambda.Start(handlernew)
}

func init() {
	//client = &http.Client{}
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
