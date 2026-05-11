"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"

export default function TokenSetup({ onSave }: { onSave: (token: string) => void }) {
  const [value, setValue] = useState("")

  return (
    <Card className="border-accent">
      <CardHeader className="pb-2">
        <CardTitle className="text-base">Enter Your Companion Token</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <p className="text-sm text-muted-foreground">
          Find your token in the companion app tray menu under &quot;Show Token&quot;. This is stored only in your browser.
        </p>
        <div className="flex gap-2">
          <input
            type="text"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            placeholder="Paste token here..."
            className="flex-1 rounded-md border border-input bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
          />
          <Button
            onClick={() => value.trim() && onSave(value.trim())}
            className="bg-primary text-primary-foreground hover:bg-primary/90"
          >
            Save
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
