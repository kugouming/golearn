package substr

// 53. 最大子数组和
// 输入：nums = [-2,1,-3,4,-1,2,1,-5,4]
// 输出：6
// 解释：连续子数组 [4,-1,2,1] 的和最大，为 6 。
/*
题目：
	给你一个整数数组 nums ，请你找出一个具有最大和的连续子数组（子数组最少包含一个元素），返回其最大和。
	子数组是数组中的一个连续部分。

思路：
	1、定义一个中间变量，用来记录最大和
	2、求最大和即为算出每个位置在前面累加过程中可能出现的最大值
	3、当前位置与前一个位置的累积和相加比较，即可得到当前位置的最大值
	4、遍历过程中需要时刻更新中间变量的最大值
	5、输入遍历完即可得到结果
*/
func maxSubArray(nums []int) int {
}

// 56. 合并区间
// 输入：intervals = [[1,3],[2,6],[8,10],[15,18]]
// 输出：[[1,6],[8,10],[15,18]]
// 解释：区间 [1,3] 和 [2,6] 重叠, 将它们合并为 [1,6].
/*
题目：
	以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。请你合并所有重叠的区间，并返回 一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间 。

思路：
	1、定义一个存储结果的变量
	2、将无序集合转换为有序集合 `sort.Slice(xx, func(x,y int) bool {return xx[x][0] < xx[y][0]})`
	3、定义一个变量，用来存储前面处理后的中间结果
	4、从0位置开始遍历集合，第0位置直接将结果赋值给中间变量，就跳过
	5、从非第一次遍历，然后判断当前集合中的值与中间变量中的值是否存在交集，存在交集则重建临时变量（merge当前集合），无重合时则将中间变量追加到结果集，然后将当前集合赋值给中间变量
	6、循环结束后，判断临时变量是否有值，有值则追加到结果集中
	7、返回结果
*/
func merge(intervals [][]int) [][]int {
}

// 189. 轮转数组
// 输入: nums = [1,2,3,4,5,6,7], k = 3
// 输出: [5,6,7,1,2,3,4]
// 解释:
// 向右轮转 1 步: [7,1,2,3,4,5,6]
// 向右轮转 2 步: [6,7,1,2,3,4,5]
// 向右轮转 3 步: [5,6,7,1,2,3,4]
/*
题目：
	给定一个整数数组 nums，将数组中的元素向右轮转 k 个位置，其中 k 是非负数。

思路：(还可以根据翻转后的数组索引重建数组)
	1、【重点】k可能大于数组的长度，这里可以先处理下k（k%len(nums)）
	2、按照位置将数组分为两部分，然后将值分别追加到新数组中
	3、操作完毕之后需要将新数组值赋值给nums上（`copy(dst, src)`）
*/
func rotate(nums []int, k int) {
	cnt := len(nums)
	if cnt <= 1 {
		return
	}
	k = k % cnt

	nnums := []int{}
	nnums = append(nnums, nums[cnt-k:]...)
	nnums = append(nnums, nums[:cnt-k]...)
	copy(nums, nnums)
}

// 238. 除自身以外数组的乘积
// 输入: nums = [1,2,3,4]
// 输出: [24,12,8,6]
/*
题目：
	给你一个整数数组 nums，返回 数组 answer ，其中 answer[i] 等于 nums 中除 nums[i] 之外其余各元素的乘积 。
	题目数据 保证 数组 nums之中任意元素的全部前缀元素和后缀的乘积都在  32 位 整数范围内。
	请 不要使用除法，且在 O(n) 时间复杂度内完成此题。

思路：
	1、【重点】将所在位置的乘积看作是当前节点的左右两部分的乘积（第一个位置的左侧默认为1）
	2、定义左侧部分map，其key为数组索引，值为左侧部分的索引乘积
	3、定义右侧部分map，其key为数组索引，值为右侧部分的索引乘积
	4、遍历数组，按照数组索引位置分别取出左右两边的乘积，然后再做乘积，得到最后结果。
*/
func productExceptSelf(nums []int) []int {
}

// 41. 缺失的第一个正数
// 给你一个未排序的整数数组 nums ，请你找出其中没有出现的最小的正整数。
//
// 请你实现时间复杂度为 O(n) 并且只使用常数级别额外空间的解决方案。
//
// 输入：nums = [1,2,0]
// 输出：3
// 解释：范围 [1,2] 中的数字都在数组中。
/*
题目：
	给你一个未排序的整数数组 nums ，请你找出其中没有出现的最小的正整数。
	请你实现时间复杂度为 O(n) 并且只使用常数级别额外空间的解决方案。

要求：
    1、最小正整数为1
    2、数组中出现的值应为索引+1，其他值不能出现
    3、判断里面没有出现的最小正整数
    4、若数组中值都满足，则返回数组长度+1的值

思路：
	1、
*/
func firstMissingPositive(nums []int) int {
}
