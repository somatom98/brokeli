import React, { useState, useEffect, useCallback } from 'react';
import { 
  ArrowUpRight, 
  ArrowDownLeft, 
  Search,
  Calendar,
  Filter,
  Loader2,
  X
} from 'lucide-react';
import { api } from './api';
import type { Account, Transaction, TransactionFilter } from './api';

const Transactions: React.FC = () => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState<TransactionFilter>({});
  const [isFilterOpen, setIsFilterOpen] = useState(false);

  const fetchTransactions = useCallback(async () => {
    setLoading(true);
    try {
      const data = await api.getTransactions(filter);
      setTransactions(data || []);
    } catch (err) {
      console.error('Error fetching transactions:', err);
    } finally {
      setLoading(false);
    }
  }, [filter]);

  const fetchAccounts = async () => {
    try {
      const accs = await api.getAccounts();
      setAccounts(accs || []);
    } catch (err) {
      console.error('Error fetching accounts:', err);
    }
  };

  useEffect(() => {
    fetchAccounts();
  }, []);

  useEffect(() => {
    fetchTransactions();
  }, [fetchTransactions]);

  const toggleAccount = (accountId: string) => {
    const current = filter.account_id || [];
    if (current.includes(accountId)) {
      setFilter({ ...filter, account_id: current.filter(id => id !== accountId) });
    } else {
      setFilter({ ...filter, account_id: [...current, accountId] });
    }
  };

  const clearFilters = () => {
    setFilter({});
  };

  return (
    <div className="w-full flex items-start justify-center p-4 md:p-8 pb-20">
      <div className="w-full max-w-5xl space-y-8">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 px-4">
          <div>
            <h1 className="text-4xl font-black text-gray-900 tracking-tight">Transactions</h1>
            <p className="text-gray-400 font-bold uppercase tracking-widest text-[10px] mt-2">Historical Ledger Activity</p>
          </div>
          
          <div className="flex items-center gap-3">
            <button 
              onClick={() => setIsFilterOpen(!isFilterOpen)}
              className={`flex items-center gap-2 px-6 py-3 rounded-2xl font-bold transition-all shadow-lg active:scale-[0.98] ${
                isFilterOpen || Object.keys(filter).length > 0 
                  ? 'bg-indigo-600 text-white shadow-indigo-200' 
                  : 'bg-white text-gray-600 hover:bg-gray-50'
              }`}
            >
              <Filter size={20} />
              Filters {Object.keys(filter).length > 0 && `(${Object.keys(filter).length})`}
            </button>
            {Object.keys(filter).length > 0 && (
              <button 
                onClick={clearFilters}
                className="p-3 text-gray-400 hover:text-rose-500 transition-colors"
                title="Clear all filters"
              >
                <X size={20} strokeWidth={3} />
              </button>
            )}
          </div>
        </div>

        {isFilterOpen && (
          <div className="bg-white/80 backdrop-blur-2xl rounded-[32px] p-8 border border-white/50 shadow-xl animate-in fade-in slide-in-from-top-4 duration-500">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div className="space-y-4">
                <div className="flex items-center gap-2 text-[10px] font-black text-gray-400 uppercase tracking-widest">
                  <Calendar size={12} />
                  Date Range
                </div>
                <div className="space-y-2">
                  <input 
                    type="date"
                    value={filter.start_date || ''}
                    onChange={(e) => setFilter({ ...filter, start_date: e.target.value })}
                    className="w-full bg-gray-50 border-none rounded-xl px-4 py-2 text-sm font-bold text-gray-700 outline-none focus:ring-2 focus:ring-indigo-500/20"
                  />
                  <input 
                    type="date"
                    value={filter.end_date || ''}
                    onChange={(e) => setFilter({ ...filter, end_date: e.target.value })}
                    className="w-full bg-gray-50 border-none rounded-xl px-4 py-2 text-sm font-bold text-gray-700 outline-none focus:ring-2 focus:ring-indigo-500/20"
                  />
                </div>
              </div>

              <div className="space-y-4">
                <div className="flex items-center gap-2 text-[10px] font-black text-gray-400 uppercase tracking-widest">
                  <Filter size={12} />
                  Type
                </div>
                <select 
                  value={filter.transaction_type || ''}
                  onChange={(e) => setFilter({ ...filter, transaction_type: e.target.value || undefined })}
                  className="w-full bg-gray-50 border-none rounded-xl px-4 py-2 text-sm font-bold text-gray-700 outline-none focus:ring-2 focus:ring-indigo-500/20 appearance-none cursor-pointer"
                >
                  <option value="">All Types</option>
                  <option value="DEBIT">Debit (Expense)</option>
                  <option value="CREDIT">Credit (Income)</option>
                </select>
              </div>

              <div className="space-y-4">
                <div className="flex items-center gap-2 text-[10px] font-black text-gray-400 uppercase tracking-widest">
                  <Search size={12} />
                  Accounts
                </div>
                <div className="flex flex-wrap gap-2">
                  {accounts.map(acc => (
                    <button
                      key={acc.id}
                      onClick={() => toggleAccount(acc.id)}
                      className={`px-3 py-1.5 rounded-xl text-[10px] font-bold uppercase tracking-wider transition-all border ${
                        filter.account_id?.includes(acc.id)
                          ? 'bg-indigo-600 border-transparent text-white shadow-md scale-105'
                          : 'bg-white border-gray-100 text-gray-400 hover:border-gray-200'
                      }`}
                    >
                      {acc.name}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </div>
        )}

        {loading ? (
          <div className="flex flex-col items-center justify-center py-20 bg-white/50 backdrop-blur-xl rounded-[48px] border border-white/50">
            <Loader2 className="animate-spin text-indigo-600 mb-4" size={48} />
            <p className="text-gray-400 font-bold uppercase tracking-widest text-xs">Fetching Ledger Data...</p>
          </div>
        ) : transactions.length === 0 ? (
          <div className="bg-white/50 backdrop-blur-xl rounded-[48px] border border-dashed border-gray-200 p-20 flex flex-col items-center text-center">
            <div className="w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mb-6">
              <Search size={32} className="text-gray-400" />
            </div>
            <h3 className="text-xl font-bold text-gray-900 mb-2">No results found</h3>
            <p className="text-gray-400 max-w-xs">Adjust your filters to see more transaction data.</p>
          </div>
        ) : (
          <div className="bg-white/90 backdrop-blur-2xl rounded-[40px] border border-white/50 shadow-sm overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full text-left">
                <thead>
                  <tr className="border-b border-gray-50">
                    <th className="px-8 py-6 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">Date</th>
                    <th className="px-8 py-6 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">Category</th>
                    <th className="px-8 py-6 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">Description</th>
                    <th className="px-8 py-6 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] text-right">Amount</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-50">
                  {transactions.map((t) => {
                    const isDebit = t.transaction_type === 'DEBIT' || parseFloat(t.amount) < 0;
                    return (
                      <tr key={t.id} className="hover:bg-gray-50/50 transition-colors group">
                        <td className="px-8 py-6">
                          <div className="flex flex-col">
                            <span className="text-sm font-bold text-gray-900">
                              {new Date(t.happened_at).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })}
                            </span>
                            <span className="text-[10px] font-bold text-gray-400 uppercase tracking-tighter">
                                {new Date(t.happened_at).toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' })}
                            </span>
                          </div>
                        </td>
                        <td className="px-8 py-6">
                          <span className="px-3 py-1.5 bg-gray-100 text-gray-500 rounded-full text-[10px] font-black uppercase tracking-widest">
                            {t.category || 'General'}
                          </span>
                        </td>
                        <td className="px-8 py-6">
                          <div className="flex flex-col">
                            <span className="text-sm font-medium text-gray-700">{t.description || 'No description'}</span>
                            <span className="text-[10px] font-bold text-gray-300">
                                Account: {accounts.find(a => a.id === t.account_id)?.name || 'Unknown'}
                            </span>
                          </div>
                        </td>
                        <td className="px-8 py-6 text-right">
                          <div className={`flex flex-col items-end ${isDebit ? 'text-rose-500' : 'text-emerald-500'}`}>
                            <div className="flex items-center gap-1.5 font-black text-lg tracking-tighter">
                              {isDebit ? <ArrowDownLeft size={16} strokeWidth={3} /> : <ArrowUpRight size={16} strokeWidth={3} />}
                              {Math.abs(parseFloat(t.amount)).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                              <span className="text-xs ml-1">{t.currency}</span>
                            </div>
                            {t.system_total_rate && parseFloat(t.system_total_rate) !== 1 && (
                                <span className="text-[9px] font-bold opacity-60">
                                    {(Math.abs(parseFloat(t.amount)) * parseFloat(t.system_total_rate)).toLocaleString(undefined, { minimumFractionDigits: 2 })} EUR
                                </span>
                            )}
                          </div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Transactions;
