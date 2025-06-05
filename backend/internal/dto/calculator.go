package dto

type CalculationRequest struct {
	Items     int   `json:"items" validate:"required,min=0" example:"251"`
	PackSizes []int `json:"pack_sizes" validate:"required,min=1,dive,min=1" swaggertype:"array,integer" example:"250,500,1000"`
}

type CalculationResponse struct {
	Allocation map[int]int `json:"allocation" swaggertype:"object,integer" example:"500:1"`
	TotalPacks int         `json:"total_packs" example:"1"`
	TotalItems int         `json:"total_items" example:"500"`
	Surplus    int         `json:"surplus" example:"249"`
}
