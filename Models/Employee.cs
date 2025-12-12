using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace Jotnal.Models;

/// <summary>
/// Сотрудник/Пользователь системы
/// </summary>
public class Employee
{
    [Key]
    public int Id { get; set; }

    /// <summary>
    /// Фамилия
    /// </summary>
    [Required]
    [MaxLength(100)]
    public string LastName { get; set; } = string.Empty;

    /// <summary>
    /// Имя
    /// </summary>
    [Required]
    [MaxLength(100)]
    public string FirstName { get; set; } = string.Empty;

    /// <summary>
    /// Отчество
    /// </summary>
    [MaxLength(100)]
    public string? MiddleName { get; set; }

    /// <summary>
    /// Полное ФИО
    /// </summary>
    [NotMapped]
    public string FullName => $"{LastName} {FirstName} {MiddleName}".Trim();

    /// <summary>
    /// Телефон
    /// </summary>
    [MaxLength(20)]
    public string? Phone { get; set; }

    /// <summary>
    /// Дата рождения
    /// </summary>
    public DateTime? BirthDate { get; set; }

    /// <summary>
    /// Табельный номер
    /// </summary>
    [Required]
    [MaxLength(50)]
    public string PersonnelNumber { get; set; } = string.Empty;

    /// <summary>
    /// Дата приема на работу
    /// </summary>
    [Required]
    public DateTime HireDate { get; set; }

    /// <summary>
    /// Работает сейчас
    /// </summary>
    public bool IsCurrentlyEmployed { get; set; } = true;

    /// <summary>
    /// Дата увольнения (если уволен)
    /// </summary>
    public DateTime? TerminationDate { get; set; }

    /// <summary>
    /// Имя пользователя Windows для авторизации
    /// </summary>
    [Required]
    [MaxLength(200)]
    public string WindowsUsername { get; set; } = string.Empty;

    /// <summary>
    /// Является ли пользователь активным в системе
    /// </summary>
    public bool IsActive { get; set; } = true;

    /// <summary>
    /// Роль пользователя в системе
    /// </summary>
    [Required]
    public UserRole Role { get; set; } = UserRole.User;

    /// <summary>
    /// Должность
    /// </summary>
    public int PositionId { get; set; }

    [ForeignKey(nameof(PositionId))]
    public virtual Position Position { get; set; } = null!;

    /// <summary>
    /// Структурное подразделение
    /// </summary>
    public int DepartmentId { get; set; }

    [ForeignKey(nameof(DepartmentId))]
    public virtual Department Department { get; set; } = null!;

    public DateTime CreatedAt { get; set; } = DateTime.Now;
    public DateTime? UpdatedAt { get; set; }

    /// <summary>
    /// Последний вход в систему
    /// </summary>
    public DateTime? LastLoginAt { get; set; }
}

/// <summary>
/// Роли пользователей в системе
/// </summary>
public enum UserRole
{
    /// <summary>
    /// Обычный пользователь
    /// </summary>
    User = 0,

    /// <summary>
    /// Администратор
    /// </summary>
    Administrator = 1,

    /// <summary>
    /// Разработчик (специальный пользователь с максимальными правами)
    /// </summary>
    Developer = 2
}
