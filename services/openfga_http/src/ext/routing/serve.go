package routing

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/openfga/openfga/cmd/run"
	"github.com/openfga/openfga/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Config struct {
	HttpClient *http.Client
}

// Start begins running the sidecar
func Start(port string, config *Config) {
	go startHTTPServer(port, config)
}

// Method that responds back with the cached values
func startHTTPServer(port string, config *Config) {
	/*
		r := chi.NewRouter()
		r.Get("/openfga/{key...}", handleValue(config))
		r.Post("/openfga/{key...}", handleValue(config))

		logrus.Infof("Starting server on %s", port)
		err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	*/
	/* works ok
	http.Handle("/", responseLogger(http.HandlerFunc(sampleHandler)))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	*/
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := runFGAServer(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	// Create a channel to receive termination signals
	sigs := make(chan os.Signal, 1)
	// Register the channel to receive signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-sigs
	fmt.Println("Received signal:", sig)
	//ensureServiceUp(t, cfg.GRPC.Addr, cfg.HTTP.Addr, nil, true)
}

func runFGAServer(ctx context.Context) error {
	//cfg := run.MustDefaultConfigWithRandomPorts()
	cfg, err := run.ReadConfig()
	if err != nil {
		panic(err)
	}
	//listen only localy inside the lambda execution context
	cfg.HTTP.Addr = fmt.Sprintf("127.0.0.1:%d", 4000)
	//allow GRPc only localy as well
	cfg.GRPC.Addr = fmt.Sprintf("127.0.0.1:%d", 8081)
	if err := cfg.Verify(); err != nil {
		return err
	}

	logger := logger.MustNewLogger(cfg.Log.Format, "DEBUG", "")
	serverCtx := &run.ServerContext{Logger: logger}
	return serverCtx.Run(ctx, cfg)
}

func sampleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	response := []byte("The received path:" + r.URL.Path)
	w.Write(response)
}

type loggingResponseWriter struct {
	status int
	body   string
	http.ResponseWriter
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) {
	w.body = string(body)
	return w.ResponseWriter.Write(body)
}

func responseLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingRW := &loggingResponseWriter{
			ResponseWriter: w,
		}
		h.ServeHTTP(loggingRW, r)
		log.Println("Status : ", loggingRW.status, "Response : ", loggingRW.body)
	})
}

func abcdHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("Response to /abcd path")
	fmt.Fprintf(w, "Response to /abcd path")
}

func handleValue(config *Config) http.HandlerFunc {
	logrus.Debug("Before Received a request on http server")
	return func(w http.ResponseWriter, r *http.Request) {

		logrus.Debug("Received a request on http server hhh")
		w.Write([]byte("Answer from http"))
	}
}
