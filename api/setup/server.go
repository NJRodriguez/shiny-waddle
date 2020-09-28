package setup

import (
	"log"
	"net/http"

	"github.com/NJRodriguez/shiny-waddle/api/controllers"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func (server *Server) Initialize(tableName string, region string) error {
	args := &controllers.APIControllerArgs{
		TableName: tableName,
		Region:    region,
	}
	log.Println("Starting API Controller...")
	apiController, err := controllers.NewAPIController(args)
	if err != nil {
		log.Fatal("Error when trying to start API Controller!")
		return err
	}
	log.Println("Registering API Routes...")
	apiController.RegisterRoutes(server.Router)
	return nil
}

func (server *Server) Run(addr string) {
	log.Println("Listening to port 80")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
