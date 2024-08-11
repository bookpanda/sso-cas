namespace backend.DTO;

public class AuthToken
{
    public required string AccessToken { get; set; }
    public required string RefreshToken { get; set; }
    public required DateTime ExpiresIn { get; set; }
}

public class Credentials
{
    public required string UserID { get; set; }
    public required string Role { get; set; }
}