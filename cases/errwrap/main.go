package main

import (
	"errors"
	"fmt"
)

// Error wrapping/inspection — sentinel errors, %w chains, errors.As.
var ErrNotFound = errors.New("not found")

type AppError struct {
	Code int
	Msg  string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("app error %d: %s", e.Code, e.Msg)
}

func lookup(id int) error {
	if id < 0 {
		return &AppError{Code: 400, Msg: "negative id"}
	}
	return fmt.Errorf("lookup %d: %w", id, ErrNotFound)
}

func main() {
	err := lookup(7)
	fmt.Println("err:", err)
	fmt.Println("is NotFound:", errors.Is(err, ErrNotFound))
	fmt.Println("is other:", errors.Is(err, errors.New("nope")))
	fmt.Println("unwrap:", errors.Unwrap(err))

	var aerr *AppError
	if errors.As(lookup(-1), &aerr) {
		fmt.Println("as code:", aerr.Code)
	}
	if errors.As(err, &aerr) {
		fmt.Println("unexpected as match")
	}

	j := errors.Join(errors.New("one"), nil, fmt.Errorf("two: %w", ErrNotFound))
	fmt.Println("join:", j)
	fmt.Println("join is NotFound:", errors.Is(j, ErrNotFound))
}
