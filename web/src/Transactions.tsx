import React, { useState, useEffect, useCallback } from 'react';
import { 
  ArrowUpRight, 
  ArrowDownLeft, 
  ArrowRightLeft,
  RotateCcw,
  Search,
  Calendar,
  Filter,
  Loader2,
  X,
  TrendingUp
} from 'lucide-react';
import { api } from './api';
import type { Account, Transaction, TransactionFilter } from './api';

interface TransactionsProps {
  refreshKey?: number;
  hideHeader?: boolean;
}

const Transactions: React.FC<TransactionsProps> = ({ refreshKey, hideHeader }) => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [filter, setFilter] = useState<TransactionFilter>({});
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const pageSize = 100;

  const fetchTransactions = useCallback(async (p: number, isInitial = false) => {
    if (isInitial) {
      setLoading(true);
      setHasMore(true);
    } else {
      setLoadingMore(true);
    }

    try {
      const data = await api.getPaginatedTransactions({
        ...filter,
        page: p,
        page_size: pageSize
      });

      const newTransactions = data.transactions || [];
      const totalCount = data.total_count || 0;

      if (isInitial) {
        setTransactions(newTransactions);
        setHasMore(newTransactions.length < totalCount);
      } else {
        setTransactions(prev => {
          const updated = [...prev, ...newTransactions];
          setHasMore(updated.length < totalCount);
          return updated;
        });
      }
    } catch (err) {
      console.error('Error fetching transactions:', err);
    } finally {
      setLoading(false);
      setLoadingMore(false);
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
    setPage(1);
    fetchTransactions(1, true);
  }, [filter, refreshKey, fetchTransactions]);

  useEffect(() => {
    if (page > 1) {
      fetchTransactions(page, false);
    }
  }, [page, fetchTransactions]);

  const observer = React.useRef<IntersectionObserver>(null);
  const lastElementRef = useCallback((node: HTMLTableRowElement | null) => {
    if (loading || loadingMore) return;
    if (observer.current) observer.current.disconnect();

    observer.current = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting && hasMore) {
        setPage(prev => prev + 1);
      }
    });

    if (node) observer.current.observe(node);
  }, [loading, loadingMore, hasMore]);

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
    <div className={`w-full h-full flex flex-col ${hideHeader ? 'pt-0' : 'p-4 md:p-8'}`}>
      <div className="w-full flex-1 flex flex-col min-h-0 space-y-4 pt-0">
        {hideHeader ? (
          <div className="flex justify-end px-4">
            <div className="flex items-center gap-3">
              <button 
                onClick={() => setIsFilterOpen(!isFilterOpen)}
                className="p-2.5 bg-glass backdrop-blur-xl rounded-xl shadow-lg hover:bg-glass-hover transition-all text-text-on-dark border border-glass-border active:scale-95 flex items-center gap-2 px-4"
              >
                <Filter size={18} strokeWidth={2.5} />
                <span className="text-xs font-bold uppercase tracking-widest">
                  Filters {Object.keys(filter).length > 0 && `(${Object.keys(filter).length})`}
                </span>
              </button>
              {Object.keys(filter).length > 0 && (
                <button 
                  onClick={clearFilters}
                  className="p-2.5 bg-glass backdrop-blur-xl rounded-xl shadow-lg hover:bg-glass-hover transition-all text-text-on-dark border border-glass-border active:scale-95"
                  title="Clear all filters"
                >
                  <X size={18} strokeWidth={2.5} />
                </button>
              )}
            </div>
          </div>
        ) : (
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 px-4">
            <div>
              <h1 className="text-4xl font-black text-text-on-dark tracking-tight">Transactions</h1>
              <p className="text-text-on-dark/40 font-bold uppercase tracking-widest text-[10px] mt-2">Historical Ledger Activity</p>
            </div>
            
            <div className="flex items-center gap-3">
              <button 
                onClick={() => setIsFilterOpen(!isFilterOpen)}
                className="p-2.5 bg-glass backdrop-blur-xl rounded-xl shadow-lg hover:bg-glass-hover transition-all text-text-on-dark border border-glass-border active:scale-95 flex items-center gap-2 px-4"
              >
                <Filter size={18} strokeWidth={2.5} />
                <span className="text-xs font-bold uppercase tracking-widest">
                  Filters {Object.keys(filter).length > 0 && `(${Object.keys(filter).length})`}
                </span>
              </button>
              {Object.keys(filter).length > 0 && (
                <button 
                  onClick={clearFilters}
                  className="p-2.5 bg-glass backdrop-blur-xl rounded-xl shadow-lg hover:bg-glass-hover transition-all text-text-on-dark border border-glass-border active:scale-95"
                  title="Clear all filters"
                >
                  <X size={18} strokeWidth={2.5} />
                </button>
              )}
            </div>
          </div>
        )}

        {isFilterOpen && (
          <div className="bg-card rounded-[32px] p-8 border border-border-pearl shadow-2xl animate-in fade-in slide-in-from-top-4 duration-500 shrink-0 mx-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div className="space-y-4">
                <div className="flex items-center gap-2 text-[10px] font-black text-text-muted uppercase tracking-widest">
                  <Calendar size={12} />
                  Date Range
                </div>
                <div className="space-y-2">
                  <input 
                    type="date"
                    value={filter.start_date || ''}
                    onChange={(e) => setFilter({ ...filter, start_date: e.target.value })}
                    className="w-full bg-card-muted border-none rounded-xl px-4 py-2 text-sm font-bold text-text-main outline-none focus:ring-2 focus:ring-accent/10"
                  />
                  <input 
                    type="date"
                    value={filter.end_date || ''}
                    onChange={(e) => setFilter({ ...filter, end_date: e.target.value })}
                    className="w-full bg-card-muted border-none rounded-xl px-4 py-2 text-sm font-bold text-text-main outline-none focus:ring-2 focus:ring-accent/10"
                  />
                </div>
              </div>

              <div className="space-y-4">
                <div className="flex items-center gap-2 text-[10px] font-black text-text-muted uppercase tracking-widest">
                  <Filter size={12} />
                  Type
                </div>
                <select 
                  value={filter.transaction_type || ''}
                  onChange={(e) => setFilter({ ...filter, transaction_type: e.target.value || undefined })}
                  className="w-full bg-card-muted border-none rounded-xl px-4 py-2 text-sm font-bold text-text-main outline-none focus:ring-2 focus:ring-accent/10 appearance-none cursor-pointer"
                >
                  <option value="">All Types</option>
                  <option value="EXPENSE">Expense</option>
                  <option value="INCOME">Income</option>
                  <option value="TRANSFER">Transfer</option>
                  <option value="REIMBURSEMENT">Reimbursement</option>
                  <option value="DEPOSIT">Deposit</option>
                  <option value="WITHDRAWAL">Withdrawal</option>
                </select>
              </div>

              <div className="space-y-4">
                <div className="flex items-center gap-2 text-[10px] font-black text-text-muted uppercase tracking-widest">
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
                          ? 'bg-accent border-transparent text-white shadow-md scale-105'
                          : 'bg-card-muted border-border-pearl text-text-muted hover:border-text-muted/20'
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
          <div className="flex-1 flex flex-col items-center justify-center bg-card/5 backdrop-blur-xl rounded-[48px] border border-white/5 m-4">
            <Loader2 className="animate-spin text-primary mb-4" size={48} />
            <p className="text-text-on-dark/40 font-bold uppercase tracking-widest text-xs">Fetching Ledger Data...</p>
          </div>
        ) : transactions.length === 0 ? (
          <div className="flex-1 flex flex-col items-center justify-center bg-card/5 backdrop-blur-xl rounded-[48px] border border-dashed border-white/5 p-20 text-center m-4">
            <div className="w-20 h-20 bg-glass-dark rounded-full flex items-center justify-center mb-6">
              <Search size={32} className="text-text-on-dark/20" />
            </div>
            <h3 className="text-xl font-bold text-text-on-dark mb-2">No results found</h3>
            <p className="text-text-on-dark/40 max-w-xs">Adjust your filters to see more transaction data.</p>
          </div>
        ) : (
          <div className="flex-1 overflow-y-auto px-4 pb-8 space-y-3 custom-scrollbar min-h-0 bg-app-bg">            <div className="grid grid-cols-1 gap-3">
              {transactions.map((t, index) => {
                const amount = Math.abs(parseFloat(t.amount));
                const isDebit = ['EXPENSE', 'WITHDRAWAL'].includes(t.transaction_type) || 
                               (t.transaction_type === 'TRANSFER' && parseFloat(t.amount) < 0) ||
                               parseFloat(t.amount) < 0;
                const isLast = index === transactions.length - 1;
                const isInterest = t.category?.toLowerCase().includes('interest') || 
                                  t.description?.toLowerCase().includes('interest');

                const typeColorClasses = 
                  t.transaction_type === 'TRANSFER' ? 'bg-neutral/20 text-neutral' :
                  t.transaction_type === 'REIMBURSEMENT' ? 'bg-primary/20 text-primary' :
                  t.transaction_type === 'WITHDRAWAL' ? 'bg-neutral/20 text-neutral' :
                  isInterest ? 'bg-primary/20 text-primary' :
                  isDebit ? 'bg-negative/20 text-negative' : 'bg-primary/20 text-primary';

                return (
                  <div 
                    key={t.id} 
                    ref={isLast ? lastElementRef : undefined}
                    className={`bg-card rounded-[24px] p-4 lg:p-5 border shadow-sm hover:shadow-xl transition-all duration-300 group flex items-center gap-3 lg:gap-5 ${
                        isDebit 
                          ? 'border-border-pearl' 
                          : 'border-primary/20'
                    }`}
                  >
                    {/* Left: Type Icon/Badge */}
                    <div className={`p-3 lg:p-4 rounded-2xl shrink-0 shadow-sm ${typeColorClasses}`}>
                      {t.transaction_type === 'TRANSFER' ? <ArrowRightLeft size={22} strokeWidth={2.5} /> :
                       t.transaction_type === 'REIMBURSEMENT' ? <RotateCcw size={22} strokeWidth={2.5} /> :
                       t.transaction_type === 'WITHDRAWAL' ? <ArrowDownLeft size={22} strokeWidth={2.5} /> :
                       isInterest ? <TrendingUp size={22} strokeWidth={2.5} /> :
                       isDebit ? <ArrowDownLeft size={22} strokeWidth={2.5} /> : 
                       <ArrowUpRight size={22} strokeWidth={2.5} />}
                    </div>

                    {/* Middle: Description & Meta */}
                    <div className="flex-1 min-w-0 py-1">
                      <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-3 mb-1.5">
                        <span className={`text-sm lg:text-base font-bold line-clamp-2 tracking-tight text-text-main`}>
                          {t.description || 'No description'}
                        </span>
                        <span className={`px-2.5 py-1 rounded-lg text-[9px] font-black uppercase tracking-[0.1em] whitespace-nowrap w-fit ${typeColorClasses.replace('bg-accent/20', 'bg-accent/10').replace('bg-negative/20', 'bg-negative/10').replace('bg-primary/20', 'bg-primary/10').replace('bg-neutral/20', 'bg-neutral/10')}`}>
                          {t.category || 'General'}
                        </span>
                      </div>
                      <div className="flex flex-wrap items-center gap-x-3 gap-y-1 text-[11px] font-bold text-text-muted">
                        <div className="flex items-center gap-1.5">
                          <span className={`w-1.5 h-1.5 rounded-full ${
                            t.transaction_type === 'TRANSFER' || t.transaction_type === 'WITHDRAWAL' || t.transaction_type === 'REIMBURSEMENT'
                              ? 'bg-neutral/40' 
                              : isDebit ? 'bg-negative/40' : 'bg-primary/40'
                          }`} />
                          <span className="uppercase tracking-wider">
                            {accounts.find(a => a.id === t.account_id)?.name || 'Unknown'}
                          </span>
                        </div>
                        <span className="text-border-pearl font-normal">|</span>
                        <span className="tracking-tight">
                          {new Date(t.happened_at).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })}
                        </span>
                      </div>
                    </div>

                    {/* Right: Amount & Type Label */}
                    <div className="text-right flex flex-col items-end gap-1 px-2 shrink-0">
                      <div className={`flex items-center gap-1 font-black text-lg lg:text-xl tracking-tighter ${
                        t.transaction_type === 'TRANSFER' ? 'text-neutral' :
                        t.transaction_type === 'REIMBURSEMENT' ? 'text-primary' :
                        t.transaction_type === 'WITHDRAWAL' ? 'text-neutral' :
                        isDebit ? 'text-negative' : 'text-primary'
                      }`}>
                        {isDebit ? '-' : '+'}{amount.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                        <span className="text-xs ml-0.5 font-bold">{t.currency}</span>
                      </div>
                      
                      <div className="flex flex-col items-end">
                        <span className={`text-[10px] font-black uppercase tracking-widest ${
                          t.transaction_type === 'TRANSFER' ? 'text-neutral/60' :
                          t.transaction_type === 'REIMBURSEMENT' ? 'text-primary/60' :
                          t.transaction_type === 'WITHDRAWAL' ? 'text-neutral/60' :
                          isDebit ? 'text-negative/60' : 'text-primary/60'
                        }`}>
                          {t.transaction_type.replace('_', ' ')}
                        </span>
                      </div>
                    </div>
                  </div>
                );
              })}
              {loadingMore && (
                <div className="py-6 flex items-center justify-center gap-2">
                  <Loader2 className="animate-spin text-primary" size={16} />
                  <span className="text-[10px] font-black text-text-muted uppercase tracking-widest">Loading more...</span>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Transactions;
