package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"

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
	internalServerError     = "Internal server error"
	idExistsError           = "Id already exists in database"
	idNotFoundError         = "Id not found in database"
	sucursalesNotFoundError = "No sucursales were found. Please load sucursales onto database"
	invalidRequestBody      = "Failed to parse the request body"
	invalidLatitude         = "Latitude must be in float64 format"
	invalidLongitude        = "Longitude must be in float64 format"
	invalidLatitudeVal      = "Latitude must be between -90 and 90"
	invalidLongitudeVal     = "Longitude must be between -180 and 180"
)

type APIController struct {
	documentsClient dynamodb.DocumentsClient
}

type APIControllerArgs struct {
	TableName string
	Region    string
}

func NewAPIController(documentsClient dynamodb.DocumentsClient) (*APIController, error) {
	validate = validator.New()
	generatedTranslator, err := RegisterErrors(validate)
	if err != nil {
		log.Fatalln("Error when trying to register api errors.")
		return nil, err
	}
	translator = generatedTranslator
	return &APIController{
		documentsClient,
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
	setJSONContentType(writer)
	pathVars := mux.Vars(r)
	lat := pathVars["lat"]
	lon := pathVars["lon"]
	position, err := validateLatLon(lat, lon)
	if err != nil {
		log.Println("Error when validating latitude/longitude")
		writer.WriteHeader(http.StatusBadRequest)
		generateErrorMessage(writer, &responses.ErrorMsg{Message: err.Error()})
		return
	}
	result, err := instance.documentsClient.ListAll()
	if err != nil {
		log.Println("Error when trying to list all items from dynamodb table.")
		writer.WriteHeader(http.StatusBadRequest)
		generateErrorMessage(writer, &responses.ErrorMsg{Message: internalServerError})
		return
	}
	if len(result) == 0 {
		log.Println("No sucursales are loaded in database!")
		writer.WriteHeader(http.StatusBadRequest)
		generateErrorMessage(writer, &responses.ErrorMsg{Message: sucursalesNotFoundError})
		return
	}
	sucursales, err := models.ToSucursalArray(result)
	if err != nil {
		log.Println("Error when trying to convert dynamodb result to sucursales array.")
		writer.WriteHeader(http.StatusBadRequest)
		generateErrorMessage(writer, &responses.ErrorMsg{Message: internalServerError})
		return
	}
	closestSucursal := &models.SucursalWithDistance{}
	for _, sucursal := range sucursales {
		distance := calcDistance(position, sucursal)
		if closestSucursal.Sucursal == nil {
			closestSucursal = &models.SucursalWithDistance{Sucursal: sucursal, Distance: distance}
		}
		if distance < closestSucursal.Distance {
			closestSucursal.Sucursal = sucursal
			closestSucursal.Distance = distance
		}
	}
	_ = json.NewEncoder(writer).Encode(&responses.ClosestSucursalResponse{Sucursal: *closestSucursal.Sucursal, DistanceInKm: closestSucursal.Distance})
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

func validateLatLon(lat string, lon string) (*models.Position, error) {
	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return nil, errors.New(invalidLatitude)
	}
	if latFloat > 90 || latFloat < -90 {
		return nil, errors.New(invalidLatitudeVal)
	}
	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return nil, errors.New(invalidLongitude)
	}
	if lonFloat > 180 || lonFloat < -180 {
		return nil, errors.New(invalidLongitudeVal)
	}
	return &models.Position{Latitude: latFloat, Longitude: lonFloat}, nil
}

func calcDistance(position *models.Position, sucursal *models.Sucursal) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * position.Latitude / 180)
	radlat2 := float64(PI * sucursal.Latitude / 180)

	theta := float64(position.Longitude - sucursal.Longitude)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515
	dist = dist * 1.609344

	return dist
}
