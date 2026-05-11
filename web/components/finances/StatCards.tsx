import type { FarmStats } from "@/types/finances"

function fmt(n: number) {
  return new Intl.NumberFormat("en-US", { style: "currency", currency: "USD", maximumFractionDigits: 0 }).format(n)
}

interface Props {
  totalIncome: number
  totalExpense: number
  netProfit: number
  stats: FarmStats
}

export default function StatCards({ totalIncome, totalExpense, netProfit, stats }: Props) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
      <StatCard icon="📈" label="Total Income" value={fmt(totalIncome)} accent="primary" />
      <StatCard icon="📉" label="Total Expenses" value={fmt(totalExpense)} accent="destructive" />
      <StatCard
        icon="💰"
        label="Net Profit"
        value={fmt(netProfit)}
        accent={netProfit >= 0 ? "primary" : "destructive"}
        large
      />
      <StatCard icon="⏱️" label="Missions" value={String(stats.missionCount)} accent="neutral" />
      <StatCard icon="🚜" label="Worked" value={`${stats.workedHectares.toFixed(1)} ha`} accent="neutral" />
      <StatCard icon="🌱" label="Sown" value={`${stats.sownHectares.toFixed(1)} ha`} accent="neutral" />
      <StatCard icon="🌾" label="Harvested" value={`${stats.threshedHectares.toFixed(1)} ha`} accent="neutral" />
      <StatCard icon="🎁" label="Bales Made" value={String(stats.baleCount)} accent="neutral" />
    </div>
  )
}

const ACCENT_CLASSES = {
  primary: { dot: "bg-primary", value: "text-primary" },
  destructive: { dot: "bg-destructive", value: "text-destructive" },
  neutral: { dot: "bg-muted-foreground/40", value: "text-foreground" },
}

function StatCard({
  icon, label, value, accent, large,
}: {
  icon: string
  label: string
  value: string
  accent: keyof typeof ACCENT_CLASSES
  large?: boolean
}) {
  const { value: valueClass } = ACCENT_CLASSES[accent]
  return (
    <div className="rounded-2xl bg-card border border-border/60 card-shadow p-4 flex flex-col gap-2">
      <div className="flex items-center gap-2">
        <span className="text-base">{icon}</span>
        <span className="text-xs text-muted-foreground">{label}</span>
      </div>
      <p className={`${large ? "text-2xl" : "text-lg"} font-semibold ${valueClass} leading-none`}>
        {value}
      </p>
    </div>
  )
}
