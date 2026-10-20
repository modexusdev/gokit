package validate

func Int(
	field string,
	value int,
	rules ...IntRule,
) error {
	for _, rule := range rules {
		if message := rule(value); message != "" {
			return newError(field, message)
		}
	}

	return nil
}
