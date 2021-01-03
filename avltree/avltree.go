package avltree

import (
	"bytes"
	"errors"

	storage "../storage"
	byteconv "../utils"
)

type Tree struct {
	root        *Node
	dbf         *storage.DbFiles
	cachedNodes caching
}

func (t *Tree) Get(key []byte) ([]byte, error) {
	r := t.has(key)
	if r == nil {
		return nil, errors.New("Error: Key not found")
	}
	if bytes.Compare(key, r.key) == 0 {
		if r.del {
			return nil, errors.New("Error: Key not found")
		}
		return r.getValue(t.dbf, t.cachedNodes), nil
	}
	return nil, nil
}

func (t *Tree) Set(key []byte, value []byte) error {
	if t.root == nil {
		t.root = makeNode(t.dbf, key, value, nil)
		t.dbf.WriteRootAddress(byteconv.Int64ToByte(t.root.adrs.nadr))
		t.dbf.Commit()
		t.cachedNodes.add(t.root)
		return nil
	}
	r := t.has(key)
	cmp := bytes.Compare(key, r.key)
	if cmp == 0 {
		if r.del {
			r.del = false
		}
		r.value = value
		r.adrs.vlen = len(r.value)
		r.adrs.vadr, _ = t.dbf.WriteData(r.value)
		t.dbf.WriteTreeElem(r.adrs.nadr, r.toBytes())
		t.dbf.Commit()
		t.cachedNodes.add(r)
		return nil
	} else {
		added := makeNode(t.dbf, key, value, r)
		if cmp < 0 {
			r.left = added
		} else {
			r.right = added
		}
		t.dbf.WriteTreeElem(r.adrs.nadr, r.toBytes())
		r.balance(t.dbf)
		for ; t.root.up != nil; t.root = t.root.up {
		}
		t.dbf.Commit()
		t.cachedNodes.add(added)
		t.dbf.WriteRootAddress(byteconv.Int64ToByte(t.root.adrs.nadr))
	}
	return nil
}

func (t *Tree) Del(key []byte) error {
	n := t.has(key)
	if n == nil {
		return errors.New("Error: Key not found")
	}
	if bytes.Compare(key, n.key) == 0 {
		if n.del {
			return errors.New("Error: Key not found")
		} else {
			n.del = true
			t.dbf.WriteTreeElem(n.adrs.nadr, n.toBytes())
			t.dbf.Commit()
			t.cachedNodes.remove(n)
			return nil
		}
	}
	return errors.New("Error: Key not found")
}

func (t *Tree) has(key []byte) *Node {
	for n := t.root; n != nil; {
		res := bytes.Compare(key, n.key)
		if res == 0 {
			return n
		} else if res == -1 {
			if n.left == nil {
				return n
			}
			n = n.left
		} else {
			if n.right == nil {
				return n
			}
			n = n.right
		}
	}
	return nil
}

func Open(address, name string, cacheCount int, rewrite bool) *Tree {
	//create cache
	var cachedNodes caching
	if cacheCount > 0 {
		cachedNodes = createQueueCache(cacheCount)
	} else if cacheCount == 0 {
		cachedNodes = &emptyCache{}
	} else {
		cachedNodes = &fullCache{}
	}
	//open
	tree := &Tree{cachedNodes: cachedNodes, dbf: nil}
	if rewrite {
		dbf, _ := storage.CreateAndOpen(address, name)
		tree.dbf = dbf
		tree.root = nil
		return tree
	}
	dbf, err := storage.Open(address, name)
	if err != nil {
		dbf, _ = storage.CreateAndOpen(address, name)
		tree.dbf = dbf
		tree.root = nil
	} else {
		tree.dbf = dbf
		data, _ := dbf.ReadFullTree()
		rootAddress := byteconv.ByteToInt64(data[:8])
		if rootAddress == int64(0) {
			tree.root = nil
			return tree
		}
		tree.root = tree.load(rootAddress, data[:], nil)
	}
	return tree
}

func (t *Tree) load(nodeAdr int64, data []byte, up *Node) *Node {
	adrs := &address{}
	length := byteconv.ByteToInt(data[nodeAdr : nodeAdr+4])
	adrs.nadr = nodeAdr
	adrs.nlen = length
	nodeData := data[nodeAdr : nodeAdr+int64(length)]
	node := &Node{up: up}
	del := false
	if nodeData[4] == 1 {
		del = true
	}
	node.del = del
	node.lheight = byteconv.ByteToInt(nodeData[29:33])
	node.rheight = byteconv.ByteToInt(nodeData[33:37])
	adrs.vadr = byteconv.ByteToInt64(nodeData[37:45])
	adrs.vlen = byteconv.ByteToInt(nodeData[45:49])
	node.adrs = adrs
	node.key = nodeData[49:]
	if left := byteconv.ByteToInt64(nodeData[5:13]); left != int64(0) {
		node.left = t.load(left, data, node)
	}
	if right := byteconv.ByteToInt64(nodeData[13:21]); right != int64(0) {
		node.right = t.load(right, data, node)
	}
	return node
}

func (t *Tree) Close() error {
	t.root = nil
	return t.dbf.Close()
}
