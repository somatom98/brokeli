package integration

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Given struct {
	s *Suite
}

func (g *Given) Account(alias string, currency string) *Given {
	return g.AccountCreated(alias, currency, alias)
}

func (g *Given) AccountCreated(name, currency, alias string) *Given {
	g.s.t.Logf("Creating account %s (%s)...", name, alias)
	openReq := map[string]interface{}{
		"name":     name,
		"currency": currency,
	}
	openBody, _ := json.Marshal(openReq)
	resp, err := g.s.client.Post(g.s.baseURL+"/accounts", "application/json", bytes.NewBuffer(openBody))
	require.NoError(g.s.t, err)
	assert.Equal(g.s.t, http.StatusCreated, resp.StatusCode)

	var openResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&openResp)
	require.NoError(g.s.t, err)
	resp.Body.Close()
	accountID := openResp["id"]
	require.NotEmpty(g.s.t, accountID)
	
	g.s.accounts[alias] = accountID
	return g
}
