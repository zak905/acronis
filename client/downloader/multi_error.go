package downloader

import (
	"strings"
)

type multiError struct {
	errors_ map[string]error
}

func newMultiError() *multiError {
	return &multiError{errors_: make(map[string]error)}
}

func (e *multiError) add(filename string, err error) {
	e.errors_[filename] = err
}

func (e *multiError) hasErrors() bool {
	return len(e.errors_) > 1
}

func (e *multiError) Error() string {
	builder := strings.Builder{}
	for filename, err := range e.errors_ {
		builder.WriteString("the following errors occured for file: " + filename + "\n")
		builder.WriteString(err.Error())
	}

	return builder.String()
}
