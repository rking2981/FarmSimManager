"use client"

import Link from "next/link"
import { useRouter } from "next/navigation"
import { useAuth } from "@/lib/auth-context"

export default function Header() {
  const { token, email, clearAuth } = useAuth()
  const router = useRouter()

  function handleLogout() {
    clearAuth()
    router.push("/login")
  }

  return (
    <header className="sticky top-0 z-40 w-full border-b border-border/60 bg-card/80 backdrop-blur-md">
      <div className="max-w-7xl mx-auto px-6 h-14 flex items-center gap-3">
        <Link href="/" className="flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center text-primary-foreground text-base shadow-sm">
            🌾
          </div>
          <span className="font-semibold text-base tracking-tight">
            FarmSim <span className="text-primary">Manager</span>
          </span>
        </Link>

        <div className="flex-1" />

        {token ? (
          <div className="flex items-center gap-3">
            <span className="text-xs text-muted-foreground hidden sm:block">{email}</span>
            <button
              onClick={handleLogout}
              className="text-xs px-3 py-1.5 rounded-lg bg-secondary text-secondary-foreground hover:bg-secondary/70 transition-colors"
            >
              Sign out
            </button>
          </div>
        ) : (
          <div className="flex items-center gap-2">
            <Link href="/login" className="text-xs px-3 py-1.5 rounded-lg bg-secondary text-secondary-foreground hover:bg-secondary/70 transition-colors">
              Sign in
            </Link>
            <Link href="/register" className="text-xs px-3 py-1.5 rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors">
              Sign up
            </Link>
          </div>
        )}
      </div>
    </header>
  )
}
