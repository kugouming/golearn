package main

import (
	"fmt"
)

/*
在上述代码中，我们定义了QuickSort函数用于执行快速排序的递归操作，它接收一个整数切片arr，以及需要排序部分的左右索引left和right。在QuickSort函数中，我们首先检查left是否小于right，如果是，则调用partition函数进行划分，并得到基准值的索引pivotIndex。接着递归地对基准值左边的子序列和右边的子序列进行快速排序。

partition函数负责根据基准值对数组进行划分。在这个实现中，我们选择最右边的元素作为基准值，然后将所有小于基准值的元素移动到基准值的左边，所有大于或等于基准值的元素移动到基准值的右边。最后，将基准值放到正确的位置，并返回其索引。

在main函数中，我们创建了一个需要排序的整数切片arr，然后调用QuickSort函数对其进行排序，并打印出排序前后的数组。

请注意，快速排序的性能与所选的基准值密切相关。在实际应用中，可能会使用更复杂的策略来选择基准值，例如随机选择或者选择中位数作为基准值，以减少最坏情况发生的概率。此外，这个示例中使用了简单的交换来移动元素，但在处理大型数据集时，可能需要使用更高效的内存管理策略。
*/

// QuickSort 实现快速排序
func QuickSort(arr []int, left, right int) {
	if left < right {
		pivotIndex := partition(arr, left, right)
		QuickSort(arr, left, pivotIndex-1)
		QuickSort(arr, pivotIndex+1, right)
	}
}

// partition 划分函数，返回基准值的索引
func partition(arr []int, left, right int) int {
	pivot := arr[right] // 选择最右边的元素作为基准值
	i := left
	for j := left; j < right; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i] // 交换元素
			i++
		}
	}
	arr[i], arr[right] = arr[right], arr[i] // 将基准值放到正确的位置
	return i
}

func main() {
	arr := []int{64, 25, 12, 22, 11}
	fmt.Println("原始数组:", arr)

	QuickSort(arr, 0, len(arr)-1)

	fmt.Println("排序后的数组:", arr)
}
