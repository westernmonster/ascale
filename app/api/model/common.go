package model

import validation "github.com/go-ozzo/ozzo-validation"

type ArgJob struct {
	Job string `json:"job"`
}

func (p *ArgJob) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.Job, validation.Required),
	)
}
