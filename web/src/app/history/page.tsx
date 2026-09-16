"use client";

import { useState, useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { api, Job } from "../../lib/api";
import { useAppContext } from "../context/AppContext";

export default function HistoryPage() {
  const { token } = useAppContext();
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("all");

  const { data: jobsResponse, isLoading, isError } = useQuery({
    queryKey: ["jobs", token],
    queryFn: () => api.listJobs(token!),
    enabled: !!token,
  });

  const jobs = useMemo<Job[]>(() => jobsResponse ?? [], [jobsResponse]);

  const filteredJobs = useMemo(() => {
    return jobs.filter((job) => {
      const matchesSearch = job.prompt_input?.template?.toLowerCase().includes(search.toLowerCase()) || 
                            job.id.includes(search);
      const matchesStatus = statusFilter === "all" || job.status === statusFilter;
      return matchesSearch && matchesStatus;
    });
  }, [jobs, search, statusFilter]);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-black pt-28 text-center text-white">
        <div className="w-12 h-12 border-4 border-purple-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
        Loading your video history...
      </div>
    );
  }

  if (isError) {
    return (
      <div className="min-h-screen bg-black pt-28 text-center text-red-400">
        Failed to load history. Please try again.
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-black text-white pt-24 pb-20 px-6">
      <div className="max-w-6xl mx-auto">
        <div className="mb-10">
          <h1 className="text-4xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-600 mb-2">
            Video Generation History
          </h1>
          <p className="text-gray-400">View and download your previously generated cinematic videos.</p>
        </div>

        <div className="flex flex-col md:flex-row gap-4 mb-8">
          <input 
            type="text" 
            placeholder="Search projects by template or ID..." 
            className="flex-1 bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-gray-500 focus:outline-none focus:border-purple-500 transition-colors"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          <select 
            className="bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-purple-500 transition-colors cursor-pointer md:w-48 appearance-none"
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
          >
            <option value="all" className="bg-gray-900">All Status</option>
            <option value="completed" className="bg-gray-900">Completed</option>
            <option value="processing" className="bg-gray-900">Processing</option>
            <option value="in_queue" className="bg-gray-900">In Queue</option>
            <option value="failed" className="bg-gray-900">Failed</option>
          </select>
        </div>

        {filteredJobs.length === 0 ? (
          <div className="text-center mt-20 text-gray-500 bg-white/5 border border-white/10 rounded-3xl p-12">
            <span className="text-4xl block mb-4">📭</span>
            <p className="text-lg">No jobs found.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
            {filteredJobs.map(item => (
              <div key={item.id} className="bg-white/5 border border-white/10 rounded-2xl overflow-hidden hover:border-white/20 transition-all hover:shadow-xl hover:shadow-purple-500/10 flex flex-col">
                <div className="aspect-video relative bg-black border-b border-white/10 flex items-center justify-center overflow-hidden group">
                  {item.status === 'completed' ? (
                    <>
                      <video src={item.gcs_video_url ?? undefined} className="w-full h-full object-cover opacity-80 group-hover:opacity-100 transition-opacity" muted loop playsInline onMouseEnter={(e) => e.currentTarget.play()} onMouseLeave={(e) => e.currentTarget.pause()} />
                      <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity bg-black/30">
                         <span className="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm border border-white/30 text-xl">▶</span>
                      </div>
                    </>
                  ) : (item.status === 'processing' || item.status === 'in_queue' || item.status === 'pending') ? (
                    <div className="flex flex-col items-center gap-3">
                      <div className="w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full animate-spin"></div>
                      <span className="text-xs text-purple-400 uppercase tracking-widest">{item.status}</span>
                    </div>
                  ) : (
                    <div className="text-red-400 text-center">
                      <span className="text-3xl block mb-2">⚠️</span>
                      <span className="text-sm">Failed</span>
                    </div>
                  )}
                  
                  {/* Status Badge */}
                  <div className={`absolute top-3 right-3 text-xs font-bold px-3 py-1 rounded-full ${
                    item.status === 'completed' ? 'bg-green-500/20 text-green-400 border border-green-500/30' :
                    item.status === 'failed' ? 'bg-red-500/20 text-red-400 border border-red-500/30' :
                    'bg-purple-500/20 text-purple-400 border border-purple-500/30'
                  }`}>
                    {item.status}
                  </div>
                </div>
                
                <div className="p-5 flex flex-col flex-1">
                  <h3 className="font-mono text-sm text-gray-300 truncate mb-2" title={item.id}>
                    ID: {item.id.split("-")[0]}...
                  </h3>
                  <div className="flex justify-between items-center text-xs text-gray-500 mb-4 pb-4 border-b border-white/10">
                    <span>{item.prompt_input?.template || 'Custom Prompt'}</span>
                    <span>{new Date(item.created_at).toLocaleDateString()}</span>
                  </div>
                  
                  <div className="mt-auto">
                    <a href={item.gcs_video_url ?? undefined} target="_blank" rel="noopener noreferrer" className={`block w-full text-center py-2.5 rounded-lg font-medium transition-all ${
                      item.status === 'completed' 
                      ? 'bg-white/10 hover:bg-white/20 text-white' 
                      : 'bg-white/5 text-gray-600 cursor-not-allowed pointer-events-none'
                    }`}>
                      Download HD
                    </a>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
