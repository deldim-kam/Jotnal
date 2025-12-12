using System.ComponentModel.DataAnnotations;

namespace Jotnal.Models;

/// <summary>
/// Должность сотрудника
/// </summary>
public class Position
{
    [Key]
    public int Id { get; set; }

    [Required]
    [MaxLength(200)]
    public string Name { get; set; } = string.Empty;

    [MaxLength(500)]
    public string? Description { get; set; }

    public bool IsActive { get; set; } = true;

    public DateTime CreatedAt { get; set; } = DateTime.Now;
    public DateTime? UpdatedAt { get; set; }

    // Navigation property
    public virtual ICollection<Employee> Employees { get; set; } = new List<Employee>();
}
