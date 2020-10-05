package main

import (
	"fmt"

	"github.com/didikprabowo/dipra"
)

type Personal struct {
	Name    string `is_valid:"required"`
	Address string `is_valid:"required"`
}

func main() {
	ps := dipra.NewValidator()
	// p := Personal{
	// 	Name:    "s",
	// 	Address: "",
	// }
	// err := ps.Validate(p)
	// if err != nil {
	// 	fmt.Println("error", err)
	// }

	fmt.Println("Validate number", ps.IsNumeric("0s"))

	fmt.Println("VAlidate email", ps.IsEmail("didik@gmail.com"))
	fmt.Println(ps.ValidateType("email"))
	// ps.
}
