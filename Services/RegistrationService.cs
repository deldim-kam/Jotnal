using Microsoft.EntityFrameworkCore;
using Jotnal.Data;
using Jotnal.Models;

namespace Jotnal.Services;

/// <summary>
/// Сервис для работы с запросами на регистрацию
/// </summary>
public class RegistrationService
{
    private readonly JotnalDbContext _context;

    public RegistrationService(JotnalDbContext context)
    {
        _context = context;
    }

    /// <summary>
    /// Создать запрос на регистрацию
    /// </summary>
    public async Task<bool> CreateRegistrationRequestAsync(RegistrationRequest request)
    {
        try
        {
            // Проверяем, есть ли уже такой пользователь
            var existingUser = await _context.Employees
                .FirstOrDefaultAsync(e => e.WindowsUsername == request.WindowsUsername);

            if (existingUser != null)
            {
                return false; // Пользователь уже существует
            }

            // Проверяем, нет ли уже активного запроса
            var existingRequest = await _context.RegistrationRequests
                .FirstOrDefaultAsync(r => r.WindowsUsername == request.WindowsUsername
                    && r.Status == RegistrationRequestStatus.Pending);

            if (existingRequest != null)
            {
                return false; // Уже есть активный запрос
            }

            _context.RegistrationRequests.Add(request);
            await _context.SaveChangesAsync();
            return true;
        }
        catch
        {
            return false;
        }
    }

    /// <summary>
    /// Получить все ожидающие запросы
    /// </summary>
    public async Task<List<RegistrationRequest>> GetPendingRequestsAsync()
    {
        return await _context.RegistrationRequests
            .Where(r => r.Status == RegistrationRequestStatus.Pending)
            .OrderBy(r => r.RequestedAt)
            .ToListAsync();
    }

    /// <summary>
    /// Одобрить запрос и создать пользователя
    /// </summary>
    public async Task<bool> ApproveRequestAsync(int requestId, int approvedByEmployeeId,
        int positionId, int departmentId)
    {
        try
        {
            var request = await _context.RegistrationRequests.FindAsync(requestId);
            if (request == null || request.Status != RegistrationRequestStatus.Pending)
            {
                return false;
            }

            // Создаем нового сотрудника
            var employee = new Employee
            {
                LastName = request.LastName,
                FirstName = request.FirstName,
                MiddleName = request.MiddleName,
                Phone = request.Phone,
                BirthDate = request.BirthDate,
                PersonnelNumber = request.PersonnelNumber,
                HireDate = request.HireDate,
                WindowsUsername = request.WindowsUsername,
                IsCurrentlyEmployed = true,
                IsActive = true,
                Role = UserRole.User,
                PositionId = positionId,
                DepartmentId = departmentId
            };

            _context.Employees.Add(employee);

            // Обновляем статус запроса
            request.Status = RegistrationRequestStatus.Approved;
            request.ApprovedByEmployeeId = approvedByEmployeeId;
            request.ProcessedAt = DateTime.Now;

            await _context.SaveChangesAsync();
            return true;
        }
        catch
        {
            return false;
        }
    }

    /// <summary>
    /// Отклонить запрос
    /// </summary>
    public async Task<bool> RejectRequestAsync(int requestId, int rejectedByEmployeeId, string reason)
    {
        try
        {
            var request = await _context.RegistrationRequests.FindAsync(requestId);
            if (request == null || request.Status != RegistrationRequestStatus.Pending)
            {
                return false;
            }

            request.Status = RegistrationRequestStatus.Rejected;
            request.ApprovedByEmployeeId = rejectedByEmployeeId;
            request.ProcessedAt = DateTime.Now;
            request.RejectionReason = reason;

            await _context.SaveChangesAsync();
            return true;
        }
        catch
        {
            return false;
        }
    }
}
