package repositories

type ServiceRiskReport struct {
	DebtCount      map[string]int64 `json:"debtCount"`
	DependentCount int64            `json:"dependentCount"`
}

type ServiceDebtReport struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
	Count int    `json:"count"`
}
