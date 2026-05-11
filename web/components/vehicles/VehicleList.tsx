"use client"

import { useState } from "react"
import Image from "next/image"
import type { Vehicle } from "@/types/vehicles"

const COMPANION_URL = "http://127.0.0.1:3847"

function fmt(n: number) {
  if (n >= 1_000_000) return `$${(n / 1_000_000).toFixed(2)}M`
  if (n >= 1_000) return `$${(n / 1_000).toFixed(0)}K`
  return `$${n.toFixed(0)}`
}

const CATEGORY_ICONS: Record<string, string> = {
  Tractor: "🚜", Harvester: "🌾", Trailer: "🚛", Tillage: "⚙️",
  Seeder: "🌱", Sprayer: "💧", Loader: "🏗️", Forage: "🎁",
  Truck: "🚚", Train: "🚂", Car: "🚗", Forestry: "🪵",
  Storage: "🏚️", Equipment: "🔧",
}

const CONDITION: Record<string, { bar: string; badge: string; label: string }> = {
  good: { bar: "bg-primary", badge: "bg-primary/15 text-primary", label: "Good" },
  fair: { bar: "bg-amber-500", badge: "bg-amber-100 text-amber-800", label: "Fair" },
  poor: { bar: "bg-destructive", badge: "bg-destructive/15 text-destructive", label: "Poor" },
}

function VehicleImage({ path, token, name }: { path: string; token: string; name: string }) {
  const [failed, setFailed] = useState(false)
  const src = `${COMPANION_URL}/api/image?path=${encodeURIComponent(path)}`

  if (failed || !path) {
    return (
      <div className="w-full h-full flex items-center justify-center text-4xl bg-secondary/50 rounded-xl">
        🚜
      </div>
    )
  }

  return (
    <img
      src={`${src}&token=${token}`}
      alt={name}
      className="w-full h-full object-contain p-2 rounded-xl"
      onError={() => setFailed(true)}
    />
  )
}

export default function VehicleList({ vehicles, token }: { vehicles: Vehicle[]; token: string }) {
  if (vehicles.length === 0) {
    return (
      <div className="rounded-2xl border border-dashed border-border p-12 text-center text-muted-foreground">
        <p className="text-3xl mb-3">🚜</p>
        <p className="font-medium">No equipment in this category</p>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
      {vehicles.map((v) => {
        const cond = CONDITION[v.condition]
        const icon = CATEGORY_ICONS[v.category] ?? "🔧"
        return (
          <div
            key={v.uniqueId}
            className="group rounded-2xl bg-card border border-border/60 card-shadow overflow-hidden hover:card-shadow-hover hover:-translate-y-0.5 transition-all duration-200"
          >
            {/* Image area */}
            <div className="relative h-36 bg-secondary/30 border-b border-border/40">
              <VehicleImage path={v.storeImagePath} token={token} name={v.name} />
              {/* Condition badge top-right */}
              <span className={`absolute top-2 right-2 text-[10px] font-semibold px-2 py-0.5 rounded-full ${cond.badge}`}>
                {cond.label}
              </span>
              {/* Mod badge top-left */}
              {v.isMod && (
                <span className="absolute top-2 left-2 text-[10px] font-semibold px-2 py-0.5 rounded-full bg-accent/30 text-accent-foreground">
                  MOD
                </span>
              )}
            </div>

            {/* Details */}
            <div className="p-3 space-y-2">
              <div>
                <p className="font-semibold text-xs leading-tight line-clamp-2">{v.name}</p>
                <p className="text-[10px] text-muted-foreground mt-0.5">
                  {icon} {v.category}
                </p>
              </div>

              {/* Damage bar */}
              <div className="space-y-1">
                <div className="flex justify-between text-[10px] text-muted-foreground">
                  <span>Damage</span>
                  <span className="tabular-nums">{(v.damage * 100).toFixed(0)}%</span>
                </div>
                <div className="h-1 rounded-full bg-secondary overflow-hidden">
                  <div
                    className={`h-full rounded-full ${cond.bar} transition-all`}
                    style={{ width: `${v.damage * 100}%` }}
                  />
                </div>
              </div>

              {/* Stats row */}
              <div className="flex justify-between text-[10px] pt-1 border-t border-border/50">
                <span className="text-muted-foreground tabular-nums">{v.operatingTimeHours.toFixed(0)}h</span>
                <span className="font-semibold tabular-nums">{fmt(v.price)}</span>
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}
