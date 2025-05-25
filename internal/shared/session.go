package authsystem

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Session struct {
	SessionId   string `json:"session_id"`
	ClientId    string `json:"client_id"`
	OrgId       string `json:"org_id"`
	AccessToken string `json:"access_token"`
}

func (s Session) AsDynamoDBItem() map[string]types.AttributeValue {

	sessionId, err := attributevalue.Marshal(s.SessionId)
	if err != nil {
		log.Fatal(err)
	}

	clientId, err := attributevalue.Marshal(s.ClientId)
	if err != nil {
		log.Fatal(err)
	}

	orgId, err := attributevalue.Marshal(s.OrgId)
	if err != nil {
		log.Fatal(err)
	}

	accessToken, err := attributevalue.Marshal(s.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	return map[string]types.AttributeValue{
		"session_id":   sessionId,
		"client_id":    clientId,
		"org_id":       orgId,
		"access_token": accessToken,
	}
}
