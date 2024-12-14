package utils

import "fmt"

type Node struct {
	Prev *Node
	Next *Node
	Data string
}

type List struct {
	Head *Node
	Tail *Node
	Size int
}

func (l *List) Add(data string, cache map[string]*Node, key string) {
	newNode := &Node{Data: data}
	cache[key] = newNode
	if l.Size >= 4 {
		l.remove(cache, key)
	}

	if l.Size == 0 {
		l.Head = newNode
		l.Tail = newNode
	} else {
		aux := l.Head
		l.Head = newNode
		newNode.Next = aux
		aux.Prev = newNode
	}
	l.Size += 1
}

func (l *List) Get(data string, cache map[string]*Node) (*string, error) {
	node := cache[data]
	if node == nil {
		return nil, fmt.Errorf("não existe no cache")
	}

	if node == l.Head {
		return &l.Head.Data, nil
	}

	if node == l.Tail {
		l.Tail = node.Prev
		l.Tail.Next = nil
		l.Head.Prev = node
		node.Prev = nil
		node.Next = l.Head
		l.Head = node
	} else {
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
		node.Next = l.Head
		node.Prev = nil
		l.Head.Prev = node
		l.Head = node
	}
	return &l.Head.Data, nil
}

func (l *List) remove(cache map[string]*Node, key string) {
	node := l.Tail
	aux := node.Prev
	node.Prev = nil
	aux.Next = nil
	l.Tail = aux
	delete(cache, key)
	l.Size -= 1
}
