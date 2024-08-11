using backend.Services.Interfaces;
using backend.Models;
using backend.Repositories.Interfaces;
using backend.Interfaces;
using backend.DTO;

namespace tiermaker_backend.Services;

public class TokenService : ITokenService
{
    private readonly ICacheRepository _cache;
    private readonly IJwtService _jwtSvc;
    private readonly IConfiguration _config;
    private readonly int _refreshTTLDays;

    public TokenService(ICacheRepository cache, IJwtService jwtSvc, IConfiguration config)
    {
        _cache = cache;
        _jwtSvc = jwtSvc;
        _config = config;
        _refreshTTLDays = _config.GetValue<int>("JWT:RefreshTTLDays");
    }

    public async Task<AuthToken> GetCredentials(User user)
    {
        var session = await _cache.GetAsync<AuthToken>(SessionKey(user.ID));
        if (session == null)
        {
            session = await CreateCredentials(user);
        }

        var isTokenValid = _jwtSvc.ValidateToken(session.AccessToken);
        if (!isTokenValid)
        {
            await _cache.RemoveAsync(SessionKey(user.ID));
            string accessToken = _jwtSvc.CreateToken(user);

            var credentials = new AuthToken
            {
                AccessToken = accessToken,
                RefreshToken = session.RefreshToken,
                ExpiresIn = DateTime.UtcNow.AddMinutes(_config.GetValue<int>("JWT:ExpirationMinutes"))
            };
            await _cache.SetAsync(SessionKey(user.ID), credentials, TimeSpan.FromDays(_refreshTTLDays));

            return credentials;
        }

        return session;
    }

    public async Task<AuthToken> CreateCredentials(User user)
    {
        string accessToken = _jwtSvc.CreateToken(user);
        string refreshToken = CreateRefreshToken();

        await _cache.SetAsync(RefreshKey(refreshToken), user, TimeSpan.FromDays(_refreshTTLDays));

        var credentials = new AuthToken
        {
            AccessToken = accessToken,
            RefreshToken = refreshToken,
            ExpiresIn = DateTime.UtcNow.AddMinutes(_config.GetValue<int>("JWT:ExpirationMinutes"))
        };

        await _cache.SetAsync(SessionKey(user.ID), credentials, TimeSpan.FromDays(_refreshTTLDays));

        return credentials;
    }

    public async Task<AuthToken?> RefreshToken(string refreshToken)
    {
        var user = await _cache.GetAsync<User>(RefreshKey(refreshToken));
        if (user == null) return null;

        await _cache.RemoveAsync(RefreshKey(refreshToken));
        await _cache.RemoveAsync(SessionKey(user.ID));

        var credentials = await CreateCredentials(user);

        return credentials;
    }

    private string CreateRefreshToken()
    {
        return Guid.NewGuid().ToString();
    }

    private string RefreshKey(string refreshToken)
    {
        return $"refresh_{refreshToken}";
    }

    private string SessionKey(string? userID)
    {
        return $"session_{userID}";
    }
}
