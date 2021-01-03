package avltree

import (
	"bytes"

	storage "../storage"
	byteconv "../utils"
)

type Node struct {
	key     []byte
	value   []byte
	left    *Node
	right   *Node
	up      *Node
	lheight int
	rheight int
	adrs    *address
	del     bool
}

type address struct {
	nadr int64
	vadr int64
	nlen int
	vlen int
}

func (n *Node) getValue(dbf *storage.DbFiles, cachedNodes caching) []byte {
	if n.value == nil {
		n.value, _ = dbf.ReadData(n.adrs.vadr, n.adrs.vlen)
		defer cachedNodes.add(n)
	}
	return n.value
}

func (n *Node) toBytes() []byte {
	data := make([]byte, 0)
	leftAdr := int64(0)
	rightAdr := int64(0)
	upAdr := int64(0)
	if n.left != nil {
		leftAdr = n.left.adrs.nadr
	}
	if n.right != nil {
		rightAdr = n.right.adrs.nadr
	}
	if n.up != nil {
		upAdr = n.up.adrs.nadr
	}
	del := byte(0)
	if n.del {
		del = byte(1)
	}
	data = append(data, del)
	data = append(data, byteconv.Int64ToByte(leftAdr)...)
	data = append(data, byteconv.Int64ToByte(rightAdr)...)
	data = append(data, byteconv.Int64ToByte(upAdr)...)
	data = append(data, byteconv.IntToByte(n.lheight)...)
	data = append(data, byteconv.IntToByte(n.rheight)...)
	data = append(data, byteconv.Int64ToByte(n.adrs.vadr)...)
	data = append(data, byteconv.IntToByte(n.adrs.vlen)...)
	data = append(data, n.key...)
	data = append(byteconv.IntToByte(4+len(data)), data...)
	return data
}

func makeNode(dbf *storage.DbFiles, key, value []byte, up *Node) *Node {
	node := &Node{key: key, value: value, lheight: 0, rheight: 0, left: nil, right: nil, up: up, del: false, adrs: &address{}}
	node.adrs.vlen = len(value)
	node.adrs.nlen = len(node.toBytes())
	vadr, _ := dbf.WriteData(node.value)
	node.adrs.vadr = vadr
	nadr, _ := dbf.WriteTreeElem(-1, node.toBytes())
	node.adrs.nadr = nadr
	return node
}

func (n *Node) heights() {
	if n.left != nil {
		n.lheight = max(n.left.lheight, n.left.rheight) + 1
	} else {
		n.lheight = 0
	}
	if n.right != nil {
		n.rheight = max(n.right.lheight, n.right.rheight) + 1
	} else {
		n.rheight = 0
	}
}

func (n *Node) balance(dbf *storage.DbFiles) {
	n.heights()
	if n.lheight-n.rheight > 1 {
		n.leftBalance(dbf)
	} else if n.lheight-n.rheight < (-1) {
		n.rightBalance(dbf)
	}
	if n.up != nil {
		n.up.balance(dbf)
	}
}

func (n *Node) leftBalance(dbf *storage.DbFiles) {
	if n.left.rheight > n.left.lheight {
		n.left.rightBalance(dbf)
	}
	if n.up != nil {
		if n.up.left != nil {
			if bytes.Compare(n.key, n.up.left.key) == 0 {
				n.up.left = n.left
			} else {
				n.up.right = n.left
			}
		} else {
			n.up.right = n.left
		}
		//write n.up
		dbf.WriteTreeElem(n.up.adrs.nadr, n.up.toBytes())
	}
	n.left.up = n.up
	n.up = n.left
	if n.left.right != nil {
		n.left.right.up = n
	}
	n.left = n.left.right
	n.up.right = n
	n.heights()
	n.up.heights()
	if n.left != nil {
		//write n.left
		dbf.WriteTreeElem(n.left.adrs.nadr, n.left.toBytes())
	}
	if n.up != nil {
		//write n.up
		dbf.WriteTreeElem(n.up.adrs.nadr, n.up.toBytes())
	}
	//write n
	dbf.WriteTreeElem(n.adrs.nadr, n.toBytes())
}

func (n *Node) rightBalance(dbf *storage.DbFiles) {
	if n.right.lheight > n.right.rheight {
		n.right.leftBalance(dbf)
	}
	if n.up != nil {
		if n.up.left != nil {
			if bytes.Compare(n.key, n.up.left.key) == 0 {
				n.up.left = n.right
			} else {
				n.up.right = n.right
			}
		} else {
			n.up.right = n.right
		}
		//write n.up
		dbf.WriteTreeElem(n.up.adrs.nadr, n.up.toBytes())
	}
	n.right.up = n.up
	n.up = n.right
	if n.right.left != nil {
		n.right.left.up = n
	}
	n.right = n.right.left
	n.up.left = n
	n.heights()
	n.up.heights()
	if n.right != nil {
		//write n.right
		dbf.WriteTreeElem(n.right.adrs.nadr, n.right.toBytes())
	}
	if n.up != nil {
		//write n.up
		dbf.WriteTreeElem(n.up.adrs.nadr, n.up.toBytes())
	}
	//write n
	dbf.WriteTreeElem(n.adrs.nadr, n.toBytes())
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
