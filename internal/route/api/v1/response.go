package v1

import (
	"eMobile/internal/crud"
	"encoding/json"
	"net/http"
)

type ResponseBase[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type ResponseBaseErr struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type ResponseBasePaginated[T any] struct {
	Message    string          `json:"message"`
	Pagination crud.Pagination `json:"next_pagination"`
	Data       []T             `json:"data"`
}

type ResponseBaseMulti[T any] struct {
	Message string `json:"message"`
	Data    []T    `json:"data"`
}

func WriteResponse[Data any](w http.ResponseWriter, code int, data Data, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := ResponseBase[Data]{
		Message: msg,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

func WriteResponsePaginated[Data any](w http.ResponseWriter, code int, pag crud.Pagination, data []Data, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := ResponseBasePaginated[Data]{
		Message:    msg,
		Data:       data,
		Pagination: pag,
	}

	json.NewEncoder(w).Encode(resp)
}

func WriteResponseErr(w http.ResponseWriter, code int, err error, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err == nil {
		return
	}

	resp := ResponseBaseErr{
		Message: msg,
		Error:   err.Error(),
	}

	json.NewEncoder(w).Encode(resp)
}
