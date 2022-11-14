package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dharan1011/be-code/internal/app"
	"github.com/dharan1011/be-code/internal/generator"
	"github.com/dharan1011/be-code/internal/lorawan"
)

const (
	DEV_EUI_LEN         = 16
	API_CALL_BATCH_SIZE = 10
	MAX_HEX_BUFFER_SIZE = API_CALL_BATCH_SIZE
)

func runCliApp() {
	lorawanApiClient, err := lorawan.NewLoRaWANApiClient()
	if err != nil {
		log.Fatal("Error creating LoRaWAN API client")
	}
	devEUIGenerator, err := generator.NewDevEUIGenerator(DEV_EUI_LEN, MAX_HEX_BUFFER_SIZE)
	if err != nil {
		log.Fatal("InternalError. Error creating DevEUI Generator")
	}
	cliApp, err := app.NewDevEUIApplication(devEUIGenerator, lorawanApiClient, API_CALL_BATCH_SIZE)
	if err != nil {
		log.Fatal("InternalError. Error creating CLI Application")
	}
	cliApp.Start()
	defer cliApp.GracefulShutdown()

	// Handling ctrl + c to intiate gracefull shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for _ = range c {
			log.Println("Greacfully Shutting down application. Waiting to finish API Call in process")
			cliApp.GracefulShutdown()
			os.Exit(1)
		}
	}()
	cliApp.Register(100)
}

func main() {
	runCliApp()
}
