package actionqueue

import (
	"fmt"
	"sync"
)

var (
	mu       sync.RWMutex
	registry = make(map[string]*Definition)
)

// Register adds a definition to the module-level registry.
// Called from generated init() functions.
func Register(def *Definition) {
	mu.Lock()
	defer mu.Unlock()

	if def.Name == "" {
		panic("actionqueue: definition name is required")
	}
	if _, exists := registry[def.Name]; exists {
		panic(fmt.Sprintf("actionqueue: definition %q already registered", def.Name))
	}
	registry[def.Name] = def
}

// Get returns a definition by name, or nil if not found.
func Get(name string) *Definition {
	mu.RLock()
	defer mu.RUnlock()
	return registry[name]
}

// All returns all registered definitions.
func All() []*Definition {
	mu.RLock()
	defer mu.RUnlock()

	defs := make([]*Definition, 0, len(registry))
	for _, def := range registry {
		defs = append(defs, def)
	}
	return defs
}
