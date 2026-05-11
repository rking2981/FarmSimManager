import type { Field } from "@/types/fields"

interface Props { fields: Field[] }

export default function FieldSummary({ fields }: Props) {
  const ready = fields.filter((f) => f.status === "ready").length
  const growing = fields.filter((f) => f.status === "growing").length
  const empty = fields.filter((f) => f.status === "empty").length
  const attention = fields.filter((f) => f.status === "needs_attention").length
  const needsLime = fields.filter((f) => f.needsLime).length
  const needsPlow = fields.filter((f) => f.needsPlow).length
  const needsFert = fields.filter((f) => f.needsFertilizer).length
  const hasWeeds = fields.filter((f) => f.hasWeeds).length
  const avgHealth = fields.length
    ? Math.round(fields.reduce((s, f) => s + f.soilHealth, 0) / fields.length)
    : 0

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
      <StatCard icon="✅" label="Ready to Harvest" value={ready} accent={ready > 0 ? "primary" : "neutral"} />
      <StatCard icon="🌱" label="Growing" value={growing} accent="neutral" />
      <StatCard icon="🪨" label="Empty / Fallow" value={empty} accent="neutral" />
      <StatCard icon="⚠️" label="Needs Attention" value={attention} accent={attention > 0 ? "destructive" : "neutral"} />
      <StatCard icon="🧪" label="Needs Lime" value={needsLime} accent={needsLime > 0 ? "warn" : "neutral"} />
      <StatCard icon="🚜" label="Needs Plow" value={needsPlow} accent={needsPlow > 0 ? "warn" : "neutral"} />
      <StatCard icon="🌿" label="Needs Fertilizer" value={needsFert} accent={needsFert > 0 ? "warn" : "neutral"} />
      <StatCard
        icon="🌍"
        label="Avg Soil Health"
        value={`${avgHealth}%`}
        accent={avgHealth >= 70 ? "primary" : avgHealth >= 40 ? "warn" : "destructive"}
      />
    </div>
  )
}

const ACCENT_CLASSES = {
  primary: "text-primary",
  destructive: "text-destructive",
  warn: "text-amber-700",
  neutral: "text-foreground",
}

function StatCard({ icon, label, value, accent }: {
  icon: string
  label: string
  value: string | number
  accent: keyof typeof ACCENT_CLASSES
}) {
  return (
    <div className="rounded-2xl bg-card border border-border/60 card-shadow p-4 flex flex-col gap-2">
      <div className="flex items-center gap-2">
        <span className="text-base">{icon}</span>
        <span className="text-xs text-muted-foreground">{label}</span>
      </div>
      <p className={`text-xl font-semibold leading-none ${ACCENT_CLASSES[accent]}`}>{value}</p>
    </div>
  )
}
