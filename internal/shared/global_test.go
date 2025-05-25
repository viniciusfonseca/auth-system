package authsystem

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	tcdynamodb "github.com/testcontainers/testcontainers-go/modules/dynamodb"
)

func TestMain(m *testing.M) {

	flag.Parse()

	ctx := context.Background()

	ctr, err := tcdynamodb.Run(ctx, "amazon/dynamodb-local:2.2.1", tcdynamodb.WithSharedDB())
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()

	if err := ctr.Terminate(ctx); err != nil {
		log.Fatal(err)
	}

	os.Exit(exitCode)
}
