package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/NJRodriguez/shiny-waddle/api/controllers/payloads/requests"
	"github.com/NJRodriguez/shiny-waddle/api/controllers/payloads/responses"
	"github.com/NJRodriguez/shiny-waddle/api/models"
	documentsMock "github.com/NJRodriguez/shiny-waddle/lib/aws/dynamodb/mocks"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type APIControllerTestSuite struct {
	suite.Suite
	controller    *APIController
	documentsMock *documentsMock.DocumentsClient
	router        *mux.Router
}

type testCaseResult struct {
	status int
	value  interface{}
}

func executeRequest(req *http.Request, router *mux.Router) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	return response
}

func (testSuite *APIControllerTestSuite) SetupTest() {
	testSuite.documentsMock = &documentsMock.DocumentsClient{}
	controller, _ := NewAPIController(testSuite.documentsMock)
	testSuite.controller = controller
	testSuite.router = mux.NewRouter()
	controller.RegisterRoutes(testSuite.router)
}

func (testSuite *APIControllerTestSuite) TestGetClosestSucursalWithInvalidParamsReturnsBadRequest() {
	request, reqErr := http.NewRequest("GET", "/sucursal/invalid/invalid", nil)

	testSuite.Require().NoError(reqErr)
	expectedResult := testCaseResult{
		http.StatusBadRequest,
		responses.ErrorMsg{
			Message: invalidLatitude,
		},
	}
	testSuite.verifyResponse(request, expectedResult)
}

func (testSuite *APIControllerTestSuite) TestGetClosestSucursalWithValidParamsReturnsClosestSucursal() {
	mockPosition := &models.Position{Latitude: 10.20, Longitude: 102.80}
	request, reqErr := http.NewRequest("GET", fmt.Sprintf("/sucursal/%f/%f", mockPosition.Latitude, mockPosition.Longitude), nil)
	mockSucursales := []models.Sucursal{
		{
			ID:        "winner",
			Address:   "123 Fake St",
			Latitude:  10.4,
			Longitude: 104.5,
		},
		{
			ID:        "failure",
			Address:   "587 Obelisco",
			Latitude:  50.32548,
			Longitude: 2.25468,
		},
		{
			ID:        "failure's cousin",
			Address:   "DEMASIADO LEJOS",
			Latitude:  80.2564,
			Longitude: 48.5648,
		},
	}
	marshaledMockSucursales := []map[string]*dynamodb.AttributeValue{}
	for _, sucursal := range mockSucursales {
		marshaledSucursal, err := dynamodbattribute.MarshalMap(sucursal)
		if err != nil {
			testSuite.Fail("Dynamodbattribute MarshalMap Failure")
		}
		marshaledMockSucursales = append(marshaledMockSucursales, marshaledSucursal)
	}
	testSuite.documentsMock.On("ListAll").Return(marshaledMockSucursales, nil).Once()
	testSuite.Require().NoError(reqErr)
	expectedDistanceInKm := calcDistance(mockPosition, &mockSucursales[0])
	expectedResult := testCaseResult{
		http.StatusOK,
		responses.ClosestSucursalResponse{
			Sucursal:     mockSucursales[0],
			DistanceInKm: expectedDistanceInKm,
		},
	}
	testSuite.verifyResponse(request, expectedResult)
}

func (testSuite *APIControllerTestSuite) TestCreateSucursalWithInvalidParamsReturnsBadRequest() {
	mockLat := 150.20
	mockLon := 2000.500
	request, reqErr := http.NewRequest("POST", "/sucursal", convertStructToBuffer(requests.PostSucursal{
		ID:        "INVALID ID",
		Address:   "",
		Latitude:  &mockLat,
		Longitude: &mockLon,
	}))

	testSuite.Require().NoError(reqErr)
	expectedResult := testCaseResult{
		http.StatusBadRequest,
		ApiError{
			Message: "Error when validating payload",
			Errors: []string{
				"ID must be in valid UUID v4 format",
				"Address is a required field",
				"Latitude must be 90 or less",
				"Longitude must be 180 or less",
			},
		},
	}
	testSuite.verifyResponse(request, expectedResult)
}

func (testSuite *APIControllerTestSuite) TestCreateSucursalWithValidParamsReturnsStatusOK() {
	mockLat := 20.252
	mockLon := 50.685
	mockUUID := uuid.NewV4()
	mockPostSucursal := requests.PostSucursal{
		ID:        mockUUID.String(),
		Address:   "123 Fake St.",
		Latitude:  &mockLat,
		Longitude: &mockLon,
	}
	request, reqErr := http.NewRequest("POST", "/sucursal", convertStructToBuffer(mockPostSucursal))

	testSuite.documentsMock.On("Create", &mockPostSucursal).Return(nil, nil).Once()

	testSuite.Require().NoError(reqErr)
	expectedResult := testCaseResult{
		http.StatusOK,
		responses.PostSucursal{
			Message: "Successfully created sucursal",
			ID:      mockUUID.String(),
		},
	}
	testSuite.verifyResponse(request, expectedResult)
}

func (testSuite *APIControllerTestSuite) verifyResponse(request *http.Request, expectedResult testCaseResult) {
	response := executeRequest(request, testSuite.router)
	parsedResult := response.Body.String()
	parsedExpectedResult := ""
	if reflect.TypeOf(expectedResult.value).Kind() != reflect.String {
		result, _ := json.Marshal(expectedResult.value)
		parsedExpectedResult = string(result)
	} else {
		parsedExpectedResult = expectedResult.value.(string)
	}

	testSuite.Require().Equal(strings.Trim(parsedExpectedResult, "\n"), strings.Trim(parsedResult, "\n"), "result mismatch")
	responseResult := response.Result()
	responseResult.Body.Close()
	testSuite.Require().Equal(expectedResult.status, responseResult.StatusCode, "status code mismatch")
}

func convertStructToBuffer(structure interface{}) *bytes.Buffer {
	marshaledStruct, _ := json.Marshal(structure)
	return bytes.NewBuffer(marshaledStruct)
}

func TestApiControllerTestSuite(t *testing.T) {
	suite.Run(t, new(APIControllerTestSuite))
}
