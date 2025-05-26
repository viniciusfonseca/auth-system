package authsystem

import (
	"context"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	tcdynamodb "github.com/testcontainers/testcontainers-go/modules/dynamodb"
)

type AuthSystem struct {
	FiberApp  *fiber.App
	DdbClient *dynamodb.Client
}

func NewAuthSystem(fiberApp *fiber.App) *AuthSystem {

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ddbClient := dynamodb.NewFromConfig(cfg)

	return &AuthSystem{
		FiberApp:  fiberApp,
		DdbClient: ddbClient,
	}
}

func (a *AuthSystem) GetFiberHandler(fn func(c *fiber.Ctx, a *AuthSystem) error) fiber.Handler {
	return func(c *fiber.Ctx) error { return fn(c, a) }
}

func GetAuthSystemTest(t *testing.T, ctx context.Context) *AuthSystem {

	ctr, err := tcdynamodb.Run(ctx, "amazon/dynamodb-local:2.2.1", tcdynamodb.WithSharedDB())

	require.NoError(t, err)

	hostPort, err := ctr.ConnectionString(ctx)

	require.NoError(t, err)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "test")),
		config.WithBaseEndpoint(hostPort),
	)

	require.NoError(t, err)

	ddbClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(hostPort)
		o.EndpointResolverV2 = &DynamoDBLocalResolver{HostAndPort: hostPort}
	})

	listTablesOutput, err := ddbClient.ListTables(ctx, &dynamodb.ListTablesInput{})

	require.NoError(t, err)

	if len(listTablesOutput.TableNames) == 0 {
		_, err = ddbClient.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String("users"),
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("client_id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("org_id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("client_id"),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String("org_id"),
					KeyType:       types.KeyTypeRange,
				},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
		})

		require.NoError(t, err)

		_, err = ddbClient.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String("sessions"),
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("session_id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("client_id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("session_id"),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String("client_id"),
					KeyType:       types.KeyTypeRange,
				},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
		})

		require.NoError(t, err)
	}

	return &AuthSystem{
		FiberApp:  fiber.New(),
		DdbClient: ddbClient,
	}

}
