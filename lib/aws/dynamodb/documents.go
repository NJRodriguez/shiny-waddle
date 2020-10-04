package dynamodb

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/pkg/errors"
)

//go:generate mockery --name DocumentsClient
type DocumentsClient interface {
	Get(key interface{}) (*dynamodb.GetItemOutput, error)
	Create(item interface{}) (*dynamodb.PutItemOutput, error)
	List(exclusiveStartKey map[string]*dynamodb.AttributeValue, limit int64) (*dynamodb.ScanOutput, error)
	ListAll() ([]map[string]*dynamodb.AttributeValue, error)
}

type documents struct {
	awsDynamodbClient dynamodbiface.DynamoDBAPI
	table             string
}

var newAwsSession = session.NewSession

// New creates a new Documents client for interacting with AWS DynamoDB.
func New(table string, awsRegion string) (*documents, error) {
	session, err := newAwsSession(&aws.Config{Region: aws.String(awsRegion)})
	if err != nil {
		return nil, errors.Wrap(err, "starting new aws sessions")
	}
	service := dynamodb.New(session)
	return &documents{service, table}, nil
}

func (instance *documents) Get(document interface{}) (*dynamodb.GetItemOutput, error) {
	item, err := dynamodbattribute.MarshalMap(document)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling interface to dynamodb readable")
	}
	args := &dynamodb.GetItemInput{
		TableName: aws.String(instance.table),
		Key:       item,
	}
	doc, err := instance.awsDynamodbClient.GetItem(args)
	if err != nil {
		log.Println("Failed to obtain item from table.")
		return nil, err
	}
	return doc, nil
}

func (instance *documents) Create(document interface{}) (*dynamodb.PutItemOutput, error) {
	item, err := dynamodbattribute.MarshalMap(document)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling interface to dynamodb readable")
	}
	condition := "attribute_not_exists(id)"
	args := &dynamodb.PutItemInput{
		TableName:           aws.String(instance.table),
		ConditionExpression: &condition,
		Item:                item,
	}
	result, err := instance.awsDynamodbClient.PutItem(args)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (instance *documents) List(exclusiveStartKey map[string]*dynamodb.AttributeValue, limit int64) (*dynamodb.ScanOutput, error) {
	result, err := instance.awsDynamodbClient.Scan(&dynamodb.ScanInput{
		TableName:         aws.String(instance.table),
		Limit:             aws.Int64(limit),
		ExclusiveStartKey: exclusiveStartKey,
	})
	if err != nil {
		log.Println("Failed to list items from table")
		return nil, err
	}
	return result, nil
}

func (instance *documents) ListAll() ([]map[string]*dynamodb.AttributeValue, error) {
	result := []map[string]*dynamodb.AttributeValue{}
	lastEvaluatedKey := map[string]*dynamodb.AttributeValue{}
	for ok := true; ok; ok = !(lastEvaluatedKey == nil) {
		if len(lastEvaluatedKey) == 0 {
			lastEvaluatedKey = nil
		}
		response, err := instance.List(lastEvaluatedKey, 10) // 10 is an arbitrary number. Considering the project scale, this can be modified in the future.
		if err != nil {
			return nil, errors.Wrap(err, "list items from dynamodb error")
		}
		result = append(result, response.Items...)
		lastEvaluatedKey = response.LastEvaluatedKey
	}
	return result, nil
}
