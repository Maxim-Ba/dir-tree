package tree

type Node struct {
	Children []*Node
	//TODO
}

// Get tree recurcive
func Get(path string) *Node {
	if path == "" {
		return nil
	}
	return &Node{}
}
func NewNode() *Node {
	return &Node{
		//TODO
	}
}
