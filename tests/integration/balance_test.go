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
		Income(100, EUR, "Initial", "Initial balance", "main").
		Expense(10, EUR, "Fun", "Cinema", "main").
		Deposit(100, EUR, "user1", "main").
		Expense(10, EUR, "Food", "Lunch", "main").
		Deposit(100, EUR, "user1", "other").
		Expense(10, EUR, "Food", "Lunch", "main").
		Transfer(50, EUR, "main", "other")

	s.Then().
		BalanceShouldBe("main", EUR, 120).
		BalanceShouldBe("other", EUR, 150).
		TransactionsDistributionShouldMatch([]map[string]string{
			{"category": "Transfer", "system_total_rate": "-0.4167"},
			{"category": "Transfer", "system_total_rate": "0.3333"},
			{"category": "Food", "system_total_rate": "0"},
			{"category": "Deposit", "system_total_rate": "0"},
			{"category": "Food", "system_total_rate": "0"},
			{"category": "Deposit", "system_total_rate": "0"},
			{"category": "Fun", "system_total_rate": "0"},
			{"category": "Initial", "system_total_rate": "0"},
		})
}
