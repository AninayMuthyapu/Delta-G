package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/your-org/blind-gpu-scheduler/pkg/extender"
	"github.com/your-org/blind-gpu-scheduler/pkg/logging"
	"github.com/your-org/blind-gpu-scheduler/pkg/spire"
)

func main() {
	logging.Init()

	sc, err := spire.NewEnvMockClient()
	if err != nil {
		log.Fatalf("failed to init SPIRE mock client: %v", err)
	}
	svc := extender.NewService(sc)

	r := mux.NewRouter()
	r.HandleFunc("/healthz", svc.Healthz).Methods("GET")
	r.HandleFunc("/filter", svc.Filter).Methods("POST")

	addr := os.Getenv("BGS_LISTEN_ADDR")
	if addr == "" {
		addr = ":8000"
	}
	log.WithField("addr", addr).Info("starting blind-gpu-scheduler extender")
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
