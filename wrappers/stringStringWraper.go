package wrappers

import (
	ct "../concurentTree"
)

type StringStringDb struct {
	avltree *ct.ConcurentTree
}

func Open(address, name string, cacheCount int, rewrite bool) *StringStringDb {
	avltree := ct.Open(address, name, cacheCount, rewrite)
	ssdb := &StringStringDb{avltree: avltree}
	return ssdb
}

func (ssdb *StringStringDb) Close() error {
	return ssdb.avltree.Close()
}

func (ssdb *StringStringDb) Get(key string) string {
	data, _ := ssdb.avltree.Get([]byte(key))
	return string(data)
}

func (ssdb *StringStringDb) Add(key, value string) error {
	return ssdb.avltree.Set([]byte(key), []byte(value))
}

func (ssdb *StringStringDb) Del(key string) error {
	return ssdb.avltree.Del([]byte(key))
}
