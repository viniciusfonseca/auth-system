package authsystemroutes

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	authsystem "github.com/viniciusfonseca/auth-system/internal/shared"
)

func TestCreateUser(t *testing.T) {

	ctx := context.Background()

	authSystem := authsystem.GetAuthSystemTest(t, ctx)
	AddAuthGroup(authSystem)

	clientId := uuid.New().String()
	clientSecret := uuid.New().String()
	orgId := uuid.New().String()

	payload := CreateUserPayload{
		ClientId:     clientId,
		OrgId:        orgId,
		ClientSecret: clientSecret,
	}
	body, err := json.Marshal(payload)

	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/auth/users", bytes.NewReader(body))

	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := authSystem.FiberApp.Test(req, -1)

	require.NoError(t, err)

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		t.Errorf("Response body: %s", body)
	}

	require.Equal(t, http.StatusCreated, resp.StatusCode)

}

func TestCreateUserAlreadyExists(t *testing.T) {

	ctx := context.Background()

	authSystem := authsystem.GetAuthSystemTest(t, ctx)
	AddAuthGroup(authSystem)

	user := authsystem.User{
		ClientId:     uuid.New().String(),
		OrgId:        uuid.New().String(),
		ClientSecret: uuid.New().String(),
	}

	_, err := authSystem.DdbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item:      user.AsDynamoDBItem(),
	})

	require.NoError(t, err)

	payload := CreateUserPayload{
		ClientId:     user.ClientId,
		OrgId:        user.OrgId,
		ClientSecret: user.ClientSecret,
	}
	body, err := json.Marshal(payload)

	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/auth/users", bytes.NewReader(body))

	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := authSystem.FiberApp.Test(req, -1)

	require.NoError(t, err)

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status code %d, got %d", http.StatusConflict, resp.StatusCode)
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		t.Errorf("Response body: %s", body)
	}

	require.Equal(t, http.StatusConflict, resp.StatusCode)

}
