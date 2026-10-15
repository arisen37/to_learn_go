package lab

func RotateLeft(nums []int, k int) []int {
	n := len(nums)
	out := make([]int, n)
	if n == 0 {
		return out
	}
	k %= n
	if k < 0 {
		k += n
	}
	copy(out, nums[k:])
	copy(out[n-k:], nums[:k])
	return out
}
