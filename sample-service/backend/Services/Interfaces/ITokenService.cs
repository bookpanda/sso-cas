using backend.DTO;
using backend.Models;

namespace backend.Services.Interfaces;

public interface ITokenService
{
    Task<AuthToken> GetCredentials(User user);
    Task<AuthToken> CreateCredentials(User user);
    Task<AuthToken?> RefreshToken(string refreshToken);
}