package handlers

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
	Error   []apiError  `json:"error"`
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ServeResponse is to unify the api responses
func ServeResponse(w http.ResponseWriter, success bool, result interface{}, err ...apiError) {
	w.Header().Set("Content-Type", "application/json")
	res := apiResponse{Success: success, Result: result, Error: err}
	if !success {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(res)
}
