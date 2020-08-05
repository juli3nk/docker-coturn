package main

import (
	"gopkg.in/go-playground/validator.v10"
)

func validatePort(fl validator.FieldLevel) bool {
	portNum := fl.Field().Int()

	if portNum > 65535 || portNum < 1 {
		return false
	}

	return true
}
