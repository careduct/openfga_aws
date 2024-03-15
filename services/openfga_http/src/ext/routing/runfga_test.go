package routing

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"testing"

	"github.com/openfga/openfga/cmd/run"
	"github.com/openfga/openfga/pkg/logger"
)

func TestMain(m *testing.M) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestHTTPServerEnabled(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := runServer(ctx); err != nil {
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

func runServer(ctx context.Context) error {
	//cfg := run.MustDefaultConfigWithRandomPorts()
	cfg, err := run.ReadConfig()
	if err != nil {
		panic(err)
	}
	//listen only localy inside the lambda execution context
	cfg.HTTP.Addr = fmt.Sprintf("127.0.0.1:%d", 8080)
	//allow GRPc only localy as well
	cfg.GRPC.Addr = fmt.Sprintf("127.0.0.1:%d", 8081)
	cfg.Datastore.Engine = "memory"
	/*
		cfg.Datastore.URI = "postgresql://postgres:aGdFrhf42df@localhost:5432/postgres"

		cmd := migrate.NewMigrateCommand()

		cmd.Flags().Set("datastore-engine", "postgres")
		cmd.Flags().Set("datastore-uri", "postgresql://postgres:aGdFrhf42df@localhost:5432/postgres")

		cmd.Execute()
	*/
	if err := cfg.Verify(); err != nil {
		return err
	}

	logger := logger.MustNewLogger(cfg.Log.Format, "DEBUG", "")
	serverCtx := &run.ServerContext{Logger: logger}
	return serverCtx.Run(ctx, cfg)
}
