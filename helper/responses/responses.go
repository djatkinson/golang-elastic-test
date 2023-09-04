package responses

import (
	"encoding/json"
	"net/http"
)

type Res struct {
	Error   bool        `json:"error"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Responses(res http.ResponseWriter, code int, msg string, payload interface{}) {
	var result Res
	if code >= 400 {
		result.Error = true
		result.Code = code
		result.Message = "Error"
		result.Data = payload
	} else {
		result.Error = false
		result.Code = code
		if msg == "" {
			msg = "Success"
		}
		result.Message = msg
		result.Data = payload
	}
	responses, _ := json.Marshal(result)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(responses)
}

func ErrorCustomStatus(res http.ResponseWriter, code int, msg string) {
	Responses(res, code, "", msg)
}
