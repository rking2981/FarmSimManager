export interface DailyFinances {
  day: number
  harvestIncome: number
  missionIncome: number
  soldWood: number
  soldBales: number
  soldWool: number
  soldMilk: number
  soldProducts: number
  soldAnimals: number
  soldBuildings: number
  soldVehicles: number
  propertyIncome: number
  fieldSelling: number
  other: number
  newVehiclesCost: number
  constructionCost: number
  fieldPurchase: number
  purchaseFuel: number
  purchaseSeeds: number
  purchaseFertilizer: number
  purchasePallets: number
  purchaseWater: number
  vehicleLeasingCost: number
  vehicleRunningCost: number
  propertyMaintenance: number
  wagePayment: number
  loanInterest: number
  newAnimalsCost: number
  productionCosts: number
  totalIncome: number
  totalExpense: number
  netProfit: number
}

export interface FarmStats {
  workedHectares: number
  sownHectares: number
  sprayedHectares: number
  threshedHectares: number
  revenue: number
  expenses: number
  missionCount: number
  baleCount: number
}

export interface FinancesResponse {
  days: DailyFinances[]
  stats: FarmStats
}
