-- FarmSim Manager Bridge
-- Collects live game data and writes it to modSettings every 30 seconds.
-- The companion app watches this file and syncs it to the cloud.

FarmSimManagerBridge = {}
FarmSimManagerBridge.modName = g_currentModName
FarmSimManagerBridge.interval = 30000  -- ms between exports
FarmSimManagerBridge.timer = 0
FarmSimManagerBridge.outputDir = getUserProfileAppPath() .. "modSettings/FS25_FarmSimManager/"
FarmSimManagerBridge.outputFile = FarmSimManagerBridge.outputDir .. "data.json"

-- ─── Utility ────────────────────────────────────────────────────────────────

local function safeStr(v)
    if v == nil then return "" end
    return tostring(v):gsub('"', '\\"')
end

local function safeNum(v)
    if v == nil or type(v) ~= "number" then return 0 end
    if v ~= v then return 0 end  -- NaN check
    return math.floor(v * 100 + 0.5) / 100
end

local function safeBool(v)
    return v == true and "true" or "false"
end

-- Minimal JSON serialiser — avoids external dependencies
local function jsonStr(s)  return '"' .. safeStr(s) .. '"' end
local function jsonNum(n)  return tostring(safeNum(n)) end
local function jsonBool(b) return safeBool(b) end

local function jsonObj(fields)
    local parts = {}
    for _, kv in ipairs(fields) do
        parts[#parts + 1] = '"' .. kv[1] .. '":' .. kv[2]
    end
    return "{" .. table.concat(parts, ",") .. "}"
end

local function jsonArr(items)
    return "[" .. table.concat(items, ",") .. "]"
end

-- ─── Data collectors ─────────────────────────────────────────────────────────

local function collectGameTime()
    local env = g_currentMission and g_currentMission.environment
    if not env then return jsonObj({}) end

    local season = "UNKNOWN"
    if env.currentSeason ~= nil then
        local seasons = {"SPRING", "SUMMER", "AUTUMN", "WINTER"}
        season = seasons[env.currentSeason + 1] or "UNKNOWN"
    end

    return jsonObj({
        {"year",        jsonNum(env.currentYear or 0)},
        {"month",       jsonNum(env.currentMonth or 0)},
        {"day",         jsonNum(env.currentDay or 0)},
        {"hour",        jsonNum(env.currentHour or 0)},
        {"minute",      jsonNum(env.currentMinute or 0)},
        {"season",      jsonStr(season)},
        {"dayLength",   jsonNum(env.dayDuration or 0)},
    })
end

local function collectCropPrices()
    local items = {}
    if not g_currentMission or not g_currentMission.economyManager then
        return jsonArr(items)
    end

    local economy = g_currentMission.economyManager
    local fillTypes = g_fillTypeManager and g_fillTypeManager:getFillTypes()
    if not fillTypes then return jsonArr(items) end

    for _, fillType in ipairs(fillTypes) do
        if fillType.pricePerLiter and fillType.pricePerLiter > 0 then
            local price = economy:getPricePerLiter(fillType.index, nil) or fillType.pricePerLiter
            items[#items + 1] = jsonObj({
                {"name",          jsonStr(fillType.name or "")},
                {"title",         jsonStr(fillType.title or fillType.name or "")},
                {"pricePerLiter", jsonNum(price)},
                {"basePrice",     jsonNum(fillType.pricePerLiter or 0)},
            })
        end
    end

    return jsonArr(items)
end

local function collectContracts()
    local items = {}
    local cm = g_currentMission and g_currentMission.contractManager
    if not cm then return jsonArr(items) end

    local contracts = cm:getContracts(false)  -- false = all farms
    if not contracts then return jsonArr(items) end

    for _, contract in ipairs(contracts) do
        if contract then
            local typeName = "UNKNOWN"
            if contract.getContractTypeName then
                typeName = contract:getContractTypeName() or "UNKNOWN"
            end

            local fieldNum = 0
            if contract.field and contract.field.fieldId then
                fieldNum = contract.field.fieldId
            end

            local rewardStr = "0"
            if contract.rewardPerHa then
                rewardStr = tostring(safeNum(contract.rewardPerHa))
            elseif contract.reward then
                rewardStr = tostring(safeNum(contract.reward))
            end

            local completion = 0
            if contract.getCompletion then
                completion = safeNum((contract:getCompletion() or 0) * 100)
            end

            items[#items + 1] = jsonObj({
                {"type",       jsonStr(typeName)},
                {"fieldId",    jsonNum(fieldNum)},
                {"reward",     rewardStr},
                {"completion", jsonNum(completion)},
                {"isActive",   jsonBool(contract.isActive or false)},
            })
        end
    end

    return jsonArr(items)
end

