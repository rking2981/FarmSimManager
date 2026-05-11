export interface Vehicle {
  uniqueId: string
  name: string
  category: string
  filename: string
  storeImagePath: string
  isMod: boolean
  ageMonths: number
  price: number
  operatingTimeHours: number
  damage: number
  condition: "good" | "fair" | "poor"
}
