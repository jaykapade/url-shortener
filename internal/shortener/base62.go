package shortener

const base62Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func EncodeBase62(num uint64) string {
	if num == 0 {
		return string(base62Digits[0])
	}
	result := ""
	for num > 0 {
		remainder := num % 62
		result = string(base62Digits[remainder]) + result
		num = num / 62
	}
	return result
}
