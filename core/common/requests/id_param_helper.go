package requests

import (
	"strconv"
)

func StringToUInt(str string) (uint, error) {
	strInt, err := strconv.Atoi(str)
	if err != nil || strInt <= 0 {
		return 0, err
	}
	return uint(strInt), nil
}
