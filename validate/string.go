package validate

func String(
	field string,
	value string,
	rules ...StringRule,
) error {
	for _, rule := range rules {
		if message := rule(value); message != "" {
			return newError(field, message)
		}
	}

	return nil
}
