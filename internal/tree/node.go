package tree

type NodeType int

const (
	ObjectOpen NodeType = iota
	ObjectClose
	ArrayOpen
	ArrayClose
	KeyValue
	ArrayElement
)

type Node struct {
	Type        NodeType
	Depth       int
	Key         string
	Value       any
	Path        string
	Collapsible bool
	Collapsed   bool
	ChildCount  int
	IsLast      bool
}
