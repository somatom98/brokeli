import React, { useEffect, useState } from 'react';
import { 
  ArrowUpRight, 
  ArrowDownLeft, 
  ArrowRightLeft, 
  Check, 
  ChevronDown,
  Loader2,
  Receipt,
  Tag,
  AlignLeft,
  Banknote
} from 'lucide-react';
import { api } from './api';
import type { Account } from './api';

const App: React.FC = () => {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [success, setSuccess] = useState(false);

  // Form State
  const [type, setType] = useState<'expense' | 'income' | 'transfer'>('expense');
  const [accountId, setAccountId] = useState('');
  const [toAccountId, setToAccountId] = useState('');
  const [amount, setAmount] = useState('');
  const [currency, setCurrency] = useState('EUR');
  const [category, setCategory] = useState('');
  const [description, setDescription] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    const fetchAccounts = async () => {
      try {
        const accs = await api.getAccounts();
        setAccounts(accs || []);
        if (accs && accs.length > 0) {
          setAccountId(accs[0].id);
          if (accs.length > 1) setToAccountId(accs[1].id);
        }
      } catch (err) {
        console.error('Error fetching accounts:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchAccounts();
  }, []);

  const themes = {
    expense: {
      primary: 'text-rose-500',
      bg: 'bg-rose-500/10',
      btn: 'bg-rose-500 hover:bg-rose-600 shadow-rose-200',
      mesh: 'from-rose-400 via-orange-300 to-amber-200',
      tab: 'text-rose-600 bg-white/90 shadow-rose-200/50'
    },
    income: {
      primary: 'text-emerald-500',
      bg: 'bg-emerald-500/10',
      btn: 'bg-emerald-500 hover:bg-emerald-600 shadow-emerald-200',
      mesh: 'from-emerald-400 via-teal-300 to-sky-300',
      tab: 'text-emerald-600 bg-white/90 shadow-emerald-200/50'
    },
    transfer: {
      primary: 'text-indigo-500',
      bg: 'bg-indigo-500/10',
      btn: 'bg-indigo-500 hover:bg-indigo-600 shadow-indigo-200',
      mesh: 'from-indigo-400 via-violet-300 to-purple-300',
      tab: 'text-indigo-600 bg-white/90 shadow-indigo-200/50'
    }
  };

  const theme = themes[type];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isSubmitting || !amount) return;
    
    setIsSubmitting(true);
    try {
      const happenedAt = new Date().toISOString();
      if (type === 'transfer') {
        await api.registerTransfer({
          from_account_id: accountId,
          from_currency: currency,
          from_amount: amount,
          to_account_id: toAccountId,
          to_currency: currency,
          to_amount: amount,
          category,
          description,
          happened_at: happenedAt
        });
      } else if (type === 'income') {
        await api.registerIncome({
          account_id: accountId,
          currency,
          amount,
          category,
          description,
          happened_at: happenedAt
        });
      } else {
        await api.registerExpense({
          account_id: accountId,
          currency,
          amount,
          category,
          description,
          happened_at: happenedAt
        });
      }
      setSuccess(true);
      setAmount('');
      setCategory('');
      setDescription('');
      setTimeout(() => setSuccess(false), 2500);
    } catch (err) {
      alert('Error: ' + err);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (loading) return (
    <div className="min-h-screen flex items-center justify-center bg-white w-full">
      <Loader2 className="animate-spin text-gray-300" size={48} strokeWidth={1} />
    </div>
  );

  return (
    <div className="min-h-screen w-full relative flex items-center justify-center p-4 antialiased overflow-hidden">
      
      {/* Animated Mesh Background */}
      <div className={`mesh-bg bg-gradient-to-tr ${theme.mesh}`} />

      {/* Glassmorphic Panel Container */}
      <div className="w-full max-w-[440px] relative z-10 transition-all duration-500 scale-95 animate-in fade-in zoom-in-95 duration-1000">
        
        {/* Main Panel */}
        <div className="bg-white/90 backdrop-blur-2xl rounded-[48px] shadow-[0_40px_100px_-20px_rgba(0,0,0,0.15)] border border-white/50 p-8 md:p-10 relative overflow-hidden group">
          
          {/* Success Overlay */}
          <div className={`absolute inset-0 z-50 flex flex-col items-center justify-center transition-all duration-700 bg-white/95 ${success ? 'translate-y-0 opacity-100' : 'translate-y-full opacity-0 pointer-events-none'}`}>
             <div className={`w-24 h-24 ${theme.btn} text-white rounded-full flex items-center justify-center mb-6 shadow-2xl animate-bounce`}>
               <Check size={48} strokeWidth={4} />
             </div>
             <h3 className="text-4xl font-black text-gray-900 tracking-tight">Saved!</h3>
             <p className="text-gray-400 font-bold mt-2 uppercase tracking-widest text-xs">Ledger Updated</p>
          </div>

          <header className="text-center mb-8 space-y-3">
            <div className={`inline-block px-4 py-1.5 rounded-full bg-gray-100/50 text-gray-400 text-[10px] font-black uppercase tracking-[0.4em]`}>
              Brøkeli Core
            </div>
            <h2 className="text-3xl font-black text-gray-900 tracking-tighter">New Entry</h2>
          </header>

          <form onSubmit={handleSubmit} className="space-y-8">
            
            {/* Type Selector (TABS) */}
            <div className="flex p-1.5 bg-gray-100/80 rounded-[28px] gap-1 shadow-inner">
              {(['expense', 'income', 'transfer'] as const).map((t) => (
                <button
                  key={t}
                  type="button"
                  onClick={() => setType(t)}
                  className={`flex-1 flex flex-col items-center py-3.5 rounded-[22px] transition-all duration-500 ease-out ${
                    type === t 
                      ? `${theme.tab} scale-[1.02] shadow-xl` 
                      : 'text-gray-400 hover:text-gray-600'
                  }`}
                >
                  {t === 'expense' && <ArrowDownLeft size={20} strokeWidth={2.5} />}
                  {t === 'income' && <ArrowUpRight size={20} strokeWidth={2.5} />}
                  {t === 'transfer' && <ArrowRightLeft size={20} strokeWidth={2.5} />}
                  <span className="text-[9px] font-black uppercase tracking-widest mt-1.5">{t}</span>
                </button>
              ))}
            </div>

            {/* Huge Amount Input */}
            <div className="text-center relative group/input">
              <div className="flex items-center justify-center gap-3">
                <span className={`text-4xl font-black transition-colors duration-500 ${theme.primary}`}>
                   {currency === 'EUR' ? '€' : currency === 'USD' ? '$' : '£'}
                </span>
                <input 
                  type="number" 
                  step="0.01" 
                  required 
                  placeholder="0.00"
                  value={amount} 
                  onChange={(e) => setAmount(e.target.value)}
                  className="w-full max-w-[200px] text-center text-7xl font-black outline-none bg-transparent placeholder:text-gray-100 transition-all caret-indigo-500"
                  style={{ color: '#111827' }}
                />
              </div>

              {/* Currency Badges */}
              <div className="flex justify-center gap-2.5 mt-4">
                {['EUR', 'USD', 'GBP'].map(c => (
                  <button
                    key={c}
                    type="button"
                    onClick={() => setCurrency(c)}
                    className={`px-4 py-1.5 rounded-full text-[10px] font-black tracking-widest transition-all duration-500 border ${
                      currency === c 
                        ? `${theme.bg} ${theme.primary} border-transparent scale-105` 
                        : 'bg-transparent text-gray-300 border-gray-100 hover:border-gray-200'
                    }`}
                  >
                    {c}
                  </button>
                ))}
              </div>
            </div>

            {/* Detailed Fields */}
            <div className="space-y-4">
              
              {/* Account Dropdowns */}
              <div className="flex flex-col gap-3">
                <div className="relative group/field">
                  <div className="absolute left-5 top-3 flex items-center gap-2 pointer-events-none">
                    <Banknote size={14} className="text-gray-300" />
                    <span className="text-[9px] font-black text-gray-300 uppercase tracking-widest">
                      {type === 'transfer' ? 'From Account' : 'Account'}
                    </span>
                  </div>
                  <select 
                    value={accountId}
                    onChange={(e) => setAccountId(e.target.value)}
                    className="w-full bg-gray-50/50 hover:bg-gray-100 border-none rounded-[24px] px-5 pt-8 pb-3.5 text-sm font-bold appearance-none outline-none transition-all cursor-pointer focus:ring-4 focus:ring-indigo-50"
                  >
                    {accounts.map(acc => <option key={acc.id} value={acc.id}>{acc.name}</option>)}
                  </select>
                  <ChevronDown size={18} className="absolute right-5 bottom-4 text-gray-300 pointer-events-none" />
                </div>

                {type === 'transfer' && (
                  <div className="relative animate-in slide-in-from-top-4 duration-500">
                    <div className="absolute left-5 top-3 flex items-center gap-2 pointer-events-none">
                      <ArrowRightLeft size={14} className="text-gray-300" />
                      <span className="text-[9px] font-black text-gray-300 uppercase tracking-widest">To Account</span>
                    </div>
                    <select 
                      value={toAccountId}
                      onChange={(e) => setToAccountId(e.target.value)}
                      className="w-full bg-gray-50/50 hover:bg-gray-100 border-none rounded-[24px] px-5 pt-8 pb-3.5 text-sm font-bold appearance-none outline-none transition-all cursor-pointer focus:ring-4 focus:ring-indigo-50"
                    >
                      {accounts.map(acc => <option key={acc.id} value={acc.id}>{acc.name}</option>)}
                    </select>
                    <ChevronDown size={18} className="absolute right-5 bottom-4 text-gray-300 pointer-events-none" />
                  </div>
                )}
              </div>

              {/* Categorization */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                <div className="relative">
                  <div className="absolute left-5 top-3 flex items-center gap-2 pointer-events-none">
                    <Tag size={14} className="text-gray-300" />
                    <span className="text-[9px] font-black text-gray-300 uppercase tracking-widest">Category</span>
                  </div>
                  <input 
                    type="text" 
                    placeholder="Shopping..."
                    value={category} 
                    onChange={(e) => setCategory(e.target.value)}
                    className="w-full bg-gray-50/50 border-none rounded-[24px] px-5 pt-8 pb-3.5 text-sm font-bold outline-none transition-all focus:ring-4 focus:ring-indigo-50 placeholder:text-gray-200"
                  />
                </div>

                <div className="relative">
                  <div className="absolute left-5 top-3 flex items-center gap-2 pointer-events-none">
                    <AlignLeft size={14} className="text-gray-300" />
                    <span className="text-[9px] font-black text-gray-300 uppercase tracking-widest">Note</span>
                  </div>
                  <input 
                    type="text" 
                    placeholder="Brief note"
                    value={description} 
                    onChange={(e) => setDescription(e.target.value)}
                    className="w-full bg-gray-50/50 border-none rounded-[24px] px-5 pt-8 pb-3.5 text-sm font-bold outline-none transition-all focus:ring-4 focus:ring-indigo-50 placeholder:text-gray-200"
                  />
                </div>
              </div>

            </div>

            {/* Submission Button */}
            <button 
              type="submit" 
              disabled={isSubmitting || !amount}
              className={`w-full ${theme.btn} text-white font-black py-6 rounded-[32px] transition-all duration-500 shadow-2xl active:scale-95 flex items-center justify-center gap-3 text-[11px] uppercase tracking-[0.3em] disabled:bg-gray-200 disabled:shadow-none`}
            >
              {isSubmitting ? (
                <Loader2 className="animate-spin" size={20} strokeWidth={4} />
              ) : (
                <>
                  <Receipt size={18} strokeWidth={3} />
                  Record Transaction
                </>
              )}
            </button>
          </form>
        </div>

        {/* Status Indicator */}
        <div className="mt-8 flex justify-center items-center gap-3 text-[10px] font-black text-gray-400 uppercase tracking-widest">
           <div className={`w-2 h-2 rounded-full ${theme.primary} shadow-[0_0_10px_currentColor] animate-pulse`} />
           Cloud Node Sync Active
        </div>
      </div>
    </div>
  );
};

export default App;
