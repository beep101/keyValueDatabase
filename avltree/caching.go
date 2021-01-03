package avltree

type caching interface {
	add(node *Node)
	remove(node *Node)
}

type queueCache struct {
	first *queueElement
	last  *queueElement
	count int
	max   int
}

type queueElement struct {
	element *Node
	next    *queueElement
}

func createQueueCache(max int) *queueCache {
	return &queueCache{first: nil, last: nil, count: 0, max: max}
}

func (q *queueCache) add(node *Node) {
	elem := &queueElement{element: node, next: nil}
	if q.count == 0 {
		q.first = elem
		q.last = elem
		q.count++
	} else if q.count < q.max {
		q.last.next = elem
		q.last = elem
		q.count++
	} else {
		q.last.next = elem
		q.last = elem
		q.first.element.value = nil
		q.first = q.first.next
	}
}

func (q *queueCache) remove(node *Node) {
	if q.first.element == node {
		if q.last.element == node {
			q.last = nil
			q.first = nil
			q.count = 0
			return
		}
		q.first = q.first.next
		q.count--
	} else {
		iterator := q.first
		for {
			if iterator.next == nil {
				return
			}
			if iterator.next.element == node {
				iterator.next = iterator.next.next
				q.count--
				return
			}
			iterator = iterator.next
		}
	}
}

type emptyCache struct{}

func (ec *emptyCache) add(node *Node) {
	node.value = nil
	return
}

func (ec *emptyCache) remove(node *Node) {
	node.value = nil
	return
}

type fullCache struct{}

func (fc *fullCache) add(node *Node) {
	return
}

func (fc *fullCache) remove(node *Node) {
	return
}
