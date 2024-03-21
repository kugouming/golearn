package main

type ListNode struct {
	Val  int
	Next *ListNode
}

type Node struct {
	Val    int
	Next   *ListNode
	Random *Node
}

// 160. 相交链表
/*
题目：
	给你两个单链表的头节点 headA 和 headB ，请你找出并返回两个单链表相交的起始节点。如果两个链表不存在相交节点，返回 null 。
思路：
*/
func getIntersectionNode(headA, headB *ListNode) *ListNode {
}

// 206. 反转链表
/*
给你单链表的头节点 head ，请你反转链表，并返回反转后的链表。
输入：head = [1,2,3,4,5]
输出：[5,4,3,2,1]
*/
func reverseList(head *ListNode) *ListNode {
	var pre *ListNode
	cur := head
	for cur != nil {
		// 断开节点
		next := cur.Next

		cur.Next = pre
		pre = cur

		// 移动遍历指针
		cur = next
	}
	return pre
}

// 234. 回文链表
// 给你一个单链表的头节点 head ，请你判断该链表是否为回文链表。如果是，返回 true ；否则，返回 false 。
func isPalindrome(head *ListNode) bool {
}

// 141. 环形链表
func hasCycle(head *ListNode) bool {
}

// 142. 环形链表 II
func detectCycle(head *ListNode) *ListNode {
}

// 21. 合并两个有序链表
func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	res := &ListNode{}
	cur := res
	fhead, shead := list1, list2
	for fhead != nil && shead != nil {
		if fhead.Val < shead.Val {
			cur.Next = fhead
			fhead = fhead.Next
		} else {
			cur.Next = shead
			shead = shead.Next
		}
		cur = cur.Next
	}
	if fhead != nil {
		cur.Next = fhead
	} else {
		cur.Next = shead
	}
	return res.Next
}

// 2. 两数相加
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
}

// 19. 删除链表的倒数第 N 个结点
func removeNthFromEnd(head *ListNode, n int) *ListNode {}

// 24. 两两交换链表中的节点
func swapPairs(head *ListNode) *ListNode {}

// 25. K 个一组翻转链表
func reverseKGroup(head *ListNode, k int) *ListNode {
}

// 138. 随机链表的复制
/*
给你一个长度为 n 的链表，每个节点包含一个额外增加的随机指针 random ，该指针可以指向链表中的任何节点或空节点。

构造这个链表的 深拷贝。 深拷贝应该正好由 n 个 全新 节点组成，其中每个新节点的值都设为其对应的原节点的值。新节点的 next 指针和 random 指针也都应指向复制链表中的新节点，并使原链表和复制链表中的这些指针能够表示相同的链表状态。复制链表中的指针都不应指向原链表中的节点 。

例如，如果原链表中有 X 和 Y 两个节点，其中 X.random --> Y 。那么在复制链表中对应的两个节点 x 和 y ，同样有 x.random --> y 。

返回复制链表的头节点。

用一个由 n 个节点组成的链表来表示输入/输出中的链表。每个节点用一个 [val, random_index] 表示：

val：一个表示 Node.val 的整数。
random_index：随机指针指向的节点索引（范围从 0 到 n-1）；如果不指向任何节点，则为  null 。
你的代码 只 接受原链表的头节点 head 作为传入参数。
*/
func copyRandomList(head *Node) *Node {
}

// 148. 排序链表
func sortList(head *ListNode) *ListNode {}

// 23. 合并 K 个升序链表
func mergeKLists(lists []*ListNode) *ListNode {
}

// 146. LRU 缓存
