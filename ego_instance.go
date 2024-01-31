package ego

const (
	DefaultName = "ego-default-instance"
)

// Instance
func Instance(names ...string) *Ego {
	return New()
}
