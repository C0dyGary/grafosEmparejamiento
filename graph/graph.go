package graph

type Edge struct {
	U int
	V int
}

type Graph struct {
	Vertices int
	Edges    []Edge
}

func Normalize(u, v int) Edge {
	if u > v {
		u, v = v, u
	}
	return Edge{U: u, V: v}
}

func MaximalMatching(graph Graph) []Edge {
	used := make([]bool, graph.Vertices)
	matching := make([]Edge, 0)

	for _, edge := range graph.Edges {
		if !used[edge.U] && !used[edge.V] {
			matching = append(matching, edge)
			used[edge.U] = true
			used[edge.V] = true
		}
	}

	return matching
}

type Blossom struct {
	graph   [][]int
	match   []int
	parent  []int
	base    []int
	used    []bool
	blossom []bool
	queue   []int
}

func NewBlossom(graph Graph) *Blossom {
	adjacency := make([][]int, graph.Vertices)

	for _, edge := range graph.Edges {
		adjacency[edge.U] = append(adjacency[edge.U], edge.V)
		adjacency[edge.V] = append(adjacency[edge.V], edge.U)
	}

	return &Blossom{
		graph:   adjacency,
		match:   make([]int, graph.Vertices),
		parent:  make([]int, graph.Vertices),
		base:    make([]int, graph.Vertices),
		used:    make([]bool, graph.Vertices),
		blossom: make([]bool, graph.Vertices),
	}
}

func (b *Blossom) lca(a, c int) int {
	usedPath := make([]bool, len(b.graph))

	for {
		a = b.base[a]
		usedPath[a] = true

		if b.match[a] == -1 {
			break
		}

		a = b.parent[b.match[a]]
	}

	for {
		c = b.base[c]

		if usedPath[c] {
			return c
		}

		c = b.parent[b.match[c]]
	}
}

func (b *Blossom) markPath(v, base, child int) {
	for b.base[v] != base {
		b.blossom[b.base[v]] = true
		b.blossom[b.base[b.match[v]]] = true

		b.parent[v] = child
		child = b.match[v]
		v = b.parent[b.match[v]]
	}
}

func (b *Blossom) findPath(root int) int {
	n := len(b.graph)

	for i := range n {
		b.used[i] = false
		b.parent[i] = -1
		b.base[i] = i
	}

	b.queue = b.queue[:0]
	b.queue = append(b.queue, root)
	b.used[root] = true

	for head := 0; head < len(b.queue); head++ {
		v := b.queue[head]

		for _, u := range b.graph[v] {
			if b.base[v] == b.base[u] || b.match[v] == u {
				continue
			}

			if u == root || (b.match[u] != -1 && b.parent[b.match[u]] != -1) {
				curBase := b.lca(v, u)

				for i := range n {
					b.blossom[i] = false
				}

				b.markPath(v, curBase, u)
				b.markPath(u, curBase, v)

				for i := range n {
					if b.blossom[b.base[i]] {
						b.base[i] = curBase

						if !b.used[i] {
							b.used[i] = true
							b.queue = append(b.queue, i)
						}
					}
				}
			} else if b.parent[u] == -1 {
				b.parent[u] = v

				if b.match[u] == -1 {
					return u
				}

				u = b.match[u]
				b.used[u] = true
				b.queue = append(b.queue, u)
			}
		}
	}

	return -1
}

func (b *Blossom) augment(v int) {
	for v != -1 {
		parent := b.parent[v]
		next := -1

		if parent != -1 {
			next = b.match[parent]
			b.match[v] = parent
			b.match[parent] = v
		}

		v = next
	}
}

func (b *Blossom) MaximumMatching() []Edge {
	n := len(b.graph)

	for i := range n {
		b.match[i] = -1
	}

	for i := range n {
		if b.match[i] == -1 {
			endpoint := b.findPath(i)

			if endpoint != -1 {
				b.augment(endpoint)
			}
		}
	}

	matching := make([]Edge, 0)

	for u := range n {
		v := b.match[u]

		if v != -1 && u < v {
			matching = append(matching, Edge{
				U: u,
				V: v,
			})
		}
	}

	return matching
}

func MaximumMatching(graph Graph) []Edge {
	blossom := NewBlossom(graph)
	return blossom.MaximumMatching()
}
