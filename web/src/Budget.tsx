import React, { useState, useEffect, useMemo } from 'react';
import { Trash2, X, PlusCircle, Save, Check, Loader2, ChevronLeft, Layout } from 'lucide-react';
import { api } from './api';
import type { Account } from './api';

const AVAILABLE_CATEGORIES = [
  'Groceries', 'Dining Out', 'Rent', 'Utilities', 
  'Metro/Bus', 'Train', 'Gas', 'Entertainment', 
  'Healthcare', 'Insurance', 'Shopping', 'Travel',
  'Subscriptions', 'Personal Care', 'Education'
];

interface BudgetItem {
  name: string;
  categories: string[];
  percentage: number;
}

interface BudgetData {
  id: string;
  name: string;
  data: {
    items: BudgetItem[];
    selectedAccounts: string[];
  };
}

const Budget: React.FC = () => {
  const [view, setView] = useState<'list' | 'edit'>('list');
  const [budgets, setBudgets] = useState<BudgetData[]>([]);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [selectedAccounts, setSelectedAccounts] = useState<string[]>([]);
  const [items, setItems] = useState<BudgetItem[]>([]);
  const [budgetName, setBudgetName] = useState('Monthly Budget');
  const [budgetId, setBudgetId] = useState<string>(crypto.randomUUID());
  const [isSaving, setIsSaving] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [success, setSuccess] = useState(false);

  const fetchBudgets = async () => {
    setIsLoading(true);
    try {
      const data = await api.getBudgets();
      setBudgets(data || []);
    } catch (err) {
      console.error('Error fetching budgets:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const fetchAccounts = async () => {
      try {
        const accs = await api.getAccounts();
        setAccounts(accs || []);
      } catch (err) {
        console.error('Error fetching accounts:', err);
      }
    };
    fetchAccounts();
    fetchBudgets();
  }, []);

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
      alert('Failed to save budget');
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
      alert('Failed to delete budget');
    }
  };

  const handleEditBudget = (budget: BudgetData) => {
    setBudgetId(budget.id);
    setBudgetName(budget.name);
    setItems(budget.data.items || []);
    setSelectedAccounts(budget.data.selectedAccounts || []);
    setView('edit');
  };

  const handleCreateNew = () => {
    setBudgetId(crypto.randomUUID());
    setBudgetName('New Budget');
    setItems([]);
    setSelectedAccounts([]);
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
    return AVAILABLE_CATEGORIES.filter(c => !assignedCategories.includes(c));
  }, [assignedCategories]);

  const totalPercentage = items.reduce((sum, item) => sum + (item.percentage || 0), 0);
  const otherPercentage = Math.max(0, 100 - totalPercentage);

  if (view === 'list') {
    return (
      <div className="w-full flex items-start justify-center p-4 md:p-8 pb-20">
        <div className="w-full max-w-4xl space-y-8">
          <div className="flex items-center justify-between px-4">
            <div>
              <h1 className="text-4xl font-black text-gray-900 tracking-tight">Budgets</h1>
              <p className="text-gray-400 font-bold uppercase tracking-widest text-[10px] mt-2">Manage your spending plans</p>
            </div>
            <button 
              onClick={handleCreateNew}
              className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white font-bold px-6 py-3 rounded-2xl transition-all shadow-lg hover:shadow-xl active:scale-[0.98]"
            >
              <PlusCircle size={20} />
              Create New
            </button>
          </div>

          {isLoading ? (
            <div className="flex flex-col items-center justify-center py-20 bg-white/50 backdrop-blur-xl rounded-[48px] border border-white/50">
              <Loader2 className="animate-spin text-indigo-600 mb-4" size={48} />
              <p className="text-gray-400 font-bold uppercase tracking-widest text-xs">Loading budgets...</p>
            </div>
          ) : budgets.length === 0 ? (
            <div className="bg-white/50 backdrop-blur-xl rounded-[48px] border border-dashed border-gray-200 p-20 flex flex-col items-center text-center">
              <div className="w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mb-6">
                <Layout size={32} className="text-gray-400" />
              </div>
              <h3 className="text-xl font-bold text-gray-900 mb-2">No budgets found</h3>
              <p className="text-gray-400 max-w-xs mb-8">You haven't created any budget plans yet. Start by creating your first one!</p>
              <button 
                onClick={handleCreateNew}
                className="text-indigo-600 font-bold hover:text-indigo-700 flex items-center gap-2"
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
                  onClick={() => handleEditBudget(budget)}
                  className="bg-white/90 backdrop-blur-2xl rounded-[40px] p-8 border border-white/50 shadow-sm hover:shadow-xl transition-all cursor-pointer group relative overflow-hidden"
                >
                  <div className="flex justify-between items-start mb-6">
                    <div>
                      <h3 className="text-2xl font-black text-gray-900 tracking-tight">{budget.name}</h3>
                      <p className="text-gray-400 font-bold uppercase tracking-widest text-[10px] mt-1">
                        {budget.data.items?.length || 0} items • {budget.data.selectedAccounts?.length || 0} accounts
                      </p>
                    </div>
                    <button 
                      onClick={(e) => handleDeleteBudget(budget.id, e)}
                      className="text-gray-300 hover:text-red-500 transition-colors p-2"
                    >
                      <Trash2 size={20} />
                    </button>
                  </div>

                  <div className="space-y-3">
                    {budget.data.items?.slice(0, 3).map((item, i) => (
                      <div key={i} className="flex items-center justify-between text-sm">
                        <span className="text-gray-600 font-medium">{item.name}</span>
                        <span className="text-gray-400 font-bold">{item.percentage}%</span>
                      </div>
                    ))}
                    {budget.data.items?.length > 3 && (
                      <div className="text-[10px] font-bold text-indigo-400 uppercase tracking-widest pt-2">
                        + {budget.data.items.length - 3} more items
                      </div>
                    )}
                  </div>

                  <div className="absolute bottom-0 left-0 h-1.5 bg-indigo-500/10 w-full">
                    <div 
                      className="h-full bg-indigo-500 transition-all duration-1000" 
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

  return (
    <div className="w-full flex items-start justify-center p-4 md:p-8 pb-20">
      <div className="bg-white/90 backdrop-blur-2xl rounded-[48px] shadow-[0_40px_100px_-20px_rgba(0,0,0,0.15)] border border-white/50 p-8 md:p-10 w-full max-w-4xl flex flex-col items-stretch my-8 relative overflow-hidden group">
        
        {/* Success Overlay */}
        <div className={`absolute inset-0 z-50 flex flex-col items-center justify-center transition-all duration-700 bg-white/95 ${success ? 'translate-y-0 opacity-100' : 'translate-y-full opacity-0 pointer-events-none'}`}>
            <div className={`w-24 h-24 bg-indigo-600 text-white rounded-full flex items-center justify-center mb-6 shadow-2xl animate-bounce`}>
              <Check size={48} strokeWidth={4} />
            </div>
            <h3 className="text-4xl font-black text-gray-900 tracking-tight">Saved!</h3>
            <p className="text-gray-400 font-bold mt-2 uppercase tracking-widest text-xs">Budget Updated</p>
        </div>

        <div className="mb-12 relative">
          <button 
            onClick={() => setView('list')}
            className="absolute -top-4 -left-4 p-4 text-gray-400 hover:text-indigo-600 transition-colors flex items-center gap-1 font-bold uppercase tracking-widest text-[10px]"
          >
            <ChevronLeft size={16} strokeWidth={3} />
            Back to List
          </button>
          
          <div className="text-center pt-8">
            <input
              type="text"
              value={budgetName}
              onChange={(e) => setBudgetName(e.target.value)}
              className="text-4xl font-black text-gray-900 tracking-tighter text-center bg-transparent focus:outline-none border-b-2 border-transparent focus:border-indigo-400 pb-2 w-full max-w-lg"
              placeholder="Budget Name"
            />
            <p className="text-gray-400 font-bold uppercase tracking-widest text-[10px] mt-4">Personal Spending Plan</p>
          </div>
        </div>
        
        {/* Accounts Section */}
        <div className="mb-10 bg-gray-50/50 p-6 rounded-3xl border border-gray-100">
          <h3 className="text-sm font-black text-gray-400 uppercase tracking-widest mb-4">Accounts to Consider</h3>
          <div className="flex flex-wrap gap-2 mb-4">
            {selectedAccounts.map(accId => {
              const acc = accounts.find(a => a.id === accId);
              return (
                <div key={accId} className="flex items-center gap-2 bg-indigo-100 text-indigo-700 px-3 py-1.5 rounded-xl text-sm font-semibold">
                  <span>{acc?.name || accId}</span>
                  <button onClick={() => handleRemoveAccount(accId)} className="hover:text-indigo-900 transition-colors">
                    <X size={14} strokeWidth={3} />
                  </button>
                </div>
              );
            })}
            {selectedAccounts.length === 0 && (
              <span className="text-gray-400 text-sm font-medium italic">No accounts selected</span>
            )}
          </div>
          
          <select 
            className="w-full md:w-auto bg-white border border-gray-200 text-gray-700 font-semibold rounded-2xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all shadow-sm"
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
            <h3 className="text-sm font-black text-gray-400 uppercase tracking-widest">Budget Items</h3>
            <span className={`text-sm font-bold ${totalPercentage > 100 ? 'text-red-500' : 'text-gray-500'}`}>
              Total: {totalPercentage}%
            </span>
          </div>

          {items.map((item, index) => (
            <div key={index} className="bg-white border border-gray-100 shadow-sm rounded-3xl p-6 flex flex-col gap-4 relative group">
              <button 
                onClick={() => handleRemoveItem(index)}
                className="absolute top-6 right-6 text-gray-300 hover:text-red-500 transition-colors"
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
                    className="w-full bg-transparent text-xl font-bold text-gray-800 focus:outline-none focus:border-b-2 border-gray-200 focus:border-indigo-400 transition-colors pb-1"
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
                    className="w-20 bg-gray-50 border border-gray-200 text-gray-700 font-bold rounded-xl px-3 py-2 text-right focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
                    placeholder="0"
                  />
                  <span className="text-gray-400 font-bold">%</span>
                </div>
              </div>

              <div className="mt-2">
                <div className="text-[10px] font-black text-gray-300 uppercase tracking-widest mb-3">Categories</div>
                <div className="flex flex-wrap gap-2">
                  {item.categories.map(cat => (
                    <div key={cat} className="flex items-center gap-1.5 bg-gray-100 text-gray-600 px-3 py-1.5 rounded-xl text-sm font-medium">
                      <span>{cat}</span>
                      <button onClick={() => handleRemoveCategoryFromItem(index, cat)} className="hover:text-gray-900 transition-colors">
                        <X size={14} strokeWidth={2.5} />
                      </button>
                    </div>
                  ))}
                  
                  <select 
                    className="bg-gray-50 border border-gray-200 text-gray-600 font-medium rounded-xl px-3 py-1.5 text-sm focus:outline-none focus:border-indigo-400 transition-all"
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

          {/* Other Item (Read-only) */}
          <div className="bg-gray-50 border border-dashed border-gray-200 rounded-3xl p-6 flex flex-col gap-4 opacity-75">
            <div className="flex flex-col md:flex-row md:items-center gap-4">
              <div className="flex-1">
                <div className="text-xl font-bold text-gray-500">Other</div>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-20 bg-gray-200/50 text-gray-500 font-bold rounded-xl px-3 py-2 text-right">
                  {otherPercentage}
                </div>
                <span className="text-gray-400 font-bold">%</span>
              </div>
            </div>
            <div className="mt-2">
              <div className="text-[10px] font-black text-gray-400 uppercase tracking-widest mb-3">Remaining Categories</div>
              <div className="flex flex-wrap gap-2">
                {unassignedCategories.length > 0 ? (
                  unassignedCategories.map(cat => (
                    <div key={cat} className="flex items-center gap-1.5 bg-gray-200/50 text-gray-500 px-3 py-1.5 rounded-xl text-sm font-medium">
                      {cat}
                    </div>
                  ))
                ) : (
                  <span className="text-gray-400 text-sm font-medium italic">None</span>
                )}
              </div>
            </div>
          </div>
        </div>

        <button 
          onClick={handleAddItem}
          className="flex items-center justify-center gap-2 w-full bg-gray-900 hover:bg-gray-800 text-white font-bold rounded-2xl py-4 transition-all shadow-lg hover:shadow-xl active:scale-[0.98] mb-4"
        >
          <PlusCircle size={20} />
          Add Budget Item
        </button>

        <button 
          onClick={handleSaveBudget}
          disabled={isSaving}
          className="flex items-center justify-center gap-2 w-full bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-2xl py-4 transition-all shadow-lg hover:shadow-xl active:scale-[0.98] disabled:bg-gray-200"
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
