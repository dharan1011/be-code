package lorawan

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dharan1011/be-code/internal/entity"
	"log"
	"net/http"
)

const (
	ENDPOINT          = "europe-west1-machinemax-dev-d524.cloudfunctions.net"
	HTTP_SCHEMA       = "http"
	HTTPS_SCHEMA      = HTTP_SCHEMA + "s"
	SCHEMA_DELIMITER  = "://"
	URL               = HTTP_SCHEMA + SCHEMA_DELIMITER + ENDPOINT
	POST_CONTENT_TYPE = "application/json"
)

type RegistrationResponse struct {
	code int
}

func (rr *RegistrationResponse) IsSuccessful() bool {
	return rr.code == 200
}

func (rr *RegistrationResponse) IsSensorAlreadyRegistered() bool {
	return rr.code == 422
}

type LoRaWanAPIClient struct {
	Endpoint string
	// lock       sync.RWMutex
	httpClient http.Client
}

func NewLoRaWANApiClient() (*LoRaWanAPIClient, error) {
	return &LoRaWanAPIClient{
		Endpoint:   URL,
		httpClient: *http.DefaultClient,
	}, nil
}

func (c *LoRaWanAPIClient) RegisterSensor(sensorId string) (*RegistrationResponse, error) {
	requestBody := entity.NewPostDevEUIRegistrationPostBody(sensorId)
	marshalledBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.New("LoRaWanAPIClientError: Error marshalling POST request body")
	}
	req, err := createDevEUIRegistrationRequest(c.Endpoint, marshalledBody)
	if err != nil {
		return nil, errors.New("LoRaWanAPIClientError: Error creating HTTP request")
	}
	resp, err := c.httpClient.Do(req)
	log.Println("Debug: Response Status Code", resp.StatusCode)
	if err != nil {
		return nil, err
	}
	return &RegistrationResponse{code: resp.StatusCode}, nil
}

func createDevEUIRegistrationRequest(endpoint string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", POST_CONTENT_TYPE)
	return req, nil
}
