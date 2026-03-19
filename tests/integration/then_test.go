package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Then struct {
	s *Suite
}

func (t *Then) BalanceShouldBe(accountAlias string, currency string, expectedAmount interface{}) *Then {
	accountID := t.s.accounts[accountAlias]
	require.NotEmpty(t.s.t, accountID, "Account alias %s not found", accountAlias)

	expectedAmountStr := fmt.Sprint(expectedAmount)
	t.s.t.Logf("Checking balance for %s should be %s %s...", accountAlias, expectedAmountStr, currency)
	assert.Eventually(t.s.t, func() bool {
		resp, err := t.s.client.Get(t.s.baseURL + "/accounts/" + accountID + "/balances")
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return false
		}

		var balances []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&balances); err != nil {
			return false
		}

		for _, b := range balances {
			if b["currency"] == currency && b["amount"] == expectedAmountStr {
				return true
			}
		}

		return false
	}, 30*time.Second, 1*time.Second, "Balance should be %s %s", expectedAmountStr, currency)

	return t
}

func (t *Then) TransactionsDistributionShouldMatch(expected []map[string]string) *Then {
	t.s.t.Log("Checking transactions distribution...")
	assert.Eventually(t.s.t, func() bool {
		resp, err := t.s.client.Get(t.s.baseURL + "/transactions")
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return false
		}

		var transactions []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
			return false
		}

		if len(transactions) < len(expected) {
			return false
		}

		for i, exp := range expected {
			actual := transactions[i]
			for k, v := range exp {
				if actual[k] != v {
					return false
				}
			}
		}

		return true
	}, 30*time.Second, 1*time.Second, "Transactions distribution did not match")

	return t
}
