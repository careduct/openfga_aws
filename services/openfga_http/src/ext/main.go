package main

import (
	"cache-layer/clients"
	"cache-layer/extension"
	"cache-layer/routing"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	extensionName   = filepath.Base(os.Args[0]) // extension name has to match the filename
	extensionClient = extension.NewClient(os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	printPrefix     = fmt.Sprintf("[%s]", extensionName)
	config          *routing.Config
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	config = &routing.Config{
		HttpClient: clients.NewHttpClient(),
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sigs
		cancel()
		println(printPrefix, "Received", s)
		println(printPrefix, "Exiting")
	}()

	res, err := extensionClient.Register(ctx, extensionName)
	if err != nil {
		panic(err)
	}

	logrus.WithFields(logrus.Fields{
		"response": res,
	}).Debug("Client registered")

	routing.Start("4000", config)

	// Will block until shutdown event is received or cancelled via the context.
	processEvents(ctx)
}

func processEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, err := extensionClient.NextEvent(ctx)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err": err,
				}).Error("Error occurred.  Exiting the extension")
				return
			}
		}
	}
}
