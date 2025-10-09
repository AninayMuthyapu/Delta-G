package spire

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type EnvMockClient struct {
	mu   sync.RWMutex
	hash map[string]string
}

func NewEnvMockClient() (*EnvMockClient, error) {
	m := &EnvMockClient{hash: map[string]string{}}
	if raw := os.Getenv("BGS_NODE_HASH_MAP"); raw != "" {
		var parsed map[string]string
		if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
			return nil, fmt.Errorf("invalid BGS_NODE_HASH_MAP JSON: %w", err)
		}
		m.hash = parsed
	}
	return m, nil
}

func (m *EnvMockClient) NodeHash(nodeName string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	h, ok := m.hash[nodeName]
	if !ok {
		return "", errors.New("node hash not found")
	}
	return h, nil
}

func (m *EnvMockClient) Update(newMap map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hash = newMap
}
