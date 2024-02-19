package main

import (
	"errors"
	"fmt"
)

func main() {
	err := B()
	fmt.Println(errors.Is(err, ErrNotAllowed))
	fmt.Println(errors.Is(err, ErrNotFound))
	fmt.Println(err.Error())
}

var ErrNotFound = errors.New("not found")
var ErrNotAllowed = errors.New("not allowed")
var ErrInternal = errors.New("internal error")

func A() error {
	return ErrNotFound
}

func C() error {
	return ErrNotAllowed
}

func B() error {
	err := A()
	if err != nil {
		err = fmt.Errorf("b: %w", err)
	}

	errC := C()
	if errC != nil {
		return errors.Join(err, errC)
	}
	return nil
}
