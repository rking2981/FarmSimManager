"use client"

import { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import { getLiveData } from "@/lib/api"
import { useAuth } from "@/lib/auth-context"

import CompanyNav from "@/components/CompanyNav"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"

interface GameTime {
  year: number; month: number; day: number
  hour: number; minute: number; season: string
}
interface CropPrice { name: string; title: string; pricePerLiter: number; basePrice: number }
const ANIMAL_PREFIXES = ["COW_", "SHEEP_", "PIG_", "HORSE_", "CHICKEN_", "GOAT_"]
const ANIMAL_ICONS: Record<string, string> = {
  COW_: "🐄", SHEEP_: "🐑", PIG_: "🐖", HORSE_: "🐴", CHICKEN_: "🐔", GOAT_: "🐐",
}
interface Contract { type: string; fieldId: number; reward: number; completion: number; isActive: boolean }
interface Animal { type: string; title: string; count: number; healthPct: number; productivityPct: number; avgAgeMonths: number }
interface AnimalPriceTier { label: string; months: number; price: number }
interface AnimalPrice { name: string; title: string; basePrice: number; tiers: AnimalPriceTier[] }
interface Worker { task: string; vehicle: string; wagePerHour: number }
interface LiveData {
  available: boolean
  exportedAt?: string
  pushedAt?: string
  gameTime?: GameTime
  cropPrices?: CropPrice[]
  animalPrices?: AnimalPrice[]
  contracts?: Contract[]
  animals?: Animal[]
  workers?: Worker[]
}

function fmt(n: number) {
  return new Intl.NumberFormat("en-US", { style: "currency", currency: "USD", maximumFractionDigits: 2 }).format(n)
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

export default function LivePage() {
  const { slotId } = useParams<{ slotId: string }>()
  const router = useRouter()
  const { token, loading: authLoading } = useAuth()
  const [data, setData] = useState<LiveData | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (authLoading) return
    if (!token) { window.location.href = "/login"; return }
    getLiveData(token, slotId)
      .then(setData)
      .catch((e) => setError(String(e)))
      .finally(() => setLoading(false))
  }, [token, authLoading, slotId, router])

  if (loading) return <div className="max-w-6xl mx-auto space-y-6"><div className="h-8 w-48 bg-secondary rounded animate-pulse" /></div>
  if (error) return <p className="text-destructive">{error}</p>

  if (!data?.available) {
    return (
      <div className="max-w-6xl mx-auto space-y-6">
        <CompanyNav slotId={slotId} />
        <div className="rounded-2xl border border-dashed border-border p-12 text-center space-y-3">
          <p className="text-4xl">📡</p>
          <p className="font-semibold">No live data available</p>
          <p className="text-sm text-muted-foreground max-w-sm mx-auto">
            Install the <strong>FS25_FarmSimManager</strong> mod in Farming Simulator 25 and load a save game. Live data exports every 30 seconds while the game is running.
          </p>
        </div>
      </div>
    )
  }

  const gt = data.gameTime!
  const seasons: Record<string, string> = { SPRING: "🌸 Spring", SUMMER: "☀️ Summer", AUTUMN: "🍂 Autumn", WINTER: "❄️ Winter" }

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <CompanyNav slotId={slotId} />

      <div className="flex items-center gap-3">
        <h2 className="text-xl font-semibold">Live Data</h2>
        <Badge className="bg-primary/15 text-primary text-xs">
          📡 Updated {data.exportedAt ? new Date(data.exportedAt).toLocaleTimeString() : "—"}
        </Badge>
      </div>

      {/* Game Time */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        {[
          { icon: "📅", label: "Season", value: seasons[gt.season] ?? gt.season },
          { icon: "🗓️", label: "In-Game Date", value: `Y${gt.year} D${gt.day}` },
          { icon: "🕐", label: "Time", value: `${String(gt.hour).padStart(2,"0")}:${String(gt.minute).padStart(2,"0")}` },
          { icon: "🌍", label: "Month", value: `Month ${gt.month}` },
        ].map((s) => (
          <div key={s.label} className="rounded-2xl bg-card border border-border/60 card-shadow p-4">
            <div className="flex items-center gap-2 mb-1">
              <span>{s.icon}</span>
              <span className="text-xs text-muted-foreground">{s.label}</span>
            </div>
            <p className="font-semibold">{s.value}</p>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Crop Prices */}
        <Card>
          <CardHeader><CardTitle className="text-base">💹 Crop Prices</CardTitle></CardHeader>
          <CardContent>
            {!data.cropPrices?.length ? (
              <p className="text-sm text-muted-foreground">No price data</p>
            ) : (
              <div className="space-y-1 max-h-72 overflow-y-auto">
                {[...(data.cropPrices ?? [])]
                  .filter((c) => !ANIMAL_PREFIXES.some(p => c.name.startsWith(p)) &&
                    !c.name.startsWith("BALE_") &&
                    !c.name.startsWith("ROUNDBALE") && !c.name.startsWith("SQUAREBALE") &&
                    !["MANURE","LIQUIDMANURE","DIGESTATE","WATER","DIESEL","DEF","ELECTRICCHARGE",
                      "METHANE","TREESAPLINGS","TREE","POPLAR","FORAGE","FORAGE_MIXING",
                      "CHAFF","STONE","OILSEEDRADISH","RICESAPLINGS"].includes(c.name))
                  .sort((a, b) => b.pricePerLiter - a.pricePerLiter)
                  .map((c) => {
                  const pct = c.basePrice > 0 ? (c.pricePerLiter / c.basePrice) * 100 : 100
                  const color = pct >= 115 ? "text-primary" : pct <= 85 ? "text-destructive" : "text-foreground"
                  return (
                    <div key={c.name} className="flex items-center justify-between py-1.5 border-b border-border/40 last:border-0 text-sm">
                      <span className="capitalize">{c.title || c.name.toLowerCase()}</span>
                      <div className="flex items-center gap-3">
                        <span className="text-xs text-muted-foreground">{fmt(c.basePrice)}/l base</span>
                        <span className={`font-semibold tabular-nums ${color}`}>{fmt(c.pricePerLiter)}/l</span>
                      </div>
                    </div>
                  )
                })}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Animals — full detail on Animals tab */}
        <Card>
          <CardHeader><CardTitle className="text-base">🐄 Animals</CardTitle></CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">See the <strong>Animals</strong> tab for full animal details, health, and productivity.</p>
          </CardContent>
        </Card>

        {/* Animal Market Prices */}
        {data.animalPrices && data.animalPrices.length > 0 && (
          <Card className="lg:col-span-2">
            <CardHeader><CardTitle className="text-base">🐾 Animal Dealer Prices</CardTitle></CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                {data.animalPrices.map((a) => {
                  const prefix = ANIMAL_PREFIXES.find(p => a.name.startsWith(p)) ?? ""
                  const icon = ANIMAL_ICONS[prefix] ?? "🐾"
                  return (
                    <div key={a.name} className="rounded-xl bg-secondary/40 p-3 space-y-2">
                      <p className="font-medium text-sm flex items-center gap-1.5">
                        <span>{icon}</span>
                        <span>{a.title}</span>
                      </p>
                      <div className="space-y-1">
                        {a.tiers.map((t) => (
                          <div key={t.label} className="flex justify-between text-xs">
                            <span className="text-muted-foreground">{t.label} ({t.months}m)</span>
                            <span className="font-semibold tabular-nums text-primary">{fmt(t.price)}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )
                })}
              </div>
            </CardContent>
          </Card>
        )}

        {/* Contracts */}
        <Card>
          <CardHeader><CardTitle className="text-base">📋 Contracts</CardTitle></CardHeader>
          <CardContent>
            {!data.contracts?.length ? (
              <p className="text-sm text-muted-foreground">No active contracts</p>
            ) : (
              <div className="space-y-2">
                {(data.contracts ?? []).filter(c => c.isActive).map((c, i) => (
                  <div key={i} className="rounded-xl bg-secondary/40 p-3 space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="font-medium capitalize">{c.type.replace(/_/g," ").toLowerCase()}</span>
                      <span className="text-primary font-semibold">{fmt(c.reward)}</span>
                    </div>
                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                      <span>Field {c.fieldId}</span>
                      <span>•</span>
                      <span className="flex-1">
                        <Bar value={c.completion} color="bg-primary" />
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Workers */}
        <Card>
          <CardHeader><CardTitle className="text-base">👷 Hired Workers</CardTitle></CardHeader>
          <CardContent>
            {!data.workers?.length ? (
              <p className="text-sm text-muted-foreground">No hired workers</p>
            ) : (
              <div className="space-y-2">
                {(data.workers ?? []).map((w, i) => (
                  <div key={i} className="flex items-center justify-between py-2 border-b border-border/40 last:border-0 text-sm">
                    <div>
                      <p className="font-medium capitalize">{w.task.replace(/_/g," ").toLowerCase()}</p>
                      {w.vehicle && <p className="text-xs text-muted-foreground">{w.vehicle}</p>}
                    </div>
                    <span className="text-xs text-muted-foreground">{fmt(w.wagePerHour)}/h</span>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
