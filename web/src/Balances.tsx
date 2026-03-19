import React, { useEffect, useState, useMemo } from 'react';
import { 
  Loader2, 
  Wallet, 
  Calendar,
  BarChart3
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
  Filler
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import { api } from './api';
import type { Account, BalancePeriod } from './api';

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
}

const Sparkline: React.FC<{ data: BalancePeriod[], color: string }> = ({ data, color }) => {
  const last6Months = useMemo(() => {
    // Ensure we have a sorted history, take last 6 points
    const sorted = [...data].sort((a, b) => new Date(a.month).getTime() - new Date(b.month).getTime());
    return sorted.slice(-6);
  }, [data]);

  const chartData = {
    labels: last6Months.map(h => h.month),
    datasets: [
      {
        data: last6Months.map(h => parseFloat(h.amount)),
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
  const [balanceHistory, setBalanceHistory] = useState<BalancePeriod[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [accs, history, transactions] = await Promise.all([
          api.getAccounts(),
          api.getBalances(),
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
          
          // Fetch history for sparkline
          const accHistory = await api.getBalancesByAccount(acc.id);

          return {
            ...acc,
            lastTransactionAt: lastTx,
            history: accHistory
          };
        }));

        // Sort by last transaction (most recent first)
        accountsWithMetadata.sort((a, b) => {
          if (!a.lastTransactionAt) return 1;
          if (!b.lastTransactionAt) return -1;
          return new Date(b.lastTransactionAt).getTime() - new Date(a.lastTransactionAt).getTime();
        });

        setAccounts(accountsWithMetadata);
        setBalanceHistory(history || []);
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
    const months = Array.from(new Set(balanceHistory.map(h => {
        const date = new Date(h.month);
        return date.toLocaleDateString(undefined, { month: 'short', year: '2-digit' });
    }))).reverse();

    const currencies = Array.from(new Set(balanceHistory.map(h => h.currency)));
    
    const datasets = currencies.map((curr, index) => {
      const data = months.map(m => {
        const h = balanceHistory.find(history => {
            const date = new Date(history.month);
            const label = date.toLocaleDateString(undefined, { month: 'short', year: '2-digit' });
            return label === m && history.currency === curr;
        });
        return h ? parseFloat(h.amount) : null;
      });

      const colors = [
        { border: 'rgb(99, 102, 241)', bg: 'rgba(99, 102, 241, 0.1)' }, // Indigo
        { border: 'rgb(16, 185, 129)', bg: 'rgba(16, 185, 129, 0.1)' }, // Emerald
        { border: 'rgb(244, 63, 94)', bg: 'rgba(244, 63, 94, 0.1)' },   // Rose
        { border: 'rgb(59, 130, 246)', bg: 'rgba(59, 130, 246, 0.1)' }, // Blue
      ];
      const color = colors[index % colors.length];

      return {
        label: curr,
        data: data,
        borderColor: color.border,
        backgroundColor: color.bg,
        fill: true,
        tension: 0.4,
        pointRadius: 0, // Hidden by default
        hoverPointRadius: 6,
        pointBackgroundColor: color.border,
        borderWidth: 3,
        pointHitRadius: 10,
      };
    });

    return {
      labels: months,
      datasets: datasets
    };
  }, [balanceHistory]);

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
          font: {
            family: 'inherit',
            weight: 'bold' as any,
            size: 11
          }
        }
      },
      datalabels: {
        display: false
      },
      tooltip: {
        backgroundColor: 'rgba(255, 255, 255, 0.9)',
        titleColor: '#111827',
        bodyColor: '#111827',
        borderColor: '#e5e7eb',
        borderWidth: 1,
        padding: 12,
        boxPadding: 6,
        usePointStyle: true,
        callbacks: {
          label: (context: any) => {
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
            weight: 'bold' as any
          },
          color: '#9ca3af'
        }
      },
      y: {
        grid: {
          color: 'rgba(0, 0, 0, 0.03)',
        },
        ticks: {
          font: {
            size: 10,
            weight: 'bold' as any
          },
          color: '#9ca3af',
          callback: (value: any) => {
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
    <div className="flex items-center justify-center p-20">
      <Loader2 className="animate-spin text-gray-300" size={48} strokeWidth={1} />
    </div>
  );

  if (error) return (
    <div className="text-center p-20">
      <p className="text-rose-500 font-bold">Failed to load balance overview</p>
    </div>
  );

  return (
    <div className="max-w-6xl mx-auto w-full space-y-8 pb-20">
      {/* Header Section */}
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 px-4">
        <div>
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-indigo-50 text-indigo-500 text-[10px] font-black uppercase tracking-widest mb-4">
            <BarChart3 size={12} />
            Financial Overview
          </div>
          <h1 className="text-5xl font-black text-gray-900 tracking-tighter">Balances</h1>
          <p className="text-gray-400 font-bold mt-2 uppercase tracking-[0.2em] text-[10px]">Net worth and account performance</p>
        </div>
      </div>

      {/* Accounts Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 px-4">
        {accounts.map(account => (
          <div key={account.id} className="bg-white/80 backdrop-blur-xl p-8 rounded-[40px] shadow-sm border border-white hover:shadow-xl transition-all duration-500 group">
            <div className="flex items-start justify-between mb-8">
              <div className="p-4 bg-gray-50 rounded-3xl group-hover:bg-indigo-50 transition-colors duration-500">
                <Wallet className="text-gray-400 group-hover:text-indigo-500 transition-colors duration-500" size={24} />
              </div>
              <div className="flex flex-col items-end">
                <span className="text-[10px] font-black text-gray-300 uppercase tracking-widest mb-1">6M Trend</span>
                {account.history && <Sparkline data={account.history} color="rgb(99, 102, 241)" />}
              </div>
            </div>
            
            <div>
              <h3 className="text-gray-400 font-bold uppercase tracking-widest text-[10px] mb-1">
                {account.lastTransactionAt 
                    ? `Last active: ${new Date(account.lastTransactionAt).toLocaleDateString()}` 
                    : account.id.slice(0, 8)}
              </h3>
              <h2 className="text-2xl font-black text-gray-900 tracking-tight mb-6">{account.name}</h2>
              
              <div className="space-y-4">
                {Object.entries(account.balance).map(([curr, amt]) => (
                  <div key={curr} className="flex items-baseline justify-between p-4 bg-gray-50/50 rounded-2xl group-hover:bg-white transition-colors duration-500">
                    <span className="text-xs font-black text-gray-400">{curr}</span>
                    <span className="text-xl font-black text-gray-900 tracking-tighter">
                      {parseFloat(amt).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* History Graph */}
      <div className="px-4">
        <div className="bg-white/80 backdrop-blur-xl p-8 md:p-12 rounded-[48px] shadow-sm border border-white">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-6 mb-12">
            <div>
              <h2 className="text-2xl font-black text-gray-900 tracking-tight">Balance History</h2>
              <p className="text-gray-400 font-bold mt-1 uppercase tracking-widest text-[9px]">Monthly progression across currencies</p>
            </div>
            <div className="flex items-center gap-4">
                <div className="flex items-center gap-2 px-4 py-2 bg-gray-50 rounded-2xl border border-gray-100">
                    <Calendar size={14} className="text-gray-400" />
                    <span className="text-[10px] font-bold text-gray-600 uppercase tracking-widest">Timeline</span>
                </div>
            </div>
          </div>
          
          <div className="h-[400px] w-full">
            <Line data={chartData} options={chartOptions} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Balances;
