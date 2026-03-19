package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type When struct {
	s *Suite
}

func (w *When) Deposit(amount interface{}, currency string, user string, accountAlias string) *When {
	accountID := w.s.accounts[accountAlias]
	require.NotEmpty(w.s.t, accountID, "Account alias %s not found", accountAlias)

	amountStr := fmt.Sprint(amount)
	w.s.t.Logf("Depositing %s %s from %s to %s...", amountStr, currency, user, accountAlias)
	depositReq := map[string]interface{}{
		"currency":    currency,
		"amount":      amountStr,
		"user":        user,
		"category":    "Deposit",
		"description": "Initial deposit",
	}
	depositBody, _ := json.Marshal(depositReq)
	resp, err := w.s.client.Post(w.s.baseURL+"/accounts/"+accountID+"/deposits", "application/json", bytes.NewBuffer(depositBody))
	require.NoError(w.s.t, err)
	assert.Equal(w.s.t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	return w
}

func (w *When) Expense(amount interface{}, currency string, category string, description string, accountAlias string) *When {
	accountID := w.s.accounts[accountAlias]
	require.NotEmpty(w.s.t, accountID, "Account alias %s not found", accountAlias)

	amountStr := fmt.Sprint(amount)
	w.s.t.Logf("Creating %s %s expense for %s...", amountStr, currency, accountAlias)
	expenseReq := map[string]interface{}{
		"account_id":  accountID,
		"currency":    currency,
		"amount":      amountStr,
		"category":    category,
		"description": description,
	}
	expenseBody, _ := json.Marshal(expenseReq)
	resp, err := w.s.client.Post(w.s.baseURL+"/expenses", "application/json", bytes.NewBuffer(expenseBody))
	require.NoError(w.s.t, err)
	assert.Equal(w.s.t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	return w
}

func (w *When) Income(amount interface{}, currency string, category string, description string, accountAlias string) *When {
	accountID := w.s.accounts[accountAlias]
	require.NotEmpty(w.s.t, accountID, "Account alias %s not found", accountAlias)

	amountStr := fmt.Sprint(amount)
	w.s.t.Logf("Creating %s %s income for %s...", amountStr, currency, accountAlias)
	incomeReq := map[string]interface{}{
		"account_id":  accountID,
		"currency":    currency,
		"amount":      amountStr,
		"category":    category,
		"description": description,
	}
	incomeBody, _ := json.Marshal(incomeReq)
	resp, err := w.s.client.Post(w.s.baseURL+"/incomes", "application/json", bytes.NewBuffer(incomeBody))
	require.NoError(w.s.t, err)
	assert.Equal(w.s.t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	return w
}

func (w *When) Transfer(amount interface{}, currency string, fromAlias string, toAlias string) *When {
	fromAccountID := w.s.accounts[fromAlias]
	require.NotEmpty(w.s.t, fromAccountID, "Account alias %s not found", fromAlias)

	toAccountID := w.s.accounts[toAlias]
	require.NotEmpty(w.s.t, toAccountID, "Account alias %s not found", toAlias)

	amountStr := fmt.Sprint(amount)
	w.s.t.Logf("Transferring %s %s from %s to %s...", amountStr, currency, fromAlias, toAlias)
	transferReq := map[string]interface{}{
		"from_account_id": fromAccountID,
		"from_currency":   currency,
		"from_amount":     amountStr,
		"to_account_id":   toAccountID,
		"to_currency":     currency,
		"to_amount":       amountStr,
		"category":        "Transfer",
		"description":     fmt.Sprintf("Transfer from %s to %s", fromAlias, toAlias),
	}
	transferBody, _ := json.Marshal(transferReq)
	resp, err := w.s.client.Post(w.s.baseURL+"/transfers", "application/json", bytes.NewBuffer(transferBody))
	require.NoError(w.s.t, err)
	assert.Equal(w.s.t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	return w
}
