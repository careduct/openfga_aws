package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/openfga/openfga/cmd/migrate"
)

// Define your event structure according to the expected input.
// For simplicity, we're using a map here, but you might want to
// define a more specific struct based on your use case.
type MyEvent struct {
	Name string `json:"name"`
}

// Define your response structure.
type MyResponse struct {
	Message string `json:"message"`
}

// This is the handler function that AWS Lambda will invoke.
func HandleLambdaEvent(ctx context.Context, event MyEvent) (MyResponse, error) {
	// Log the received event (you might want to do more complex processing)
	mcmd := migrate.NewMigrateCommand()
	mcmd.Flags().Set("datastore-engine", os.Getenv("OPENFGA_DATASTORE_ENGINE"))
	mcmd.Flags().Set("datastore-uri", os.Getenv("OPENFGA_DATASTORE_URI"))

	mcmd.Execute()

	// Return a response
	return MyResponse{Message: fmt.Sprintf("Hello, %s!", event.Name)}, nil
}

func main() {
	// Tell AWS Lambda to start your handler function
	lambda.Start(HandleLambdaEvent)
}
