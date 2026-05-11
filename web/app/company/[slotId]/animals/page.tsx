"use client"

import { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import { useAuth } from "@/lib/auth-context"
import { getAnimals } from "@/lib/api"
import CompanyNav from "@/components/CompanyNav"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

interface Animal {
  subType: string
  name: string
  category: string
  numAnimals: number
  ageMonths: number
  healthPct: number
  reproductionPct: number
  basePrice: number
  estimatedValue: number
}

const CATEGORY_ICONS: Record<string, string> = {
  Cow: "🐄", Sheep: "🐑", Pig: "🐖", Horse: "🐴", Chicken: "🐔", Goat: "🐐", Animal: "🐾",
}

function fmt(n: number) {
  if (n >= 1_000_000) return `$${(n / 1_000_000).toFixed(2)}M`
  if (n >= 1_000) return `$${(n / 1_000).toFixed(0)}K`
  return `$${n.toFixed(0)}`
}

function Bar({ value, color }: { value: number; color: string }) {
  return (
    <div className="flex items-center gap-2">
      <div className="flex-1 h-1.5 rounded-full bg-secondary overflow-hidden">
        <div className={`h-full rounded-full ${color}`} style={{ width: `${Math.min(value, 100)}%` }} />
      </div>
      <span className="text-xs tabular-nums w-8 text-right">{value.toFixed(0)}%</span>
    </div>
  )
}

function ageLabel(months: number) {
  const years = Math.floor(months / 12)
  const rem = months % 12
  if (years === 0) return `${months}m`
  if (rem === 0) return `${years}y`
  return `${years}y ${rem}m`
}

export default function AnimalsPage() {
  const { slotId } = useParams<{ slotId: string }>()
  const router = useRouter()
  const { token, loading: authLoading } = useAuth()
  const [animals, setAnimals] = useState<Animal[]>([])
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (authLoading) return
    if (!token) { window.location.href = "/login"; return }
    getAnimals(token, slotId)
      .then(setAnimals)
      .catch((e) => setError(String(e)))
      .finally(() => setLoading(false))
  }, [slotId, token, authLoading, router])

  const totalValue = animals.reduce((s, a) => s + a.estimatedValue, 0)
  const totalCount = animals.reduce((s, a) => s + a.numAnimals, 0)

  if (loading) return (
    <div className="max-w-6xl mx-auto space-y-6">
      <div className="h-8 w-48 bg-secondary rounded animate-pulse" />
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        {Array.from({ length: 4 }).map((_, i) => <div key={i} className="h-24 bg-secondary rounded-2xl animate-pulse" />)}
      </div>
    </div>
  )

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <CompanyNav slotId={slotId} />

      <div className="flex items-center gap-3">
        <h2 className="text-xl font-semibold">Animals</h2>
        <span className="text-sm text-muted-foreground">{totalCount} owned</span>
      </div>

      {error && <p className="text-sm text-destructive">{error}</p>}

      {/* Summary cards */}
      {animals.length > 0 && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
          <StatCard icon="🐾" label="Total Animals" value={String(totalCount)} />
          <StatCard icon="🏷️" label="Fleet Value" value={fmt(totalValue)} accent="primary" />
          <StatCard icon="❤️" label="Avg Health" value={`${(animals.reduce((s,a) => s + a.healthPct, 0) / animals.length).toFixed(0)}%`} accent={animals.reduce((s,a) => s + a.healthPct, 0) / animals.length >= 70 ? "primary" : "destructive"} />
          <StatCard icon="⚡" label="Avg Productivity" value={`${(animals.reduce((s,a) => s + a.reproductionPct, 0) / animals.length).toFixed(0)}%`} />
        </div>
      )}

      {animals.length === 0 ? (
        <div className="rounded-2xl border border-dashed border-border p-12 text-center space-y-3">
          <p className="text-4xl">🐄</p>
          <p className="font-semibold">No animals owned</p>
          <p className="text-sm text-muted-foreground">Purchase animals from the animal dealer to see them here.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {animals.map((a) => {
            const icon = CATEGORY_ICONS[a.category] ?? "🐾"
            const healthColor = a.healthPct >= 70 ? "bg-primary" : a.healthPct >= 40 ? "bg-amber-500" : "bg-destructive"
            const prodColor = a.reproductionPct >= 70 ? "bg-primary" : a.reproductionPct >= 40 ? "bg-amber-500" : "bg-destructive"
            return (
              <div key={a.subType} className="rounded-2xl bg-card border border-border/60 card-shadow overflow-hidden">
                <div className="h-1.5 bg-gradient-to-r from-primary/60 to-accent/60" />
                <div className="p-4 space-y-3">
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-2.5">
                      <div className="w-9 h-9 rounded-xl bg-secondary flex items-center justify-center text-xl">
                        {icon}
                      </div>
                      <div>
                        <p className="font-semibold text-sm">{a.name}</p>
                        <p className="text-xs text-muted-foreground">{a.category}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-semibold text-primary">{a.numAnimals}x</p>
                      <p className="text-xs text-muted-foreground">{ageLabel(a.ageMonths)}</p>
                    </div>
                  </div>

                  <div className="space-y-1.5">
                    <div>
                      <p className="text-[10px] text-muted-foreground mb-1">Health</p>
                      <Bar value={a.healthPct} color={healthColor} />
                    </div>
                    <div>
                      <p className="text-[10px] text-muted-foreground mb-1">Productivity</p>
                      <Bar value={a.reproductionPct} color={prodColor} />
                    </div>
                  </div>

                  <div className="flex justify-between pt-2 border-t border-border/50 text-xs text-muted-foreground">
                    <span>Est. value</span>
                    <span className="font-semibold text-foreground">{fmt(a.estimatedValue)}</span>
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}

function StatCard({ icon, label, value, accent }: { icon: string; label: string; value: string; accent?: "primary" | "destructive" }) {
  const color = accent === "primary" ? "text-primary" : accent === "destructive" ? "text-destructive" : "text-foreground"
  return (
    <div className="rounded-2xl bg-card border border-border/60 card-shadow p-4">
      <div className="flex items-center gap-2 mb-1">
        <span>{icon}</span>
        <span className="text-xs text-muted-foreground">{label}</span>
      </div>
      <p className={`text-xl font-semibold ${color}`}>{value}</p>
    </div>
  )
}
