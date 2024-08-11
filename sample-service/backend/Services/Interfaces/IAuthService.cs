using backend.DTO;

namespace backend.Services.Interfaces;

public interface IAuthService
{
    Task<IEnumerable<AuthToken>> AuthenticateSSO(string userID);
}