package cache

import "fmt"

func GetUserPattern(id uint) string {
	return fmt.Sprintf("users:%d", id)
}

func GetUsers() string {
	return fmt.Sprintf("users:*")
}
