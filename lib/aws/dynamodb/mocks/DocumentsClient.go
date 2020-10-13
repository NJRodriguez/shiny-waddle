// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	mock "github.com/stretchr/testify/mock"
)

// DocumentsClient is an autogenerated mock type for the DocumentsClient type
type DocumentsClient struct {
	mock.Mock
}

// Create provides a mock function with given fields: item
func (_m *DocumentsClient) Create(item interface{}) (*dynamodb.PutItemOutput, error) {
	ret := _m.Called(item)

	var r0 *dynamodb.PutItemOutput
	if rf, ok := ret.Get(0).(func(interface{}) *dynamodb.PutItemOutput); ok {
		r0 = rf(item)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.PutItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(item)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: key
func (_m *DocumentsClient) Get(key interface{}) (*dynamodb.GetItemOutput, error) {
	ret := _m.Called(key)

	var r0 *dynamodb.GetItemOutput
	if rf, ok := ret.Get(0).(func(interface{}) *dynamodb.GetItemOutput); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.GetItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: exclusiveStartKey, limit
func (_m *DocumentsClient) List(exclusiveStartKey map[string]*dynamodb.AttributeValue, limit int64) (*dynamodb.ScanOutput, error) {
	ret := _m.Called(exclusiveStartKey, limit)

	var r0 *dynamodb.ScanOutput
	if rf, ok := ret.Get(0).(func(map[string]*dynamodb.AttributeValue, int64) *dynamodb.ScanOutput); ok {
		r0 = rf(exclusiveStartKey, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.ScanOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(map[string]*dynamodb.AttributeValue, int64) error); ok {
		r1 = rf(exclusiveStartKey, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAll provides a mock function with given fields:
func (_m *DocumentsClient) ListAll() ([]map[string]*dynamodb.AttributeValue, error) {
	ret := _m.Called()

	var r0 []map[string]*dynamodb.AttributeValue
	if rf, ok := ret.Get(0).(func() []map[string]*dynamodb.AttributeValue); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]map[string]*dynamodb.AttributeValue)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}