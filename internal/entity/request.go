package entity

type RequestBody struct {
	DevEUI map[string]string `json:"deveui"`
}

/*
Requirements are unclear around the POST body schema. Left it open to fix it
*/
func NewPostDevEUIRegistrationPostBody(deveui string) *RequestBody {
	return &RequestBody{}
}
