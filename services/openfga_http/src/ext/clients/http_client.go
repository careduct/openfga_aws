package clients

import (
	"fmt"
	"net/http"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

func NewHttpClient() *http.Client {
	client := httptrace.WrapClient(&http.Client{}, httptrace.RTWithResourceNamer(func(h *http.Request) string {
		return fmt.Sprintf("%s %s://%s%s", h.Method, h.URL.Scheme, h.URL.Host, h.URL.Path)
	}))

	return client
}
