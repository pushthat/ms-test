package validator

import (
	"gopkg.in/go-playground/validator.v9"
)

// NodeRequestValidator is a Validator compatible struct used to validate node request
type NodeRequestValidator struct {
	Validator *validator.Validate
}

// ContainerRequestValidator is a Validator compatible struct used to validate node request
type ContainerRequestValidator struct {
	Validator *validator.Validate
}

func (cv *ContainerRequestValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func (nv *NodeRequestValidator) Validate(i interface{}) error {
	return nv.Validator.Struct(i)
}
