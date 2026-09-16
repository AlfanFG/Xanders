"use client";
import { useAppContext } from "../context/AppContext";
import { useQuery } from "@tanstack/react-query";
import { api } from "../../lib/api";
import { useRouter } from "next/navigation";

export default function ProfilePage() {
  const { user, token, logout } = useAppContext();
  const router = useRouter();

  const { data: historyResponse } = useQuery({
    queryKey: ["creditHistory", token],
    queryFn: () => api.getCreditHistory(token!),
    enabled: !!token,
  });

	const transactions = historyResponse || [];
	const creditRequestHref = `mailto:alfanfaturahman10@gmail.com?subject=${encodeURIComponent("Xanders credit request")}&body=${encodeURIComponent(`Hello, I would like to request more Xanders video-generation credits.\n\nAccount: ${user?.email ?? ""}\nCurrent balance: ${user?.credit_balance ?? 0}`)}`;

  if (!user) return null;

  return (
    <div className="min-h-screen bg-black text-white pt-24 pb-20 px-6">
      <div className="max-w-6xl mx-auto">
        <div className="mb-10">
          <h1 className="text-4xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-600 mb-2">
            My Profile
          </h1>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          <div className="lg:col-span-4 flex flex-col gap-6">
            <div className="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-8 flex flex-col items-center text-center shadow-xl">
              <div className="w-24 h-24 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 mb-4 shadow-lg"></div>
              <h2 className="text-xl font-bold break-all mb-2">{user.email}</h2>
              <div className="bg-purple-500/20 text-purple-400 text-xs font-bold px-3 py-1 rounded-full border border-purple-500/30 uppercase tracking-wide mb-8">
                {user.plan} Plan
              </div>
              
              <div className="w-full grid grid-cols-2 gap-4 mb-8 text-left">
                <div className="bg-black/40 rounded-2xl p-4 border border-white/5">
                  <span className="block text-xs text-gray-500 mb-1 font-medium">Credits</span>
                  <span className="text-2xl font-bold text-white">{user.credit_balance}</span>
                </div>
                <div className="bg-black/40 rounded-2xl p-4 border border-white/5">
                  <span className="block text-xs text-gray-500 mb-1 font-medium">Joined</span>
                  <span className="text-sm font-bold text-white block mt-1.5">
                    {new Date(user.created_at).toLocaleDateString()}
                  </span>
                </div>
              </div>
              
			  <a
				  href={creditRequestHref}
				  className="w-full py-3 rounded-xl font-bold text-center text-white bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 shadow-lg mb-3 transition-all"
			  >
				  Request More Credits
			  </a>
              <button 
                className="w-full py-3 rounded-xl font-medium text-red-400 bg-red-500/10 border border-red-500/20 hover:bg-red-500/20 transition-all"
                onClick={() => {
                  logout();
                  router.push("/login");
                }}
              >
                Sign Out
              </button>
            </div>
          </div>

          <div className="lg:col-span-8 flex flex-col gap-8">
            <div className="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-8 shadow-xl">
              <h3 className="text-xl font-bold mb-6">Credit Usage History</h3>
              <div className="h-48 border-b border-white/10 relative flex items-end justify-between px-4 pb-4">
                {transactions.length === 0 ? (
                  <div className="absolute inset-0 flex items-center justify-center text-gray-500">No transaction history.</div>
                ) : (
                  <div className="w-full h-full flex items-end justify-around gap-2 pt-4">
                    {transactions.slice(0, 7).map((t, i) => {
                      const h = Math.min(Math.abs(t.amount) * 5, 100);
                      return (
                        <div key={i} className="w-full max-w-[40px] flex flex-col justify-end h-full relative group">
                          <div className="w-full rounded-t-sm transition-all duration-300" style={{ height: `${h}%`, backgroundColor: t.amount < 0 ? '#ef4444' : '#10b981' }}></div>
                          <div className="absolute -top-8 left-1/2 -translate-x-1/2 bg-black text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap z-10 pointer-events-none">
                            {t.amount}
                          </div>
                        </div>
                      );
                    })}
                  </div>
                )}
              </div>
            </div>
            
            <div className="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-8 shadow-xl">
              <h3 className="text-xl font-bold mb-6">Recent Activity</h3>
              {transactions.length === 0 ? (
                <div className="text-gray-500">No recent activity.</div>
              ) : (
                <ul className="space-y-4">
                  {transactions.slice(0, 5).map(t => (
                    <li key={t.id} className="flex items-center gap-4 p-4 rounded-2xl bg-black/40 border border-white/5">
                      <div className="w-12 h-12 rounded-xl bg-white/5 flex items-center justify-center text-2xl flex-shrink-0">
                        {t.amount < 0 ? '🎬' : '💳'}
                      </div>
                      <div className="flex-1">
                        <h4 className="font-bold text-gray-200 capitalize">{t.action.replace("_", " ")}</h4>
                        <span className="text-xs text-gray-500">{new Date(t.created_at).toLocaleDateString()}</span>
                      </div>
                      <div className={`font-bold ${t.amount < 0 ? 'text-red-400' : 'text-green-400'}`}>
                        {t.amount > 0 ? '+' : ''}{t.amount} Credits
                      </div>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
