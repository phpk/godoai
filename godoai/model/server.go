package model

import (
	"net/http"
)

type DownserverStucct struct {
	Path string `json:"path"`
}

func DownServerHandler(w http.ResponseWriter, r *http.Request) {

}
