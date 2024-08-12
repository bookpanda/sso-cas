
using System.Security.Claims;
using backend.Models;

namespace backend.Interfaces;

public interface IJwtService
{
    string CreateToken(User user);
    ClaimsPrincipal? ValidateToken(string token);
}