
package errors

import (
	"fmt"
	"log"
	"net/http"
)

func handleError(w http.ResponseWriter, err error) {
	log.Println(err)

	w.WriteHeader(http.StatusInternalServerError) 
	fmt.Fprintf(w, err.Error())
}
