using backend.Models;

namespace backend.Repositories.Interfaces;

public interface IUserRepository
{
    Task<User?> FindOne(long id);
    Task<User> Create(User todo);
}