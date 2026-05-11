"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"

interface Props {
  slotId: string
  farmName?: string
}

const tabs = [
  { label: "Finances", icon: "💰", href: (id: string) => `/company/${id}` },
  { label: "Fields", icon: "🌾", href: (id: string) => `/company/${id}/fields` },
  { label: "Equipment", icon: "🚜", href: (id: string) => `/company/${id}/vehicles` },
]

export default function CompanyNav({ slotId, farmName }: Props) {
  const pathname = usePathname()

  return (
    <div className="space-y-4">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-sm">
        <Link href="/" className="text-muted-foreground hover:text-foreground transition-colors">
          All Companies
        </Link>
        <span className="text-border">/</span>
        <span className="font-medium">{farmName ?? `Slot ${slotId}`}</span>
      </div>

      {/* Tab bar */}
      <div className="flex gap-1 p-1 bg-secondary/60 rounded-xl w-fit">
        {tabs.map((tab) => {
          const href = tab.href(slotId)
          const active = pathname === href
          return (
            <Link
              key={tab.label}
              href={href}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                active
                  ? "bg-card text-foreground card-shadow"
                  : "text-muted-foreground hover:text-foreground"
              }`}
            >
              <span>{tab.icon}</span>
              {tab.label}
            </Link>
          )
        })}
      </div>
    </div>
  )
}
