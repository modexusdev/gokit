package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeJSON(
	r *http.Request,
	dst any,
) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf(
			"decode json: %w",
			err,
		)
	}

	return nil
}

func NewJSONRequest(
	ctx context.Context,
	method string,
	url string,
	body any,
) (*http.Request, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf(
			"encode json: %w",
			err,
		)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create request: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	return req, nil
}
