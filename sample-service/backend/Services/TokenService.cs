using backend.Services.Interfaces;
using backend.Models;
using backend.Repositories.Interfaces;
using backend.Interfaces;
using backend.DTO;
using System.Security.Claims;
using backend.Config;
using Microsoft.Extensions.Options;
using Microsoft.EntityFrameworkCore.Metadata.Internal;

namespace backend.Services;

public class TokenService : ITokenService
{
    private readonly ICacheRepository _cache;
    private readonly IJwtService _jwtSvc;
    private readonly JwtConfig _config;

    public TokenService(ICacheRepository cache, IJwtService jwtSvc, IOptions<JwtConfig> config)
    {
        _cache = cache;
        _jwtSvc = jwtSvc;
        _config = config.Value;
    }

    public async Task<AuthToken> GetCredentials(User user, DateTime refreshExpiry)
    {
        var session = await _cache.GetAsync<AuthToken>(SessionKey(user.ID));
        if (session == null)
        {
            session = await CreateCredentials(user, refreshExpiry);
        }

        var claims = _jwtSvc.ValidateToken(session.AccessToken);
        if (claims == null)
        {
            await RemoveSessionCache(user.ID);
            // await _cache.RemoveAsync(SessionKey(user.ID));
            string accessToken = _jwtSvc.CreateToken(user);

            var credentials = new AuthToken
            {
                AccessToken = accessToken,
                RefreshToken = session.RefreshToken,
                ExpiresIn = DateTime.UtcNow.AddSeconds(_config.AccessTTL)
            };
            await _cache.SetAsync(SessionKey(user.ID), credentials, TimeSpan.FromSeconds(_config.AccessTTL));

            return credentials;
        }

        return session;
    }

    public async Task<AuthToken?> RefreshToken(string refreshToken)
    {
        var user = await _cache.GetAsync<User>(RefreshKey(refreshToken));
        if (user == null) return null;

        // await _cache.RemoveAsync(SessionKey(user.ID));
        await RemoveSessionCache(user.ID);

        string accessToken = _jwtSvc.CreateToken(user);
        var credentials = new AuthToken
        {
            AccessToken = accessToken,
            RefreshToken = refreshToken,
            ExpiresIn = DateTime.UtcNow.AddMinutes(_config.AccessTTL)
        };

        await _cache.SetAsync(SessionKey(user.ID), credentials, TimeSpan.FromSeconds(_config.AccessTTL));

        return credentials;
    }

    public async Task<Credentials?> ValidateToken(string accessToken)
    {
        var claims = _jwtSvc.ValidateToken(accessToken);
        if (claims == null) return null;

        var userID = claims.FindFirstValue("userID") ?? throw new InvalidOperationException("User ID is missing");

        var session = await _cache.GetAsync<AuthToken>(SessionKey(userID));
        if (session == null) return null;

        return new Credentials
        {
            UserID = userID,
            Role = "user",
        };
    }

    public async Task<AuthToken?> GetSessionCache(string? userID)
    {
        return await _cache.GetAsync<AuthToken>(SessionKey(userID));
    }

    public async Task RemoveSessionCache(string? userID)
    {
        await _cache.RemoveAsync(SessionKey(userID));
    }

    public async Task RemoveRefreshCache(string? refreshToken)
    {
        await _cache.RemoveAsync(RefreshKey(refreshToken));
    }

    private async Task<AuthToken> CreateCredentials(User user, DateTime refreshExpiry)
    {
        string accessToken = _jwtSvc.CreateToken(user);
        string refreshToken = CreateRefreshToken();

        await _cache.SetAsync(RefreshKey(refreshToken), user, refreshExpiry - DateTime.Now);

        var credentials = new AuthToken
        {
            AccessToken = accessToken,
            RefreshToken = refreshToken,
            ExpiresIn = DateTime.UtcNow.AddSeconds(_config.AccessTTL)
        };

        await _cache.SetAsync(SessionKey(user.ID), credentials, TimeSpan.FromSeconds(_config.AccessTTL));

        return credentials;
    }

    private string CreateRefreshToken()
    {
        return Guid.NewGuid().ToString();
    }

    private string RefreshKey(string? refreshToken)
    {
        return $"refresh_{refreshToken}";
    }

    private string SessionKey(string? userID)
    {
        return $"session_{userID}";
    }
}
