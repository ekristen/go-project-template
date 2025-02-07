package registry

var registry = make(map[string]WithID)

// Register a route with the registry.
func Register(r WithID) {
	if _, ok := registry[r.ID()]; ok {
		panic("route already registered")
	}

	registry[r.ID()] = r
}

// Get a route from the registry by ID.
func Get(id string) WithID {
	if r, ok := registry[id]; ok {
		return r
	}
	return nil
}

// GetRegistry returns the entire registry.
func GetRegistry() map[string]WithID {
	return registry
}

// GetRoutesWithPermissions returns all routes that have permissions defined.
func GetRoutesWithPermissions() []WithPermission {
	var routes []WithPermission
	for _, r := range registry {
		if rp, ok := r.(WithPermission); ok {
			routes = append(routes, rp)
		}
	}
	return routes
}
