
---comment
---@param n number
---@return number
function Fib(n)
    if n < 2 then return n end
    return Fib(n-2)+Fib(n-1)
end
