//////////
//单链表 -- 线性表
package datas

//定义节点
type Node struct {
	Data int
	Next *Node
}

/*
* 返回第一个节点
* h 头结点
 */
func GetFirst(h *Node) *Node {
	if h.Next == nil {
		return nil
	}
	return h.Next
}

/*
* 返回最后一个节点
* h 头结点
 */
func GetLast(h *Node) *Node {
	if h.Next == nil {
		return nil
	}
	i := h
	for i.Next != nil {
		i = i.Next
		if i.Next == nil {
			return i
		}
	}
	return nil
}

//取长度
func GetLength(h *Node) int {
	var i int = 0
	n := h
	for n.Next != nil {
		i++
		n = n.Next
	}
	return i
}

//插入一个节点
//h: 头结点
//d:要插入的节点
//p:要插入的位置
func Insert(h, d *Node, p int) bool {
	if h.Next == nil {
		h.Next = d
		return true
	}
	i := 0
	n := h
	for n.Next != nil {
		i++
		if i == p {
			if n.Next.Next == nil {
				n.Next = d
				return true
			} else {
				d.Next = n.Next
				n.Next = d.Next
				return true
			}
		}
		n = n.Next
		if n.Next == nil {
			n.Next = d
			return true
		}
	}
	return false
}

//取出指定节点
func GetLoc(h *Node, p int) *Node {
	if p < 0 || p > GetLength(h) {
		return nil
	}
	var i int = 0
	n := h
	for n.Next != nil {
		i++
		n = n.Next
		if i == p {
			return n
		}
	}
	return nil
}
