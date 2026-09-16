"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAppContext } from "./context/AppContext";
import { useRef, useState } from "react";
import { useMotionValueEvent, useScroll, motion, Variants } from "motion/react";

export default function Navbar() {
	const { theme, setTheme, lang, setLang, user, logout } = useAppContext();
	const [hidden, setHidden] = useState<boolean>(false);
	const [showDropdown, setShowDropdown] = useState<boolean>(false);
	const lastYRef = useRef(0);
	const { scrollY } = useScroll();
	const router = useRouter();

	useMotionValueEvent(scrollY, "change", (y: number) => {
		const difference = y - lastYRef.current;

		if (Math.abs(difference) > 180) {
			setHidden(difference > 0);
			lastYRef.current = y;
			setShowDropdown(false);
		}
	});

	const toggleTheme = () => {
		setTheme(theme === "dark" ? "light" : "dark");
	};

	const handleLogout = () => {
		logout();
		setShowDropdown(false);
		router.push("/login");
	};

	return (
		<motion.div
			animate={hidden ? "hidden" : "visible"}
			initial={"visible"}
			whileHover={hidden ? "peeking" : "visible"}
			onFocusCapture={hidden ? () => setHidden(false) : undefined}
			variants={
				{
					visible: { y: "0%" },
					hidden: { y: "-90%" },
					peeking: { y: "0%", cursor: "pointer" },
				} as Variants
			}
			transition={{ duration: 0.2 }}
			className="fixed top-0 left-0 w-full z-50 p-4"
		>
			<div className="max-w-6xl mx-auto flex items-center justify-between h-16 bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl px-6 shadow-2xl">
				<div className="text-2xl font-extrabold tracking-tight">
					<Link href="/" className="bg-clip-text text-transparent bg-gradient-to-r from-cyan-400 via-blue-500 to-purple-600">
						Xanders
					</Link>
				</div>
				
				<div className="flex items-center gap-6">
					<div className="flex items-center gap-2">
						<button 
							onClick={toggleTheme} 
							className="w-8 h-8 rounded-full bg-white/5 border border-white/10 flex items-center justify-center text-sm hover:bg-white/10 transition-colors"
						>
							{theme === "dark" ? "☀️" : "🌙"}
						</button>
						<select
							className="bg-white/5 border border-white/10 rounded-lg px-2 py-1 text-sm text-white focus:outline-none focus:border-purple-500 cursor-pointer"
							value={lang}
							onChange={(e) => setLang(e.target.value)}
						>
							<option value="en" className="bg-black">EN</option>
							<option value="id" className="bg-black">ID</option>
						</select>
					</div>

					{user ? (
						<div className="flex items-center gap-6">
							<Link href="/history" className="text-sm font-medium text-gray-300 hover:text-white transition-colors">
								History
							</Link>
							
							<div className="bg-purple-500/20 text-purple-400 text-xs font-bold px-3 py-1.5 rounded-full border border-purple-500/30">
								{user.credit_balance} Credits
							</div>
							
							<div className="relative">
								<div 
									className="w-9 h-9 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 cursor-pointer shadow-lg border-2 border-transparent hover:border-white/50 transition-colors"
									onClick={() => setShowDropdown(!showDropdown)}
								></div>
								
								{showDropdown && (
									<div className="absolute top-[120%] right-0 min-w-[150px] bg-black/80 backdrop-blur-xl border border-white/10 rounded-xl p-2 flex flex-col gap-1 shadow-2xl z-50">
										<Link href="/profile" className="px-3 py-2 text-sm text-gray-300 hover:text-white hover:bg-white/10 rounded-lg transition-colors" onClick={() => setShowDropdown(false)}>
											Profile
										</Link>
										<div 
											className="px-3 py-2 text-sm text-red-400 hover:text-red-300 hover:bg-red-400/10 rounded-lg cursor-pointer transition-colors"
											onClick={handleLogout}
										>
											Sign Out
										</div>
									</div>
								)}
							</div>
						</div>
					) : (
						<div className="flex items-center gap-4">
							<Link href="/login" className="text-sm font-medium text-gray-300 hover:text-white transition-colors">Login</Link>
							<Link href="/register">
								<button className="px-4 py-2 text-sm font-bold text-white bg-gradient-to-r from-blue-600 to-purple-600 rounded-xl shadow-lg hover:shadow-purple-500/25 transition-all">
									Register
								</button>
							</Link>
						</div>
					)}
				</div>
			</div>
		</motion.div>
	);
}
