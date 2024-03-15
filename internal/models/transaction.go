package models

type TopUp struct {
	Id          int64
	Date        string
	Amount      float64
	Description string
	IncomeId    int
}
