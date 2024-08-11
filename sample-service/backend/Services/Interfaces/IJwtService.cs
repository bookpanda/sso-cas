using backend.Models;

namespace backend.Interfaces;

public interface IJwtService
{
    string CreateToken(User user);
    bool ValidateToken(string token);
}