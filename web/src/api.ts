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
  system_total_rate?: string;
}

export interface BalancePeriod {
  month: string;
  currency: string;
  amount: string;
}

export interface TransactionFilter {
  start_date?: string;
  end_date?: string;
  account_id?: string[];
  transaction_type?: string;
}

export interface BudgetItem {
  name: string;
  categories: string[];
  percentage: number;
}

export interface BudgetData {
  id: string;
  name: string;
  data: {
    items: BudgetItem[];
    selectedAccounts: string[];
  };
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
      const accounts = Object.entries(data).map(([id, details]) => {
        const d = details as { name?: string, balance?: Record<string, string> };
        return {
          id: id,
          name: d.name || `Account ${id.slice(-4)}`,
          balance: d.balance || {}
        };
      });
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
  getTransactions: async (filter?: TransactionFilter): Promise<Transaction[]> => {
    const query = new URLSearchParams();
    if (filter) {
      if (filter.start_date) query.append('start_date', filter.start_date);
      if (filter.end_date) query.append('end_date', filter.end_date);
      if (filter.account_id) {
        filter.account_id.forEach(id => query.append('account_id', id));
      }
      if (filter.transaction_type) query.append('transaction_type', filter.transaction_type);
    }
    const res = await fetch(`/api/transactions?${query.toString()}`);
    if (!res.ok) throw new Error('Failed to fetch transactions');
    const data = await res.json();
    return Array.isArray(data) ? data : [];
  },
  registerExpense: async (data: { account_id: string, currency: string, amount: string, category?: string, description?: string, happened_at?: string }) => {
    const res = await fetch('/api/expenses', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to register expense');
    return res.status;
  },
  registerIncome: async (data: { account_id: string, currency: string, amount: string, category?: string, description?: string, happened_at?: string }) => {
    const res = await fetch('/api/incomes', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to register income');
    return res.status;
  },
  registerTransfer: async (data: { from_account_id: string, from_currency: string, from_amount: string, to_account_id: string, to_currency: string, to_amount: string, category?: string, description?: string, happened_at?: string }) => {
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
  getBudgets: async (): Promise<BudgetData[]> => {
    const res = await fetch('/api/budgets');
    if (!res.ok) throw new Error('Failed to fetch budgets');
    const data = await res.json();
    return Array.isArray(data) ? data : [];
  },
  saveBudget: async (data: { id?: string, name: string, data: { items: BudgetItem[], selectedAccounts: string[] } }) => {
    const res = await fetch('/api/budgets', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to save budget');
    return res.status;
  },
  deleteBudget: async (id: string) => {
    const res = await fetch(`/api/budgets/${id}`, {
      method: 'DELETE',
    });
    if (!res.ok) throw new Error('Failed to delete budget');
    return res.status;
  },
  getCategories: async (): Promise<string[]> => {
    const res = await fetch('/api/budgets/categories');
    if (!res.ok) throw new Error('Failed to fetch categories');
    const data = await res.json();
    return Array.isArray(data) ? data : [];
  },
};
