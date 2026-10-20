package validate

type StringRule func(string) string
type IntRule func(int) string

func Required() StringRule {
	return func(value string) string {
		if value == "" {
			return "is required"
		}

		return ""
	}
}

func MinLength(length int) StringRule {
	return func(value string) string {
		if len([]rune(value)) < length {
			return "is too short"
		}

		return ""
	}
}

func MaxLength(length int) StringRule {
	return func(value string) string {
		if len([]rune(value)) > length {
			return "is too long"
		}

		return ""
	}
}

func OneOf(values ...string) StringRule {
	return func(value string) string {
		for _, allowed := range values {
			if value == allowed {
				return ""
			}
		}

		return "contains an invalid value"
	}
}

func Min(minimum int) IntRule {
	return func(value int) string {
		if value < minimum {
			return "is too small"
		}

		return ""
	}
}

func Max(maximum int) IntRule {
	return func(value int) string {
		if value > maximum {
			return "is too large"
		}

		return ""
	}
}
