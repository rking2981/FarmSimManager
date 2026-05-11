"use client"

import Link from "next/link"
import type { Company } from "@/types/company"

function fmt(n: number) {
  return new Intl.NumberFormat("en-US", { style: "currency", currency: "USD", maximumFractionDigits: 0 }).format(n)
}

const DIFFICULTY_COLORS: Record<string, string> = {
  EASY: "text-primary bg-primary/10",
  NORMAL: "text-accent-foreground bg-accent/20",
  HARD: "text-destructive bg-destructive/10",
}

export default function CompanyCard({ company }: { company: Company }) {
  const netWorth = company.money - company.loanAmount
  const diffColor = DIFFICULTY_COLORS[company.difficulty ?? "NORMAL"] ?? DIFFICULTY_COLORS.NORMAL

  return (
    <Link href={`/company/${company.id ?? company.slotId}`} className="block group">
      <div className="relative rounded-2xl bg-card card-shadow border border-border/60 overflow-hidden transition-all duration-200 group-hover:card-shadow-hover group-hover:-translate-y-0.5 group-hover:border-primary/30">
        {/* Top accent strip */}
        <div className="h-1.5 w-full bg-gradient-to-r from-primary/70 via-primary to-accent/70" />

        <div className="p-5">
          {/* Header */}
          <div className="flex items-start justify-between gap-3 mb-4">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center text-xl shrink-0">
                🏡
              </div>
              <div>
                <h3 className="font-semibold text-sm leading-tight">{company.farmName}</h3>
                <p className="text-xs text-muted-foreground mt-0.5">{company.mapTitle}</p>
              </div>
            </div>
            <div className="flex flex-col items-end gap-1.5 shrink-0">
              <span className="text-[10px] font-medium text-muted-foreground bg-secondary px-2 py-0.5 rounded-full">
                Slot {company.slotId}
              </span>
              {company.difficulty && (
                <span className={`text-[10px] font-medium px-2 py-0.5 rounded-full ${diffColor}`}>
                  {company.difficulty.charAt(0) + company.difficulty.slice(1).toLowerCase()}
                </span>
              )}
            </div>
          </div>

          {/* Financial stats */}
          <div className="grid grid-cols-3 gap-2 mb-4">
            <Stat label="Cash" value={fmt(company.money)} color="text-primary" />
            <Stat label="Loan" value={fmt(company.loanAmount)} color="text-destructive" />
            <Stat
              label="Net Worth"
              value={fmt(netWorth)}
              color={netWorth >= 0 ? "text-foreground font-semibold" : "text-destructive font-semibold"}
            />
          </div>

          {/* Footer */}
          <div className="flex items-center justify-between pt-3 border-t border-border/60 text-xs text-muted-foreground">
            <span>{(company.playTimeHours ?? 0).toFixed(1)}h played</span>
            <span className="flex items-center gap-1">
              <span>{company.lastSaved}</span>
              <span className="text-primary/50 group-hover:text-primary transition-colors ml-1">→</span>
            </span>
          </div>
        </div>
      </div>
    </Link>
  )
}

function Stat({ label, value, color }: { label: string; value: string; color?: string }) {
  return (
    <div className="bg-background rounded-xl px-3 py-2">
      <p className="text-[10px] text-muted-foreground mb-0.5">{label}</p>
      <p className={`text-xs font-medium truncate ${color ?? ""}`}>{value}</p>
    </div>
  )
}
