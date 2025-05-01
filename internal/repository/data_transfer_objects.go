package repository

type CreatePersonDTO struct {
	Name       string `json:"name" validate:"required,alpha" example:"John"`
	Surname    string `json:"surname" validate:"required,alpha" example:"Smith"`
	Patronymic string `json:"patronymic" validate:"omitempty,alpha" example:"William"`
}

type FullPersonDTO struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" validate:"required,alpha" example:"John"`
	Surname     string `json:"surname" validate:"required,alpha" example:"Smith"`
	Patronymic  string `json:"patronymic" validate:"omitempty,alpha" example:"William"`
	Age         uint   `json:"age" validate:"omitempty,gte=0" example:"30"`
	Gender      string `json:"gender" validate:"omitempty,alpha" example:"male"`
	Nationality string `json:"nationality" validate:"omitempty,alpha" example:"US"`
}

type PersonParams struct {
	Name        string `form:"name" example:"John"`
	Surname     string `form:"surname" example:"Smith"`
	Patronymic  string `form:"patronymic" example:"William"`
	Age         uint   `form:"age" example:"30"`
	AgeFrom     uint   `form:"age_from" example:"18"`
	AgeTo       uint   `form:"age_to" example:"40"`
	Gender      string `form:"gender" example:"male"`
	Nationality string `form:"nationality" example:"US"`
}
