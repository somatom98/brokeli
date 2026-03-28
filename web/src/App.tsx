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
  Banknote,
  PlusSquare,
  ArrowDownToLine,
  ArrowUpFromLine,
  Menu,
  X,
  XCircle,
  Home,
  PieChart,
  BarChart3,
  Calendar
  } from 'lucide-react';
  import { api } from './api';
  import type { Account } from './api';
  import Budget from './Budget';
  import Transactions from './Transactions';
  import Balances from './Balances';

  const App: React.FC = () => {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [categories, setCategories] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  // App Navigation State
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [currentView, setCurrentView] = useState<'home' | 'budget' | 'balances'>('home');
  const [transactionsRefreshKey, setTransactionsRefreshKey] = useState(0);

  // Form State
  const [type, setType] = useState<'expense' | 'income' | 'transfer' | 'openAccount' | 'deposit' | 'withdraw'>('expense');
  const [accountId, setAccountId] = useState('');
  const [toAccountId, setToAccountId] = useState('');
  const [amount, setAmount] = useState('');
  const [currency, setCurrency] = useState('EUR');
  const [category, setCategory] = useState('');
  const [description, setDescription] = useState('');
  const [accountName, setAccountName] = useState('');
  const [happenedAtDateTime, setHappenedAtDateTime] = useState(() => {
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    return now.toISOString().slice(0, 16);
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

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
      setError(true);
      setErrorMessage('Failed to fetch accounts');
      setTimeout(() => setError(false), 3500);
    } finally {
      setLoading(false);
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
    Promise.all([fetchAccounts(), fetchCategories()]);
  }, []);

  const themes = {
    expense: {
      primary: 'text-negative',
      bg: 'bg-negative/10',
      btn: 'bg-negative hover:bg-negative/90 shadow-negative/20',
      mesh: 'from-negative/5 via-app-bg to-app-bg',
      tab: 'text-negative bg-card/90 shadow-negative/10'
    },
    income: {
      primary: 'text-primary',
      bg: 'bg-primary/10',
      btn: 'bg-primary hover:bg-primary/90 shadow-primary/20',
      mesh: 'from-primary/5 via-app-bg to-app-bg',
      tab: 'text-primary bg-card/90 shadow-primary/10'
    },
    transfer: {
      primary: 'text-neutral',
      bg: 'bg-neutral/10',
      btn: 'bg-neutral hover:bg-neutral/90 shadow-neutral/20',
      mesh: 'from-neutral/5 via-app-bg to-app-bg',
      tab: 'text-neutral bg-card/90 shadow-neutral/10'
    },
    openAccount: {
      primary: 'text-accent',
      bg: 'bg-accent/10',
      btn: 'bg-accent hover:bg-accent/90 shadow-accent/20',
      mesh: 'from-accent/5 via-app-bg to-app-bg',
      tab: 'text-accent bg-card/90 shadow-accent/10'
    },
    deposit: {
      primary: 'text-neutral',
      bg: 'bg-neutral/10',
      btn: 'bg-neutral hover:bg-neutral/90 shadow-neutral/20',
      mesh: 'from-neutral/5 via-app-bg to-app-bg',
      tab: 'text-neutral bg-card/90 shadow-neutral/10'
    },
    withdraw: {
      primary: 'text-neutral',
      bg: 'bg-neutral/10',
      btn: 'bg-neutral hover:bg-neutral/90 shadow-neutral/20',
      mesh: 'from-neutral/5 via-app-bg to-app-bg',
      tab: 'text-neutral bg-card/90 shadow-neutral/10'
    }
  };

  const theme = themes[type];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isSubmitting) return;
    if (type !== 'openAccount' && !amount) return;
    if (type === 'openAccount' && !accountName) return;
    
    setIsSubmitting(true);
    try {
      const happenedAt = new Date(happenedAtDateTime).toISOString();
      
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
      } else if (type === 'expense') {
        await api.registerExpense({
          account_id: accountId,
          currency,
          amount,
          category,
          description,
          happened_at: happenedAt
        });
      } else if (type === 'openAccount') {
        await api.openAccount({
          name: accountName,
          currency
        });
        await fetchAccounts();
      } else if (type === 'deposit') {
        await api.deposit(accountId, {
          currency,
          amount,
          happened_at: happenedAt
        });
      } else if (type === 'withdraw') {
        await api.withdraw(accountId, {
          currency,
          amount,
          happened_at: happenedAt
        });
      }
      setSuccess(true);
      setAmount('');
      setCategory('');
      setDescription('');
      setAccountName('');
      setHappenedAtDateTime(() => {
        const now = new Date();
        now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
        return now.toISOString().slice(0, 16);
      });
      fetchCategories();
      setTransactionsRefreshKey(prev => prev + 1);
      setTimeout(() => setSuccess(false), 2500);
    } catch (err) {
      setError(true);
      setErrorMessage(String(err));
      setTimeout(() => setError(false), 3500);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (loading) return (
    <div className="min-h-screen flex items-center justify-center bg-app-bg w-full">
      <Loader2 className="animate-spin text-text-muted/20" size={48} strokeWidth={1} />
    </div>
  );

  return (
    <div className="min-h-screen lg:h-screen w-full relative flex flex-col lg:overflow-hidden antialiased bg-transparent text-text-main">
      
      {/* Animated Mesh Background */}
      <div className={`fixed inset-0 z-0 mesh-bg bg-gradient-to-tr ${theme.mesh} transition-colors duration-1000`} />

      {/* Top Navigation Bar */}
      <div className="fixed top-4 left-4 z-30">
        <button 
          onClick={() => setIsSidebarOpen(true)}
          className="p-2.5 bg-glass backdrop-blur-xl rounded-xl shadow-lg hover:bg-glass-hover transition-all text-text-on-dark border border-glass-border active:scale-95"
          title="Menu"
        >
          <Menu size={18} strokeWidth={2.5} />
        </button>
      </div>

      {/* Sidebar Overlay */}
      {isSidebarOpen && (
        <div 
          className="fixed inset-0 bg-transparent z-40 transition-opacity"
          onClick={() => setIsSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <div 
        className={`fixed top-0 left-0 h-full w-80 bg-glass backdrop-blur-xl shadow-2xl z-50 transform transition-transform duration-500 ease-in-out ${
          isSidebarOpen ? 'translate-x-0' : '-translate-x-full'
        } flex flex-col border-r border-glass-border`}
      >
        <div className="p-8 flex justify-between items-center border-b border-glass-border">
          <h2 className="text-2xl font-black text-text-on-dark tracking-tighter">Brøkeli</h2>
          <button 
            onClick={() => setIsSidebarOpen(false)}
            className="p-2 bg-glass-dark rounded-full hover:bg-glass-dark-hover transition-colors text-text-on-dark"
          >
            <X size={20} strokeWidth={2.5} />
          </button>
        </div>
        <div className="flex-1 p-6 space-y-2">
          <button
            onClick={() => { setCurrentView('home'); setIsSidebarOpen(false); }}
            className={`w-full flex items-center gap-4 px-6 py-4 rounded-[20px] transition-all font-black uppercase tracking-[0.2em] text-[10px] group border ${
              currentView === 'home' 
                ? 'bg-accent/20 backdrop-blur-md text-text-on-dark border-accent/30 shadow-lg scale-[1.02]' 
                : 'text-text-on-dark border-transparent hover:bg-glass-dark-hover hover:scale-[1.02]'
            }`}
          >
            <Home size={18} strokeWidth={3} />
            <span>Ledger</span>
          </button>
          <button
            onClick={() => { setCurrentView('balances'); setIsSidebarOpen(false); }}
            className={`w-full flex items-center gap-4 px-6 py-4 rounded-[20px] transition-all font-black uppercase tracking-[0.2em] text-[10px] group border ${
              currentView === 'balances' 
                ? 'bg-accent/20 backdrop-blur-md text-text-on-dark border-accent/30 shadow-lg scale-[1.02]' 
                : 'text-text-on-dark border-transparent hover:bg-glass-dark-hover hover:scale-[1.02]'
            }`}
          >
            <PieChart size={18} strokeWidth={3} />
            <span>Balances</span>
          </button>
          <button
            onClick={() => { setCurrentView('budget'); setIsSidebarOpen(false); }}
            className={`w-full flex items-center gap-4 px-6 py-4 rounded-[20px] transition-all font-black uppercase tracking-[0.2em] text-[10px] group border ${
              currentView === 'budget' 
                ? 'bg-accent/20 backdrop-blur-md text-text-on-dark border-accent/30 shadow-lg scale-[1.02]' 
                : 'text-text-on-dark border-transparent hover:bg-glass-dark-hover hover:scale-[1.02]'
            }`}
          >
            <BarChart3 size={18} strokeWidth={3} />
            <span>Budget</span>
          </button>
        </div>
      </div>

      {/* Main Content Area */}
      <main className="flex-1 w-full relative z-10 lg:overflow-hidden bg-transparent">
        {currentView === 'home' ? (
          <div className="w-full h-auto lg:h-full max-w-[1800px] mx-auto flex flex-col lg:flex-row items-stretch justify-start lg:justify-center gap-8 lg:gap-12 transition-all duration-500 animate-in fade-in zoom-in-100 duration-1000 px-6 pb-10 lg:px-10 lg:pb-10 pt-20 lg:pt-4">
            <div className="w-full lg:w-[440px] flex flex-col shrink-0 min-h-0">
              {/* Alignment Spacer: Matches Transactions Filter height + gap */}
              <div className="h-[54px] shrink-0 hidden lg:block" />
              
              {/* Main Panel Content (WHITE CARD) */}
              <div className="bg-card rounded-[40px] shadow-sm hover:shadow-xl border border-border-pearl p-8 relative overflow-hidden group flex flex-col flex-1 min-h-0 transition-all duration-300">
                
                {/* Success Overlay */}
                <div className={`absolute inset-0 z-50 flex flex-col items-center justify-center transition-all duration-700 bg-app-bg/95 ${success ? 'translate-y-0 opacity-100' : 'translate-y-full opacity-0 pointer-events-none'}`}>
                   <div className="w-20 h-20 bg-primary text-white rounded-full flex items-center justify-center mb-4 shadow-2xl animate-bounce">
                     <Check size={40} strokeWidth={4} />
                   </div>
                   <h3 className="text-3xl font-black text-text-main tracking-tight">Saved!</h3>
                </div>


                {/* Error Overlay */}
                <div className={`absolute inset-0 z-50 flex flex-col items-center justify-center transition-all duration-700 bg-app-bg/95 ${error ? 'translate-y-0 opacity-100' : 'translate-y-full opacity-0 pointer-events-none'}`}>
                   <div className="w-20 h-20 bg-negative text-white rounded-full flex items-center justify-center mb-4 shadow-2xl animate-bounce">
                     <XCircle size={40} strokeWidth={4} />
                   </div>
                   <h3 className="text-3xl font-black text-text-main tracking-tight text-center px-6">Error</h3>
                   <p className="text-negative font-bold mt-2 uppercase tracking-widest text-[10px] text-center px-10 leading-relaxed">{errorMessage}</p>
                </div>

                <form onSubmit={handleSubmit} className="space-y-6 overflow-y-auto pr-2 custom-scrollbar flex-1">
                  
                  {/* Type Selector (TABS) */}
                  <div className="grid grid-cols-3 p-1 bg-card-muted rounded-[24px] gap-1 shadow-inner shrink-0">
                    {(['expense', 'income', 'transfer', 'openAccount', 'deposit', 'withdraw'] as const).map((t) => (
                      <button
                        key={t}
                        type="button"
                        onClick={() => setType(t)}
                        className={`flex flex-col items-center py-3 rounded-[20px] transition-all duration-500 ease-out ${
                          type === t
                            ? `text-accent scale-[1.05]`
                            : 'text-text-muted hover:text-text-main'
                        }`}                      >
                        {t === 'expense' && <ArrowDownLeft size={18} strokeWidth={2.5} />}
                        {t === 'income' && <ArrowUpRight size={18} strokeWidth={2.5} />}
                        {t === 'transfer' && <ArrowRightLeft size={18} strokeWidth={2.5} />}
                        {t === 'openAccount' && <PlusSquare size={18} strokeWidth={2.5} />}
                        {t === 'deposit' && <ArrowDownToLine size={18} strokeWidth={2.5} />}
                        {t === 'withdraw' && <ArrowUpFromLine size={18} strokeWidth={2.5} />}
                        <span className="text-[8px] font-black uppercase tracking-widest mt-1">
                          {t === 'openAccount' ? 'Account' : t}
                        </span>
                      </button>
                    ))}
                  </div>

                  {/* Currency Badges for openAccount */}
                  {type === 'openAccount' && (
                    <div className="flex justify-center items-center gap-2 mt-2">
                      {['EUR', 'DKK'].map(c => (
                        <button
                          key={c}
                          type="button"
                          onClick={() => setCurrency(c)}
                          className={`px-3 py-1 rounded-full text-[9px] font-black tracking-widest transition-all duration-500 border ${
                            currency === c 
                              ? `bg-accent/10 text-accent border-accent/20 scale-105` 
                              : 'bg-transparent text-text-muted/30 border-border-pearl hover:border-text-muted/20'
                          }`}
                        >
                          {c}
                        </button>
                      ))}
                    </div>
                  )}

                  {/* Huge Amount Input */}
                  {type !== 'openAccount' && (
                    <div className="text-center relative group/input">
                      <div className="flex items-center justify-center gap-2">
                        <span className={`text-3xl font-black transition-colors duration-500 text-accent-secondary`}>
                           {currency === 'EUR' ? '€' : currency === 'DKK' ? 'kr' : currency}
                        </span>
                        <input 
                          type="number" 
                          step="0.01" 
                          required 
                          placeholder="0.00"
                          value={amount} 
                          onChange={(e) => setAmount(e.target.value)}
                          className="w-full max-w-[160px] text-center text-6xl font-black outline-none bg-transparent placeholder:text-border-pearl transition-all caret-accent text-text-main"
                        />
                      </div>

                      {/* Currency Badges */}
                      <div className="flex justify-center items-center gap-2 mt-2">
                        {['EUR', 'DKK'].map(c => (
                          <button
                            key={c}
                            type="button"
                            onClick={() => setCurrency(c)}
                            className={`px-3 py-1 rounded-full text-[9px] font-black tracking-widest transition-all duration-500 border ${
                              currency === c 
                                ? `bg-accent text-white border-accent scale-105 shadow-lg shadow-accent/20` 
                                : 'bg-transparent text-text-muted border-border-pearl hover:border-text-muted/20'
                            }`}
                          >
                            {c}
                          </button>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Detailed Fields */}
                  <div className="space-y-3">
                    
                    {/* Account Dropdowns */}
                    {type !== 'openAccount' && (
                      <div className="flex flex-col gap-2">
                        <div className="relative group/field">
                          <div className="absolute left-4 top-2.5 flex items-center gap-2 pointer-events-none">
                            <Banknote size={12} className="text-text-muted/40" />
                            <span className="text-[8px] font-black text-text-muted/40 uppercase tracking-widest">
                              {type === 'transfer' ? 'From' : 'Account'}
                            </span>
                          </div>
                          <select 
                            value={accountId}
                            onChange={(e) => setAccountId(e.target.value)}
                            className="w-full bg-card-muted hover:bg-border-pearl border-none rounded-[20px] px-4 pt-7 pb-3 text-sm font-bold appearance-none outline-none transition-all cursor-pointer focus:ring-4 focus:ring-accent/5 text-text-main"
                          >
                            {accounts.map(acc => <option key={acc.id} value={acc.id}>{acc.name}</option>)}
                          </select>
                          <ChevronDown size={16} className="absolute right-4 bottom-3.5 text-text-muted/40 pointer-events-none" />
                        </div>

                        {type === 'transfer' && (
                          <div className="relative animate-in slide-in-from-top-2 duration-500">
                            <div className="absolute left-4 top-2.5 flex items-center gap-2 pointer-events-none">
                              <ArrowRightLeft size={12} className="text-text-muted/40" />
                              <span className="text-[8px] font-black text-text-muted/40 uppercase tracking-widest">To</span>
                            </div>
                            <select 
                              value={toAccountId}
                              onChange={(e) => setToAccountId(e.target.value)}
                              className="w-full bg-card-muted hover:bg-border-pearl border-none rounded-[20px] px-4 pt-7 pb-3 text-sm font-bold appearance-none outline-none transition-all cursor-pointer focus:ring-4 focus:ring-accent/5 text-text-main"
                            >
                              {accounts.map(acc => <option key={acc.id} value={acc.id}>{acc.name}</option>)}
                            </select>
                            <ChevronDown size={16} className="absolute right-4 bottom-3.5 text-text-muted/40 pointer-events-none" />
                          </div>
                        )}
                      </div>
                    )}

                    {/* Date & Time Selection */}
                    {type !== 'openAccount' && (
                      <div className="relative">
                        <div className="absolute left-4 top-2.5 flex items-center gap-2 pointer-events-none">
                          <Calendar size={12} className="text-text-muted/40" />
                          <span className="text-[8px] font-black text-text-muted/40 uppercase tracking-widest">Date & Time</span>
                        </div>
                        <input 
                          type="datetime-local" 
                          required
                          value={happenedAtDateTime} 
                          onChange={(e) => setHappenedAtDateTime(e.target.value)}
                          className="w-full bg-card-muted border-none rounded-[20px] px-4 pt-7 pb-3 text-sm font-bold outline-none transition-all focus:ring-4 focus:ring-accent/5 text-text-main"
                        />
                      </div>
                    )}

                    {/* Categorization */}
                    {(type === 'income' || type === 'expense' || type === 'transfer') && (
                      <div className="grid grid-cols-1 gap-2">
                        <div className="relative">
                          <div className="absolute left-4 top-2.5 flex items-center gap-2 pointer-events-none">
                            <Tag size={12} className="text-text-muted/40" />
                            <span className="text-[8px] font-black text-text-muted/40 uppercase tracking-widest">Category</span>
                          </div>
                          <input 
                            type="text" 
                            placeholder="Shopping..."
                            list="category-suggestions"
                            value={category} 
                            onChange={(e) => setCategory(e.target.value)}
                            className="w-full bg-card-muted border-none rounded-[20px] px-4 pt-7 pb-3 text-sm font-bold outline-none transition-all focus:ring-4 focus:ring-accent/5 text-text-main placeholder:text-text-muted/20"
                          />
                          <datalist id="category-suggestions">
                            {categories.map(cat => (
                              <option key={cat} value={cat} />
                            ))}
                          </datalist>
                        </div>

                        <div className="relative">
                          <div className="absolute left-4 top-2.5 flex items-center gap-2 pointer-events-none">
                            <AlignLeft size={12} className="text-text-muted/40" />
                            <span className="text-[8px] font-black text-text-muted/40 uppercase tracking-widest">Note</span>
                          </div>
                          <input 
                            type="text" 
                            placeholder="Brief note"
                            value={description} 
                            onChange={(e) => setDescription(e.target.value)}
                            className="w-full bg-card-muted border-none rounded-[20px] px-4 pt-7 pb-3 text-sm font-bold outline-none transition-all focus:ring-4 focus:ring-accent/5 text-text-main placeholder:text-text-muted/20"
                          />
                        </div>
                      </div>
                    )}

                    {/* Account Name */}
                    {type === 'openAccount' && (
                      <div className="relative group/field">
                        <div className="absolute left-4 top-2.5 flex items-center gap-2 pointer-events-none">
                          <Tag size={12} className="text-text-muted/40" />
                          <span className="text-[8px] font-black text-text-muted/40 uppercase tracking-widest">Account Name</span>
                        </div>
                        <input 
                          type="text" 
                          required
                          placeholder="e.g. Main Wallet"
                          value={accountName} 
                          onChange={(e) => setAccountName(e.target.value)}
                          className="w-full bg-card-muted border-none rounded-[20px] px-4 pt-7 pb-3 text-sm font-bold outline-none transition-all focus:ring-4 focus:ring-accent/5 text-text-main placeholder:text-text-muted/20"
                        />
                      </div>
                    )}

                  </div>

                  {/* Submission Button */}
                  <button 
                    type="submit" 
                    disabled={isSubmitting || (type !== 'openAccount' && !amount) || (type === 'openAccount' && !accountName)}
                    className={`w-full bg-accent text-white font-black py-5 rounded-[28px] transition-all duration-500 shadow-2xl shadow-accent/20 hover:bg-accent/90 active:scale-95 flex items-center justify-center gap-3 text-[10px] uppercase tracking-[0.3em] disabled:bg-card-muted disabled:text-text-muted/20 disabled:shadow-none shrink-0`}
                  >
                    {isSubmitting ? (
                      <Loader2 className="animate-spin" size={18} strokeWidth={4} />
                    ) : (
                      <>
                        <Receipt size={16} strokeWidth={3} />
                        {type === 'openAccount' ? 'Open Account' : type === 'deposit' ? 'Deposit' : type === 'withdraw' ? 'Withdraw' : 'Record'}
                      </>
                    )}
                  </button>
                </form>
              </div>
            </div>

            <div className="flex-1 min-w-0 h-auto lg:h-full">
              <Transactions refreshKey={transactionsRefreshKey} hideHeader />
            </div>
          </div>
        ) : currentView === 'balances' ? (
          <div className="w-full h-auto lg:h-full overflow-y-auto custom-scrollbar pt-20 lg:pt-16">
            <div className="w-full relative z-10 animate-in fade-in zoom-in-95 duration-500">
              <Balances />
            </div>
          </div>
        ) : (
          <div className="w-full h-auto lg:h-full overflow-y-auto custom-scrollbar pt-20 lg:pt-16">
            <div className="w-full relative z-10 animate-in fade-in zoom-in-95 duration-500">
              <Budget />
            </div>
          </div>
        )}
      </main>
    </div>
  );
};

export default App;
