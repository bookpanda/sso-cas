using backend.DTO;
using backend.Services.Interfaces;

namespace TodoApi.Services;

public class AuthService : IAuthService
{
    private readonly IUserService _userService;
    private readonly ITokenService _tokenService;
    private readonly ILogger<AuthService> _logger;

    public AuthService(IUserService userService, ITokenService tokenService, ILogger<AuthService> logger)
    {
        _userService = userService;
        _tokenService = tokenService;
        _logger = logger;
    }

    public async Task<IEnumerable<AuthToken>> AuthenticateSSO(string userID)
    {
        var user = await _userService.FindOne(userID);
        throw new NotImplementedException();
    }
}