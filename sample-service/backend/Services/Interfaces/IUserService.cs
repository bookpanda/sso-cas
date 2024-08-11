using backend.DTO;

namespace backend.Services.Interfaces;

public interface IUserService
{
    Task<UserDTO?> GetUser(long id);
    Task<UserDTO> CreateUser(UserDTO user);
}