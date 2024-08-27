package sys

import (
	"godoai/libs"
	"net/http"
)

func ConfigHandle(w http.ResponseWriter, r *http.Request) {

	libs.Success(w, "success", "The config set success!")
}
