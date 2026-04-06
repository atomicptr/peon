package err

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type Error struct {
	Code    Code
	Message string
	File    string
	Line    int
	Err     error
}

func (e Error) Error() string {
	err := e.GetMessage()

	if e.Err != nil {
		return fmt.Sprintf("%s: %v", err, e.Err)
	}

	return err
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) GetMessage() string {
	loc := e.LocationShort()

	if loc != "" {
		loc = fmt.Sprintf(" (%s)", loc)
	}

	return fmt.Sprintf("E%d %s%s", e.Code, e.Message, loc)
}

func (e Error) LocationShort() string {
	if e.File == "" || e.Line == -1 {
		return ""
	}

	path := filepath.ToSlash(e.File)
	const marker = "/pkg/"

	_, after, ok := strings.Cut(path, marker)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%s:%d", after, e.Line)
}

func New(code Code, message string, err error) Error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = ""
		line = -1
	}

	return Error{
		Code:    code,
		Message: message,
		Err:     err,
		File:    file,
		Line:    line,
	}
}

func (c Code) Is(err error) bool {
	if e, ok := errors.AsType[Error](err); ok {
		return e.Code == c
	}

	return false
}
