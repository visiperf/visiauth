package renderer

import (
	"encoding/json"
	"errors"
	"net/http"
)

func Success(v interface{}, statusCode int, rw http.ResponseWriter) {
	data, err := json.Marshal(map[string]interface{}{"data": v})
	if err != nil {
		Error(err, rw)
		return
	}

	render(data, statusCode, "application/json", rw)
}

func Error(err error, rw http.ResponseWriter) {
	data, _ := json.Marshal(map[string]interface{}{"error": err})

	var sce interface {
		error
		StatusCode() int
	}

	sc := http.StatusInternalServerError
	if errors.As(err, &sce) {
		sc = sce.StatusCode()
	}

	render(data, sc, "application/json", rw)
}

func render(data []byte, statusCode int, contentType string, rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", contentType)
	rw.WriteHeader(statusCode)
	rw.Write(data)
}
