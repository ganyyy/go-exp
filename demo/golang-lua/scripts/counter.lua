
---counter
---@param start number
---@param step number
---@return fun():number
local function counter(start, step)
    return function()
        start = start + step
        return start
    end
end

---@type fun(): number
NXT = counter(1, 10)


print("lua", user:String())

-- user:SetName("456")
print("lua", NXT())

-- print("lua", user.Name)