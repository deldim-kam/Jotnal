using Microsoft.EntityFrameworkCore;
using Jotnal.Models;

namespace Jotnal.Data;

/// <summary>
/// Контекст базы данных для системы журнала смен
/// </summary>
public class JotnalDbContext : DbContext
{
    public DbSet<Position> Positions { get; set; } = null!;
    public DbSet<Department> Departments { get; set; } = null!;
    public DbSet<Employee> Employees { get; set; } = null!;
    public DbSet<RegistrationRequest> RegistrationRequests { get; set; } = null!;

    public JotnalDbContext(DbContextOptions<JotnalDbContext> options) : base(options)
    {
    }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        // Настройка иерархической структуры для Department
        modelBuilder.Entity<Department>()
            .HasOne(d => d.ParentDepartment)
            .WithMany(d => d.ChildDepartments)
            .HasForeignKey(d => d.ParentDepartmentId)
            .OnDelete(DeleteBehavior.Restrict);

        // Уникальный индекс для Windows username
        modelBuilder.Entity<Employee>()
            .HasIndex(e => e.WindowsUsername)
            .IsUnique();

        // Уникальный индекс для табельного номера
        modelBuilder.Entity<Employee>()
            .HasIndex(e => e.PersonnelNumber)
            .IsUnique();

        // Добавление начальных данных
        SeedData(modelBuilder);
    }

    private void SeedData(ModelBuilder modelBuilder)
    {
        // Добавляем базовые должности
        modelBuilder.Entity<Position>().HasData(
            new Position
            {
                Id = 1,
                Name = "Администратор системы",
                Description = "Администратор системы с полными правами",
                IsActive = true,
                CreatedAt = DateTime.Now
            },
            new Position
            {
                Id = 2,
                Name = "Разработчик",
                Description = "Разработчик системы",
                IsActive = true,
                CreatedAt = DateTime.Now
            }
        );

        // Добавляем корневое подразделение
        modelBuilder.Entity<Department>().HasData(
            new Department
            {
                Id = 1,
                Name = "Головная организация",
                Description = "Корневое подразделение",
                IsActive = true,
                CreatedAt = DateTime.Now
            }
        );

        // Добавляем специального пользователя "Разработчик"
        // Этот пользователь имеет минимальные следы в системе и максимальные права
        modelBuilder.Entity<Employee>().HasData(
            new Employee
            {
                Id = 1,
                LastName = "Системный",
                FirstName = "Разработчик",
                MiddleName = "",
                PersonnelNumber = "DEV-001",
                WindowsUsername = "DEVELOPER",
                HireDate = DateTime.Now,
                IsCurrentlyEmployed = true,
                IsActive = true,
                Role = UserRole.Developer,
                PositionId = 2,
                DepartmentId = 1,
                CreatedAt = DateTime.Now
            }
        );
    }
}
