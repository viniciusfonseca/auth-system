package authsystemroutes

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofiber/fiber/v2"
	authsystem "github.com/viniciusfonseca/auth-system/internal/shared"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserPayload struct {
	ClientId     string `json:"client_id"`
	OrgId        string `json:"org_id"`
	ClientSecret string `json:"client_secret"`
}

func CreateUser(c *fiber.Ctx, a *authsystem.AuthSystem) error {

	payload := CreateUserPayload{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(payload.ClientSecret), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := authsystem.User{
		ClientId:     payload.ClientId,
		OrgId:        payload.OrgId,
		ClientSecret: string(hash),
	}

	if getItemOutput, err := a.DdbClient.GetItem(c.Context(), &dynamodb.GetItemInput{
		TableName: aws.String("users"),
		Key:       user.GetKey(),
	}); err == nil && getItemOutput.Item != nil {
		return fiber.ErrConflict
	} else if err != nil {
		return err
	}

	_, err = a.DdbClient.PutItem(c.Context(), &dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item:      user.AsDynamoDBItem(),
	})

	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}
