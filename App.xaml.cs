using System.Windows;
using Microsoft.EntityFrameworkCore;
using Jotnal.Data;

namespace Jotnal;

public partial class App : Application
{
    protected override void OnStartup(StartupEventArgs e)
    {
        base.OnStartup(e);

        // Инициализация базы данных при запуске
        using (var context = new JotnalDbContext(GetDbContextOptions()))
        {
            // Создаем базу данных, если она не существует
            context.Database.EnsureCreated();
        }
    }

    public static DbContextOptions<JotnalDbContext> GetDbContextOptions()
    {
        var optionsBuilder = new DbContextOptionsBuilder<JotnalDbContext>();
        optionsBuilder.UseSqlite("Data Source=jotnal.db");
        return optionsBuilder.Options;
    }
}
