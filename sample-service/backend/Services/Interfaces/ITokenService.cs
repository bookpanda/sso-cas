using backend.DTO;
using backend.Models;

namespace backend.Services.Interfaces;

public interface ITokenService
{
    Task<AuthToken> GetCredentials(User user, DateTime refreshExpiry);
    Task<AuthToken> CreateCredentials(User user, DateTime refreshExpiry);
    Task<AuthToken?> RefreshToken(string refreshToken);
    Task<Credentials?> ValidateToken(string accessToken);
}