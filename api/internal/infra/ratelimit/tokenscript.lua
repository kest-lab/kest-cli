-- Token bucket rate limiting using Redis
-- KEYS[1]: tokens key
-- KEYS[2]: timestamp key
-- ARGV[1]: rate (tokens per second)
-- ARGV[2]: capacity (burst size)
-- ARGV[3]: current timestamp (unix seconds)
-- ARGV[4]: requested tokens

local rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

-- Calculate TTL as 2x the time to fill the bucket
local fill_time = capacity / rate
local ttl = math.floor(fill_time * 2)
if ttl < 1 then
    ttl = 1
end

-- Get current tokens, default to capacity (full bucket)
local last_tokens = tonumber(redis.call("get", KEYS[1]))
if last_tokens == nil then
    last_tokens = capacity
end

-- Get last refresh time
local last_refreshed = tonumber(redis.call("get", KEYS[2]))
if last_refreshed == nil then
    last_refreshed = 0
end

-- Calculate tokens to add based on time elapsed
local delta = math.max(0, now - last_refreshed)
local filled_tokens = math.min(capacity, last_tokens + (delta * rate))

-- Check if we have enough tokens
local allowed = filled_tokens >= requested
local new_tokens = filled_tokens

if allowed then
    new_tokens = filled_tokens - requested
end

-- Update Redis with new values
redis.call("setex", KEYS[1], ttl, new_tokens)
redis.call("setex", KEYS[2], ttl, now)

-- Return 1 if allowed, 0 if not
if allowed then
    return 1
else
    return 0
end
