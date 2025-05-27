package repositories

type ServiceRiskReport struct {
	DebtCount      map[string]int64
	DependentCount int64
}
