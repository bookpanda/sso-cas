using backend.Models;
using backend.Repositories.Interfaces;

namespace backend.Repositories;

public class UserRepository : IUserRepository
{
    private readonly ApplicationDbContext _context;

    public UserRepository(ApplicationDbContext context)
    {
        _context = context;
    }

    public async Task<User?> FindOne(long id)
    {
        return await _context.Users.FindAsync(id);
    }

    public async Task<User> Create(User user)
    {
        _context.Users.Add(user);
        await _context.SaveChangesAsync();

        return user;
    }
}