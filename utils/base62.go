package utils

const Base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func EncodeBase62(num int) string {
	if num == 0 {
		return string(Base62Chars[0])
	}

	result := ""
	for num > 0 {
		remainder := num % 62
		result = string(Base62Chars[remainder]) + result
		num /= 62
	}
	return result
}

func DecodeBase62(str string) int {
	num := 0
	for _, char := range str {
		num *= 62
		num += indexOf(char)
	}
	return num
}

func indexOf(char rune) int {
	for i, c := range Base62Chars {
		if c == char {
			return i
		}
	}
	return -1 // Not found
}
