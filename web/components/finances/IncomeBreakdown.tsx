"use client"

import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer, Legend } from "recharts"
import type { DailyFinances } from "@/types/finances"

const COLORS = ["#5A7A3A", "#D4A017", "#7EB8D4", "#C8703A", "#8B6914", "#4A9A6A", "#B89432", "#6EB0C4"]

function sum(days: DailyFinances[], key: keyof DailyFinances) {
  return days.reduce((s, d) => s + (d[key] as number), 0)
}

function fmtMoney(n: number) {
  if (n >= 1_000_000) return `$${(n / 1_000_000).toFixed(2)}M`
  if (n >= 1_000) return `$${(n / 1_000).toFixed(1)}K`
  return `$${n.toFixed(0)}`
}

export default function IncomeBreakdown({ days }: { days: DailyFinances[] }) {
  const raw = [
    { name: "Harvest", value: sum(days, "harvestIncome") },
    { name: "Missions", value: sum(days, "missionIncome") },
    { name: "Wood", value: sum(days, "soldWood") },
    { name: "Bales", value: sum(days, "soldBales") },
    { name: "Wool", value: sum(days, "soldWool") },
    { name: "Milk", value: sum(days, "soldMilk") },
    { name: "Products", value: sum(days, "soldProducts") },
    { name: "Animals", value: sum(days, "soldAnimals") },
    { name: "Buildings", value: sum(days, "soldBuildings") },
    { name: "Vehicles", value: sum(days, "soldVehicles") },
    { name: "Property", value: sum(days, "propertyIncome") },
    { name: "Fields", value: sum(days, "fieldSelling") },
    { name: "Other", value: days.reduce((s, d) => s + (d.other > 0 ? d.other : 0), 0) },
  ].filter((e) => e.value > 0)

  if (raw.length === 0) return <p className="text-sm text-muted-foreground text-center py-8">No income recorded yet.</p>

  return (
    <div className="space-y-3">
      <ResponsiveContainer width="100%" height={220}>
        <PieChart>
          <Pie data={raw} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={80} label={false}>
            {raw.map((_, i) => (
              <Cell key={i} fill={COLORS[i % COLORS.length]} />
            ))}
          </Pie>
          <Tooltip
            formatter={(v) => fmtMoney(Number(v))}
            contentStyle={{ backgroundColor: "#FDF6E3", border: "1px solid #E8D5B0", borderRadius: 8 }}
          />
          <Legend />
        </PieChart>
      </ResponsiveContainer>
      <table className="w-full text-sm">
        <tbody>
          {raw.map((e, i) => (
            <tr key={e.name} className="border-b border-border last:border-0">
              <td className="py-1 flex items-center gap-2">
                <span className="w-3 h-3 rounded-full inline-block shrink-0" style={{ backgroundColor: COLORS[i % COLORS.length] }} />
                {e.name}
              </td>
              <td className="py-1 text-right text-primary font-medium">{fmtMoney(e.value)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
