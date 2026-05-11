"use client"

import { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import { useAuth } from "@/lib/auth-context"
import { getVehicles } from "@/lib/api"
import type { Vehicle } from "@/types/vehicles"
import CompanyNav from "@/components/CompanyNav"
import VehicleSummary from "@/components/vehicles/VehicleSummary"
import VehicleList from "@/components/vehicles/VehicleList"

const ALL_CATEGORIES = "All"

export default function VehiclesPage() {
  const { slotId } = useParams<{ slotId: string }>()
  const router = useRouter()
  const { token } = useAuth()
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [category, setCategory] = useState<string>(ALL_CATEGORIES)
  const [sortBy, setSortBy] = useState<"name" | "damage" | "price" | "hours">("damage")
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!token) { router.push("/login"); return }
    getVehicles(token, slotId)
      .then(setVehicles)
      .catch((e) => setError(String(e)))
      .finally(() => setLoading(false))
  }, [slotId, token, router])

  const categories = [ALL_CATEGORIES, ...Array.from(new Set(vehicles.map((v) => v.category))).sort()]
  const filtered = vehicles
    .filter((v) => category === ALL_CATEGORIES || v.category === category)
    .sort((a, b) => {
      switch (sortBy) {
        case "damage": return b.damage - a.damage
        case "price": return b.price - a.price
        case "hours": return b.operatingTimeHours - a.operatingTimeHours
        default: return a.name.localeCompare(b.name)
      }
    })

  if (loading) return <LoadingState />
  if (error) return <p className="text-destructive">{error}</p>

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <CompanyNav slotId={slotId} />
      <div className="flex items-center justify-between flex-wrap gap-3">
        <div className="flex items-center gap-3">
          <h2 className="text-xl font-semibold">Equipment</h2>
          <span className="text-sm text-muted-foreground">{vehicles.length} owned</span>
        </div>
        <select
          value={sortBy}
          onChange={(e) => setSortBy(e.target.value as typeof sortBy)}
          className="text-sm rounded-lg border border-border bg-card px-3 py-1.5 text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        >
          <option value="damage">Sort: Most Damaged</option>
          <option value="price">Sort: Highest Value</option>
          <option value="hours">Sort: Most Hours</option>
          <option value="name">Sort: Name</option>
        </select>
      </div>
      <VehicleSummary vehicles={vehicles} />
      <div className="flex gap-2 flex-wrap">
        {categories.map((cat) => (
          <button
            key={cat}
            onClick={() => setCategory(cat)}
            className={`px-3 py-1 rounded-full text-xs font-medium transition-colors ${
              category === cat
                ? "bg-primary text-primary-foreground"
                : "bg-secondary text-secondary-foreground hover:bg-secondary/70"
            }`}
          >
            {cat}
            {cat !== ALL_CATEGORIES && (
              <span className="ml-1.5 opacity-60">{vehicles.filter((v) => v.category === cat).length}</span>
            )}
          </button>
        ))}
      </div>
      <VehicleList vehicles={filtered} />
    </div>
  )
}

function LoadingState() {
  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <div className="h-8 w-48 bg-secondary rounded animate-pulse" />
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="h-24 bg-secondary rounded-2xl animate-pulse" />
        ))}
      </div>
    </div>
  )
}
