package models

import "gorm.io/gorm"

type Aluno struct {
	gorm.Model
	Nome string `kson:"nome"`
	CPF  string `kson:"cpf"`
	RG   string `kson:"rg"`
}
