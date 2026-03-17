package integration

import (
	"testing"
)

func TestBalanceScenario(t *testing.T) {
	s := NewSuite(t)

	s.Given().
		Account("main", EUR).
		Account("other", EUR)

	s.When().
		Deposit(100, EUR, "system", "main").
		Expense(10, EUR, "Fun", "Cinema", "main").
		Deposit(100, EUR, "user1", "main").
		Expense(10, EUR, "Food", "Lunch", "main").
		Deposit(100, EUR, "user1", "other").
		Expense(10, EUR, "Food", "Lunch", "main")

	s.Then().
		BalanceShouldBe("main", EUR, 170).
		TransactionsDistributionShouldMatch([]map[string]string{
			{"category": "Food", "system_total_rate": "0.5"},
			{"category": "Deposit", "system_total_rate": "0"},
			{"category": "Food", "system_total_rate": "0.5"},
			{"category": "Deposit", "system_total_rate": "0.5"},
			{"category": "Fun", "system_total_rate": "1"},
			{"category": "Deposit", "system_total_rate": "1"},
		}).
		AccountTransactionsCountShouldBeAtLeast("main", 5)
}
