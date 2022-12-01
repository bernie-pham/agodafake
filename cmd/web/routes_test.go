package main

import (
	"fmt"
	"testing"

	"github.com/bernie-pham/agodafake/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app_config *config.AppConfig

	mux := routes(app_config)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but it is %T", v))
	}

}
