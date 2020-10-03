package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/NJRodriguez/shiny-waddle/api/controllers/payloads/requests"
	"github.com/NJRodriguez/shiny-waddle/api/controllers/payloads/responses"
	"github.com/NJRodriguez/shiny-waddle/api/models"
	"github.com/NJRodriguez/shiny-waddle/lib/aws/dynamodb"
	"github.com/aws/aws-sdk-go/aws/awserr"
	dynamodbSdk "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

const (
	internalServerError = "Internal server error"
	idExistsError       = "Id already exists in database"
	idNotFoundError     = "Id not found in database"
	invalidRequestBody  = "Failed to parse the request body"
)

type APIController struct {
	documentsClient dynamodb.DocumentsClient
}

type APIControllerArgs struct {
	TableName string
	Region    string
}

func NewAPIController(args *APIControllerArgs) (*APIController, error) {
	client, err := dynamodb.New(args.TableName, args.Region)
	if err != nil {
		log.Fatalln("Error when trying to start DynamoDB Client.")
		return nil, err
	}
	validate = validator.New()
	translator, err = RegisterErrors(validate)
	if err != nil {
		log.Fatalln("Error when trying to register api errors.")
		return nil, err
	}
	return &APIController{
		documentsClient: client,
	}, nil
}

func (instance *APIController) RegisterRoutes(router *mux.Router) {

	//Sucursales routes
	router.HandleFunc("/sucursal", instance.CreateSucursal).Methods("POST")
	router.HandleFunc("/sucursal/{id}", instance.GetSucursal).Methods("GET")
	router.HandleFunc("/sucursal/{lat}/{lon}", instance.GetClosestSucursal).Methods("GET")
}

func (instance *APIController) CreateSucursal(writer http.ResponseWriter, r *http.Request) {
	setJSONContentType(writer)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error when trying to read request body: %s", err)
		writer.WriteHeader(http.StatusBadRequest)
		generateErrorMessage(writer, &responses.ErrorMsg{Message: invalidRequestBody})
		return
	}
	valErrs, err := ValidateRequest(body, &requests.PostSucursal{})
	if err != nil {
		log.Printf("Error when trying to validate requests: %s", err)
		writer.WriteHeader(http.StatusBadRequest)
		generateErrorMessage(writer, &responses.ErrorMsg{Message: invalidRequestBody})
		return
	}
	if valErrs != nil {
		log.Println("Validation error in payload.")
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(valErrs)
		return
	}
	sucursal, err := deserializePostSucursalRequest(body)
	if err != nil {
		log.Printf("Error when trying to deserialize request body: %s", err)
		writer.WriteHeader(http.StatusBadRequest)
		generateErrorMessage(writer, &responses.ErrorMsg{Message: internalServerError})
	}
	_, err = instance.documentsClient.Create(sucursal)
	if err != nil {
		log.Printf("Error when trying to create Sucursal: %s", err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodbSdk.ErrCodeConditionalCheckFailedException:
				writer.WriteHeader(http.StatusConflict)
				generateErrorMessage(writer, &responses.ErrorMsg{Message: idExistsError})
				return
			default:
				writer.WriteHeader(http.StatusBadRequest)
				generateErrorMessage(writer, &responses.ErrorMsg{Message: internalServerError})
			}
		}
		return
	}
	_ = json.NewEncoder(writer).Encode(responses.PostSucursal{Message: "Successfully created sucursal", ID: sucursal.ID})
}

func (instance *APIController) GetSucursal(writer http.ResponseWriter, r *http.Request) {
	setJSONContentType(writer)
	pathVars := mux.Vars(r)
	id := pathVars["id"]
	log.Printf("Got the following id from path variables: %s", id)
	sucursalKey := models.SucursalKey{
		ID: id,
	}
	result, err := instance.documentsClient.Get(sucursalKey)
	if err != nil {
		log.Println("Error when trying to get Sucursal from db.")
		writer.WriteHeader(http.StatusBadRequest)
		generateErrorMessage(writer, &responses.ErrorMsg{Message: internalServerError})
		return
	}
	if result.Item == nil {
		log.Println("Sucursal does not exist in db.")
		writer.WriteHeader(http.StatusBadRequest)
		generateErrorMessage(writer, &responses.ErrorMsg{Message: idNotFoundError})
		return
	}
	sucursal := models.Sucursal{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &sucursal)
	if err != nil {
		log.Println("Error when trying to parse Sucursal Object.")
		writer.WriteHeader(http.StatusBadRequest)
		generateErrorMessage(writer, &responses.ErrorMsg{Message: internalServerError})
		return
	}
	_ = json.NewEncoder(writer).Encode(sucursal)
}

func (instance *APIController) GetClosestSucursal(writer http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(writer).Encode("Get closest sucursal not implemented yet!")
}

func generateErrorMessage(writer http.ResponseWriter, msg *responses.ErrorMsg) {
	_ = json.NewEncoder(writer).Encode(msg)
}

func setJSONContentType(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func deserializePostSucursalRequest(request []byte) (*requests.PostSucursal, error) {
	decoder := json.NewDecoder(bytes.NewReader(request))
	decoder.DisallowUnknownFields()
	var postRequest requests.PostSucursal
	err := decoder.Decode(&postRequest)
	if err != nil {
		return nil, errors.Wrap(err, "deserializing post sucursal request")
	}
	return &postRequest, nil
}
