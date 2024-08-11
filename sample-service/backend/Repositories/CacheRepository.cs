using System.Text;
using Microsoft.Extensions.Caching.Distributed;
using Newtonsoft.Json;
using backend.Repositories.Interfaces;

namespace backend.Repositories;

public class CacheRepository : ICacheRepository
{
    private readonly IDistributedCache _cache;

    public CacheRepository(IDistributedCache cache)
    {
        _cache = cache;
    }

    public async Task SetAsync<T>(string key, T value, TimeSpan expiration)
    {
        var encodedValue = Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(value));
        await _cache.SetAsync(key, encodedValue, new DistributedCacheEntryOptions
        {
            AbsoluteExpirationRelativeToNow = expiration
        });
    }

    public async Task<T?> GetAsync<T>(string key)
    {
        var encodedData = await _cache.GetAsync(key);
        if (encodedData == null)
        {
            return default;
        }

        var jsonData = Encoding.UTF8.GetString(encodedData);
        return JsonConvert.DeserializeObject<T>(jsonData);
    }

    public async Task RemoveAsync(string key)
    {
        await _cache.RemoveAsync(key);
    }
}