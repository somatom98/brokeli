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
    
    // API returns a map { "uuid": { balance: {...} } }
    // We convert it to an array for the UI
    if (data && typeof data === 'object' && !Array.isArray(data)) {
      return Object.entries(data).map(([id, details]: [string, any]) => ({
        id: id,
        name: `Account ${id.slice(-4)}`, // Fallback name since projection doesn't have names yet
        balance: details.balance || {}
      }));
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
};
