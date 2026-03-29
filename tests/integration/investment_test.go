package integration

import (
	"testing"
)

func TestInvestmentScenario(t *testing.T) {
	s := NewSuite(t)

	s.Given().
		Account("main", EUR)

	s.When().
		Deposit(1000, EUR, "user1", "main").
		// Buy 10 units of AAPL at 50 EUR each, with 5 EUR fee (Total 505 EUR)
		Investment("AAPL", 10, 50, EUR, 5, EUR, "main").
		// Buy 5 units of MSFT at 100 USD each, with 10 USD fee (Total 510 USD)
		Investment("MSFT", 5, 100, USD, 10, USD, "main")

	s.Then().
		// Liquidity check: 1000 - 505 = 495 EUR
		BalanceTypeShouldBe("main", EUR, 495, "LIQUIDITY").
		// Investment check: 500 EUR (fees are not credited to investment balance)
		BalanceTypeShouldBe("main", EUR, 500, "INVESTMENT").
		// Liquidity check for USD: 0 - 510 = -510 USD
		BalanceTypeShouldBe("main", USD, -510, "LIQUIDITY").
		// Investment check for USD: 500 USD
		BalanceTypeShouldBe("main", USD, 500, "INVESTMENT").
		TransactionsDistributionShouldMatch([]map[string]string{
			{"category": "Investments", "description": "MSFT", "amount": "-510", "currency": "USD"},
			{"category": "Investments", "description": "AAPL", "amount": "-505", "currency": "EUR"},
			{"category": "Deposit", "amount": "1000", "currency": "EUR"},
		})
}

func TestMixedCurrencyInvestmentScenario(t *testing.T) {
	s := NewSuite(t)

	s.Given().
		Account("main", EUR)

	s.When().
		Deposit(1000, EUR, "user1", "main").
		Deposit(100, DKK, "user1", "main").
		// Buy 10 units of TSLA at 50 EUR each, with 20 DKK fee
		Investment("TSLA", 10, 50, EUR, 20, DKK, "main")

	s.Then().
		// EUR Liquidity: 1000 - 500 = 500
		BalanceTypeShouldBe("main", EUR, 500, "LIQUIDITY").
		// EUR Investment: 500
		BalanceTypeShouldBe("main", EUR, 500, "INVESTMENT").
		// DKK Liquidity: 100 - 20 = 80
		BalanceTypeShouldBe("main", DKK, 80, "LIQUIDITY").
		TransactionsDistributionShouldMatch([]map[string]string{
			{"category": "Investments", "description": "TSLA (Fee)", "amount": "-20", "currency": "DKK", "transaction_type": "INVESTMENT"},
			{"category": "Investments", "description": "TSLA", "amount": "-500", "currency": "EUR", "transaction_type": "INVESTMENT"},
			{"category": "Deposit", "amount": "100", "currency": "DKK"},
			{"category": "Deposit", "amount": "1000", "currency": "EUR"},
		})
}
