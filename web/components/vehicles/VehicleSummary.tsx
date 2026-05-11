import type { Vehicle } from "@/types/vehicles"

function fmt(n: number) {
  if (n >= 1_000_000) return `$${(n / 1_000_000).toFixed(2)}M`
  if (n >= 1_000) return `$${(n / 1_000).toFixed(0)}K`
  return `$${n.toFixed(0)}`
}

export default function VehicleSummary({ vehicles }: { vehicles: Vehicle[] }) {
  const totalValue = vehicles.reduce((s, v) => s + v.price, 0)
  const totalHours = vehicles.reduce((s, v) => s + v.operatingTimeHours, 0)
  const poor = vehicles.filter((v) => v.condition === "poor").length
  const fair = vehicles.filter((v) => v.condition === "fair").length
  const mods = vehicles.filter((v) => v.isMod).length

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
      <StatCard icon="💵" label="Fleet Value" value={fmt(totalValue)} accent="primary" />
      <StatCard icon="⏱️" label="Total Hours" value={`${totalHours.toFixed(0)}h`} accent="neutral" />
      <StatCard icon="🔧" label="Needs Repair" value={poor} accent={poor > 0 ? "destructive" : "neutral"} />
      <StatCard icon="🧩" label="Mod Vehicles" value={mods} accent="neutral" />
    </div>
  )
}

const ACCENTS = {
  primary: "text-primary",
  destructive: "text-destructive",
  warn: "text-amber-700",
  neutral: "text-foreground",
}

function StatCard({ icon, label, value, accent }: {
  icon: string; label: string; value: string | number; accent: keyof typeof ACCENTS
}) {
  return (
    <div className="rounded-2xl bg-card border border-border/60 card-shadow p-4 flex flex-col gap-2">
      <div className="flex items-center gap-2">
        <span className="text-base">{icon}</span>
        <span className="text-xs text-muted-foreground">{label}</span>
      </div>
      <p className={`text-xl font-semibold leading-none ${ACCENTS[accent]}`}>{value}</p>
    </div>
  )
}
