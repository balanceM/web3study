package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

func main() {
	// arr := [...]int{1, 2, 3, 4, 4, 3, 1}
	// value, finded := singleNumber(arr)
	// if finded {
	// 	fmt.Printf("%d 只出现了一次", value)
	// } else {
	// 	fmt.Println("没有任何元素只出现一次")
	// }

	// str := "{a(sdf)]a}"
	// fmt.Println(isValid(str))

	// strs := []string{"flower", "flow", "aflight"}
	// fmt.Println(longestCommonPrefix(strs))

	// digits := []int{4, 3, 2, 1}
	// fmt.Println(plusOne(digits))

	// arr := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	// fmt.Println(removeDuplicates(arr))

	// arr := []int{2, 7, 11, 15}
	// fmt.Println(twoSum(arr, 14))

	input := [][2]int{[2]int{3, 5}, [2]int{1, 2}, [2]int{2, 4}}
	fmt.Println(mergeIntervals(input))
}

// 只出现一次的数字
func singleNumber(nums []int) int {
	counter := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		counter[nums[i]]++
	}
	for k, v := range counter {
		if v == 1 {
			return k
		}
	}
	return -1
}

// 字符串是否有效
func isValid(s string) bool {
	rune_map := map[rune]rune{'}': '{', ')': '(', ']': '['}
	count_map := map[rune]int{'{': 0, '(': 0, '[': 0}
	for _, v := range s {
		switch v {
		case '(', '[', '{':
			count_map[v]++
		case ')', ']', '}':
			if count_map[rune_map[v]] == 0 {
				return false
			}
			count_map[rune_map[v]]--
		}
	}
	for _, count := range count_map {
		if count > 0 {
			return false
		}
	}
	return true
}

// 最长公共前缀
func longestCommonPrefix(strs []string) string {
	longestPre := []rune(strs[0])
	for _, str := range strs[1:] {
		for i, r := range str {
			if i > len(longestPre)-1 {
				break
			}
			if longestPre[i] != r {
				longestPre = longestPre[0:i]
			}
		}
	}
	return string(longestPre)
}

// 加一
func plusOne(digits []int) []int {
	length := len(digits)
	total := 0
	for i, v := range digits {
		total += v * int(math.Pow(10, float64(length-1-i)))
	}
	total++
	res := []int{}
	for _, v := range strconv.FormatInt(int64(total), 10) {
		res = append(res, int(v-'0'))
	}
	return res
}

// 删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	baseValue := nums[0]
	changeIndex := 1
	for _, v := range nums {
		if v != baseValue {
			nums[changeIndex] = v
			baseValue = v
			changeIndex++
		}
	}
	return len(nums[:changeIndex])
}

// 合并区间
func mergeIntervals(intervals [][2]int) [][2]int {
	// 按starti排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	// 合并
	newIntervals := [][2]int{intervals[0]}
	for _, curr := range intervals[1:] {
		lastInterval := &newIntervals[len(newIntervals)-1]
		if curr[0] <= lastInterval[1] {
			if curr[1] > lastInterval[1] {
				lastInterval[1] = curr[1]
			}
		} else {
			newIntervals = append(newIntervals, curr)
		}
	}
	//
	return newIntervals
}

// 两数之和
func twoSum(nums []int, target int) []int {
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{nums[i], nums[j]}
			}
		}
	}
	return []int{}
}
