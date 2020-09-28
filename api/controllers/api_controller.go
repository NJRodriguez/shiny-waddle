package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NJRodriguez/shiny-waddle/lib/aws/dynamodb"
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
	router.HandleFunc("/sucursal", instance.GetSucursal).Methods("GET")
	router.HandleFunc("/sucursal", instance.CreateSucursal).Methods("POST")
	router.HandleFunc("/sucursal/{lat}/{lon}", instance.GetClosestSucursal).Methods("GET")
}

func (instance *APIController) CreateSucursal(writer http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(writer).Encode("Create not implemented yet!")
}

func (instance *APIController) GetSucursal(writer http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(writer).Encode("Get sucursal not implemented yet!")
}

func (instance *APIController) GetClosestSucursal(writer http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(writer).Encode("Get closest sucursal not implemented yet!")
}
