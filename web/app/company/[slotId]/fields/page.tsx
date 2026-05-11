"use client"

import { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import { useAuth } from "@/lib/auth-context"
import { getFields } from "@/lib/api"
import type { Field } from "@/types/fields"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import CompanyNav from "@/components/CompanyNav"
import FieldGrid from "@/components/fields/FieldGrid"
import FieldSummary from "@/components/fields/FieldSummary"
import FieldTable from "@/components/fields/FieldTable"

export default function FieldsPage() {
  const { slotId } = useParams<{ slotId: string }>()
  const router = useRouter()
  const { token } = useAuth()
  const [fields, setFields] = useState<Field[]>([])
  const [filter, setFilter] = useState<"all" | "owned">("owned")
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!token) { window.location.href = "/login"; return }
    getFields(token, slotId)
      .then(setFields)
      .catch((e) => setError(String(e)))
      .finally(() => setLoading(false))
  }, [slotId, token, router])

  const displayed = filter === "owned" ? fields.filter((f) => f.owned) : fields
  const owned = fields.filter((f) => f.owned)

  if (loading) return <LoadingState />
  if (error) return <p className="text-destructive">{error}</p>

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <CompanyNav slotId={slotId} />
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <h2 className="text-xl font-semibold">Fields</h2>
          <Badge variant="secondary">{owned.length} owned</Badge>
        </div>
        <div className="flex gap-2">
          {(["owned", "all"] as const).map((f) => (
            <button
              key={f}
              onClick={() => setFilter(f)}
              className={`px-3 py-1 rounded-md text-sm transition-colors ${
                filter === f
                  ? "bg-primary text-primary-foreground"
                  : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
              }`}
            >
              {f === "owned" ? "My Fields" : "All Fields"}
            </button>
          ))}
        </div>
      </div>
      <FieldSummary fields={displayed} />
      <Card>
        <CardHeader><CardTitle className="text-base">Field Status Grid</CardTitle></CardHeader>
        <CardContent><FieldGrid fields={displayed} /></CardContent>
      </Card>
      <Card>
        <CardHeader><CardTitle className="text-base">Field Details</CardTitle></CardHeader>
        <CardContent><FieldTable fields={displayed} /></CardContent>
      </Card>
    </div>
  )
}

function LoadingState() {
  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <div className="h-8 w-48 bg-secondary rounded animate-pulse" />
      <div className="h-24 bg-secondary rounded-lg animate-pulse" />
      <div className="h-48 bg-secondary rounded-lg animate-pulse" />
    </div>
  )
}
