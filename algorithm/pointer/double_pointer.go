package pointer

import "sort"

// 283. 移动零
// 输入: nums = [0,1,0,3,12]
// 输出: [1,3,12,0,0]
func moveZeroes(nums []int) {
	if len(nums) <= 1 {
		return
	}

	// 此处必须从0开始，否则个别情况会更改原始数据位置
	l, r := 0, 0
	for r < len(nums) {
		if nums[r] != 0 {
			// 交换数据
			nums[l], nums[r] = nums[r], nums[l]
			// 左指针只在出现数据变动时才移动
			l++
		}
		// 右指针一直移动
		r++
	}
}

// 11. 盛最多水的容器
// 输入：[1,8,6,2,5,4,8,3,7]
// 输出：49
func maxArea(height []int) int {
	l, r := 0, len(height)-1

	var ans int
	for l < r {
		ans = max(ans, min(height[l], height[r])*(r-l))
		if height[l] < height[r] {
			l++
		} else {
			r--
		}
	}
	return ans
}

// [X] 15. 三数之和
// 输入：nums = [-1,0,1,2,-1,-4]
// 输出：[[-1,-1,2],[-1,0,1]]
// 解释：
//
//	nums[0] + nums[1] + nums[2] = (-1) + 0 + 1 = 0 。
//	nums[1] + nums[2] + nums[4] = 0 + 1 + (-1) = 0 。
//	nums[0] + nums[3] + nums[4] = (-1) + 2 + (-1) = 0 。
//
// 不同的三元组是 [-1,0,1] 和 [-1,-1,2] 。
// 注意，输出的顺序和三元组的顺序并不重要。
func threeSum(nums []int) [][]int {
	// 将输入结果进行排序，得到一个递增切片
	sort.Ints(nums)

	// 定义结果集
	var ans [][]int
	cnt := len(nums)

	// 遍历次数为 n-2（l&r占用的位置）
	for i := 0; i < cnt-2; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		// 当遍历的最小值都大于0时就没有必要再遍历了
		if nums[i] > 0 {
			break
		}

		r := cnt - 1
		target := -nums[i]
		for l := i + 1; l < cnt; l++ {
			if l > i+1 && nums[l] == nums[l-1] {
				continue
			}

			if l == r {
				break
			}

			// nval := nums[l] + nums[r]
			if nums[l]+nums[r] > target {
				r--
			}

			// 这里切记一定不要用中间变量 nval，否则会出问题
			if nums[l]+nums[r] == target {
				ans = append(ans, []int{nums[i], nums[l], nums[r]})
			}
		}
	}
	return ans
}

// 公共函数
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
