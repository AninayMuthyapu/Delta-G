package spire

// Client defines methods to retrieve node attestation evidence (hash).
// In the MVP we provide a mock implementation backed by an env-provided map.

type Client interface {
	// NodeHash returns the attestation hash for a given node name.
	NodeHash(nodeName string) (string, error)
}
