package handle

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
)

// RPCErr prints an error and returns to the caller
func RPCErr(w http.ResponseWriter, err error, status int, msgs ...string) bool {
	if err != nil {
		log.Println(err, msgs)
		erpc.ResponseHandler(w, status)
		return true
	}

	return false
}
