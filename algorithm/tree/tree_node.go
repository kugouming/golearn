package main

import "math"

type TreeNode struct {
	Val         int
	Left, Right *TreeNode
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// 104. 二叉树的最大深度
func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	lLen := maxDepth(root.Left)
	rLen := maxDepth(root.Right)

	return max(lLen, rLen) + 1
}

// 226. 翻转二叉树
func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}

	left := invertTree(root.Left)
	right := invertTree(root.Right)

	root.Left = right
	root.Right = left

	return root
}

// 101. 对称二叉树
func isSymmetric(root *TreeNode) bool {
	var checkDuiChen func(left, right *TreeNode) bool
	checkDuiChen = func(left, right *TreeNode) bool {
		// 边界条件处理
		if left == nil && right == nil {
			return true
		}

		if left == nil || right == nil {
			return false
		}

		return (left.Val == right.Val) && checkDuiChen(left.Left, right.Right) && checkDuiChen(left.Right, right.Left)
	}
	// 按照对称的方式进行遍历
	return checkDuiChen(root, root)
}

// 543. 二叉树的直径
func diameterOfBinaryTree(root *TreeNode) int {
	var ans int
	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		// 下面计算边的时候进行了+1，故当为nil时，需要将提前+1的值减掉
		if node == nil {
			return -1
		}

		// 计算左侧节点最长路径，并加上到当前节点的边
		left := dfs(node.Left) + 1
		// 计算右侧节点最长路径，并加上到当前节点的边
		right := dfs(node.Right) + 1

		// 计算当前节点的最长路径，并于哨兵值比较
		ans = max(ans, left+right)

		// 求得当前节点下的最长路径
		return max(left, right)
	}
	dfs(root)

	return ans
}

// 108. 将有序数组转换为二叉搜索树
func sortedArrayToBST(nums []int) *TreeNode {
	// 1. 高度平衡
	// 2. 遵循 node.Left < node.Val < node.Right

	var helper func(nums []int, l, r int) *TreeNode

	helper = func(nums []int, l, r int) *TreeNode {
		if len(nums) == 0 {
			return nil
		}

		// 计算当前节点的值
		mid := (l + r) / 2

		// 新建根节点
		node := &TreeNode{Val: nums[mid]}
		// 递归创建左节点
		node.Left = helper(nums, l, mid-1)
		// 递归创建右节点
		node.Right = helper(nums, mid+1, r)

		return node
	}

	return helper(nums, 0, len(nums)-1)
}

// 98. 验证二叉搜索树
/*
	有效 二叉搜索树定义如下：
		- 节点的左子树只包含 小于 当前节点的数。
		- 节点的右子树只包含 大于 当前节点的数。
		- 所有左子树和右子树自身必须也是二叉搜索树。
*/
func isValidBST(root *TreeNode) bool {
	// 左边的节点值最小，右边节点值最大，故可以考虑借用math中的最小/大值作为参考
	var helper func(node *TreeNode, min, max int) bool
	helper = func(node *TreeNode, min, max int) bool {
		if node == nil {
			return true
		}

		if node.Val <= min || node.Val >= max {
			return false
		}

		// 分节点判断各自的取值范围，均符合则有效
		return helper(node.Left, min, node.Val) && helper(node.Right, node.Val, max)
	}

	return helper(root, math.MinInt64, math.MaxInt64)
}

// 230. 二叉搜索树中第K小的元素
//
/* ************************
效果图如下：

        3
	  /  \
     /     \
    1      4
  /  \
 /	  \
nil	   2

// 原值：[3,1,4,null,2]
*************************/
//
// 前序：[3 1 2 4] - 根 左 右
// 中序：[1 2 3 4] - 左 根 右
// 后序：[2 1 4 3] - 左 右 根
func kthSmallest(root *TreeNode, k int) int {
	var ret = []int{}
	var helper func(node *TreeNode)
	helper = func(node *TreeNode) {
		if node == nil {
			return
		}

		ret = append(ret, node.Val)
		helper(node.Left)
		helper(node.Right)
	}
	helper(root)
	return ret[k-1]
}

// 199. 二叉树的右视图
func rightSideView(root *TreeNode) []int {
	ret := []int{}
	var helper func(node *TreeNode, depth int)
	helper = func(node *TreeNode, depth int) {
		if node == nil {
			return
		}

		// 这里需要注意下深度，由于当前的depth是在+1前的，所以下面可以直接用，无需做-1操作
		if len(ret) <= depth {
			ret = append(ret, node.Val)
		} else {
			ret[depth] = node.Val
		}

		// 这里是从做到右遍历，所以会存在上面重新赋值的情况
		helper(node.Left, depth+1)
		helper(node.Right, depth+1)
	}
	helper(root, 0)
	return ret
}

// 114. 二叉树展开为链表
func flatten(root *TreeNode) {
	var helper func(node *TreeNode) []*TreeNode

	// 前序遍历二叉树节点，并将节点收集到列表中
	helper = func(node *TreeNode) []*TreeNode {
		list := []*TreeNode{}
		if node == nil {
			return list
		}
		list = append(list, node)
		list = append(list, helper(node.Left)...)
		list = append(list, helper(node.Right)...)
		return list
	}
	lists := helper(root)

	// 按照顺序更改节点指针（由于将二叉树展开为链表之后会破坏二叉树的结构，因此在前序遍历结束之后更新每个节点的左右子节点的信息，将二叉树展开为单链表。）
	for i := 1; i < len(lists); i++ {
		pre, cur := lists[i-1], lists[i]
		pre.Left = nil
		pre.Right = cur
	}
}
