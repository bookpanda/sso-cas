using backend.DTO;

namespace backend.Services.Interfaces;

public interface IUserService
{
    Task<UserDTO?> FindOne(long id);
    Task<UserDTO> Create(UserDTO user);
}