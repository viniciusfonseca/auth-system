package authsystemroutes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	authsystem "github.com/viniciusfonseca/auth-system/internal/shared"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func TestCreateSession(t *testing.T) {

	ctx := context.Background()

	authSystem := authsystem.GetAuthSystemTest(t, ctx)
	AddAuthGroup(authSystem)

	clientId := uuid.New().String()
	clientSecret := uuid.New().String()
	orgId := uuid.New().String()

	clientSecretHash, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := authsystem.User{
		ClientId:     clientId,
		OrgId:        orgId,
		ClientSecret: string(clientSecretHash),
	}

	_, err = authSystem.DdbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item:      user.AsDynamoDBItem(),
	})

	require.NoError(t, err)

	payload := CreateSessionPayload{
		ClientId:     clientId,
		OrgId:        orgId,
		ClientSecret: clientSecret,
	}
	body, err := json.Marshal(payload)

	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/auth/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)

	resp, err := authSystem.FiberApp.Test(req, -1)
	require.NoError(t, err)

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("%+v\n", string(body))
	}

	require.Equal(t, http.StatusCreated, resp.StatusCode)
}
