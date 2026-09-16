"use client";
import {
	createContext,
	useContext,
	useState,
	useEffect,
	useSyncExternalStore,
	ReactNode,
} from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api, User } from "../../lib/api";
import { getToken, removeToken, saveToken } from "../../lib/auth";
import { usePathname } from "next/navigation";

type AppContextType = {
	theme: "dark" | "light";
	setTheme: (theme: "dark" | "light") => void;
	lang: string;
	setLang: (lang: string) => void;
	user: User | null;
	setUser: (user: User | null) => void;
	token: string | null;
	setToken: (token: string | null) => void;
	logout: () => void;
	isLoading: boolean;
};

const AppContext = createContext<AppContextType | undefined>(undefined);

type Theme = "dark" | "light";

const subscribeToStorage = (callback: () => void) => {
	window.addEventListener("storage", callback);
	return () => window.removeEventListener("storage", callback);
};

const getStoredTheme = (): Theme =>
	localStorage.getItem("theme") === "light" ? "light" : "dark";

const getStoredLanguage = () => localStorage.getItem("lang") ?? "en";

export function AppProvider({ children }: { children: ReactNode }) {
	const storedToken = useSyncExternalStore(subscribeToStorage, getToken, () => null);
	const storedTheme = useSyncExternalStore<Theme>(
		subscribeToStorage,
		getStoredTheme,
		() => "dark",
	);
	const storedLanguage = useSyncExternalStore(subscribeToStorage, getStoredLanguage, () => "en");
	const [themeOverride, setThemeOverride] = useState<Theme | undefined>(undefined);
	const [languageOverride, setLanguageOverride] = useState<string | undefined>(undefined);
	const [tokenOverride, setTokenOverride] = useState<string | null | undefined>(undefined);
	const theme = themeOverride ?? storedTheme;
	const lang = languageOverride ?? storedLanguage;
	const token = tokenOverride === undefined ? storedToken : tokenOverride;
	const queryClient = useQueryClient();
	const pathname = usePathname();

	const { data: userResponse, isLoading } = useQuery({
		queryKey: ["user", token],
		queryFn: () => api.getMe(token!),
		enabled: !!token && pathname !== "/login", // Only fetch user if token exists and not on login page
		retry: false, // Don't retry auth errors
	});

	const user = userResponse ?? null;

	// Manual setter for optimistic UI updates during login/register
	const setUser = (newUser: User | null) => {
		if (newUser) {
			queryClient.setQueryData(["user", token], newUser);
		} else {
			queryClient.removeQueries({ queryKey: ["user"] });
		}
	};

	const setToken = (newToken: string | null) => {
		if (newToken) {
			saveToken(newToken);
		} else {
			removeToken();
		}
		setTokenOverride(newToken);
	};

	const logout = () => {
		setToken(null);
		setUser(null);
	};

	useEffect(() => {
		document.documentElement.setAttribute("data-theme", theme);
	}, [theme]);

	const handleThemeChange = (newTheme: Theme) => {
		setThemeOverride(newTheme);
		localStorage.setItem("theme", newTheme);
		document.documentElement.setAttribute("data-theme", newTheme);
	};

	const handleLangChange = (newLang: string) => {
		setLanguageOverride(newLang);
		localStorage.setItem("lang", newLang);
	};

	return (
		<AppContext.Provider
			value={{
				theme,
				setTheme: handleThemeChange,
				lang,
				setLang: handleLangChange,
				user,
				setUser,
				token,
				setToken,
				logout,
				isLoading,
			}}
		>
			{children}
		</AppContext.Provider>
	);
}

export function useAppContext() {
	const context = useContext(AppContext);
	if (!context)
		throw new Error("useAppContext must be used within AppProvider");
	return context;
}
