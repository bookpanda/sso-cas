using backend.DTO;
using backend.Services.Interfaces;

namespace TodoApi.Services;

public class AuthService : IAuthService
{
    private readonly ILogger<AuthService> _logger;

    public AuthService(ILogger<AuthService> logger)
    {
        // _repo = repo;
        _logger = logger;
    }

    public async Task<IEnumerable<AuthToken>> AuthenticateSSO()
    {
        throw new NotImplementedException();
    }

    public async Task<IEnumerable<AuthToken>> RefreshToken()
    {
        throw new NotImplementedException();
    }

    public async Task<Credentials?> Validate()
    {
        throw new NotImplementedException();
    }
}