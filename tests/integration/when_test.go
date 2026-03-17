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
		"currency": currency,
		"amount":   amountStr,
		"user":     user,
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
