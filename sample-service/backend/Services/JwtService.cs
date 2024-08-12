using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;
using Microsoft.IdentityModel.Tokens;
using backend.Interfaces;
using backend.Models;
using backend.Config;
using Microsoft.Extensions.Options;

namespace backend.Services;

public class JwtService : IJwtService
{
    private readonly JwtConfig _config;
    private readonly SymmetricSecurityKey _key;

    public JwtService(IOptions<JwtConfig> config)
    {
        _config = config.Value;
        _key = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(
            _config.Secret ?? throw new InvalidOperationException("JWT Secret is missing")));
    }

    public string CreateToken(User user)
    {
        var claims = new List<Claim>
        {
            new Claim("email", user.Email ?? throw new InvalidOperationException("User Email is missing")),
            new Claim("userID", user.ID ?? throw new InvalidOperationException("User Name is missing")),
        };

        var creds = new SigningCredentials(_key, SecurityAlgorithms.HmacSha512Signature);
        var accessTTL = _config.AccessTTL;

        var tokenDescriptor = new SecurityTokenDescriptor
        {
            Subject = new ClaimsIdentity(claims),
            IssuedAt = DateTime.Now,
            Expires = DateTime.Now.AddSeconds(accessTTL),
            SigningCredentials = creds,
            Issuer = _config.Issuer
        };

        var tokenHandler = new JwtSecurityTokenHandler();

        var token = tokenHandler.CreateToken(tokenDescriptor);

        return tokenHandler.WriteToken(token);
    }

    public ClaimsPrincipal? ValidateToken(string token)
    {
        var tokenHandler = new JwtSecurityTokenHandler();
        try
        {
            var claims = tokenHandler.ValidateToken(token, new TokenValidationParameters
            {
                ValidateIssuer = true,
                ValidIssuer = _config.Issuer,
                ValidateAudience = true,
                ValidateLifetime = true,
                ClockSkew = TimeSpan.Zero,
                ValidateIssuerSigningKey = true,
                IssuerSigningKey = _key,
            }, out var validatedToken);

            return claims;
        }
        catch
        {
            return null;
        }
    }
}
