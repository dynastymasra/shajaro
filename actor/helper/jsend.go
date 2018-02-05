package helper

//Jsend used to format JSON with jsend rules
type Jsend struct {
	Status  string      `json:"status" binding:"required"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// FailResponse is used to return response with JSON format if failure
func FailResponse(msg string) Jsend {
	return Jsend{Status: "failed", Message: msg}
}

// SuccessResponse used to return response with JSON format success
func SuccessResponse() Jsend {
	return Jsend{Status: "success"}
}

// ObjectResponse used to return response JSON format if have data value
func ObjectResponse(data interface{}) Jsend {
	return Jsend{Status: "success", Data: data}
}
