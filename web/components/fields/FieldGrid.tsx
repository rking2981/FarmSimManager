import type { Field } from "@/types/fields"

const STATUS_COLORS: Record<string, string> = {
  ready: "bg-primary text-primary-foreground",
  growing: "bg-accent text-accent-foreground",
  empty: "bg-secondary text-secondary-foreground",
  needs_attention: "bg-destructive/80 text-white",
}

const STATUS_LABELS: Record<string, string> = {
  ready: "Ready",
  growing: "Growing",
  empty: "Empty",
  needs_attention: "Attention",
}

function fruitLabel(fruit: string) {
  if (!fruit || fruit === "UNKNOWN") return "—"
  return fruit.charAt(0) + fruit.slice(1).toLowerCase().replace(/_/g, " ")
}

interface Props { fields: Field[] }

export default function FieldGrid({ fields }: Props) {
  if (fields.length === 0) {
    return <p className="text-sm text-muted-foreground text-center py-8">No fields to display.</p>
  }

  return (
    <div className="grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 lg:grid-cols-10 gap-2">
      {fields.map((f) => (
        <div
          key={f.id}
          className={`rounded-lg p-2 text-center text-xs font-medium ${STATUS_COLORS[f.status]} relative`}
          title={`Field ${f.id} — ${fruitLabel(f.fruitType)} — ${STATUS_LABELS[f.status]}`}
        >
          <div className="font-bold">{f.id}</div>
          <div className="truncate opacity-80 text-[10px]">{fruitLabel(f.fruitType)}</div>
          {(f.hasWeeds || f.needsLime) && (
            <span className="absolute top-0.5 right-0.5 w-1.5 h-1.5 rounded-full bg-white/70" />
          )}
        </div>
      ))}
    </div>
  )
}
