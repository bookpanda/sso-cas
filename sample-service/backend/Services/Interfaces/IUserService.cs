using backend.DTO;

namespace backend.Services.Interfaces;

public interface IUserService
{
    Task<UserDTO?> FindOne(string id);
    Task<UserDTO> Create(UserDTO user);
}