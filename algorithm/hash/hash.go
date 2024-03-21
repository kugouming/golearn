package hash

// 1. 两数之和
// 输入：nums = [2,7,11,15], target = 9
// 输出：[0,1]
func twoSum(nums []int, target int) []int {
	// hash 方法，需要确保给定的集合中不存在重复元素
	hFunc := func(nums []int, target int) []int {
		hmap := make(map[int]int, len(nums))
		for i, v := range nums {
			if j, has := hmap[target-v]; has {
				return []int{i, j}
			}
			hmap[v] = i
		}
		return nil
	}
	hFunc(nums, target)

	// 迭代方法
	iFunc := func(nums []int, target int) []int {
		llen := len(nums)
		for i := 0; i < llen; i++ {
			for j := llen; j > i; j-- {
				if nums[i]+nums[j] == target && i != j {
					return []int{i, j}
				}
			}
		}
		return nil
	}

	return iFunc(nums, target)
}

// 49. 字母异位词分组
// 输入: strs = ["eat", "tea", "tan", "ate", "nat", "bat"]
// 输出: [["bat"],["nat","tan"],["ate","eat","tea"]]
func groupAnagrams(strs []string) [][]string {
	hmap := map[[26]int][]string{}
	for _, v := range strs {
		tmp := [26]int{}
		for _, s := range v {
			// 需要借助字符串的ASCII值，转化为有限的数组
			tmp[s-'a']++
		}
		hmap[tmp] = append(hmap[tmp], v)
	}

	ans := make([][]string, 0, len(hmap))
	for _, val := range hmap {
		ans = append(ans, val)
	}
	return ans
}

// 128. 最长连续序列
// 输入：nums = [100,4,200,1,3,2]
// 输出：4
func longestConsecutive(nums []int) int {
	hmap := make(map[int]bool, len(nums))
	for _, v := range nums {
		hmap[v] = true
	}

	var ret int
	for _, v := range nums {
		times := 1
		for tmp := v; hmap[tmp-1]; tmp-- {
			times++
		}
		ret = max(ret, times)
	}
	return ret
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
