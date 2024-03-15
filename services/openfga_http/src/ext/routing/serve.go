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
	//to speed up disable the playground
	cfg.Playground.Enabled = false
	//to speed up disable the metrics server
	cfg.Metrics.Enabled = false

	//get the postgres connection
	cfg.Datastore.Engine = os.Getenv("OPENFGA_DATASTORE_ENGINE")
	cfg.Datastore.URI = os.Getenv("OPENFGA_DATASTORE_URI")

	if err := cfg.Verify(); err != nil {
		return err
	}

	logger := logger.MustNewLogger(cfg.Log.Format, "DEBUG", "")
	serverCtx := &run.ServerContext{Logger: logger}
	return serverCtx.Run(ctx, cfg)
}
