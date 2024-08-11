using backend.DTO;
using backend.Services.Interfaces;

namespace TodoApi.Services;

public class UserService : IUserService
{
    private readonly ILogger<UserService> _logger;

    public UserService(ILogger<UserService> logger)
    {
        // _repo = repo;
        _logger = logger;
    }

    public async Task<UserDTO?> GetUser(long id)
    {
        throw new NotImplementedException();
    }

    public async Task<UserDTO> CreateUser(UserDTO user)
    {
        throw new NotImplementedException();
    }
}