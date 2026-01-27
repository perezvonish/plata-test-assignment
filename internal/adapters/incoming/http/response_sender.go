package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type HeaderType string
type HeaderValue string

var (
	contentTypeHeader    = HeaderType("Content-Type")
	contentTypeValueJson = HeaderValue("application/json")
)

type body[T any] struct {
	Data  T      `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type SendResponseParams[T any] struct {
	Status int
	Data   T
	Error  error
}

func SendResponse[T any](w http.ResponseWriter, params SendResponseParams[T]) {
	var res body[T]

	if params.Error != nil {
		res = createErrorResponse[T](params.Error.Error())
	} else {
		res = createSuccessResponse[T](params.Data)
		if params.Status == 0 {
			params.Status = http.StatusOK
		}
	}

	w.Header().Set(string(contentTypeHeader), string(contentTypeValueJson))
	w.WriteHeader(params.Status)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("response error: %v", err)
	}
}

func createSuccessResponse[T any](data T) body[T] {
	return body[T]{
		Data: data,
	}
}

func createErrorResponse[T any](message string) body[T] {
	return body[T]{
		Error: message,
	}
}
