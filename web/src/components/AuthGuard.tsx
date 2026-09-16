"use client";
import { useEffect, useSyncExternalStore } from "react";
import { useRouter } from "next/navigation";
import { useAppContext } from "../app/context/AppContext";

export default function AuthGuard({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAppContext();
  const router = useRouter();
  const mounted = useSyncExternalStore(
    () => () => undefined,
    () => true,
    () => false,
  );

  useEffect(() => {
    if (mounted && !isLoading && !user) {
      router.push("/login");
    }
  }, [user, isLoading, router, mounted]);

  if (!mounted || isLoading) {
    return (
      <div style={{ display:'flex', justifyContent:'center', alignItems:'center', height:'100vh' }}>
        <div style={{ color:'var(--text-secondary)' }}>Loading...</div>
      </div>
    );
  }

  if (!user) return null;
  
  return <>{children}</>;
}
