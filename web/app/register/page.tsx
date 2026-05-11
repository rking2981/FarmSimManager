"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { register } from "@/lib/api"
import { useAuth } from "@/lib/auth-context"
import AuthForm from "@/components/AuthForm"

export default function RegisterPage() {
  const router = useRouter()
  const { setAuth } = useAuth()
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(email: string, password: string) {
    setLoading(true)
    setError(null)
    try {
      const { token, companionToken } = await register(email, password)
      setAuth(token, companionToken, email)
      router.push("/")
    } catch (e) {
      setError(String(e).replace("Error: ", ""))
    } finally {
      setLoading(false)
    }
  }

  return <AuthForm mode="register" onSubmit={handleSubmit} error={error} loading={loading} />
}
