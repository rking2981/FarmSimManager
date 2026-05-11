"use client"

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  ReferenceLine,
} from "recharts"
import type { DailyFinances } from "@/types/finances"

interface Props {
  days: DailyFinances[]
}

function fmtMoney(n: number) {
  if (Math.abs(n) >= 1_000_000) return `$${(n / 1_000_000).toFixed(1)}M`
  if (Math.abs(n) >= 1_000) return `$${(n / 1_000).toFixed(0)}K`
  return `$${n.toFixed(0)}`
}

export default function FinancialChart({ days }: Props) {
  const data = days.map((d) => ({
    day: `Day ${d.day}`,
    Income: Math.round(d.totalIncome),
    Expenses: Math.round(d.totalExpense),
    "Net Profit": Math.round(d.netProfit),
  }))

  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={data} margin={{ top: 4, right: 16, left: 8, bottom: 4 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#E8D5B0" />
        <XAxis dataKey="day" tick={{ fontSize: 12 }} />
        <YAxis tickFormatter={fmtMoney} tick={{ fontSize: 12 }} width={64} />
        <Tooltip
          formatter={(value) => fmtMoney(Number(value))}
          contentStyle={{ backgroundColor: "#FDF6E3", border: "1px solid #E8D5B0", borderRadius: 8 }}
        />
        <Legend />
        <ReferenceLine y={0} stroke="#3B2A1A" strokeOpacity={0.3} />
        <Bar dataKey="Income" fill="#5A7A3A" radius={[3, 3, 0, 0]} />
        <Bar dataKey="Expenses" fill="#D4A017" radius={[3, 3, 0, 0]} />
        <Bar dataKey="Net Profit" fill="#7EB8D4" radius={[3, 3, 0, 0]} />
      </BarChart>
    </ResponsiveContainer>
  )
}
