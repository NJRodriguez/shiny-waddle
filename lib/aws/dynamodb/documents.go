package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type documents struct {
	awsDynamodbClient dynamodbiface.DynamoDBAPI
	table             string
}

var newAwsSession = session.NewSession

// New creates a new Documents client for interacting with AWS DynamoDB.
func New(table string, awsRegion string) (*documents, error) {
	session, err := newAwsSession(&aws.Config{Region: aws.String(awsRegion)})
	if err != nil {
		return nil, err
	}
	service := dynamodb.New(session)
	return &documents{service, table}, nil
}

func (instance *documents) Get(key map[string]*dynamodb.AttributeValue) (*dynamodb.GetItemOutput, error) {
	args := &dynamodb.GetItemInput{
		TableName: aws.String(instance.table),
		Key:       key,
	}
	item, err := instance.awsDynamodbClient.GetItem(args)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (instance *documents) Create(item map[string]*dynamodb.AttributeValue) (*dynamodb.PutItemOutput, error) {
	args := &dynamodb.PutItemInput{
		TableName: aws.String(instance.table),
		Item:      item,
	}
	result, err := instance.awsDynamodbClient.PutItem(args)
	if err != nil {
		return nil, err
	}
	return result, err
}
