using System.Net;
using backend.DTO;
using backend.Exceptions;
using backend.Models;
using backend.Repositories.Interfaces;
using backend.Services.Interfaces;

namespace backend.Services;

public class UserService : IUserService
{
    private readonly IUserRepository _repo;
    private readonly ILogger<UserService> _logger;

    public UserService(IUserRepository repo, ILogger<UserService> logger)
    {
        _repo = repo;
        _logger = logger;
    }

    public async Task<User?> FindOne(string casid)
    {
        var user = await _repo.FindOne(casid);
        if (user == null)
        {
            _logger.LogInformation($"User with CASID {casid} not found");
            throw new ServiceException("User not found", HttpStatusCode.NotFound);
        }

        return user;
    }

    public async Task<User> Create(CreateUserDTO user)
    {
        var newUser = new User
        {
            CASID = user.CASID,
        };

        try
        {
            await _repo.Create(newUser);
        }
        catch (Exception ex)
        {
            throw new ServiceException("Error creating user", HttpStatusCode.InternalServerError, ex);
        }

        return newUser;
    }
}