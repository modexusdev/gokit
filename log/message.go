package log

import "net/http"

type messageWriter interface {
	SetLogMessage(message string)
}

func Message(
	w http.ResponseWriter,
	message string,
) {
	if writer, ok := w.(messageWriter); ok {
		writer.SetLogMessage(message)
	}
}
