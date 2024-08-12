using System.Text.Json.Serialization;

namespace backend.DTO;

public class UserDTO
{
    [JsonPropertyName("id")]
    public required string ID { get; set; }
    [JsonPropertyName("email")]
    public required string Email { get; set; }
}

public class CreateUserDTO
{
    [JsonPropertyName("email")]
    public required string Email { get; set; }
}
