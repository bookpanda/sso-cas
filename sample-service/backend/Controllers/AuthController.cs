using Microsoft.AspNetCore.Mvc;
using backend.Services.Interfaces;
using backend.Config;
using backend.DTO;
using backend.Exceptions;
using Microsoft.Extensions.Options;

namespace backend.Controllers;

[Route("auth")]
[ApiController]
public class AuthController : ControllerBase
{
    private readonly SSOConfig _config;
    private readonly IAuthService _authSvc;
    private readonly ITokenService _tokenSvc;
    private readonly HttpClient _httpClient;
    private readonly ILogger<AuthController> _logger;

    public AuthController(IOptions<SSOConfig> config, IAuthService authSvc, ITokenService tokenSvc, HttpClient httpClient, ILogger<AuthController> logger)
    {
        _config = config.Value;
        _authSvc = authSvc;
        _tokenSvc = tokenSvc;
        _httpClient = httpClient;
        _logger = logger;
    }

    [HttpGet("auth-sso")]
    public async Task<IActionResult> AuthenticateSSO([FromQuery] string ticket)
    { // only after redirect from CAS
        try
        {
            var response = await _httpClient.GetAsync(_config.Authority + "/api/v1/auth/validate-st?ticket=" + ticket);
            if (!response.IsSuccessStatusCode) return Unauthorized("Invalid ticket");

            var session = await response.Content.ReadFromJsonAsync<SessionCAS>();
            if (session == null) return Unauthorized("Invalid session");

            var authToken = await _authSvc.AuthenticateSSO(session.UserID);

            return Ok(authToken);
        }
        catch (ServiceException ex)
        {
            _logger.LogError(ex, "Error authenticating SSO");
            return StatusCode((int)ex.StatusCode, ex.Message);
        }
    }

    [HttpPost("refresh-token")]
    public async Task<IActionResult> RefreshToken([FromBody] RefreshTokenBody request)
    {
        if (!ModelState.IsValid) return BadRequest(ModelState);

        try
        {
            var credentials = await _tokenSvc.RefreshToken(request.RefreshToken);
            if (credentials == null) return Unauthorized("Invalid refresh token");

            return Ok(credentials);
        }
        catch (ServiceException ex)
        {
            _logger.LogError(ex, "Error refreshing token");
            return StatusCode((int)ex.StatusCode, ex.Message);
        }
    }

    [HttpGet("validate")]
    public async Task<IActionResult> Validate([FromHeader(Name = "Authorization")] string authHeader)
    {
        var accessToken = authHeader?.Split(" ").Last();
        if (string.IsNullOrEmpty(accessToken))
        {
            return BadRequest("Authorization header is missing.");
        }

        try
        {
            var creds = await _tokenSvc.ValidateToken(accessToken);
            return Ok(creds);
        }
        catch (ServiceException ex)
        {
            _logger.LogError(ex, "Error validating token");
            return StatusCode((int)ex.StatusCode, ex.Message);
        }
    }
}