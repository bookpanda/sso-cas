using backend.DTO;
using backend.Exceptions;
using backend.Services.Interfaces;

namespace backend.Services;

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

    public async Task<AuthToken> AuthenticateSSO(SessionCAS sessionCAS)
    {
        Models.User? user = null;
        try
        {
            user = await _userService.FindOne(sessionCAS.UserID);
        }
        catch (ServiceException ex)
        {
            if (ex.StatusCode == System.Net.HttpStatusCode.NotFound)
            {
                _logger.LogInformation($"User with id {sessionCAS.UserID} not found, creating new user");

                var newUser = new CreateUserDTO { CASID = sessionCAS.UserID };
                var createdUser = await _userService.Create(newUser);

                var newAuthToken = await _tokenService.GetCredentials(createdUser, sessionCAS.ExpiresAt);
                return newAuthToken;
            }
        }

        var authToken = await _tokenService.GetCredentials(user!, sessionCAS.ExpiresAt);
        return authToken;
    }
}