package app

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/dharan1011/be-code/internal"
	"github.com/dharan1011/be-code/internal/generator"
	"github.com/dharan1011/be-code/internal/lorawan"
)

type DevEUIApplication struct {
	devEUIGenerator  *generator.DevEUIGenerator
	lorawanClient    *lorawan.LoRaWanAPIClient
	MaxBatchSize     int
	registeredDevEUI chan string
	wg               *sync.WaitGroup
	stop             bool
	printChannel     chan string
}

func NewDevEUIApplication(g *generator.DevEUIGenerator, lrwc *lorawan.LoRaWanAPIClient, maxBatchSize int) (*DevEUIApplication, error) {
	if g == nil {
		return nil, errors.New("DevEUIApplicationError: Input DevEUIGenerator cannot be nil.")
	}
	if lrwc == nil {
		return nil, errors.New("DevEUIApplicationError: Input LoRaWAN client cannot be nil.")
	}
	if maxBatchSize <= 0 {
		return nil, errors.New("DevEUIApplicationError: Batch size cannot be 0. Recomended batch size 10")

	}
	return &DevEUIApplication{devEUIGenerator: g,
		MaxBatchSize:     maxBatchSize,
		registeredDevEUI: make(chan string, maxBatchSize),
		printChannel:     make(chan string),
		wg:               new(sync.WaitGroup),
		stop:             false,
		lorawanClient:    lrwc,
	}, nil
}

func (m *DevEUIApplication) Start() {
	m.devEUIGenerator.Run()
	go m.printDevEUI()
}

func (m *DevEUIApplication) printDevEUI() {
	for id := range m.printChannel {
		fmt.Println("Registed DevEUI:", id)
	}
}

func (m *DevEUIApplication) Register(size int) {
	m.registerDevEUI(size)
}
func (m *DevEUIApplication) registerDevEUI(size int) {
	batchSize := m.MaxBatchSize
	sensorsToRegisterCount := size
	for sensorsToRegisterCount > 0 {
		for i := 0; i < internal.Min(sensorsToRegisterCount, batchSize); i++ {
			m.wg.Add(1)
			go func() {
				defer m.wg.Done()
				if m.stop {
					// Gracefull shutdown initiated
					return
				} else {
					for {
						generatedId := m.devEUIGenerator.GetDevEUI()
						res, err := m.lorawanClient.RegisterSensor(generatedId)
						if err != nil {
							log.Println("DevEUIApplicationError: Error making call to REST API.", err)
							return
						}
						if res.IsSuccessful() {
							m.printChannel <- generatedId
							break
						} else if res.IsSensorAlreadyRegistered() && !m.stop {
							log.Println(generatedId, "Already used. Retrying")
							generatedId = m.devEUIGenerator.GetDevEUI()
						} else {
							log.Println("DevEUIApplicationError: Something went wrong, unexcepted response status code.", err)
							return
						}
					}
				}
			}()
		}
		m.wg.Wait()
		sensorsToRegisterCount -= batchSize
	}
}

func (m *DevEUIApplication) GracefulShutdown() {
	m.stop = true
	m.wg.Wait()
}
