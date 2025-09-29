package errors

import (
	"strconv"
)

func HttpCode(code int) int {
	codeStr := strconv.Itoa(code)
	if len(codeStr) < 3 {
		return code
	}
	firstThreeStr := codeStr[:3]
	firstThree, _ := strconv.Atoi(firstThreeStr)
	return firstThree
}
