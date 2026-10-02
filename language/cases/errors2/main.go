package main

import (
	"errors"
	"fmt"
)

type MyErr struct{ Code int }
func (e *MyErr) Error() string { return fmt.Sprintf("myerr%d", e.Code) }

func wrap() error { return fmt.Errorf("outer: %w", &MyErr{Code: 7}) }
func joinErrs() error {
	return errors.Join(errors.New("a"), errors.New("b"))
}

func main() {
	err := wrap()
	fmt.Println(err.Error())
	u := errors.Unwrap(err)
	fmt.Println(u)
	fmt.Printf("T %T\n", u)
	var me *MyErr
	if errors.As(err, &me) {
		fmt.Println("as", me.Code)
	}
	e1 := errors.New("x")
	fmt.Println(errors.Is(e1, e1))
	fmt.Println(joinErrs().Error())
	// Unwrap returning chain
	fmt.Println(errors.Is(fmt.Errorf("w: %w", e1), e1))
}
