import type { Company } from "@/types/company"
import type { FinancesResponse } from "@/types/finances"
import type { Field } from "@/types/fields"
import type { Vehicle } from "@/types/vehicles"

const COMPANION_URL = "http://127.0.0.1:3847"

export async function pingCompanion(): Promise<boolean> {
  try {
    const res = await fetch(`${COMPANION_URL}/api/health`, { signal: AbortSignal.timeout(2000) })
    return res.ok
  } catch {
    return false
  }
}

export async function fetchCompanies(token: string): Promise<Company[]> {
  const res = await fetch(`${COMPANION_URL}/api/companies`, {
    headers: { Authorization: `Bearer ${token}` },
    signal: AbortSignal.timeout(5000),
  })
  if (!res.ok) throw new Error(`companion returned ${res.status}`)
  return res.json()
}

export async function fetchFinances(token: string, slotId: string): Promise<FinancesResponse> {
  const res = await fetch(`${COMPANION_URL}/api/companies/${slotId}/finances`, {
    headers: { Authorization: `Bearer ${token}` },
    signal: AbortSignal.timeout(5000),
  })
  if (!res.ok) throw new Error(`companion returned ${res.status}`)
  return res.json()
}

export async function fetchFields(token: string, slotId: string): Promise<Field[]> {
  const res = await fetch(`${COMPANION_URL}/api/companies/${slotId}/fields`, {
    headers: { Authorization: `Bearer ${token}` },
    signal: AbortSignal.timeout(5000),
  })
  if (!res.ok) throw new Error(`companion returned ${res.status}`)
  return res.json()
}

export async function fetchVehicles(token: string, slotId: string): Promise<Vehicle[]> {
  const res = await fetch(`${COMPANION_URL}/api/companies/${slotId}/vehicles`, {
    headers: { Authorization: `Bearer ${token}` },
    signal: AbortSignal.timeout(5000),
  })
  if (!res.ok) throw new Error(`companion returned ${res.status}`)
  return res.json()
}
