import React, { useEffect, useState, useMemo } from 'react';
import { 
  Loader2, 
  Wallet, 
  Calendar
} from 'lucide-react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
  type TooltipItem
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import { api } from './api';
import type { Account, BalancePeriod, AccountDistribution } from './api';
import { getCSSVariableValue } from './utils/colors';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

interface AccountWithMetadata extends Account {
  lastTransactionAt?: string;
  history?: BalancePeriod[];
  distributions?: AccountDistribution[];
}

const Sparkline: React.FC<{ data: BalancePeriod[], color: string }> = ({ data, color }) => {
  const chartData = {
    labels: data.map(h => h.month),
    datasets: [
      {
        data: data.map(h => parseFloat(h.amount)),
        borderColor: color,
        borderWidth: 2,
        pointRadius: 0,
        tension: 0.4,
        fill: false,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: { 
      legend: { display: false }, 
      tooltip: { enabled: false },
      datalabels: { display: false }
    },
    scales: {
      x: { display: false },
      y: { display: false },
    },
    elements: {
        point: {
            radius: 0
        }
    }
  };

  return (
    <div className="h-12 w-24">
      <Line data={chartData} options={options} />
    </div>
  );
};

const Balances: React.FC = () => {
  const [accounts, setAccounts] = useState<AccountWithMetadata[]>([]);
  const [liquidityHistory, setLiquidityHistory] = useState<BalancePeriod[]>([]);
  const [investmentHistory, setInvestmentHistory] = useState<BalancePeriod[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [timeFilter, setTimeFilter] = useState<'ytd' | 'year' | '5years' | 'all'>('year');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [accs, liqHistory, invHistory, transactions] = await Promise.all([
          api.getAccounts(),
          api.getBalances('LIQUIDITY'),
          api.getBalances('INVESTMENT'),
          api.getTransactions()
        ]);

        // Map last transaction time to accounts
        const accountsWithMetadata: AccountWithMetadata[] = await Promise.all(accs.map(async (acc) => {
          const accTransactions = transactions.filter(t => t.account_id === acc.id);
          const lastTx = accTransactions.length > 0 
            ? accTransactions.reduce((latest, current) => 
                new Date(current.happened_at) > new Date(latest.happened_at) ? current : latest
              ).happened_at
            : undefined;
          
          // Fetch history for sparkline (Total Wealth = Liquidity + Investment)
          const [accHistory, distributions] = await Promise.all([
            api.getBalancesByAccount(acc.id),
            api.getAccountDistributions(acc.id, 'LIQUIDITY')
          ]);

          return {
            ...acc,
            lastTransactionAt: lastTx,
            history: accHistory,
            distributions: distributions
          };
        }));

        // Sort by last transaction (most recent first)
        accountsWithMetadata.sort((a, b) => {
          if (!a.lastTransactionAt) return 1;
          if (!b.lastTransactionAt) return -1;
          return new Date(b.lastTransactionAt).getTime() - new Date(a.lastTransactionAt).getTime();
        });

        setAccounts(accountsWithMetadata);
        setLiquidityHistory(liqHistory || []);
        setInvestmentHistory(invHistory || []);
      } catch (err) {
        console.error('Error fetching balance data:', err);
        setError(true);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const chartData = useMemo(() => {
    const now = new Date();
    
    const filterByTime = (history: BalancePeriod[]) => history.filter(h => {
      const date = new Date(h.month);
      switch (timeFilter) {
        case 'ytd':
          return date >= new Date(now.getFullYear(), 0, 1);
        case 'year':
          return date >= new Date(now.getFullYear() - 1, now.getMonth(), now.getDate());
        case '5years':
          return date >= new Date(now.getFullYear() - 5, now.getMonth(), now.getDate());
        case 'all':
        default:
          return true;
      }
    });

    const filteredLiquidity = filterByTime(liquidityHistory);
    const filteredInvestment = filterByTime(investmentHistory);

    const allFiltered = [...filteredLiquidity, ...filteredInvestment];
    if (allFiltered.length === 0) return { labels: [], datasets: [] };

    // Find earliest date to start the timeline
    const earliestDate = allFiltered.reduce((min, h) => {
        const d = new Date(h.month);
        return d < min ? d : min;
    }, new Date());

    const months: string[] = [];
    const current = new Date(earliestDate.getFullYear(), earliestDate.getMonth(), 1);
    const end = new Date(now.getFullYear(), now.getMonth(), 1);

    while (current <= end) {
        months.push(current.toLocaleDateString(undefined, { month: 'short', year: '2-digit' }));
        current.setMonth(current.getMonth() + 1);
    }

    const currencies = Array.from(new Set(allFiltered.map(h => h.currency)));
    const datasets: any[] = [];

    const primary = getCSSVariableValue('--color-primary');
    const accentSecondary = getCSSVariableValue('--color-accent-secondary');

    currencies.forEach((curr) => {
      // 1. Liquidity Carrying Forward
      let lastLiq = 0;
      const sortedLiq = [...liquidityHistory]
        .filter(h => h.currency === curr)
        .sort((a, b) => new Date(a.month).getTime() - new Date(b.month).getTime());

      const liqData = months.map(mLabel => {
        const h = sortedLiq.find(history => {
            const date = new Date(history.month);
            return date.toLocaleDateString(undefined, { month: 'short', year: '2-digit' }) === mLabel;
        });
        if (h) lastLiq = parseFloat(h.amount);
        return lastLiq;
      });

      datasets.push({
        label: `${curr} Liquidity`,
        data: liqData,
        borderColor: primary,
        backgroundColor: `${primary}33`,
        fill: true,
        tension: 0.4,
        pointRadius: 0,
        hoverPointRadius: 6,
        pointBackgroundColor: primary,
        borderWidth: 3,
        pointHitRadius: 10,
      });

      // 2. Investment Carrying Forward
      let lastInv = 0;
      const sortedInv = [...investmentHistory]
        .filter(h => h.currency === curr)
        .sort((a, b) => new Date(a.month).getTime() - new Date(b.month).getTime());

      const invData = months.map(mLabel => {
        const h = sortedInv.find(history => {
            const date = new Date(history.month);
            return date.toLocaleDateString(undefined, { month: 'short', year: '2-digit' }) === mLabel;
        });
        if (h) lastInv = parseFloat(h.amount);
        return lastInv;
      });

      if (invData.some(d => d !== 0)) {
        datasets.push({
          label: `${curr} Investments`,
          data: invData,
          borderColor: accentSecondary,
          backgroundColor: `${accentSecondary}33`,
          fill: true,
          tension: 0.4,
          pointRadius: 0,
          hoverPointRadius: 6,
          pointBackgroundColor: accentSecondary,
          borderWidth: 3,
          pointHitRadius: 10,
        });
      }
    });

    return {
      labels: months,
      datasets: datasets
    };
  }, [liquidityHistory, investmentHistory, timeFilter]);

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
        intersect: false,
        mode: 'index' as const,
    },
    plugins: {
      legend: {
        display: true,
        position: 'top' as const,
        labels: {
          usePointStyle: true,
          padding: 20,
          color: getCSSVariableValue('--color-text-main'),
          font: {
            family: 'inherit',
            weight: 'bold' as const,
            size: 11
          }
        }
      },
      datalabels: {
        display: false
      },
      tooltip: {
        backgroundColor: getCSSVariableValue('--color-card'),
        titleColor: getCSSVariableValue('--color-text-main'),
        bodyColor: getCSSVariableValue('--color-text-main'),
        borderColor: getCSSVariableValue('--color-border-pearl'),
        borderWidth: 1,
        padding: 12,
        boxPadding: 6,
        usePointStyle: true,
        callbacks: {
          label: (context: TooltipItem<'line'>) => {
            let label = context.dataset.label || '';
            if (label) label += ': ';
            if (context.parsed.y !== null) {
              label += new Intl.NumberFormat(undefined, { style: 'currency', currency: context.dataset.label }).format(context.parsed.y);
            }
            return label;
          }
        }
      }
    },
    scales: {
      x: {
        grid: {
          display: false
        },
        ticks: {
          font: {
            size: 10,
            weight: 'bold' as const
          },
          color: getCSSVariableValue('--color-text-muted')
        }
      },
      y: {
        grid: {
          color: getCSSVariableValue('--color-border-pearl'),
        },
        ticks: {
          font: {
            size: 10,
            weight: 'bold' as const
          },
          color: getCSSVariableValue('--color-text-muted'),
          callback: (value: string | number) => {
            return value.toLocaleString();
          }
        }
      }
    },
    elements: {
        point: {
            radius: 0,
            hoverRadius: 6
        }
    }
  };

  if (loading) return (
    <div className="flex items-center justify-center p-20 bg-transparent w-full h-full">
      <Loader2 className="animate-spin text-accent" size={48} strokeWidth={1} />
    </div>
  );

  if (error) return (
    <div className="text-center p-20 bg-transparent w-full h-full">
      <p className="text-negative font-bold">Failed to load balance overview</p>
    </div>
  );

  return (
    <div className="max-w-6xl mx-auto w-full space-y-8 pb-20 bg-transparent">
      {/* Header Section */}
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 px-4">
        <div>
          <h1 className="text-5xl font-black text-text-on-dark tracking-tighter">Balances</h1>
          <p className="text-text-on-dark/40 font-bold mt-2 uppercase tracking-[0.2em] text-[10px]">Net worth and account performance</p>
        </div>
      </div>

      {/* History Graph */}
      <div className="px-4">
        <div className="bg-card p-8 md:p-12 rounded-[48px] shadow-lg border border-border-pearl">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-6 mb-12">
            <div>
              <h2 className="text-2xl font-black text-text-main tracking-tight">Balance History</h2>
              <p className="text-text-muted/40 font-bold mt-1 uppercase tracking-widest text-[9px]">Progression across currencies</p>
            </div>
            <div className="flex items-center gap-4">
                <Calendar size={16} className="text-text-muted/20" />
                <div className="flex items-center gap-2 bg-card-muted p-1.5 rounded-2xl border border-border-pearl">
                  {(['ytd', 'year', '5years', 'all'] as const).map((f) => (
                    <button
                      key={f}
                      onClick={() => setTimeFilter(f)}
                      className={`px-4 py-1.5 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all duration-300 ${
                        timeFilter === f 
                          ? 'bg-accent text-white shadow-lg shadow-accent/20' 
                          : 'text-text-muted/40 hover:text-text-muted hover:bg-card'
                      }`}
                    >
                      {f === '5years' ? '5Y' : f}
                    </button>
                  ))}
                </div>
            </div>
          </div>
          
          <div className="h-[400px] w-full">
            <Line data={chartData} options={chartOptions} />
          </div>
        </div>
      </div>

      {/* Accounts Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 px-4">
        {accounts.map(account => (
          <div key={account.id} className="bg-card p-8 rounded-[40px] shadow-lg border border-border-pearl hover:shadow-2xl transition-all duration-500 group">
            <div className="flex items-start justify-between mb-8">
              <div className="p-4 bg-card-muted rounded-3xl group-hover:bg-accent/10 transition-colors duration-500">
                <Wallet className="text-text-muted group-hover:text-accent transition-colors duration-500" size={24} />
              </div>
              <div className="flex flex-col items-end">
                <span className="text-[10px] font-black text-text-muted/40 uppercase tracking-widest mb-1">6M Trend</span>
                {account.history && (() => {
                  // Get current month and create a list of last 6 months in YYYY-MM format
                  const now = new Date();
                  const months: string[] = [];
                  for (let i = 5; i >= 0; i--) {
                    const d = new Date(now.getFullYear(), now.getMonth() - i, 1);
                    months.push(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`);
                  }

                  // Map existing history to these months, carrying forward balances
                  const sortedHistory = [...account.history].sort((a, b) => new Date(a.month).getTime() - new Date(b.month).getTime());
                  
                  let runningBalance = "0";
                  const last6Points = months.map(m => {
                    const existing = sortedHistory.find(h => h.month.startsWith(m));
                    if (existing) {
                      runningBalance = existing.amount;
                    } else {
                      // Carry forward: find the last available balance before this month
                      const lastBefore = [...sortedHistory].reverse().find(h => h.month < m);
                      if (lastBefore) {
                        runningBalance = lastBefore.amount;
                      } else {
                        // If no data before this month, it stays at 0 (or last carried value)
                        // but if we are iterating chronologically, we should just use runningBalance
                      }
                    }
                    return {
                      month: m,
                      amount: runningBalance,
                      currency: account.history?.[0]?.currency || ""
                    };
                  });

                  let color = getCSSVariableValue('--color-text-main'); // Default black
                  if (last6Points.length >= 2) {
                    const first = parseFloat(last6Points[0].amount);
                    const last = parseFloat(last6Points[last6Points.length - 1].amount);
                    if (last > first) color = getCSSVariableValue('--color-primary'); // Green
                    else if (last < first) color = getCSSVariableValue('--color-negative'); // Red
                  }
                  
                  return <Sparkline data={last6Points} color={color} />;
                })()}
              </div>
            </div>
            
            <div>
              <h3 className="text-text-muted/40 font-bold uppercase tracking-widest text-[10px] mb-1">
                {account.lastTransactionAt 
                    ? `Last active: ${new Date(account.lastTransactionAt).toLocaleDateString()}` 
                    : account.id.slice(0, 8)}
              </h3>
              <h2 className="text-2xl font-black text-text-main tracking-tight mb-6">{account.name}</h2>
              
              <div className="space-y-4">
                {Object.entries(account.balance).map(([curr, amt]) => {
                  const latestDist = account.distributions?.find(d => d.currency === curr);
                  const systemAmt = parseFloat(latestDist?.system_amount || '0');
                  const otherAmt = parseFloat(latestDist?.other_amount || '0');

                  return (
                    <div key={curr} className="flex flex-col gap-2 p-4 bg-card-muted rounded-2xl group-hover:bg-card transition-colors duration-500 border border-transparent group-hover:border-border-pearl">
                      <div className="flex items-baseline justify-between">
                        <span className="text-xs font-black text-text-muted">{curr}</span>
                        <div className="flex flex-col items-end">
                          <span className="text-xl font-black text-text-main tracking-tighter">
                            {parseFloat(amt).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                          </span>
                          {(() => {
                            const total = systemAmt + otherAmt;
                            const rate = total !== 0 ? systemAmt / total : 0;
                            if (rate === 1 || total === 0) return null;
                            return (
                              <span className="text-[10px] font-bold text-accent italic">
                                Rate: {(rate * 100).toFixed(1)}%
                              </span>
                            );
                          })()}
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Balances;