local function collectAnimals()
    local items = {}
    local am = g_currentMission and g_currentMission.animalManager
    if not am then return jsonArr(items) end

    local clusters = nil
    if am.getClusters then clusters = am:getClusters() end
    if not clusters then return jsonArr(items) end

    -- Group by animal type
    local byType = {}
    for _, cluster in ipairs(clusters) do
        if cluster then
            local typeName = "UNKNOWN"
            if cluster.animalType and cluster.animalType.name then
                typeName = cluster.animalType.name
            elseif cluster.getAnimalTypeName then
                typeName = cluster:getAnimalTypeName() or "UNKNOWN"
            end

            if not byType[typeName] then
                byType[typeName] = {count = 0, health = 0, productivity = 0, n = 0}
            end

            local count = cluster.numAnimals or cluster:getNumAnimals() or 0
            local health = 0
            local productivity = 0

            if cluster.getHealthFactor then
                health = safeNum((cluster:getHealthFactor() or 0) * 100)
            end
            if cluster.getOutputFactor then
                productivity = safeNum((cluster:getOutputFactor() or 0) * 100)
            end

            byType[typeName].count = byType[typeName].count + count
            byType[typeName].health = byType[typeName].health + health
            byType[typeName].productivity = byType[typeName].productivity + productivity
            byType[typeName].n = byType[typeName].n + 1
        end
    end

    for typeName, data in pairs(byType) do
        local avgHealth = data.n > 0 and data.health / data.n or 0
        local avgProd   = data.n > 0 and data.productivity / data.n or 0
        items[#items + 1] = jsonObj({
            {"type",           jsonStr(typeName)},
            {"count",          jsonNum(data.count)},
            {"healthPct",      jsonNum(avgHealth)},
            {"productivityPct",jsonNum(avgProd)},
        })
    end

    return jsonArr(items)
end

local function collectWorkers()
    local items = {}
    local hm = g_currentMission and g_currentMission.helperManager
    if not hm then return jsonArr(items) end

    local helpers = nil
    if hm.getHelpers then helpers = hm:getHelpers() end
    if not helpers then return jsonArr(items) end

    for _, helper in ipairs(helpers) do
        if helper then
            local task = "IDLE"
            if helper.getTaskName then
                task = helper:getTaskName() or "IDLE"
            end

            local vehicleName = ""
            if helper.vehicle and helper.vehicle.getName then
                vehicleName = helper.vehicle:getName() or ""
            end

            local wage = 0
            if helper.getPrice then
                wage = safeNum(helper:getPrice() or 0)
            end

            items[#items + 1] = jsonObj({
                {"task",        jsonStr(task)},
                {"vehicle",     jsonStr(vehicleName)},
                {"wagePerHour", jsonNum(wage)},
            })
        end
    end

    return jsonArr(items)
end

local function collectFarmInfo()
    local farms = g_farmManager and g_farmManager:getFarms()
    if not farms then return jsonArr({}) end

    local items = {}
    for _, farm in ipairs(farms) do
        if farm and farm.farmId ~= 0 then  -- skip spectator farm
            items[#items + 1] = jsonObj({
                {"farmId",   jsonNum(farm.farmId or 0)},
                {"name",     jsonStr(farm.name or "")},
                {"money",    jsonNum(farm.money or 0)},
                {"loan",     jsonNum(farm.loan or 0)},
                {"color",    jsonNum(farm.color or 0)},
            })
        end
    end

    return jsonArr(items)
end

-- ─── Export ──────────────────────────────────────────────────────────────────

local function export()
    -- Ensure output directory exists
    createFolder(FarmSimManagerBridge.outputDir)

    local gameTime   = collectGameTime()
    local cropPrices = collectCropPrices()
    local contracts  = collectContracts()
    local animals    = collectAnimals()
    local workers    = collectWorkers()
    local farms      = collectFarmInfo()

    local json = jsonObj({
        {"exportedAt",  jsonStr(getDate("%Y-%m-%dT%H:%M:%S"))},
        {"gameTime",    gameTime},
        {"farms",       farms},
        {"cropPrices",  cropPrices},
        {"contracts",   contracts},
        {"animals",     animals},
        {"workers",     workers},
    })

    local file = io.open(FarmSimManagerBridge.outputFile, "w")
    if file then
        file:write(json)
        file:close()
    else
        print("FarmSimManager: failed to write " .. FarmSimManagerBridge.outputFile)
    end
end

-- ─── Lifecycle ───────────────────────────────────────────────────────────────

function FarmSimManagerBridge:update(dt)
    self.timer = self.timer + dt
    if self.timer >= self.interval then
        self.timer = 0
        export()
    end
end

function FarmSimManagerBridge:onMissionLoaded(mission)
    -- Export immediately on load, then every 30s
    self.timer = self.interval
end

function FarmSimManagerBridge:onMissionStart()
    export()
    self.timer = 0
    print("FarmSimManager: data bridge active, exporting to " .. self.outputFile)
end

-- Register with the game engine
local function init()
    Mission00.onMissionLoaded = Utils.appendedFunction(
        Mission00.onMissionLoaded,
        function(mission)
            FarmSimManagerBridge:onMissionLoaded(mission)
        end
    )

    Mission00.onStartMission = Utils.appendedFunction(
        Mission00.onStartMission,
        function()
            FarmSimManagerBridge:onMissionStart()
        end
    )

    addModEventListener(FarmSimManagerBridge)
end

init()
