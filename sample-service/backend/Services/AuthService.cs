using System.Net;
using backend.DTO;
using backend.Exceptions;
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

    public async Task<AuthToken> AuthenticateSSO(string userID)
    {
        var user = await _userService.FindOne(userID);
        if (user == null)
        {
            _logger.LogInformation($"User with id {userID} not found, creating new user");

            var newUser = new CreateUserDTO { Email = userID };
            var createdUser = await _userService.Create(newUser);

            var newAuthToken = await _tokenService.GetCredentials(createdUser);
            return newAuthToken;
        }

        var authToken = await _tokenService.GetCredentials(user);
        return authToken;
    }
}