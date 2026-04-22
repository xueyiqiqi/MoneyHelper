package dto

type CreateBillRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Category    string  `json:"category" binding:"required"`
	Remarks     string  `json:"description"`
	Date        string  `json:"bill_date"`
	IsPersonal  bool    `json:"is_personal"`
	SpaceID     *uint   `json:"space_id,omitempty"`
}

type BillResponse struct {
	ID          uint    `json:"id"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Remarks     string  `json:"description"`
	Date        string  `json:"bill_date"`
	CreatorName string  `json:"creator_name,omitempty"`
}