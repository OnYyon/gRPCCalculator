package rest

import (
	"fmt"
	"net/http"

	"github.com/OnYyon/gRPCCalculator/internal/services/manager"
)

type restAPI struct {
	manager *manager.Manager
}

func (a *restAPI) AddNewExpression(w http.ResponseWriter, r *http.Request) {
	fmt.Println("yesss")
}
