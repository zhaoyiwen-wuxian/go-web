package utils

import "regexp"

func ValidatePhoneNumber(phone string) bool {
	phoneRegex := `^(\+?\d{1,4}[\s-]?)?(\(?\d{1,3}\)?[\s-]?)?[\d\s\-]{7,15}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}

// 校验邮箱
func ValidateEmail(email string) bool {
	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
