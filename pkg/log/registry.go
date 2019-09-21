package log

import (
"fmt"
"sync"
)

// Register is the registry of subsystems in operation
var Register = make(Registry)

var globalLevel string

var mutex sync.Mutex

// Registry is the centralised store that links all the loggers so they  can
// be accessed programmatically
type Registry map[string]*Logger

// Add appends a new subsystem to its map for access and introspection
func (r *Registry) Add(s *Logger) {
	mutex.Lock()
	_, ok := (*r)[s.Name]
	if ok {
		fmt.Println(s.Name, "subsystem already registered")
	} else {
		(*r)[s.Name] = s
	}
	mutex.Unlock()
}

// List returns a string slice containing all the registered loggers
func (r *Registry) List() (out []string) {
	mutex.Lock()
	for _, x := range *r {
		out = append(out, x.Name)
	}
	mutex.Unlock()
	return
}

// Get returns the subsystem.
// This could then be used to close or set its level eg `*r.Get("subsystem").
// SetLevel("debug")`
func (r *Registry) Get(name string) (out *Logger) {
	var ok bool
	mutex.Lock()
	if out, ok = (*r)[name]; ok {
		mutex.Unlock()
	}
	return
}

// GetGlobalLevel returns the global level
func (r *Registry) GetGlobalLevel() string {
	mutex.Lock()
	out := globalLevel
	mutex.Lock()
	return out
}

// SetAllLevels sets the level in all registered loggers
func (r *Registry) SetAllLevels(level string) {
	level = sanitizeLoglevel(level)
	mutex.Lock()
	globalLevel = level
	mutex.Unlock()
	loggers := r.List()
	for _, x := range loggers {
		r.Get(x).SetLevel(level)
	}
}

