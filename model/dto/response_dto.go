package dto

type Paginate struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalRows  int `json:"totalRows"`
	TotalPages int `json:"totalPages"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SingleResponse struct {
	Status Status `json:"status"`
	Data   any    `json:"data"`
}

type PaginateResponse struct {
	Status   Status   `json:"status"`
	Data     []any    `json:"data"`
	Paginate Paginate `json:"paginate"`
}

type ReportPaginateResponse struct {
	Status      Status   `json:"status"`
	Data        []any    `json:"data"`
	TotalIncome int64    `json:"totalIncome"`
	Paginate    Paginate `json:"paginate"`
}

type PaymentResponse struct {
	OrderId           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	Status            Status `json:"status"`
}
