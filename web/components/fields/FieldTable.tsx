import type { Field } from "@/types/fields"

function fruitLabel(fruit: string) {
  if (!fruit || fruit === "UNKNOWN") return "—"
  return fruit.charAt(0) + fruit.slice(1).toLowerCase().replace(/_/g, " ")
}

function groundLabel(g: string) {
  return g.replace(/_/g, " ").charAt(0) + g.replace(/_/g, " ").slice(1).toLowerCase()
}

function SoilBar({ value }: { value: number }) {
  const color = value >= 70 ? "bg-primary" : value >= 40 ? "bg-amber-500" : "bg-destructive"
  return (
    <div className="flex items-center gap-2 min-w-[80px]">
      <div className="flex-1 h-1.5 rounded-full bg-secondary overflow-hidden">
        <div className={`h-full rounded-full ${color} transition-all`} style={{ width: `${value}%` }} />
      </div>
      <span className="text-[11px] text-muted-foreground w-7 text-right tabular-nums">{value}%</span>
    </div>
  )
}

const STATUS_PILL: Record<string, string> = {
  ready: "bg-primary/15 text-primary",
  growing: "bg-amber-100 text-amber-800",
  empty: "bg-secondary text-muted-foreground",
  needs_attention: "bg-destructive/15 text-destructive",
}

interface Props { fields: Field[] }

export default function FieldTable({ fields }: Props) {
  const sorted = [...fields].sort((a, b) => {
    const order = { ready: 0, needs_attention: 1, growing: 2, empty: 3 }
    return (order[a.status] ?? 9) - (order[b.status] ?? 9)
  })

  return (
    <div className="overflow-x-auto -mx-1">
      <table className="w-full text-sm">
        <thead>
          <tr className="text-left text-xs text-muted-foreground border-b border-border/60">
            <th className="pb-3 pl-2 font-medium w-14">Field</th>
            <th className="pb-3 font-medium">Crop</th>
            <th className="pb-3 font-medium hidden sm:table-cell">Planned</th>
            <th className="pb-3 font-medium">Status</th>
            <th className="pb-3 font-medium hidden md:table-cell">Ground</th>
            <th className="pb-3 font-medium hidden md:table-cell">Growth</th>
            <th className="pb-3 font-medium w-36">Soil</th>
            <th className="pb-3 font-medium">Issues</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border/40">
          {sorted.map((f, i) => (
            <tr
              key={f.id}
              className={`transition-colors hover:bg-secondary/40 ${i % 2 === 0 ? "" : "bg-secondary/20"}`}
            >
              <td className="py-3 pl-2 font-semibold text-muted-foreground tabular-nums">#{f.id}</td>
              <td className="py-3 font-medium">{fruitLabel(f.fruitType)}</td>
              <td className="py-3 text-muted-foreground hidden sm:table-cell">{fruitLabel(f.plannedFruit)}</td>
              <td className="py-3">
                <span className={`px-2.5 py-1 rounded-full text-[11px] font-medium ${STATUS_PILL[f.status]}`}>
                  {f.status.replace("_", " ")}
                </span>
              </td>
              <td className="py-3 text-xs text-muted-foreground hidden md:table-cell">{groundLabel(f.groundType)}</td>
              <td className="py-3 text-xs tabular-nums hidden md:table-cell">
                {f.fruitType !== "UNKNOWN" && f.fruitType ? `${f.growthState} / ${f.lastGrowthState}` : "—"}
              </td>
              <td className="py-3">
                <SoilBar value={f.soilHealth} />
              </td>
              <td className="py-3">
                <div className="flex gap-1 flex-wrap">
                  {f.needsLime && <Pill label="Lime" />}
                  {f.needsPlow && <Pill label="Plow" />}
                  {f.needsFertilizer && <Pill label="Fert" />}
                  {f.hasWeeds && <Pill label="Weeds" />}
                  {f.hasStones && <Pill label="Stones" variant="neutral" />}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function Pill({ label, variant = "danger" }: { label: string; variant?: "danger" | "neutral" }) {
  return (
    <span className={`text-[10px] px-1.5 py-0.5 rounded-md font-medium ${
      variant === "danger"
        ? "bg-destructive/15 text-destructive"
        : "bg-secondary text-muted-foreground"
    }`}>
      {label}
    </span>
  )
}
