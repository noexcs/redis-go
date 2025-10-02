package maps

// Map is an interface for map data structures
type Map interface {
	// Put inserts or updates a key-value pair
	Put(key string, value string)

	// Get retrieves a value by key
	Get(key string) (value string, ok bool)

	// Contains checks if a key exists
	Contains(key string) bool
}
