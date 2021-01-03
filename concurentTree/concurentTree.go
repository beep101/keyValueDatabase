package concurentTree

import (
	"sync"

	avl "../avltree"
)

type ConcurentTree struct {
	tree   *avl.Tree
	locker sync.RWMutex
}

func Open(address, name string, cacheCount int, rewrite bool) *ConcurentTree {
	tree := avl.Open(address, name, cacheCount, rewrite)
	return &ConcurentTree{tree: tree}
}

func (ct *ConcurentTree) Close() error {
	ct.locker.Lock()
	return ct.tree.Close()
}

func (ct *ConcurentTree) Get(key []byte) ([]byte, error) {
	ct.locker.RLock()
	defer ct.locker.RUnlock()
	return ct.tree.Get(key)
}

func (ct *ConcurentTree) Set(key, value []byte) error {
	ct.locker.Lock()
	defer ct.locker.Unlock()
	return ct.tree.Set(key, value)
}

func (ct *ConcurentTree) Del(key []byte) error {
	ct.locker.Lock()
	defer ct.locker.Unlock()
	return ct.tree.Del(key)
}
