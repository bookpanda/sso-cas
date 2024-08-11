using System.Net;
using backend.DTO;
using backend.Exceptions;
using backend.Models;
using backend.Parsers;
using backend.Repositories.Interfaces;
using backend.Services.Interfaces;

namespace TodoApi.Services;

public class UserService : IUserService
{
    private readonly IUserRepository _repo;
    private readonly ILogger<UserService> _logger;

    public UserService(IUserRepository repo, ILogger<UserService> logger)
    {
        _repo = repo;
        _logger = logger;
    }

    public async Task<UserDTO?> FindOne(string id)
    {
        var user = await _repo.FindOne(id);
        if (user == null)
        {
            _logger.LogInformation($"User with id {id} not found");
            throw new ServiceException("User not found", HttpStatusCode.NotFound);
        }

        return UserParser.ModelToDTO(user);
    }

    public async Task<UserDTO> Create(UserDTO user)
    {
        var newUser = new User
        {
            Email = user.Email,
        };

        try
        {
            await _repo.Create(newUser);
        }
        catch (Exception ex)
        {
            throw new ServiceException("Error creating user", HttpStatusCode.InternalServerError, ex);
        }

        return UserParser.ModelToDTO(newUser);
    }
}