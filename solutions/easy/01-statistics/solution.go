package lab

import "errors"

func Stats(nums []int) (min, max, sum int, mean float64, err error) {
	if len(nums) == 0 {
		return 0, 0, 0, 0, errors.New("empty input")
	}
	min, max = nums[0], nums[0]
	for _, n := range nums {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
		sum += n
	}
	mean = float64(sum) / float64(len(nums))
	return
}
