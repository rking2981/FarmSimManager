export interface Company {
  id?: string
  slotId: string
  farmName: string
  mapTitle: string
  money: number
  loanAmount: number
  playTimeHours: number
  lastSaved: string
  difficulty?: string
  creationDate?: string
  mods?: { name: string; title: string; version: string }[]
}
