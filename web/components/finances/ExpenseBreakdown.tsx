"use client"

import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer, Legend } from "recharts"
import type { DailyFinances } from "@/types/finances"

const COLORS = ["#C8703A", "#D4A017", "#A05020", "#B89432", "#8B4513", "#D4701A", "#C89020", "#A06030"]

function sum(days: DailyFinances[], key: keyof DailyFinances) {
  return days.reduce((s, d) => s + (d[key] as number), 0)
}

function fmtMoney(n: number) {
  if (n >= 1_000_000) return `$${(n / 1_000_000).toFixed(2)}M`
  if (n >= 1_000) return `$${(n / 1_000).toFixed(1)}K`
  return `$${n.toFixed(0)}`
}

export default function ExpenseBreakdown({ days }: { days: DailyFinances[] }) {
  const raw = [
    { name: "Vehicles", value: sum(days, "newVehiclesCost") },
    { name: "Construction", value: sum(days, "constructionCost") },
    { name: "Field Purchase", value: sum(days, "fieldPurchase") },
    { name: "Leasing", value: sum(days, "vehicleLeasingCost") },
    { name: "Running Costs", value: sum(days, "vehicleRunningCost") },
    { name: "Seeds", value: sum(days, "purchaseSeeds") },
    { name: "Fertilizer", value: sum(days, "purchaseFertilizer") },
    { name: "Fuel", value: sum(days, "purchaseFuel") },
    { name: "Pallets", value: sum(days, "purchasePallets") },
    { name: "Water", value: sum(days, "purchaseWater") },
    { name: "Wages", value: sum(days, "wagePayment") },
    { name: "Maintenance", value: sum(days, "propertyMaintenance") },
    { name: "Animals", value: sum(days, "newAnimalsCost") },
    { name: "Production", value: sum(days, "productionCosts") },
    { name: "Loan Interest", value: sum(days, "loanInterest") },
  ].filter((e) => e.value > 0)

  if (raw.length === 0) return <p className="text-sm text-muted-foreground text-center py-8">No expenses recorded yet.</p>

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
              <td className="py-1 text-right text-destructive font-medium">{fmtMoney(e.value)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
