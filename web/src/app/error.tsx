"use client";
import { useEffect } from "react";
import Link from "next/link";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    // Optionally log the error to an error reporting service
    console.error(error);
  }, [error]);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '100vh', textAlign: 'center', padding: '20px' }}>
      <h2 style={{ fontSize: '2rem', marginBottom: '16px', color: '#ef4444' }}>Something went wrong!</h2>
      <p style={{ color: 'var(--text-secondary)', marginBottom: '24px', maxWidth: '500px' }}>
        {error.message || "An unexpected error occurred. Please try again."}
      </p>
      <div style={{ display: 'flex', gap: '16px' }}>
        <button
          className="primary-btn"
          onClick={() => reset()}
        >
          Try again
        </button>
        <Link href="/">
          <button className="secondary-btn">Go back home</button>
        </Link>
      </div>
    </div>
  );
}
