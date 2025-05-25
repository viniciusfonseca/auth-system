package authsystemroutes

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	authsystem "github.com/viniciusfonseca/auth-system/internal/shared"
	"golang.org/x/crypto/bcrypt"
)

type CreateSessionPayload struct {
	ClientId     string `json:"client_id"`
	OrgId        string `json:"org_id"`
	ClientSecret string `json:"client_secret"`
}

func (p CreateSessionPayload) GetKey() map[string]types.AttributeValue {

	clientId, err := attributevalue.Marshal(p.ClientId)
	if err != nil {
		log.Fatal(err)
	}

	orgId, err := attributevalue.Marshal(p.OrgId)
	if err != nil {
		log.Fatal(err)
	}

	return map[string]types.AttributeValue{
		"client_id": clientId,
		"org_id":    orgId,
	}
}

func CreateSession(c *fiber.Ctx, a *authsystem.AuthSystem) error {

	payload := CreateSessionPayload{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	user, err := a.DdbClient.GetItem(c.Context(), &dynamodb.GetItemInput{
		TableName: aws.String("users"),
		Key:       payload.GetKey(),
	})

	if err != nil {
		return err
	}

	passwordHash := user.Item["client_secret"].(*types.AttributeValueMemberS).Value

	if err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(payload.ClientSecret)); err != nil {
		return fiber.ErrUnauthorized
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)

	session := authsystem.Session{
		SessionId:   uuid.New().String(),
		ClientId:    payload.ClientId,
		OrgId:       payload.OrgId,
		AccessToken: accessToken.Raw,
	}

	_, err = a.DdbClient.PutItem(c.Context(), &dynamodb.PutItemInput{
		TableName: aws.String("sessions"),
		Item:      session.AsDynamoDBItem(),
	})

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(session)
}
