"use client"

import { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import { useAuth } from "@/lib/auth-context"
import { getFinances } from "@/lib/api"
import type { FinancesResponse } from "@/types/finances"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import FinancialChart from "@/components/finances/FinancialChart"
import IncomeBreakdown from "@/components/finances/IncomeBreakdown"
import ExpenseBreakdown from "@/components/finances/ExpenseBreakdown"
import StatCards from "@/components/finances/StatCards"
import CompanyNav from "@/components/CompanyNav"

export default function CompanyFinancesPage() {
  const { slotId } = useParams<{ slotId: string }>()
  const router = useRouter()
  const { token } = useAuth()
  const [data, setData] = useState<FinancesResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!token) { window.location.href = "/login"; return }
    getFinances(token, slotId)
      .then(setData)
      .catch((e) => setError(String(e)))
      .finally(() => setLoading(false))
  }, [slotId, token, router])

  if (loading) return <LoadingState />
  if (error) return <p className="text-destructive">{error}</p>
  if (!data) return null

  const totalIncome = data.days.reduce((s, d) => s + d.totalIncome, 0)
  const totalExpense = data.days.reduce((s, d) => s + d.totalExpense, 0)
  const netProfit = totalIncome - totalExpense
  const stats = data.stats ?? {
    workedHectares: 0, sownHectares: 0, sprayedHectares: 0,
    threshedHectares: 0, revenue: 0, expenses: 0, missionCount: 0, baleCount: 0,
  }

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <CompanyNav slotId={slotId} />
      <div className="flex items-center gap-3">
        <h2 className="text-xl font-semibold">Financial Dashboard</h2>
        <Badge variant="secondary">{data.days.length} days recorded</Badge>
      </div>
      <StatCards totalIncome={totalIncome} totalExpense={totalExpense} netProfit={netProfit} stats={stats} />
      <Card>
        <CardHeader><CardTitle className="text-base">Income vs Expenses — By Day</CardTitle></CardHeader>
        <CardContent><FinancialChart days={data.days} /></CardContent>
      </Card>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader><CardTitle className="text-base">Income Breakdown</CardTitle></CardHeader>
          <CardContent><IncomeBreakdown days={data.days} /></CardContent>
        </Card>
        <Card>
          <CardHeader><CardTitle className="text-base">Expense Breakdown</CardTitle></CardHeader>
          <CardContent><ExpenseBreakdown days={data.days} /></CardContent>
        </Card>
      </div>
    </div>
  )
}

function LoadingState() {
  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <div className="h-8 w-48 bg-secondary rounded animate-pulse" />
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="h-24 bg-secondary rounded-lg animate-pulse" />
        ))}
      </div>
      <div className="h-72 bg-secondary rounded-lg animate-pulse" />
    </div>
  )
}
