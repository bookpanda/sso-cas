using Microsoft.AspNetCore.Mvc;
using backend.Services.Interfaces;
using backend.Config;
using backend.DTO;

namespace backend.Controllers;

[Route("auth")]
[ApiController]
public class AuthController : ControllerBase
{
    private readonly SSOConfig _config;
    private readonly IAuthService _authSvc;
    private readonly ITokenService _tokenSvc;
    private readonly HttpClient _httpClient;

    public AuthController(SSOConfig config, IAuthService authSvc, ITokenService tokenSvc, HttpClient httpClient)
    {
        _config = config;
        _authSvc = authSvc;
        _tokenSvc = tokenSvc;
        _httpClient = httpClient;
    }

    [HttpPost("auth-sso")]
    public async Task<IActionResult> AuthenticateSSO([FromQuery] string ticket)
    {
        // send request to SSO service
        // get user data   
        // create user
        // create token
        // return token
        var response = await _httpClient.GetAsync(_config.Authority + "/api/v1/validate-st?ticket=" + ticket);
        if (!response.IsSuccessStatusCode) return Unauthorized("Invalid ticket");

        var session = await response.Content.ReadFromJsonAsync<SessionCAS>();
        if (session == null) return Unauthorized("Invalid session");

        var authToken = await _authSvc.AuthenticateSSO(session.UserID);

        return Ok(authToken);
    }

    [HttpPost("refresh-token")]
    public async Task<IActionResult> RefreshToken([FromBody] RefreshTokenBody request)
    {
        if (!ModelState.IsValid) return BadRequest(ModelState);

        var credentials = await _tokenSvc.RefreshToken(request.RefreshToken);
        if (credentials == null) return Unauthorized("Invalid refresh token");

        return Ok(credentials);
    }
}