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

    public async Task<AuthToken> AuthenticateSSO(string userCASID)
    {
        Models.User? user = null;
        try
        {
            user = await _userService.FindOne(userCASID);
        }
        catch (ServiceException ex)
        {
            if (ex.StatusCode == System.Net.HttpStatusCode.NotFound)
            {
                _logger.LogInformation($"User with id {userCASID} not found, creating new user");

                var newUser = new CreateUserDTO { CASID = userCASID };
                var createdUser = await _userService.Create(newUser);

                var newAuthToken = await _tokenService.GetCredentials(createdUser);
                return newAuthToken;
            }
        }

        var authToken = await _tokenService.GetCredentials(user!);
        return authToken;
    }
}