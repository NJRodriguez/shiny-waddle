package models

import (
	dynamodbSdk "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

// Sucursal defines the properties belonging to the Sucursal object in DynamoDB.
type Sucursal struct {
	// Used for DynamoDB lookup.
	ID string `json:"id"`
	// Address of Sucursal.
	Address string `json:"address"`
	// Latitude of Sucursal in decimal degrees.
	Latitude float64 `json:"latitude"`
	// Longitude of Sucursal in decimal degrees.
	Longitude float64 `json:"longitude"`
}

type SucursalWithDistance struct {
	Sucursal *Sucursal
	Distance float64
}

type SucursalKey struct {
	ID string `json:"id"`
}

func ToSucursal(dynamodbItem map[string]*dynamodbSdk.AttributeValue) (*Sucursal, error) {
	sucursal := Sucursal{}
	err := dynamodbattribute.UnmarshalMap(dynamodbItem, &sucursal)
	if err != nil {
		return nil, errors.Wrap(err, "unable to convert to sucursal object")
	}
	return &sucursal, nil
}

func ToSucursalArray(dynamodbArray []map[string]*dynamodbSdk.AttributeValue) ([]*Sucursal, error) {
	result := []*Sucursal{}
	for _, item := range dynamodbArray {
		sucursal := Sucursal{}
		err := dynamodbattribute.UnmarshalMap(item, &sucursal)
		if err != nil {
			return nil, errors.Wrap(err, "unable to convert to sucursal object")
		}
		result = append(result, &sucursal)
	}
	return result, nil
}
