package models

type TopUp struct {
	Id          int64
	Date        string
	Amount      int
	Description string
	IncomeId    int
}
