package client

import (
	"gopkg.in/validator.v2"
)

type NamespaceStat struct {
	NrObjects           int64   `json:"nrObjects" validate:"nonzero"`
	ReadRequestPerHour  int64   `json:"readRequestPerHour" validate:"nonzero"`
	SpaceAvailable      float64 `json:"spaceAvailable,omitempty"`
	SpaceUsed           float64 `json:"spaceUsed,omitempty"`
	WriteRequestPerHour int64   `json:"writeRequestPerHour" validate:"nonzero"`
}

func (s NamespaceStat) Validate() error {

	return validator.Validate(s)
}