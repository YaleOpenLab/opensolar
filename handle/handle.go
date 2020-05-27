package handle

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
)

// RPCErr prints an error and returns to the caller
func RPCErr(w http.ResponseWriter, err error, status int, msgs ...string) bool {
	var retmsg bool
	if len(msgs) == 2 {
		retmsg = true
	}

	if err != nil {
		log.Println(err, msgs[0])
		if retmsg {
			erpc.ResponseHandler(w, erpc.StatusBadRequest, msgs[1])
		} else {
			erpc.ResponseHandler(w, status)
		}
		return true
	}

	return false
}
