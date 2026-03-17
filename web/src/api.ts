export interface Account {
  id: string;
  name: string;
  balance: Record<string, string>;
}

export interface Transaction {
  id: string;
  account_id: string;
  currency: string;
  amount: string;
  category: string;
  description: string;
  happened_at: string;
  transaction_type: string;
}

export interface BalancePeriod {
  month: string;
  currency: string;
  amount: string;
}

export const api = {
  getAccounts: async (): Promise<Account[]> => {
    const res = await fetch('/api/accounts');
    if (!res.ok) throw new Error('Failed to fetch accounts');
    const data = await res.json();
    console.log('Raw accounts data:', data);
    
    // API returns a map { "uuid": { name: "...", balance: {...} } }
    // We convert it to an array for the UI
    if (data && typeof data === 'object' && !Array.isArray(data)) {
      const accounts = Object.entries(data).map(([id, details]: [string, any]) => ({
        id: id,
        name: details.name || `Account ${id.slice(-4)}`,
        balance: details.balance || {}
      }));
      console.log('Mapped accounts:', accounts);
      return accounts;
    }
    
    return Array.isArray(data) ? data : [];
  },
  getBalances: async (): Promise<BalancePeriod[]> => {
    const res = await fetch('/api/balances');
    if (!res.ok) throw new Error('Failed to fetch balances');
    const data = await res.json();
    return Array.isArray(data) ? data : [];
  },
  getTransactions: async (): Promise<Transaction[]> => {
    const res = await fetch('/api/transactions');
    if (!res.ok) throw new Error('Failed to fetch transactions');
    const data = await res.json();
    return Array.isArray(data) ? data : [];
  },
  registerExpense: async (data: any) => {
    const res = await fetch('/api/expenses', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to register expense');
    return res.status;
  },
  registerIncome: async (data: any) => {
    const res = await fetch('/api/incomes', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to register income');
    return res.status;
  },
  registerTransfer: async (data: any) => {
    const res = await fetch('/api/transfers', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to register transfer');
    return res.status;
  },
  openAccount: async (data: { name: string, currency: string }) => {
    const res = await fetch('/api/accounts', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to open account');
    return res.status;
  },
  deposit: async (accountId: string, data: { currency: string, amount: string }) => {
    const res = await fetch(`/api/accounts/${accountId}/deposits`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to deposit');
    return res.status;
  },
  withdraw: async (accountId: string, data: { currency: string, amount: string }) => {
    const res = await fetch(`/api/accounts/${accountId}/withdrawals`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to withdraw');
    return res.status;
  },
};
