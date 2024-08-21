package utility

import (
	"errors"
	"slices"
	"strings"
)

type Error struct {
	errs []error
}

func (e Error) Error() string {
	s := new(strings.Builder)
	s.WriteString("multiple errors:")
	for _, err := range e.errs {
		s.WriteString("\n\t")
		s.WriteString(err.Error())
	}
	return s.String()
}

func (e Error) Errors() []error {
	return slices.Clone(e.errs)
}

func (e Error) Unwrap() []error {
	return e.errs
}

func NewMultiErr(errs ...error) error {
	var n int
	var errFirst error
	for _, e := range errs {
		switch e := e.(type) {
		case nil:
			continue
		case Error:
			n += len(e.errs)
			if errFirst == nil && len(e.errs) > 0 {
				errFirst = e.errs[0]
			}
		default:
			n++
			if errFirst == nil {
				errFirst = e
			}
		}
	}
	if n <= 1 {
		return errFirst // nil if n == 0
	}

	dst := make([]error, 0, n)
	for _, e := range errs {
		switch e := e.(type) {
		case nil:
			continue
		case Error:
			dst = append(dst, e.errs...)
		default:
			dst = append(dst, e)
		}
	}
	return Error{errs: dst}
}

func (e Error) Is(target error) bool {
	for _, err := range e.errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e Error) As(target any) bool {
	for _, err := range e.errs {
		if ok := errors.As(err, target); ok {
			return true
		}
	}
	return false
}

func Range(err error, fn func(error) bool) bool {
	if err == nil {
		return true
	}
	if !fn(err) {
		return false
	}
	switch err := err.(type) {
	case interface{ Unwrap() error }:
		if err := err.Unwrap(); err != nil {
			if !Range(err, fn) {
				return false
			}
		}
	case interface{ Unwrap() []error }:
		for _, err := range err.Unwrap() {
			if !Range(err, fn) {
				return false
			}
		}
	}
	return true
}
