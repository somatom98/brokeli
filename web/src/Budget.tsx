import React, { useState, useEffect, useMemo } from 'react';
import { Trash2, X, PlusCircle } from 'lucide-react';
import { api } from './api';
import type { Account } from './api';

const AVAILABLE_CATEGORIES = [
  'Groceries', 'Dining Out', 'Rent', 'Utilities', 
  'Metro/Bus', 'Train', 'Gas', 'Entertainment', 
  'Healthcare', 'Insurance', 'Shopping', 'Travel',
  'Subscriptions', 'Personal Care', 'Education'
];

interface BudgetItem {
  id: string;
  name: string;
  categories: string[];
  percentage: number;
}

const Budget: React.FC = () => {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [selectedAccounts, setSelectedAccounts] = useState<string[]>([]);
  const [items, setItems] = useState<BudgetItem[]>([]);
  const [budgetName, setBudgetName] = useState('Monthly Budget');

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
      id: Math.random().toString(36).substring(2, 9),
      name: 'New Item',
      categories: [],
      percentage: 0,
    };
    setItems([...items, newItem]);
  };

  const handleRemoveItem = (id: string) => {
    setItems(items.filter(item => item.id !== id));
  };

  const handleUpdateItem = (id: string, updates: Partial<BudgetItem>) => {
    setItems(items.map(item => item.id === id ? { ...item, ...updates } : item));
  };

  const handleAddCategoryToItem = (itemId: string, category: string) => {
    setItems(items.map(item => {
      if (item.id === itemId && !item.categories.includes(category)) {
        return { ...item, categories: [...item.categories, category] };
      }
      return item;
    }));
  };

  const handleRemoveCategoryFromItem = (itemId: string, category: string) => {
    setItems(items.map(item => {
      if (item.id === itemId) {
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

  return (
    <div className="w-full flex items-start justify-center p-4 md:p-8 pb-20">
      <div className="bg-white/90 backdrop-blur-2xl rounded-[48px] shadow-[0_40px_100px_-20px_rgba(0,0,0,0.15)] border border-white/50 p-8 md:p-10 w-full max-w-4xl flex flex-col items-stretch my-8">
        <div className="mb-12 text-center">
          <input
            type="text"
            value={budgetName}
            onChange={(e) => setBudgetName(e.target.value)}
            className="text-4xl font-black text-gray-900 tracking-tighter text-center bg-transparent focus:outline-none border-b-2 border-transparent focus:border-indigo-400 pb-2 w-full max-w-lg"
            placeholder="Budget Name"
          />
          <p className="text-gray-400 font-bold uppercase tracking-widest text-[10px] mt-4">Personal Spending Plan</p>
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

          {items.map(item => (
            <div key={item.id} className="bg-white border border-gray-100 shadow-sm rounded-3xl p-6 flex flex-col gap-4 relative group">
              <button 
                onClick={() => handleRemoveItem(item.id)}
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
                    onChange={(e) => handleUpdateItem(item.id, { name: e.target.value })}
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
                    onChange={(e) => handleUpdateItem(item.id, { percentage: parseFloat(e.target.value) || 0 })}
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
                      <button onClick={() => handleRemoveCategoryFromItem(item.id, cat)} className="hover:text-gray-900 transition-colors">
                        <X size={14} strokeWidth={2.5} />
                      </button>
                    </div>
                  ))}
                  
                  <select 
                    className="bg-gray-50 border border-gray-200 text-gray-600 font-medium rounded-xl px-3 py-1.5 text-sm focus:outline-none focus:border-indigo-400 transition-all"
                    onChange={(e) => {
                      if (e.target.value) {
                        handleAddCategoryToItem(item.id, e.target.value);
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
          className="flex items-center justify-center gap-2 w-full bg-gray-900 hover:bg-gray-800 text-white font-bold rounded-2xl py-4 transition-all shadow-lg hover:shadow-xl active:scale-[0.98]"
        >
          <PlusCircle size={20} />
          Add Budget Item
        </button>
      </div>
    </div>
  );
};

export default Budget;
