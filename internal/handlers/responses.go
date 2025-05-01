package handlers

type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// swagger:response successfullyDeletedResponse
type SuccessfullyDeletedResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Person successfully deleted"`
	Status  string `json:"status" example:"ok"`
}

// swagger:response badRequestResponse
type BadRequestResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid ID format"`
	Status  string `json:"status" example:"error"`
}

// swagger:response internalServerErrorResponse
type InternalServerErrorResponse struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Some error message"`
	Status  string `json:"status" example:"error"`
}

// swagger:response personResponse
type PersonResponse struct {
	ID          int    `json:"ID" example:"42"`
	CreatedAt   string `json:"CreatedAt" example:"2025-05-01T03:35:20.374096098+03:00"`
	UpdatedAt   string `json:"UpdatedAt" example:"2025-05-01T03:35:20.374096098+03:00"`
	DeletedAt   string `json:"DeletedAt" example:""`
	Name        string `json:"Name" example:"John"`
	Surname     string `json:"Surname" example:"Smith"`
	Patronymic  string `json:"Patronymic" example:""`
	Age         int    `json:"Age" example:"49"`
	Gender      string `json:"Gender" example:"female"`
	Nationality string `json:"Nationality" example:"SK"`
}
