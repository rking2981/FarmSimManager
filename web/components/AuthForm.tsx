"use client"

import { useState } from "react"
import Link from "next/link"
import { Button } from "@/components/ui/button"

interface Props {
  mode: "login" | "register"
  onSubmit: (email: string, password: string) => Promise<void>
  error: string | null
  loading: boolean
}

export default function AuthForm({ mode, onSubmit, error, loading }: Props) {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    await onSubmit(email, password)
  }

  return (
    <div className="min-h-[80vh] flex items-center justify-center">
      <div className="w-full max-w-sm">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="w-14 h-14 rounded-2xl bg-primary flex items-center justify-center text-3xl mx-auto mb-4 shadow-md">
            🌾
          </div>
          <h1 className="text-2xl font-bold tracking-tight">
            {mode === "login" ? "Welcome back" : "Create account"}
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            {mode === "login"
              ? "Sign in to your FarmSim Manager account"
              : "Start managing your farms across devices"}
          </p>
        </div>

        <div className="rounded-2xl bg-card border border-border/60 card-shadow p-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-1.5">
              <label className="text-sm font-medium">Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@example.com"
                required
                className="w-full rounded-xl border border-input bg-background px-3.5 py-2.5 text-sm outline-none focus:ring-2 focus:ring-ring transition-shadow"
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-sm font-medium">Password</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder={mode === "register" ? "Min. 8 characters" : "••••••••"}
                required
                minLength={mode === "register" ? 8 : undefined}
                className="w-full rounded-xl border border-input bg-background px-3.5 py-2.5 text-sm outline-none focus:ring-2 focus:ring-ring transition-shadow"
              />
            </div>

            {error && (
              <p className="text-sm text-destructive bg-destructive/10 rounded-lg px-3 py-2">{error}</p>
            )}

            <Button
              type="submit"
              disabled={loading}
              className="w-full rounded-xl bg-primary text-primary-foreground hover:bg-primary/90 py-2.5"
            >
              {loading ? "Please wait..." : mode === "login" ? "Sign In" : "Create Account"}
            </Button>
          </form>

          <p className="text-center text-sm text-muted-foreground mt-4">
            {mode === "login" ? (
              <>Don&apos;t have an account?{" "}
                <Link href="/register" className="text-primary font-medium hover:underline">Sign up</Link>
              </>
            ) : (
              <>Already have an account?{" "}
                <Link href="/login" className="text-primary font-medium hover:underline">Sign in</Link>
              </>
            )}
          </p>
        </div>
      </div>
    </div>
  )
}
