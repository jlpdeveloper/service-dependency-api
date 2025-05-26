package repositories

type ServiceRiskReport struct {
	DebtCount      map[string]int
	DependentCount int
}
