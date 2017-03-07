package govuegui

// Option for an element is represented by a function, which takes
// strings as input. For example Class() could be a function where
// css classes can added to a element.
// Class("active","box")
type Option func(vals ...string)

// Element represents a simple html element
type Element struct {
	Key     string
	Name    string
	Options []Option
}

// Container is the way elements are grouped. Every Element
type Container struct {
	Elements []*Element
}

// Form wrapps one ore more Containers
type Form struct {
	Containers []*Container
}
