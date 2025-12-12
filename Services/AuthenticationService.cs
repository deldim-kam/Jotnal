using System.Security.Principal;
using Microsoft.EntityFrameworkCore;
using Jotnal.Data;
using Jotnal.Models;

namespace Jotnal.Services;

/// <summary>
/// Сервис авторизации через Windows
/// </summary>
public class AuthenticationService
{
    private readonly JotnalDbContext _context;
    private Employee? _currentUser;

    public AuthenticationService(JotnalDbContext context)
    {
        _context = context;
    }

    /// <summary>
    /// Текущий авторизованный пользователь
    /// </summary>
    public Employee? CurrentUser => _currentUser;

    /// <summary>
    /// Авторизация через Windows username
    /// </summary>
    public async Task<AuthenticationResult> AuthenticateAsync()
    {
        try
        {
            // Получаем имя текущего пользователя Windows
            var windowsIdentity = WindowsIdentity.GetCurrent();
            var username = windowsIdentity.Name;

            // Извлекаем только имя пользователя без домена
            var shortUsername = username.Contains('\\')
                ? username.Split('\\')[1]
                : username;

            // Проверяем специального пользователя "Разработчик"
            if (shortUsername.Equals("DEVELOPER", StringComparison.OrdinalIgnoreCase))
            {
                _currentUser = await _context.Employees
                    .Include(e => e.Position)
                    .Include(e => e.Department)
                    .FirstOrDefaultAsync(e => e.WindowsUsername == "DEVELOPER");

                if (_currentUser != null)
                {
                    // НЕ обновляем LastLoginAt для разработчика (минимальные следы)
                    return new AuthenticationResult
                    {
                        IsSuccess = true,
                        User = _currentUser,
                        Message = "Вход выполнен как Разработчик"
                    };
                }
            }

            // Ищем пользователя в базе данных
            _currentUser = await _context.Employees
                .Include(e => e.Position)
                .Include(e => e.Department)
                .FirstOrDefaultAsync(e => e.WindowsUsername == shortUsername && e.IsActive);

            if (_currentUser == null)
            {
                return new AuthenticationResult
                {
                    IsSuccess = false,
                    Message = $"Пользователь '{shortUsername}' не найден в системе. Запросите регистрацию у администратора."
                };
            }

            if (!_currentUser.IsCurrentlyEmployed)
            {
                return new AuthenticationResult
                {
                    IsSuccess = false,
                    Message = "Пользователь уволен и не может войти в систему."
                };
            }

            // Обновляем время последнего входа
            _currentUser.LastLoginAt = DateTime.Now;
            await _context.SaveChangesAsync();

            return new AuthenticationResult
            {
                IsSuccess = true,
                User = _currentUser,
                Message = $"Добро пожаловать, {_currentUser.FullName}!"
            };
        }
        catch (Exception ex)
        {
            return new AuthenticationResult
            {
                IsSuccess = false,
                Message = $"Ошибка авторизации: {ex.Message}"
            };
        }
    }

    /// <summary>
    /// Сменить пользователя (выход)
    /// </summary>
    public void Logout()
    {
        _currentUser = null;
    }

    /// <summary>
    /// Проверка прав администратора
    /// </summary>
    public bool IsAdministrator()
    {
        return _currentUser?.Role == UserRole.Administrator || _currentUser?.Role == UserRole.Developer;
    }

    /// <summary>
    /// Проверка прав разработчика
    /// </summary>
    public bool IsDeveloper()
    {
        return _currentUser?.Role == UserRole.Developer;
    }
}

/// <summary>
/// Результат авторизации
/// </summary>
public class AuthenticationResult
{
    public bool IsSuccess { get; set; }
    public Employee? User { get; set; }
    public string Message { get; set; } = string.Empty;
}
