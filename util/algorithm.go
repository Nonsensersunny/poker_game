package util

import (
	"fmt"
	"math"
)

func SplitStep(max, num int) []int {
	var (
		result []int
		step   = max / num
	)

	for i := 1; i <= num; i++ {
		result = append(result, step*i)
	}
	return result
}

func Extremes(nums []int) (max, min int) {
	max, min = math.MinInt32, math.MaxInt32
	for _, v := range nums {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return
}

func IsContinuous(nums []int) bool {
	for i := 0; i < len(nums)-1; i++ {
		if nums[i]-nums[i+1] != -1 {
			return false
		}
	}
	return true
}

func GenUserGameKey(name string) string {
	return fmt.Sprintf("%v_game", name)
}

func GenGameRemainKey(game string) string {
	return fmt.Sprintf("%v_remain", game)
}

func GenGameRecordKey(game string) string {
	return fmt.Sprintf("%v_record", game)
}

func SliceContains(s []string, t string) bool {
	for _, v := range s {
		if v == t {
			return true
		}
	}
	return false
}
