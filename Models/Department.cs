using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace Jotnal.Models;

/// <summary>
/// Структурное подразделение с иерархической структурой
/// </summary>
public class Department
{
    [Key]
    public int Id { get; set; }

    [Required]
    [MaxLength(200)]
    public string Name { get; set; } = string.Empty;

    [MaxLength(500)]
    public string? Description { get; set; }

    // Иерархия - ссылка на родительское подразделение
    public int? ParentDepartmentId { get; set; }

    [ForeignKey(nameof(ParentDepartmentId))]
    public virtual Department? ParentDepartment { get; set; }

    public bool IsActive { get; set; } = true;

    public DateTime CreatedAt { get; set; } = DateTime.Now;
    public DateTime? UpdatedAt { get; set; }

    // Navigation properties
    public virtual ICollection<Department> ChildDepartments { get; set; } = new List<Department>();
    public virtual ICollection<Employee> Employees { get; set; } = new List<Employee>();
}
