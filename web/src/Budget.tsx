import React, { useState, useEffect, useMemo, useRef } from 'react';
import { Trash2, X, PlusCircle, Save, Check, Loader2, ChevronLeft, Layout, XCircle, Pencil, ChevronDown, ChevronRight, Calendar } from 'lucide-react';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Pie } from 'react-chartjs-2';
import ChartDataLabels from 'chartjs-plugin-datalabels';
import { api } from './api';

ChartJS.register(ArcElement, Tooltip, Legend, ChartDataLabels);
import type { Account, Transaction, TransactionFilter, BudgetItem, BudgetData } from './api';
import { getCSSVariableValue } from './utils/colors';

// Fallback for crypto.randomUUID() if not in a secure context
const generateId = () => {
  try {
    return crypto.randomUUID();
  } catch {
    return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
  }
};

const TransactionList: React.FC<{ transactions: Transaction[], accounts: Account[] }> = ({ transactions, accounts }) => {
  if (transactions.length === 0) return <div className="text-[10px] font-bold text-text-muted/40 uppercase tracking-widest py-4 px-2 italic">No transactions</div>;

  return (
    <div className="mt-4 border-t border-border-pearl pt-4 space-y-2 animate-in fade-in slide-in-from-top-2 duration-300">
      {transactions.map((t) => {
        const amount = Math.abs(parseFloat(t.amount));
        const isInvestment = t.transaction_type === 'INVESTMENT';
        const isDebit = isInvestment || 
                       ['EXPENSE', 'WITHDRAWAL'].includes(t.transaction_type) ||
                       (t.transaction_type === 'TRANSFER' && parseFloat(t.amount) < 0) ||
                       parseFloat(t.amount) < 0;
        const rate = (parseFloat(t.system_total_rate || '1') || 1);
        const systemAmount = amount * rate;

        const textColor = isInvestment ? 'text-accent-secondary' : (!isDebit ? 'text-accent-secondary' : 'text-black');

        return (
          <div key={t.id} className="flex items-center justify-between py-2 px-3 hover:bg-card-muted rounded-xl transition-colors">
            <div className="flex flex-col">
              <span className="text-[11px] font-bold text-text-main leading-tight">
                {t.description || t.category || 'Transaction'}
              </span>
              <span className="text-[9px] font-bold text-text-muted/40 uppercase tracking-tighter">
                {new Date(t.happened_at).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })} • {accounts.find(a => a.id === t.account_id)?.name || 'Unknown'}
              </span>
            </div>
            <div className={`flex items-center gap-1.5 font-black text-xs tracking-tighter ${textColor}`}>
              {isDebit ? '-' : '+'}
              {amount.toLocaleString(undefined, { minimumFractionDigits: 2 })}
              <span className="text-[8px] opacity-70 ml-0.5">{t.currency}</span>
              {rate !== 1 && (
                <span className="text-[9px] text-text-muted/60 font-bold ml-1 italic">
                  (System: {systemAmount.toLocaleString(undefined, { minimumFractionDigits: 2 })} {t.currency})
                </span>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
};
interface ProcessedBudgetItem extends BudgetItem {
  actualSpent: number;
  actualPercentage: number;
  transactions: Transaction[];
  prevMonthPercentage?: number;
  avg12MonthsPercentage?: number;
}

const Budget: React.FC = () => {
  const [view, setView] = useState<'list' | 'edit' | 'view'>('list');
  const [budgets, setBudgets] = useState<BudgetData[]>([]);
  const [selectedBudget, setSelectedBudget] = useState<BudgetData | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [isFetchingTransactions, setIsFetchingTransactions] = useState(false);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [categories, setCategories] = useState<string[]>([]);
  const [selectedAccounts, setSelectedAccounts] = useState<string[]>([]);
  const [expandedSections, setExpandedSections] = useState<string[]>([]);
  const [selectedMonth, setSelectedMonth] = useState<string>(() => {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
  });
  const [timeframe, setTimeframe] = useState<'month' | 'year'>('month');

  const monthInputRef = useRef<HTMLInputElement>(null);

  const handleMonthClick = () => {
    if (monthInputRef.current) {
      try {
        monthInputRef.current.showPicker();
      } catch {
        monthInputRef.current.click();
      }
    }
  };

  const toggleSection = (sectionId: string) => {
    setExpandedSections(prev => 
      prev.includes(sectionId) ? prev.filter(s => s !== sectionId) : [...prev, sectionId]
    );
  };
  const [items, setItems] = useState<BudgetItem[]>([]);
  const [otherPercentage, setOtherPercentage] = useState<number>(0);
  const [budgetName, setBudgetName] = useState('Monthly Budget');
  const [budgetId, setBudgetId] = useState<string>(generateId());
  const [isSaving, setIsSaving] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');

  const fetchBudgets = async () => {
    setIsLoading(true);
    try {
      const data = await api.getBudgets();
      setBudgets(data || []);
    } catch (err) {
      console.error('Error fetching budgets:', err);
      setError(true);
      setErrorMessage('Failed to fetch budgets');
      setTimeout(() => setError(false), 3500);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchCategories = async () => {
    try {
      const cats = await api.getCategories();
      setCategories(cats || []);
    } catch (err) {
      console.error('Error fetching categories:', err);
    }
  };

  useEffect(() => {
    const fetchAccounts = async () => {
      try {
        const accs = await api.getAccounts();
        setAccounts(accs || []);
      } catch (err) {
        console.error('Error fetching accounts:', err);
        setError(true);
        setErrorMessage('Failed to fetch accounts');
        setTimeout(() => setError(false), 3500);
      }
    };
    fetchAccounts();
    fetchBudgets();
    fetchCategories();
  }, []);

  useEffect(() => {
    if (view === 'view' && selectedBudget) {
      const fetchTransactions = async () => {
        setIsFetchingTransactions(true);
        try {
          const [year, month] = selectedMonth.split('-').map(Number);
          // Fetch last 13 months (current month + 12 previous months)
          const currentEnd = new Date(year, month, 0, 23, 59, 59, 999);
          const historicalStart = new Date(year, month - 13, 1).toISOString();
          
          const filter: TransactionFilter = {
            start_date: historicalStart,
            end_date: currentEnd.toISOString(),
            account_id: selectedBudget.data.selectedAccounts
          };

          const data = await api.getTransactions(filter);
          setTransactions(data || []);
        } catch (err) {
          console.error('Error fetching transactions:', err);
        } finally {
          setIsFetchingTransactions(false);
        }
      };
      fetchTransactions();
    }
  }, [view, selectedBudget, selectedMonth]);

  const budgetStats = useMemo(() => {
    if (!selectedBudget) return { 
      totalSpending: 0, 
      totalIncome: 0, 
      totalOutcome: 0, 
      items: [], 
      incomeTransactions: [], 
      outcomeTransactions: [],
      otherItem: {
        name: 'Others',
        actualSpent: 0,
        actualPercentage: 0,
        transactions: [] as Transaction[],
        categories: [] as string[],
        percentage: 0
      },
      expectedSavingsPercentage: 0
    };

    const selectedAccIds = selectedBudget.data.selectedAccounts || [];

    const [year, month] = selectedMonth.split('-').map(Number);
    const start = timeframe === 'year' ? new Date(year, month - 13, 1) : new Date(year, month - 1, 1);
    const end = new Date(year, month, 0, 23, 59, 59, 999);

    const currentTransactions = transactions.filter(t => {
      const happenedAt = new Date(t.happened_at);
      const isAccountMatch = selectedAccIds.includes(t.account_id);
      const isDateMatch = happenedAt >= start && happenedAt <= end;
      return isAccountMatch && isDateMatch;
    });

    const incomeTransactions = currentTransactions
      .filter(t => t.transaction_type !== 'DEPOSIT' && t.transaction_type !== 'TRANSFER' && t.transaction_type !== 'INVESTMENT' &&
                   (['INCOME'].includes(t.transaction_type) || 
                   (t.transaction_type === '' && parseFloat(t.amount) > 0)));

    const totalIncome = incomeTransactions
      .reduce((sum, t) => {
        const rate = parseFloat(t.system_total_rate || '1') || 1;
        return sum + (parseFloat(t.amount) * rate);
      }, 0);

    const investmentTransactions = currentTransactions
      .filter(t => t.transaction_type === 'INVESTMENT');

    const totalInvestments = investmentTransactions
      .reduce((sum, t) => {
        const rate = parseFloat(t.system_total_rate || '1') || 1;
        const amount = Math.abs(parseFloat(t.amount));
        return sum + (amount * rate);
      }, 0);

    const outcomeTransactions = currentTransactions
      .filter(t => t.transaction_type !== 'WITHDRAWAL' && t.transaction_type !== 'TRANSFER' &&
                   (['EXPENSE', 'REIMBURSEMENT', 'INVESTMENT'].includes(t.transaction_type) || 
                   (t.transaction_type === '' && parseFloat(t.amount) < 0)));

    const totalOutcome = outcomeTransactions
      .reduce((sum, t) => {
        const rate = parseFloat(t.system_total_rate || '1') || 1;
        const amount = parseFloat(t.amount);
        return sum + ((-amount) * rate);
      }, 0);

    const totalSpending = totalOutcome;

    const assignedCategories = new Set((selectedBudget.data.items || []).flatMap(i => i.categories));
    
    // Helper to calculate percentage of income for a category over a period
    const calcPeriodStats = (cats: string[] | null, startDate: Date, endDate: Date) => {
      const periodTransactions = transactions.filter(t => {
        const d = new Date(t.happened_at);
        const isAccountMatch = selectedAccIds.includes(t.account_id);
        return isAccountMatch && d >= startDate && d <= endDate;
      });

      const periodIncome = periodTransactions
        .filter(t => t.transaction_type !== 'DEPOSIT' && t.transaction_type !== 'TRANSFER' &&
                     (['INCOME'].includes(t.transaction_type) || (t.transaction_type === '' && parseFloat(t.amount) > 0)))
        .reduce((sum, t) => sum + (parseFloat(t.amount) * (parseFloat(t.system_total_rate || '1') || 1)), 0);

      const periodSpent = periodTransactions
        .filter(t => t.transaction_type !== 'WITHDRAWAL' && t.transaction_type !== 'TRANSFER' &&
                     (['EXPENSE', 'REIMBURSEMENT', 'INVESTMENT'].includes(t.transaction_type) || (t.transaction_type === '' && parseFloat(t.amount) < 0)) &&
                     (cats === null ? !assignedCategories.has(t.category) : cats.includes(t.category)))
        .reduce((sum, t) => sum + ((-parseFloat(t.amount)) * (parseFloat(t.system_total_rate || '1') || 1)), 0);

      return periodIncome > 0 ? (periodSpent / periodIncome) * 100 : 0;
    };

    const prevMonthStart = new Date(year, month - 2, 1);
    const prevMonthEnd = new Date(year, month - 1, 0, 23, 59, 59, 999);
    const last12Start = new Date(year, month - 13, 1);
    const last12End = new Date(year, month - 1, 0, 23, 59, 59, 999);

    const otherTransactions = outcomeTransactions.filter(t => !assignedCategories.has(t.category));
    const otherSpent = otherTransactions.reduce((sum, t) => sum + ((-parseFloat(t.amount)) * (parseFloat(t.system_total_rate || '1') || 1)), 0);

    const totalAssignedPercentage = (selectedBudget.data.items || []).reduce((sum, i) => sum + (i.percentage || 0), 0);
    const budgetOtherPercentage = selectedBudget.data.otherPercentage || 0;
    const expectedSavingsPercentage = Math.max(0, 100 - totalAssignedPercentage - budgetOtherPercentage);

    const items = (selectedBudget.data.items || []).map(item => {
      const itemTransactions = currentTransactions.filter(t => 
        t.transaction_type !== 'WITHDRAWAL' && t.transaction_type !== 'TRANSFER' &&
        (['EXPENSE', 'REIMBURSEMENT', 'INVESTMENT'].includes(t.transaction_type) || (t.transaction_type === '' && parseFloat(t.amount) < 0)) &&
        item.categories.includes(t.category)
      );
      const itemSpent = itemTransactions.reduce((sum, t) => sum + ((-parseFloat(t.amount)) * (parseFloat(t.system_total_rate || '1') || 1)), 0);
      
      return {
        ...item,
        actualSpent: itemSpent,
        actualPercentage: totalIncome > 0 ? (itemSpent / totalIncome) * 100 : 0,
        transactions: itemTransactions,
        prevMonthPercentage: calcPeriodStats(item.categories, prevMonthStart, prevMonthEnd),
        avg12MonthsPercentage: calcPeriodStats(item.categories, last12Start, last12End)
      };
    }).sort((a, b) => b.actualSpent - a.actualSpent);

    return {
      totalSpending,
      totalIncome,
      totalOutcome,
      totalInvestments,
      incomeTransactions,
      outcomeTransactions,
      investmentTransactions,
      otherItem: {
        name: 'Others',
        actualSpent: otherSpent,
        actualPercentage: totalIncome > 0 ? (otherSpent / totalIncome) * 100 : 0,
        transactions: otherTransactions,
        categories: [] as string[],
        percentage: budgetOtherPercentage,
        prevMonthPercentage: calcPeriodStats(null, prevMonthStart, prevMonthEnd),
        avg12MonthsPercentage: calcPeriodStats(null, last12Start, last12End)
      },
      expectedSavingsPercentage,
      items
    };
  }, [selectedBudget, transactions, selectedMonth, timeframe]);

  const renderBudgetItemRow = (item: ProcessedBudgetItem, id: string) => (
    <div key={id} className="space-y-3 p-4 hover:bg-card-muted rounded-3xl transition-colors cursor-pointer group/item" onClick={() => toggleSection(id)}>
      <div className="flex justify-between items-start">
        <div className="flex-1">
          <div className="flex items-center gap-2">
            <span className="text-xl font-bold text-text-main">{item.name}</span>
            {expandedSections.includes(id) ? <ChevronDown size={16} className="text-text-muted/40" /> : <ChevronRight size={16} className="text-text-muted/20 group-hover/item:text-accent transition-colors" />}
          </div>
          {item.categories && item.categories.length > 0 && (
            <div className="flex flex-wrap gap-1 mt-1.5">
              {item.categories.slice(0, 3).map((c: string) => (
                <span key={c} className="text-[8px] font-bold bg-card-muted text-text-muted/40 px-2 py-0.5 rounded-full uppercase tracking-wider border border-border-pearl">{c}</span>
              ))}
              {item.categories.length > 3 && (
                  <span className="text-[8px] font-bold text-text-muted/20 px-1 py-0.5 uppercase tracking-wider">+{item.categories.length - 3} more</span>
              )}
            </div>
          )}
        </div>
        <div className="text-right flex flex-col items-end gap-3">
          <div className="flex flex-col items-end">
            <div className={`text-lg font-bold ${(item.actualSpent ?? 0) < 0 ? 'text-accent-secondary' : 'text-black'}`}>{item.actualSpent?.toLocaleString(undefined, { minimumFractionDigits: 2 })}</div>
            <div className={`text-[10px] font-bold uppercase tracking-widest ${item.actualPercentage > item.percentage ? 'text-accent-secondary' : 'text-text-muted/40'}`}>
              {item.actualPercentage?.toFixed(1)}% vs {item.percentage}%
            </div>
          </div>
          
          <div className="flex gap-4 border-t border-border-pearl pt-2">
            <div className="flex flex-col items-end">
              <span className="text-[7px] font-black text-text-muted/20 uppercase tracking-widest">Prev. Month</span>
              <span className={`text-[10px] font-bold ${(item.prevMonthPercentage ?? 0) > item.percentage ? 'text-accent-secondary' : 'text-text-muted/40'}`}>{item.prevMonthPercentage?.toFixed(1)}% <span className="text-[8px] font-normal opacity-50">inc</span></span>
            </div>
            <div className="flex flex-col items-end">
              <span className="text-[7px] font-black text-text-muted/20 uppercase tracking-widest">12M Avg</span>
              <span className={`text-[10px] font-bold ${(item.avg12MonthsPercentage ?? 0) > item.percentage ? 'text-accent-secondary' : 'text-text-muted/40'}`}>{item.avg12MonthsPercentage?.toFixed(1)}% <span className="text-[8px] font-normal opacity-50">inc</span></span>
            </div>
          </div>
        </div>
      </div>
      <div className="h-3 bg-glass-dark rounded-full overflow-hidden flex relative border border-border-pearl shadow-inner">
        <div 
          className={`h-full transition-all duration-1000 ${item.actualPercentage > item.percentage ? 'bg-accent' : 'bg-glass-dark-hover'}`}
          style={{ width: `${Math.min(100, item.actualPercentage)}%` }}
        />
        {/* Target line */}
        <div 
          className="absolute top-0 bottom-0 w-1 z-10 bg-accent-secondary"
          style={{ left: `${item.percentage}%` }}
          title={`Target: ${item.percentage}%`}
        />
      </div>
      {expandedSections.includes(id) && <TransactionList transactions={item.transactions || []} accounts={accounts} />}
    </div>
  );

  const handleAddAccount = (accountId: string) => {
    if (accountId && !selectedAccounts.includes(accountId)) {
      setSelectedAccounts([...selectedAccounts, accountId]);
    }
  };

  const handleRemoveAccount = (accountId: string) => {
    setSelectedAccounts(selectedAccounts.filter(id => id !== accountId));
  };

  const handleAddItem = () => {
    const newItem: BudgetItem = {
      name: `New Item ${items.length + 1}`,
      categories: [],
      percentage: 0,
    };
    setItems([...items, newItem]);
  };

  const handleRemoveItem = (index: number) => {
    setItems(items.filter((_, i) => i !== index));
  };

  const handleUpdateItem = (index: number, updates: Partial<BudgetItem>) => {
    setItems(items.map((item, i) => i === index ? { ...item, ...updates } : item));
  };

  const handleSaveBudget = async () => {
    if (isSaving) return;
    setIsSaving(true);
    try {
      await api.saveBudget({
        id: budgetId,
        name: budgetName,
        data: {
          items,
          selectedAccounts,
          otherPercentage,
        }
      });
      setSuccess(true);
      setTimeout(() => {
        setSuccess(false);
        fetchBudgets();
        setView('list');
      }, 2000);
    } catch (err) {
      console.error('Error saving budget:', err);
      setError(true);
      setErrorMessage('Failed to save budget');
      setTimeout(() => setError(false), 3500);
    } finally {
      setIsSaving(false);
    }
  };

  const handleDeleteBudget = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (!confirm('Are you sure you want to delete this budget?')) return;
    try {
      await api.deleteBudget(id);
      fetchBudgets();
    } catch (err) {
      console.error('Error deleting budget:', err);
      setError(true);
      setErrorMessage('Failed to delete budget');
      setTimeout(() => setError(false), 3500);
    }
  };

  const handleEditBudget = (budget: BudgetData) => {
    setBudgetId(budget.id);
    setBudgetName(budget.name);
    setItems(budget.data.items || []);
    setSelectedAccounts(budget.data.selectedAccounts || []);
    setOtherPercentage(budget.data.otherPercentage || 0);
    setView('edit');
  };

  const handleViewBudget = (budget: BudgetData) => {
    setSelectedBudget(budget);
    setView('view');
  };

  const handleCreateNew = () => {
    setBudgetId(generateId());
    setBudgetName('New Budget');
    setItems([]);
    setSelectedAccounts([]);
    setOtherPercentage(0);
    setView('edit');
  };

  const handleAddCategoryToItem = (index: number, category: string) => {
    setItems(items.map((item, i) => {
      if (i === index && !item.categories.includes(category)) {
        return { ...item, categories: [...item.categories, category] };
      }
      return item;
    }));
  };

  const handleRemoveCategoryFromItem = (index: number, category: string) => {
    setItems(items.map((item, i) => {
      if (i === index) {
        return { ...item, categories: item.categories.filter(c => c !== category) };
      }
      return item;
    }));
  };

  const assignedCategories = useMemo(() => {
    const assigned = new Set<string>();
    items.forEach(item => {
      item.categories.forEach(c => assigned.add(c));
    });
    return Array.from(assigned);
  }, [items]);

  const unassignedCategories = useMemo(() => {
    return categories.filter(c => !assignedCategories.includes(c));
  }, [categories, assignedCategories]);

  const totalPercentage = items.reduce((sum, item) => sum + (item.percentage || 0), 0);

  if (view === 'list') {
    return (
      <div className="w-full flex items-start justify-center p-4 md:p-8 pb-20 bg-transparent">
        <div className="w-full max-w-4xl space-y-8">
          <div className="flex items-center justify-between px-4">
            <div>
              <h1 className="text-4xl font-black text-text-on-dark tracking-tight">Budgets</h1>
              <p className="text-text-on-dark/40 font-bold uppercase tracking-widest text-[10px] mt-2">Manage your spending plans</p>
            </div>
            <button 
              onClick={handleCreateNew}
              className="flex items-center gap-2 bg-accent hover:bg-accent/90 text-white font-bold px-6 py-3 rounded-2xl transition-all shadow-lg hover:shadow-xl active:scale-[0.98]"
            >
              <PlusCircle size={20} />
              Create New
            </button>
          </div>

          {isLoading ? (
            <div className="flex flex-col items-center justify-center py-20 bg-card/5 rounded-[48px] border border-white/5">
              <Loader2 className="animate-spin text-accent mb-4" size={48} />
              <p className="text-text-on-dark/40 font-bold uppercase tracking-widest text-xs">Loading budgets...</p>
            </div>
          ) : budgets.length === 0 ? (
            <div className="bg-card/5 rounded-[48px] border border-dashed border-white/5 p-20 flex flex-col items-center text-center">
              <div className="w-20 h-20 bg-glass-dark rounded-full flex items-center justify-center mb-6">
                <Layout size={32} className="text-text-on-dark/20" />
              </div>
              <h3 className="text-xl font-bold text-text-on-dark mb-2">No budgets found</h3>
              <p className="text-text-on-dark/40 max-w-xs mb-8">You haven't created any budget plans yet. Start by creating your first one!</p>
              <button 
                onClick={handleCreateNew}
                className="text-accent font-bold hover:text-accent flex items-center gap-2"
              >
                <PlusCircle size={18} />
                Create your first budget
              </button>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {budgets.map((budget) => (
                <div 
                  key={budget.id}
                  onClick={() => handleViewBudget(budget)}
                  className="bg-card rounded-[40px] p-8 border border-border-pearl shadow-lg hover:shadow-xl transition-all cursor-pointer group relative overflow-hidden"
                >
                  <div className="flex justify-between items-start mb-6">
                    <div>
                      <h3 className="text-2xl font-black text-text-main tracking-tight">{budget.name}</h3>
                      <p className="text-text-muted font-bold uppercase tracking-widest text-[10px] mt-1">
                        {budget.data.items?.length || 0} items • {budget.data.selectedAccounts?.length || 0} accounts
                      </p>
                    </div>
                    <div className="flex items-center gap-1">
                      <button 
                        onClick={(e) => {
                          e.stopPropagation();
                          handleEditBudget(budget);
                        }}
                        className="text-text-muted/20 hover:text-accent transition-colors p-2"
                        title="Edit Budget"
                      >
                        <Pencil size={20} />
                      </button>
                      <button 
                        onClick={(e) => handleDeleteBudget(budget.id, e)}
                        className="text-text-muted/20 hover:text-negative transition-colors p-2"
                        title="Delete Budget"
                      >
                        <Trash2 size={20} />
                      </button>
                    </div>
                  </div>

                  <div className="space-y-3">
                    {budget.data.items?.slice(0, 3).map((item, i) => (
                      <div key={i} className="flex items-center justify-between text-sm">
                        <span className="text-text-muted font-medium">{item.name}</span>
                        <span className="text-text-muted/40 font-bold">{item.percentage}%</span>
                      </div>
                    ))}
                    {budget.data.items?.length > 3 && (
                      <div className="text-[10px] font-bold text-accent uppercase tracking-widest pt-2">
                        + {budget.data.items.length - 3} more items
                      </div>
                    )}
                  </div>

                  <div className="absolute bottom-0 left-0 h-1.5 bg-accent/10 w-full">
                    <div 
                      className="h-full bg-accent transition-all duration-1000" 
                      style={{ width: `${Math.min(100, (budget.data.items?.reduce((s, i) => s + (i.percentage || 0), 0) || 0))}%` }} 
                    />
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    );
  }

  if (view === 'view' && selectedBudget) {
    return (
      <div className="w-full flex items-start justify-center p-4 md:p-8 pb-20 bg-transparent">
        <div className="w-full max-w-4xl space-y-8">
          <div className="flex items-center justify-between px-4">
            <div className="flex flex-col">
              <button 
                onClick={() => setView('list')}
                className="text-text-on-dark/30 hover:text-accent transition-colors flex items-center gap-1 font-bold uppercase tracking-widest text-[10px] mb-2"
              >
                <ChevronLeft size={16} strokeWidth={3} />
                Back to List
              </button>
              <h1 className="text-4xl font-black text-text-on-dark tracking-tight">{selectedBudget.name}</h1>
            </div>
            
            <div className="flex items-center gap-3">
              <div className="flex bg-glass backdrop-blur-xl rounded-xl border border-glass-border p-1 overflow-hidden">
                <button
                  onClick={() => setTimeframe('month')}
                  className={`px-4 py-1.5 rounded-lg text-[10px] font-black uppercase tracking-widest transition-all ${timeframe === 'month' ? 'bg-accent text-white shadow-lg' : 'text-text-on-dark/40 hover:text-text-on-dark'}`}
                >
                  Month
                </button>
                <button
                  onClick={() => setTimeframe('year')}
                  className={`px-4 py-1.5 rounded-lg text-[10px] font-black uppercase tracking-widest transition-all ${timeframe === 'year' ? 'bg-accent text-white shadow-lg' : 'text-text-on-dark/40 hover:text-text-on-dark'}`}
                >
                  Year
                </button>
              </div>
              
              <div className="relative">
                <button 
                  onClick={handleMonthClick}
                  className="p-2.5 bg-glass backdrop-blur-xl rounded-xl shadow-lg hover:bg-glass-hover transition-all text-text-on-dark border border-glass-border active:scale-95 flex items-center gap-2 px-4"
                >
                  <Calendar size={18} strokeWidth={2.5} />
                  <span className="text-xs font-bold uppercase tracking-widest">
                    {new Date(selectedMonth + '-01').toLocaleDateString(undefined, { month: 'long', year: 'numeric' })}
                  </span>
                </button>
                <input 
                  ref={monthInputRef}
                  type="month" 
                  value={selectedMonth}
                  onChange={(e) => setSelectedMonth(e.target.value)}
                  className="absolute inset-0 opacity-0 pointer-events-none"
                />
              </div>
            </div>
          </div>

          <div className="bg-card rounded-[48px] p-10 border border-border-pearl shadow-2xl relative overflow-hidden">
            {isFetchingTransactions ? (
              <div className="flex flex-col items-center justify-center py-20">
                <Loader2 className="animate-spin text-accent mb-4" size={48} />
                <p className="text-text-muted font-bold uppercase tracking-widest text-xs">Analyzing transactions...</p>
              </div>
            ) : (
              <div className="flex flex-col gap-12">
                {/* 1. Total Income */}
                <div className="cursor-pointer group/stat" onClick={() => toggleSection('income')}>
                  <div className="flex items-center justify-between mb-2">
                    <div className="text-[10px] font-black text-text-muted uppercase tracking-widest">Total Income</div>
                    {expandedSections.includes('income') ? <ChevronDown size={14} className="text-text-muted" /> : <ChevronRight size={14} className="text-text-muted/20 group-hover/stat:text-accent transition-colors" />}
                  </div>
                  <div className="text-5xl font-black text-black tracking-tighter">
                    + {budgetStats.totalIncome?.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                  </div>
                  {expandedSections.includes('income') && <TransactionList transactions={budgetStats.incomeTransactions || []} accounts={accounts} />}
                </div>

                {/* Total Investments */}
                <div className="cursor-pointer group/stat text-right" onClick={() => toggleSection('investments')}>
                  <div className="flex items-center justify-end gap-2 mb-1">
                    <div className="text-[8px] font-black text-accent uppercase tracking-widest">Total Investments</div>
                    {expandedSections.includes('investments') ? <ChevronDown size={12} className="text-accent" /> : <ChevronRight size={12} className="text-accent/20 group-hover/stat:text-accent transition-colors" />}
                  </div>
                  <div className="text-2xl font-black text-accent tracking-tighter">
                    {budgetStats.totalInvestments?.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                  </div>
                  {expandedSections.includes('investments') && <div className="text-left"><TransactionList transactions={budgetStats.investmentTransactions || []} accounts={accounts} /></div>}
                </div>

                {/* 2. Total Outcome with Expenses Breakdown */}
                <div className="cursor-pointer group/stat" onClick={() => toggleSection('breakdown')}>
                  <div className="flex items-center justify-between mb-2">
                    <div className="text-[10px] font-black text-text-muted uppercase tracking-widest">Total Outcome</div>
                    {expandedSections.includes('breakdown') ? <ChevronDown size={14} className="text-text-muted" /> : <ChevronRight size={14} className="text-text-muted/20 group-hover/stat:text-accent transition-colors" />}
                  </div>
                  <div className="flex items-center justify-end mb-4">
                    <div className="text-5xl font-black text-black tracking-tighter">
                      - {budgetStats.totalOutcome?.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </div>
                  </div>
                  
                  {expandedSections.includes('breakdown') && (
                    <div className="space-y-4 pt-6 border-t border-border-pearl animate-in fade-in slide-in-from-top-2 duration-300" onClick={(e) => e.stopPropagation()}>
                      <div className="text-[10px] font-black text-text-muted px-4 pb-2">Expenses Breakdown</div>
                      <div className="space-y-4">
                        {budgetStats.items?.map((item, i) => renderBudgetItemRow(item, `item-${i}`))}
                        {(budgetStats.otherItem.actualSpent > 0 || budgetStats.otherItem.transactions.length > 0) && renderBudgetItemRow(budgetStats.otherItem, 'item-others')}
                        {budgetStats.items?.length === 0 && budgetStats.otherItem.actualSpent === 0 && (
                            <div className="text-center py-10 text-text-muted/40 italic">No expenses recorded</div>
                        )}
                      </div>
                    </div>
                  )}
                </div>

                {/* 3. Net Gain/Loss */}
                <div className="pt-10 border-t border-border-pearl">
                  <div className="text-[10px] font-black text-text-muted uppercase tracking-widest mb-2">Net Gain/Loss</div>
                  {(() => {
                    const netGain = budgetStats.totalIncome - budgetStats.totalOutcome;

                    return (
                      <div className={`text-7xl font-black tracking-tighter ${netGain > 0 ? 'text-accent-secondary' : 'text-black'}`}>
                        {netGain.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                      </div>
                    );
                  })()}
                </div>

                {/* 4. Distribution Chart */}
                <div className="flex flex-col items-center">
                   <div className="w-full max-w-xl aspect-square p-8">
                     <Pie 
                        data={{
                          labels: [
                            ...budgetStats.items.map(i => i.name),
                            budgetStats.otherItem.name,
                            'Net Savings'
                          ],
                          datasets: [{
                            data: [
                              ...budgetStats.items.map(i => i.actualSpent),
                              budgetStats.otherItem.actualSpent,
                              Math.max(0, budgetStats.totalIncome - budgetStats.totalOutcome)
                            ],
                            backgroundColor: [
                              ...budgetStats.items.map((_, i) => {
                                const varName = `--color-chart-${(i % 12) + 1}`;
                                return getCSSVariableValue(varName);
                              }),
                              getCSSVariableValue('--color-chart-others'), // Specific green for 'Others'
                              getCSSVariableValue('--color-accent-secondary') // Yellow for Net Savings
                            ],
                            borderWidth: 0,
                          }]
                        }}
                        options={{
                          layout: {
                            padding: 40
                          },
                          plugins: {
                            legend: {
                              display: false,
                            },
                            datalabels: {
                              color: getCSSVariableValue('--color-text-main'),
                              anchor: 'end',
                              align: 'end',
                              offset: 10,
                              display: 'auto',
                              font: {
                                size: 9,
                                weight: 'bold',
                                family: 'Inter, system-ui, sans-serif'
                              },
                              formatter: (value, context) => {
                                const total = (context.dataset.data as number[]).reduce((a, b) => a + b, 0);
                                const percentage = (value / total) * 100;
                                if (percentage < 1) return null;
                                
                                const label = context.chart.data.labels![context.dataIndex] as string;
                                return `${label}\n${percentage.toFixed(1)}%`;
                              },
                              textAlign: 'center'
                            }
                          }
                        }}
                     />
                   </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full flex items-start justify-center p-4 md:p-8 pb-20 bg-transparent">
      <div className="bg-card rounded-[48px] shadow-2xl border border-border-pearl p-8 md:p-10 w-full max-w-4xl flex flex-col items-stretch my-8 relative overflow-hidden group">
        
        {/* Success Overlay */}
        <div className={`absolute inset-0 z-50 flex flex-col items-center justify-center transition-all duration-700 bg-card/95 ${success ? 'translate-y-0 opacity-100' : 'translate-y-full opacity-0 pointer-events-none'}`}>
            <div className="w-24 h-24 bg-accent text-white rounded-full flex items-center justify-center mb-6 shadow-2xl animate-bounce">
              <Check size={48} strokeWidth={4} />
            </div>
            <h3 className="text-4xl font-black text-text-main tracking-tight">Saved!</h3>
            <p className="text-text-muted font-bold mt-2 uppercase tracking-widest text-xs">Budget Updated</p>
        </div>

        {/* Error Overlay */}
        <div className={`absolute inset-0 z-50 flex flex-col items-center justify-center transition-all duration-700 bg-card/95 ${error ? 'translate-y-0 opacity-100' : 'translate-y-full opacity-0 pointer-events-none'}`}>
            <div className="w-24 h-24 bg-negative text-white rounded-full flex items-center justify-center mb-6 shadow-2xl animate-bounce">
              <XCircle size={48} strokeWidth={4} />
            </div>
            <h3 className="text-4xl font-black text-text-main tracking-tight">Error</h3>
            <p className="text-negative font-bold mt-2 uppercase tracking-widest text-xs text-center px-10">{errorMessage}</p>
        </div>

        <div className="mb-12 relative">
          <button 
            onClick={() => setView('list')}
            className="absolute -top-4 -left-4 p-4 text-text-muted hover:text-accent transition-colors flex items-center gap-1 font-bold uppercase tracking-widest text-[10px]"
          >
            <ChevronLeft size={16} strokeWidth={3} />
            Back to List
          </button>
          
          <div className="text-center pt-8">
            <input
              type="text"
              value={budgetName}
              onChange={(e) => setBudgetName(e.target.value)}
              className="text-4xl font-black text-text-main tracking-tighter text-center bg-transparent focus:outline-none border-b-2 border-transparent focus:border-accent pb-2 w-full max-w-lg"
              placeholder="Budget Name"
            />
            <p className="text-text-muted font-bold uppercase tracking-widest text-[10px] mt-4">Personal Spending Plan</p>
          </div>
        </div>
        
        {/* Accounts Section */}
        <div className="mb-10 bg-card-muted p-6 rounded-3xl border border-border-pearl">
          <h3 className="text-sm font-black text-text-muted uppercase tracking-widest mb-4">Accounts to Consider</h3>
          <div className="flex flex-wrap gap-2 mb-4">
            {selectedAccounts.map(accId => {
              const acc = accounts.find(a => a.id === accId);
              return (
                <div key={accId} className="flex items-center gap-2 bg-accent/10 text-accent px-3 py-1.5 rounded-xl text-sm font-semibold">
                  <span>{acc?.name || accId}</span>
                  <button onClick={() => handleRemoveAccount(accId)} className="hover:text-accent transition-colors">
                    <X size={14} strokeWidth={3} />
                  </button>
                </div>
              );
            })}
            {selectedAccounts.length === 0 && (
              <span className="text-text-muted/40 text-sm font-medium italic">No accounts selected</span>
            )}
          </div>
          
          <select 
            className="w-full md:w-auto bg-card border border-border-pearl text-text-main font-semibold rounded-2xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-accent/10 focus:border-accent transition-all shadow-sm"
            onChange={(e) => handleAddAccount(e.target.value)}
            value=""
          >
            <option value="" disabled>+ Add Account</option>
            {accounts.filter(acc => !selectedAccounts.includes(acc.id)).map(acc => (
              <option key={acc.id} value={acc.id}>{acc.name}</option>
            ))}
          </select>
        </div>

        {/* Budget Items Section */}
        <div className="space-y-6 mb-8">
          <div className="flex items-center justify-between px-2">
            <h3 className="text-sm font-black text-text-muted uppercase tracking-widest">Budget Items</h3>
            <span className={`text-sm font-bold ${totalPercentage > 100 ? 'text-negative' : 'text-text-muted/40'}`}>
              Total: {totalPercentage}%
            </span>
          </div>

          {items.map((item, index) => (
            <div key={index} className="bg-card border border-border-pearl shadow-sm rounded-3xl p-6 flex flex-col gap-4 relative group">
              <button 
                onClick={() => handleRemoveItem(index)}
                className="absolute top-6 right-6 text-text-muted/20 hover:text-negative transition-colors"
                title="Remove Item"
              >
                <Trash2 size={20} />
              </button>
              
              <div className="flex flex-col md:flex-row md:items-center gap-4 pr-8">
                <div className="flex-1">
                  <input
                    type="text"
                    value={item.name}
                    onChange={(e) => handleUpdateItem(index, { name: e.target.value })}
                    className="w-full bg-transparent text-xl font-bold text-text-main focus:outline-none focus:border-b-2 border-border-pearl focus:border-accent transition-colors pb-1"
                    placeholder="Item Name (e.g., Transport)"
                  />
                </div>
                <div className="flex items-center gap-2">
                  <input
                    type="number"
                    min="0"
                    max="100"
                    value={item.percentage || ''}
                    onChange={(e) => handleUpdateItem(index, { percentage: parseFloat(e.target.value) || 0 })}
                    className="w-20 bg-card-muted border border-border-pearl text-text-main font-bold rounded-xl px-3 py-2 text-right focus:outline-none focus:ring-2 focus:ring-accent/10 focus:border-accent transition-all"
                    placeholder="0"
                  />
                  <span className="text-text-muted/40 font-bold">%</span>
                </div>
              </div>

              <div className="mt-2">
                <div className="text-[10px] font-black text-text-muted uppercase tracking-widest mb-3">Categories</div>
                <div className="flex flex-wrap gap-2">
                  {item.categories.map(cat => (
                    <div key={cat} className="flex items-center gap-1.5 bg-card-muted text-text-main px-3 py-1.5 rounded-xl text-sm font-medium">
                      <span>{cat}</span>
                      <button onClick={() => handleRemoveCategoryFromItem(index, cat)} className="hover:text-accent transition-colors">
                        <X size={14} strokeWidth={2.5} />
                      </button>
                    </div>
                  ))}
                  
                  <select 
                    className="bg-card-muted border border-border-pearl text-text-main font-medium rounded-xl px-3 py-1.5 text-sm focus:outline-none focus:border-accent transition-all"
                    onChange={(e) => {
                      if (e.target.value) {
                        handleAddCategoryToItem(index, e.target.value);
                      }
                    }}
                    value=""
                  >
                    <option value="" disabled>+ Add Category</option>
                    {unassignedCategories.map(cat => (
                      <option key={cat} value={cat}>{cat}</option>
                    ))}
                  </select>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Other Item */}
        <div className="bg-card-muted border border-dashed border-border-pearl rounded-3xl p-6 flex flex-col gap-4 mb-4">
          <div className="flex flex-col md:flex-row md:items-center gap-4">
            <div className="flex-1">
              <div className="text-xl font-bold text-text-muted">Others</div>
              <p className="text-[10px] font-bold text-text-muted uppercase tracking-widest mt-1">Budget for all remaining categories</p>
            </div>
            <div className="flex items-center gap-2">
              <input
                type="number"
                min="0"
                max="100"
                value={otherPercentage || ''}
                onChange={(e) => setOtherPercentage(parseFloat(e.target.value) || 0)}
                className="w-20 bg-card border border-border-pearl text-text-main font-bold rounded-xl px-3 py-2 text-right focus:outline-none focus:ring-2 focus:ring-accent/10 focus:border-accent transition-all"
                placeholder="0"
              />
              <span className="text-text-muted/40 font-bold">%</span>
            </div>
          </div>
          <div className="mt-2">
            <div className="text-[10px] font-black text-text-muted uppercase tracking-widest mb-3">Remaining Categories</div>
            <div className="flex flex-wrap gap-2">
              {unassignedCategories.length > 0 ? (
                unassignedCategories.map(cat => (
                  <div key={cat} className="flex items-center gap-1.5 bg-card text-text-main px-3 py-1.5 rounded-xl text-sm font-medium">
                    {cat}
                  </div>
                ))
              ) : (
                <span className="text-text-muted/40 text-sm font-medium italic">None</span>
              )}
            </div>
          </div>
        </div>

        {/* Expected Savings (Calculated) */}
        <div className="bg-accent/5 border border-accent/10 rounded-3xl p-6 flex flex-col gap-1 mb-8">
          <div className="text-[10px] font-black text-accent/60 uppercase tracking-widest">Calculated Expected Savings</div>
          <div className="flex items-baseline gap-2">
            <div className="text-4xl font-black text-accent">
              {(100 - items.reduce((sum, i) => sum + (i.percentage || 0), 0) - (otherPercentage || 0)).toFixed(1)}%
            </div>
            <div className="text-sm font-bold text-accent/40">of total income</div>
          </div>
        </div>

        <button 
          onClick={handleAddItem}
          className="flex items-center justify-center gap-2 w-full bg-accent hover:bg-accent/90 text-white font-bold rounded-2xl py-4 transition-all shadow-lg hover:shadow-xl active:scale-[0.98] mb-4"
        >
          <PlusCircle size={20} />
          Add Budget Item
        </button>

        <button 
          onClick={handleSaveBudget}
          disabled={isSaving}
          className="flex items-center justify-center gap-2 w-full bg-accent hover:bg-accent/90 text-white font-bold rounded-2xl py-4 transition-all shadow-lg hover:shadow-xl active:scale-[0.98] disabled:bg-card-muted disabled:text-text-muted/20"
        >
          {isSaving ? (
            <Loader2 className="animate-spin" size={20} strokeWidth={4} />
          ) : (
            <>
              <Save size={20} />
              Save Budget
            </>
          )}
        </button>
      </div>
    </div>
  );
};

export default Budget;
