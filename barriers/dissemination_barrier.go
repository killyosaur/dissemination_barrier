package barriers

import (
	"math/bits"
)

type Node struct {
	send   chan bool
	sense  bool
	parity int
	id     int
	flags  [][]bool
}

func pow2(n int) int {
	return 1 << n
}

func ceil_log2(n int) int {
	res := 31 - bits.TrailingZeros32(uint32(n))
	if (1 << res) != n {
		res++
	}

	return res
}

type NodeCollection struct {
	Nodes []*Node
}

func Init(threads int) *NodeCollection {
	collection := &NodeCollection{Nodes: make([]*Node, threads)}

	for i := 0; i < threads; i++ {
		node := &Node{parity: 0, sense: true, id: i, send: make(chan bool), flags: make([][]bool, 2)}
		node.flags[0] = make([]bool, ceil_log2(threads))
		node.flags[1] = make([]bool, ceil_log2(threads))

		collection.Nodes = append(collection.Nodes, node)
	}

	return collection
}

func (coll *NodeCollection) Barrier(current int) {
	threads := len(coll.Nodes)

	rounds := ceil_log2(threads)

	node := coll.Nodes[current]

	for i := 0; i < rounds; i++ {
		sent_to := (current + pow2(i)) % threads
		sent_from := (current + threads - pow2(i)) % threads

		coll.Nodes[sent_to].send <- node.sense
		node.flags[node.parity][i] = <-coll.Nodes[sent_from].send
	}

	if node.parity == 1 {
		node.sense = !node.sense
	}

	node.parity = 1 - node.parity
}
