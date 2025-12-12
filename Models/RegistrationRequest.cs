using System.ComponentModel.DataAnnotations;

namespace Jotnal.Models;

/// <summary>
/// Запрос на регистрацию в системе
/// </summary>
public class RegistrationRequest
{
    [Key]
    public int Id { get; set; }

    [Required]
    [MaxLength(100)]
    public string LastName { get; set; } = string.Empty;

    [Required]
    [MaxLength(100)]
    public string FirstName { get; set; } = string.Empty;

    [MaxLength(100)]
    public string? MiddleName { get; set; }

    [MaxLength(20)]
    public string? Phone { get; set; }

    public DateTime? BirthDate { get; set; }

    [Required]
    [MaxLength(50)]
    public string PersonnelNumber { get; set; } = string.Empty;

    [Required]
    public DateTime HireDate { get; set; }

    [Required]
    [MaxLength(200)]
    public string WindowsUsername { get; set; } = string.Empty;

    [Required]
    [MaxLength(200)]
    public string RequestedPosition { get; set; } = string.Empty;

    [Required]
    [MaxLength(200)]
    public string RequestedDepartment { get; set; } = string.Empty;

    [MaxLength(1000)]
    public string? Comments { get; set; }

    public RegistrationRequestStatus Status { get; set; } = RegistrationRequestStatus.Pending;

    public DateTime RequestedAt { get; set; } = DateTime.Now;

    public int? ApprovedByEmployeeId { get; set; }

    public DateTime? ProcessedAt { get; set; }

    [MaxLength(500)]
    public string? RejectionReason { get; set; }
}

/// <summary>
/// Статус запроса на регистрацию
/// </summary>
public enum RegistrationRequestStatus
{
    /// <summary>
    /// Ожидает рассмотрения
    /// </summary>
    Pending = 0,

    /// <summary>
    /// Одобрен
    /// </summary>
    Approved = 1,

    /// <summary>
    /// Отклонен
    /// </summary>
    Rejected = 2
}
