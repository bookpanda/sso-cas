namespace backend.DTO;

public class UserDTO
{
    public required string ID { get; set; }
    public required string Email { get; set; }
}

public class CreateUserDTO
{
    public required string Email { get; set; }
}
