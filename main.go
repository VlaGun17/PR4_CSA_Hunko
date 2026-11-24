package main

import(
	"net/http"
	"log"
	"PR4_CSA_Hunko/internal/handler"
	"PR4_CSA_Hunko/internal/repository"
)

func main(){
	deviceRepo := repository.NewMemoryDeviceRepository()
	router := handler.NewHandler(deviceRepo)

	log.Println("Server is starting on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}