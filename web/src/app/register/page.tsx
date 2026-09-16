"use client";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { api } from "../../lib/api";
import { useAppContext } from "../context/AppContext";

export default function RegisterPage() {
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const router = useRouter();
	const { setToken, setUser } = useAppContext();

	const mutation = useMutation({
		mutationFn: () => api.register(email, password),
		onSuccess: (response) => {
			if (response) {
				setToken(response.token);
				setUser(response.user);
				router.push("/");
			}
		},
	});

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		if (email && password) {
			mutation.mutate();
		}
	};

	return (
		<div className="min-h-screen bg-black flex items-center justify-center p-6 selection:bg-purple-500/30">
			<div className="w-full max-w-md bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-8 shadow-2xl relative overflow-hidden">
				<div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-purple-600 via-pink-500 to-red-500"></div>
				
				<div className="mb-8 text-center">
					<h1 className="text-3xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-white to-gray-400 mb-2">
						Create Account
					</h1>
					<p className="text-gray-400 text-sm">Join today and get 15 free credits instantly.</p>
				</div>

				<form className="space-y-6" onSubmit={handleSubmit}>
					<div className="space-y-2">
						<label className="text-sm font-medium text-gray-300">Email Address</label>
						<input
							type="email"
							placeholder="you@company.com"
							className="w-full bg-black/40 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-gray-600 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500 transition-all"
							value={email}
							onChange={(e) => setEmail(e.target.value)}
							required
						/>
					</div>
					<div className="space-y-2">
						<label className="text-sm font-medium text-gray-300">Password</label>
						<input
							type="password"
							placeholder="••••••••"
							className="w-full bg-black/40 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-gray-600 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500 transition-all"
							value={password}
							onChange={(e) => setPassword(e.target.value)}
							required
						/>
					</div>

					{mutation.isError && (
						<div className="text-red-400 text-sm bg-red-400/10 p-3 rounded-lg border border-red-400/20">
							{mutation.error.message || "Registration failed. Try a different email."}
						</div>
					)}

					<button
						type="submit"
						className="w-full py-3 rounded-xl font-bold text-white bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 shadow-lg hover:shadow-pink-500/25 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
						disabled={mutation.isPending}
					>
						{mutation.isPending ? "Creating Account..." : "Create Account"}
					</button>
				</form>

				<p className="mt-8 text-center text-sm text-gray-400">
					Already have an account?{" "}
					<Link href="/login" className="text-purple-400 hover:text-purple-300 font-medium transition-colors">
						Login here
					</Link>
				</p>
			</div>
		</div>
	);
}
