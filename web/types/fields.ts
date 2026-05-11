export interface Field {
  id: number
  fruitType: string
  plannedFruit: string
  growthState: number
  lastGrowthState: number
  groundType: string
  sprayType: string
  sprayLevel: number
  limeLevel: number
  plowLevel: number
  weedState: number
  stoneLevel: number
  owned: boolean
  soilHealth: number
  status: "ready" | "growing" | "empty" | "needs_attention"
  needsLime: boolean
  needsPlow: boolean
  needsFertilizer: boolean
  hasWeeds: boolean
  hasStones: boolean
}
