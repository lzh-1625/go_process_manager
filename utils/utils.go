package utils

import "strconv"

func Unwarp[T any](result T, err error) T {
	if err != nil {
		panic(err)
	}
	return result
}

func UnwarpIgnore[T any](result T, _ error) T {
	return result
}

func GetIntByString(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	} else {
		return i
	}
}
