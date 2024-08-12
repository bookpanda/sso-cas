using backend.Config;
using backend.Repositories;
using backend.Repositories.Interfaces;
using backend.Services.Interfaces;
using backend.Services;
using backend.Interfaces;

namespace Microsoft.Extensions.DependencyInjection;

public static class MyConfigServiceCollectionExtensions
{
    public static IServiceCollection AddConfig(this IServiceCollection services, IConfiguration config)
    {
        services.Configure<JwtConfig>(config.GetSection(JwtConfig.Jwt));
        services.Configure<SSOConfig>(config.GetSection(SSOConfig.SSO));

        return services;
    }
    public static IServiceCollection AddMyDependencyGroup(
         this IServiceCollection services)
    {
        services.AddScoped<ICacheRepository, CacheRepository>();
        services.AddScoped<IUserRepository, UserRepository>();

        services.AddScoped<IAuthService, AuthService>();
        services.AddScoped<IJwtService, JwtService>();
        services.AddScoped<ITokenService, TokenService>();
        services.AddScoped<IUserService, UserService>();

        services.AddScoped<HttpClient>();

        return services;
    }
}
