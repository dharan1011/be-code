package generator

import (
	"crypto/rand"
	"errors"
	"math/big"
	"sync"
)

const ALLOWED_CHARS = "ABCDEF0123456789"

type DevEUIGenerator struct {
	stopChannel         chan bool
	idChannel           chan string
	IdLen               int
	ChannelBufferLength int
	createdHex          map[string]struct{}
	lock                sync.Mutex
}

func NewDevEUIGenerator(idLen, channelBufferLength int) (*DevEUIGenerator, error) {
	if channelBufferLength <= 0 {
		return nil, errors.New("DevEUIGeneratorError: Invalid channel error buffer length. Channel Buffer length >= 1")
	}
	if idLen <= 0 {
		return nil, errors.New("DevEUIGeneratorError: Invalid DevEUI id length. Recommended DevEUI id length is 16")

	}
	return &DevEUIGenerator{
		IdLen:               idLen,
		ChannelBufferLength: channelBufferLength,
		stopChannel:         make(chan bool),
		idChannel:           make(chan string, channelBufferLength),
		createdHex:          make(map[string]struct{}),
	}, nil
}

func (d *DevEUIGenerator) Run() {
	go d.runDevEUIGenerator()
}

func (d *DevEUIGenerator) runDevEUIGenerator() {
	devEUIChannel := d.idChannel
	closeChannel := d.stopChannel
	for {
		select {
		case _ = <-closeChannel:
			close(devEUIChannel)
			return
		default:
			hex, _ := generateHexString(d.IdLen)
			devEUIChannel <- hex
		}
	}
}

func (d *DevEUIGenerator) GetDevEUI() string {
	// Acquire lock
	d.lock.Lock()
	defer d.lock.Unlock()
	// Check & get unique DevEUI
	hex := <-d.idChannel
	for ; HasIdAlreadyGenerated(hex, d.createdHex); hex = <-d.idChannel {
	}
	return hex
}

func (d *DevEUIGenerator) Stop() {
	d.stopChannel <- true
}

/*
Internal Functions to generate hex
*/
func generateHexString(length int) (string, error) {

	max := big.NewInt(int64(len(ALLOWED_CHARS)))
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = ALLOWED_CHARS[n.Int64()]
	}
	return string(b), nil
}

/*
Internal function to get short code give DevEUI as input
*/
func getShortCode(hex string) string {
	return hex[len(hex)-5:]
}

/*
Internal function to check if map contains input key
*/
func HasIdAlreadyGenerated(key string, set map[string]struct{}) bool {
	_, ok := set[getShortCode(key)]
	return ok
}
