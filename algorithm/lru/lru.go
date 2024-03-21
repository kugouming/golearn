package lru

/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */

type LinkNode struct {
	key  int       // 存储节点Key，后续删除时会用到
	val  int       // 节点的value
	pre  *LinkNode // 前一个节点
	next *LinkNode // 后一个节点
}

type LRUCache struct {
	capacity int               // 容量
	size     int               // 当前的数量
	hmap     map[int]*LinkNode // 记录key对应的节点
	head     *LinkNode         // 头节点指针
	tail     *LinkNode         // 尾节点指针
}

func Constructor(capacity int) LRUCache {
	l := LRUCache{
		capacity: capacity,
		size:     0,
		hmap:     make(map[int]*LinkNode, capacity),
		head:     &LinkNode{}, // 初始空头节点，减少边界处理
		tail:     &LinkNode{}, // 初始空尾节点，减少边界处理
	}

	// 需要建立好链接关系
	l.head.next = l.tail
	l.tail.pre = l.head

	return l
}

// 查询节点值，且需更新使用频率（移动节点到头节点）
func (this *LRUCache) Get(key int) int {
	if v, has := this.hmap[key]; has {
		this.moveToHead(v)
		return v.val
	}
	return -1
}

// 添加节点值，并更新内部映射关系
func (this *LRUCache) Put(key int, value int) {
	if node, has := this.hmap[key]; has {
		node.val = value
		this.moveToHead(node)
	} else {
		tnode := &LinkNode{
			key: key,
			val: value,
		}
		this.hmap[key] = tnode
		this.addToHead(tnode)
		this.size++
		if this.size > this.capacity {
			dnode := this.removeTail()
			delete(this.hmap, dnode.key)
			this.size--
		}
	}
}

// 添加到头节点（这里自己画图是最好理解的）
func (this *LRUCache) addToHead(node *LinkNode) {
	node.pre = this.head
	node.next = this.head.next

	this.head.next.pre = node
	this.head.next = node
}

// 移除节点（这里自己画图是最好理解的）
func (this *LRUCache) removeNode(node *LinkNode) {
	node.pre.next = node.next
	node.next.pre = node.pre
}

// 移动节点到头（这里自己画图是最好理解的）
func (this *LRUCache) moveToHead(node *LinkNode) {
	this.removeNode(node)
	this.addToHead(node)
}

// 移除尾部节点（需要返回节点值，用于清理map）
func (this *LRUCache) removeTail() *LinkNode {
	node := this.tail.pre
	this.removeNode(node)
	return node
}
