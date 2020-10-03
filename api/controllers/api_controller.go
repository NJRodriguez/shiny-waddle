package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/NJRodriguez/shiny-waddle/api/controllers/payloads/requests"
	"github.com/NJRodriguez/shiny-waddle/api/models"
	"github.com/NJRodriguez/shiny-waddle/lib/aws/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate *validator.Validate
var translator ut.Translator

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
	writer.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when trying to read request body.")
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(err)
		return
	}
	valErrs, err := ValidateRequest(body, &requests.PostSucursal{})
	if err != nil {
		log.Println("Error when trying to validate requests.")
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(err)
		return
	}
	if valErrs != nil {
		log.Println("Validation error in payload.")
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(valErrs)
		return
	}
	_ = json.NewEncoder(writer).Encode("Create sucursal not implemented yet!")
}

func (instance *APIController) GetSucursal(writer http.ResponseWriter, r *http.Request) {
	setJSONContentType(writer)
	pathVars := mux.Vars(r)
	id := pathVars["id"]
	log.Printf("Got the following id from path variables: %s", id)
	sucursalKey := models.SucursalKey{
		ID: id,
	}
	key, err := dynamodbattribute.MarshalMap(sucursalKey)
	if err != nil {
		log.Println("Error when trying to parse Sucursal Key.")
		_ = json.NewEncoder(writer).Encode(err)
		return
	}
	result, err := instance.documentsClient.Get(key)
	if err != nil {
		log.Println("Error when trying to get Sucursal from db.")
		_ = json.NewEncoder(writer).Encode(err)
		return
	}
	sucursal := models.Sucursal{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &sucursal)
	if err != nil {
		log.Println("Error when trying to parse Sucursal Object.")
		_ = json.NewEncoder(writer).Encode(err)
		return
	}
	_ = json.NewEncoder(writer).Encode(sucursal)
}

func (instance *APIController) GetClosestSucursal(writer http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(writer).Encode("Get closest sucursal not implemented yet!")
}

func setJSONContentType(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}
