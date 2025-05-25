package authsystem

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	ClientId     string
	OrgId        string
	ClientSecret string
}

func (u User) GetKey() map[string]types.AttributeValue {

	clientId, err := attributevalue.Marshal(u.ClientId)
	if err != nil {
		log.Fatal(err)
	}

	orgId, err := attributevalue.Marshal(u.OrgId)
	if err != nil {
		log.Fatal(err)
	}

	return map[string]types.AttributeValue{
		"client_id": clientId,
		"org_id":    orgId,
	}

}

func (u User) AsDynamoDBItem() map[string]types.AttributeValue {

	clientId, err := attributevalue.Marshal(u.ClientId)
	if err != nil {
		log.Fatal(err)
	}

	orgId, err := attributevalue.Marshal(u.OrgId)
	if err != nil {
		log.Fatal(err)
	}

	clientSecret, err := attributevalue.Marshal(u.ClientSecret)
	if err != nil {
		log.Fatal(err)
	}

	return map[string]types.AttributeValue{
		"client_id":     clientId,
		"org_id":        orgId,
		"client_secret": clientSecret,
	}
}
