"use client"

import { createContext, useContext, useEffect, useState, ReactNode } from "react"

interface AuthState {
  token: string | null
  companionToken: string | null
  email: string | null
  loading: boolean
  setAuth: (token: string, companionToken: string, email: string) => void
  clearAuth: () => void
}

const AuthContext = createContext<AuthState>({
  token: null, companionToken: null, email: null, loading: true,
  setAuth: () => {}, clearAuth: () => {},
})

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null)
  const [companionToken, setCompanionToken] = useState<string | null>(null)
  const [email, setEmail] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const t = localStorage.getItem("api_token")
    const ct = localStorage.getItem("companion_cloud_token")
    const em = localStorage.getItem("user_email")
    if (t) { setToken(t); setCompanionToken(ct); setEmail(em) }
    setLoading(false)
  }, [])

  function setAuth(t: string, ct: string, em: string) {
    localStorage.setItem("api_token", t)
    localStorage.setItem("companion_cloud_token", ct)
    localStorage.setItem("user_email", em)
    setToken(t); setCompanionToken(ct); setEmail(em)
  }

  function clearAuth() {
    localStorage.removeItem("api_token")
    localStorage.removeItem("companion_cloud_token")
    localStorage.removeItem("user_email")
    setToken(null); setCompanionToken(null); setEmail(null)
  }

  return (
    <AuthContext.Provider value={{ token, companionToken, email, loading, setAuth, clearAuth }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}
