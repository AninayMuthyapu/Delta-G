package extender

import (
	"encoding/json"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/your-org/blind-gpu-scheduler/pkg/spire"
)

const DefaultRequiredAnnotation = "attestation-hash.my-company.com/required-hash"

type Service struct {
	spireClient      spire.Client
	requiredAnnotKey string
}

func NewService(sc spire.Client) *Service {
	key := os.Getenv("BGS_REQUIRED_ANNOTATION")
	if key == "" {
		key = DefaultRequiredAnnotation
	}
	return &Service{spireClient: sc, requiredAnnotKey: key}
}

func (s *Service) Filter(w http.ResponseWriter, r *http.Request) {
	var args ExtenderArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		writeError(w, http.StatusBadRequest, "decode args: "+err.Error())
		return
	}
	if args.Pod == nil {
		writeError(w, http.StatusBadRequest, "missing pod in args")
		return
	}

	required := args.Pod.GetAnnotations()[s.requiredAnnotKey]
	if required == "" {
		writeJSON(w, http.StatusOK, &ExtenderFilterResult{
			Error: "missing required attestation annotation",
		})
		return
	}

	candidateNames := getCandidateNames(&args)
	var (
		okNames []string
		failed = map[string]string{}
	)
	for _, name := range candidateNames {
		live, err := s.spireClient.NodeHash(name)
		if err != nil {
			failed[name] = "no attestation hash"
			continue
		}
		if live == required {
			okNames = append(okNames, name)
		} else {
			failed[name] = "attestation mismatch"
		}
	}

	resp := &ExtenderFilterResult{NodeNames: okNames, FailedNodes: failed}
	writeJSON(w, http.StatusOK, resp)
}

func getCandidateNames(args *ExtenderArgs) []string {
	if args.NodeNames != nil && len(args.NodeNames) > 0 {
		return args.NodeNames
	}
	var names []string
	if args.Nodes != nil {
		for _, n := range args.Nodes.Items {
			names = append(names, n.Name)
		}
	}
	return names
}

func writeError(w http.ResponseWriter, code int, msg string) {
	log.Error(msg)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(&ExtenderFilterResult{Error: msg})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// Healthz returns 200 OK. Useful for k8s probes.
func (s *Service) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
