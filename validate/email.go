package validate

import "net/mail"

func Email(
	field string,
	value string,
	rules ...StringRule,
) error {
	if err := String(field, value, rules...); err != nil {
		return err
	}

	if value == "" {
		return nil
	}

	address, err := mail.ParseAddress(value)
	if err != nil || address.Address != value {
		return newError(
			field,
			"must be a valid email address",
		)
	}

	return nil
}
