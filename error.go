package main

import "fmt"

type ErrUnknownControlName struct {
    controlName string
}

func (e *ErrUnknownControlName) Error() string {
    return fmt.Sprintf("There is no builder associated with the control name/type of '%v'", e.controlName)
}