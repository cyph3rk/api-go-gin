package models

import (
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

type Aluno struct {
	gorm.Model
	Nome string `kson:"nome" validate:"nonzero"`
	CPF  string `kson:"cpf" validate:"len=11, regexp=^[0-9]*$"`
	RG   string `kson:"rg" validate:"len=9, regexp=^[0-9]*$"`
}

func ValidaDadosDeAluno(aluno *Aluno) error {
	if err := validator.Validate(aluno); err != nil {
		return err
	}
	return nil
}
