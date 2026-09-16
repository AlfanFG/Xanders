"use client";
import { useState } from "react";
import { useAppContext } from "../context/AppContext";

export default function SettingsPage() {
  const { theme, setTheme, lang, setLang, user } = useAppContext();
  const [activeTab, setActiveTab] = useState("appearance");
  
  if (!user) return null;

  return (
    <div className="min-h-screen bg-black text-white pt-24 pb-20 px-6">
      <div className="max-w-6xl mx-auto">
        <div className="mb-10">
          <h1 className="text-4xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-600 mb-2">
            Settings
          </h1>
          <p className="text-gray-400">Manage your account preferences and application settings.</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          <div className="lg:col-span-3 flex flex-col gap-2">
            {[
              { id: 'appearance', icon: '🎨', label: 'Appearance' },
              { id: 'language', icon: '🌐', label: 'Language' },
              { id: 'notifications', icon: '🔔', label: 'Notifications' }
            ].map(tab => (
              <button 
                key={tab.id}
                className={`w-full text-left px-5 py-3 rounded-xl font-medium transition-all flex items-center gap-3 ${
                  activeTab === tab.id 
                  ? 'bg-purple-500/20 text-purple-300 border border-purple-500/30' 
                  : 'bg-transparent text-gray-400 hover:bg-white/5 hover:text-white border border-transparent'
                }`}
                onClick={() => setActiveTab(tab.id)}
              >
                <span>{tab.icon}</span> {tab.label}
              </button>
            ))}
          </div>

          <div className="lg:col-span-9">
            <div className="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-8 shadow-xl min-h-[400px]">
              
              {activeTab === 'appearance' && (
                <div className="animate-in fade-in duration-300">
                  <h2 className="text-2xl font-bold mb-2">Appearance</h2>
                  <p className="text-gray-400 mb-8">Customize how Xanders looks on your device.</p>
                  
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 max-w-2xl">
                    <div 
                      className={`rounded-2xl border-2 p-1 cursor-pointer transition-all ${theme === 'dark' ? 'border-purple-500 bg-purple-500/10' : 'border-white/10 hover:border-white/30 bg-black/40'}`}
                      onClick={() => setTheme('dark')}
                    >
                      <div className="h-32 rounded-xl bg-black border border-white/10 flex items-center justify-center">
                        <span className="text-white font-medium">Dark UI</span>
                      </div>
                      <div className="p-3 text-center font-medium">Dark Mode</div>
                    </div>
                    
                    <div 
                      className={`rounded-2xl border-2 p-1 cursor-pointer transition-all ${theme === 'light' ? 'border-purple-500 bg-purple-500/10' : 'border-white/10 hover:border-white/30 bg-black/40'}`}
                      onClick={() => setTheme('light')}
                    >
                      <div className="h-32 rounded-xl bg-gray-100 border border-gray-200 flex items-center justify-center">
                        <span className="text-black font-medium">Light UI</span>
                      </div>
                      <div className="p-3 text-center font-medium">Light Mode</div>
                    </div>
                  </div>
                </div>
              )}

              {activeTab === 'language' && (
                <div className="animate-in fade-in duration-300">
                  <h2 className="text-2xl font-bold mb-2">Language Preferences</h2>
                  <p className="text-gray-400 mb-8">Select your preferred language for the interface.</p>
                  
                  <div className="max-w-md">
                    <select 
                      className="w-full bg-black/40 border border-white/10 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-purple-500 transition-colors" 
                      value={lang} 
                      onChange={(e) => setLang(e.target.value)}
                    >
                      <option value="en" className="bg-gray-900">English (US)</option>
                      <option value="id" className="bg-gray-900">Bahasa Indonesia</option>
                      <option value="es" className="bg-gray-900">Español</option>
                    </select>
                  </div>
                </div>
              )}

              {activeTab === 'notifications' && (
                <div className="animate-in fade-in duration-300">
                  <h2 className="text-2xl font-bold mb-2">Email Notifications</h2>
                  <p className="text-gray-400 mb-8">Choose what updates you want to receive.</p>
                  
                  <div className="space-y-6 max-w-2xl">
                    <div className="flex items-center justify-between p-4 rounded-2xl bg-black/40 border border-white/5">
                      <div>
                        <h4 className="font-bold text-white mb-1">Video Completion</h4>
                        <p className="text-sm text-gray-500">Get notified when your AI video finishes rendering.</p>
                      </div>
                      <label className="relative inline-flex items-center cursor-pointer">
                        <input type="checkbox" className="sr-only peer" defaultChecked />
                        <div className="w-11 h-6 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-500"></div>
                      </label>
                    </div>
                    
                    <div className="flex items-center justify-between p-4 rounded-2xl bg-black/40 border border-white/5">
                      <div>
                        <h4 className="font-bold text-white mb-1">Low Credit Alert</h4>
                        <p className="text-sm text-gray-500">Receive an email when you have less than 50 credits left.</p>
                      </div>
                      <label className="relative inline-flex items-center cursor-pointer">
                        <input type="checkbox" className="sr-only peer" defaultChecked />
                        <div className="w-11 h-6 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-500"></div>
                      </label>
                    </div>
                  </div>
                </div>
              )}
              
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
