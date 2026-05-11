// Cloud API client — talks to the Railway backend

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080"

function authHeaders(token: string) {
  return { Authorization: `Bearer ${token}`, "Content-Type": "application/json" }
}

export async function register(email: string, password: string) {
  const res = await fetch(`${API_URL}/auth/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json() as Promise<{ token: string; companionToken: string }>
}

export async function login(email: string, password: string) {
  const res = await fetch(`${API_URL}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json() as Promise<{ token: string; companionToken: string }>
}

export async function getMe(token: string) {
  const res = await fetch(`${API_URL}/auth/me`, { headers: authHeaders(token) })
  if (!res.ok) throw new Error(await res.text())
  return res.json() as Promise<{ id: string; email: string; companionToken: string; companionLastSeen: string | null }>
}

export async function getCompanies(token: string) {
  const res = await fetch(`${API_URL}/api/companies`, { headers: authHeaders(token) })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function getFinances(token: string, companyId: string) {
  const res = await fetch(`${API_URL}/api/companies/${companyId}/finances`, { headers: authHeaders(token) })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function getFields(token: string, companyId: string) {
  const res = await fetch(`${API_URL}/api/companies/${companyId}/fields`, { headers: authHeaders(token) })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function getVehicles(token: string, companyId: string) {
  const res = await fetch(`${API_URL}/api/companies/${companyId}/vehicles`, { headers: authHeaders(token) })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}
