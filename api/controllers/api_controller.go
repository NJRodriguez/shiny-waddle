package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NJRodriguez/shiny-waddle/api/models"
	"github.com/NJRodriguez/shiny-waddle/lib/aws/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
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
		log.Fatal("Error when trying to parse Sucursal Key!")
		_ = json.NewEncoder(writer).Encode(err)
	}
	result, err := instance.documentsClient.Get(key)
	if err != nil {
		log.Fatal("Error when trying to get Sucursal from db!")
		_ = json.NewEncoder(writer).Encode(err)
	}
	sucursal := models.Sucursal{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &sucursal)
	if err != nil {
		log.Fatal("Error when trying to parse Sucursal Object!!")
		_ = json.NewEncoder(writer).Encode(err)
	}
	_ = json.NewEncoder(writer).Encode(sucursal)
}

func (instance *APIController) GetClosestSucursal(writer http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(writer).Encode("Get closest sucursal not implemented yet!")
}

func setJSONContentType(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}
