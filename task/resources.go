package task

// ResourcesSet represents a set of task resources.
// A resource is represented by its unique key
type ResourcesSet struct {
	set map[string]bool
}

// NewResourcesSet creates new instance of resources set
func NewResourcesSet() *ResourcesSet {
	return &ResourcesSet{make(map[string]bool)}
}

// Add given key to resources set
func (r *ResourcesSet) Add(keys []string) {
	for _, key := range(keys) {
		r.set[key] = true
	}
}

// ContainsAny checks whether one of given keys
// is currently in set of resources.
func (r *ResourcesSet) ContainsAny(keys []string) bool {
	found := false
	for _, key := range(keys) {
		_, found = r.set[key]
		if found {
			break
		}
	}

	return found
}

// Remove removes given keys from dependency set
func (r *ResourcesSet) Remove(keys []string) {
	for _, key := range(keys) {
		delete(r.set, key)
	}
}
