"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { useAuth } from "@/lib/auth-context"
import { getCompanies } from "@/lib/api"
import CompanyCard from "@/components/CompanyCard"

export default function Home() {
  const { token, email, loading } = useAuth()
  const router = useRouter()
  const [companies, setCompanies] = useState<any[]>([])
  const [fetching, setFetching] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (loading) return
    if (!token) { window.location.href = "/login"; return }

    setFetching(true)
    getCompanies(token)
      .then(setCompanies)
      .catch((e) => setError(String(e)))
      .finally(() => setFetching(false))
  }, [token, loading, router])

  if (loading || fetching) return <LoadingState />

  return (
    <div className="space-y-8">
      {/* Hero */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-primary/90 to-primary px-8 py-10 text-primary-foreground shadow-lg">
        <div className="absolute inset-0 opacity-10" style={{backgroundImage: "url(\"data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='1'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E\")"}} />
        <div className="relative">
          <p className="text-sm font-medium text-primary-foreground/70 mb-1">Welcome back, {email?.split("@")[0]}</p>
          <h1 className="text-3xl font-bold tracking-tight mb-2">Your Farms</h1>
          <p className="text-primary-foreground/80 text-sm max-w-md">
            Track finances, fields, and equipment across all your Farming Simulator 25 companies.
          </p>
        </div>
      </div>

      {error && (
        <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">{error}</div>
      )}

      {companies.length > 0 ? (
        <div>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold">Companies</h2>
            <span className="text-sm text-muted-foreground">{companies.length} save{companies.length !== 1 ? "s" : ""} synced</span>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {companies.map((c) => (
              <CompanyCard key={c.id} company={{ ...c, slotId: c.slotId ?? c.id }} />
            ))}
          </div>
        </div>
      ) : (
        <div className="rounded-2xl border border-dashed border-border p-12 text-center space-y-3">
          <p className="text-4xl">🌱</p>
          <p className="font-semibold">No companies synced yet</p>
          <p className="text-sm text-muted-foreground max-w-sm mx-auto">
            Set up the companion app on your PC and configure it with your companion token to start syncing your farms.
          </p>
          <CompanionSetupInstructions />
        </div>
      )}
    </div>
  )
}

function CompanionSetupInstructions() {
  const { companionToken } = useAuth()
  const [copied, setCopied] = useState(false)

  function copy() {
    if (companionToken) {
      navigator.clipboard.writeText(companionToken)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  return (
    <div className="mt-4 text-left rounded-xl bg-secondary/60 p-4 space-y-3 max-w-md mx-auto">
      <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">Companion Setup</p>
      <ol className="text-sm space-y-2 text-muted-foreground list-decimal list-inside">
        <li>Download and run the FarmSim companion app on your gaming PC</li>
        <li>Open the config file and set <code className="text-xs bg-background px-1 py-0.5 rounded">cloudToken</code> to your companion token</li>
        <li>Set <code className="text-xs bg-background px-1 py-0.5 rounded">cloudApiUrl</code> to your Railway API URL</li>
        <li>Restart the companion — it will sync automatically</li>
      </ol>
      {companionToken && (
        <div className="mt-2">
          <p className="text-xs text-muted-foreground mb-1">Your companion token:</p>
          <div className="flex items-center gap-2">
            <code className="text-xs bg-background px-2 py-1.5 rounded-lg border border-border flex-1 truncate">
              {companionToken}
            </code>
            <button
              onClick={copy}
              className="text-xs px-2 py-1.5 rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors shrink-0"
            >
              {copied ? "Copied!" : "Copy"}
            </button>
          </div>
        </div>
      )}
    </div>
  )
}

function LoadingState() {
  return (
    <div className="space-y-8">
      <div className="h-40 rounded-2xl bg-secondary animate-pulse" />
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="h-48 rounded-2xl bg-secondary animate-pulse" />
        ))}
      </div>
    </div>
  )
}
