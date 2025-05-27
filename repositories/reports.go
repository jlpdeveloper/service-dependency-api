package repositories

type ServiceRiskReport struct {
	DebtCount      map[string]int64 `json:"debtCount"`
	DependentCount int64            `json:"dependentCount"`
}
