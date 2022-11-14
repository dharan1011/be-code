## Requirements
1. The application must return every DevEUI that it registers with the LoRaWAN provider
(e.g., if the application is killed it must wait for in-flight requests to finish otherwise, we
would have registered those DevEUIs but would not be using them)
2. It must handle user interrupts gracefully (SIGINT)
3. It must register exactly 100 DevEUIs (no more) with the provider (to avoid paying for
DevEUIs that we do not use)
4. It should make multiple requests concurrently (but there must never be more than 10
requests in-flight, to avoid throttling)

## Instrucation to RUN APP
```
go build -o deveui
./deveui
```