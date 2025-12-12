using System.Windows;
using Jotnal.Data;
using Jotnal.Services;

namespace Jotnal;

public partial class MainWindow : Window
{
    private readonly JotnalDbContext _context;
    private readonly AuthenticationService _authService;

    public MainWindow()
    {
        InitializeComponent();

        _context = new JotnalDbContext(App.GetDbContextOptions());
        _authService = new AuthenticationService(_context);

        Loaded += MainWindow_Loaded;
    }

    private async void MainWindow_Loaded(object sender, RoutedEventArgs e)
    {
        await AuthenticateUser();
    }

    private async Task AuthenticateUser()
    {
        var result = await _authService.AuthenticateAsync();

        if (result.IsSuccess && result.User != null)
        {
            UserNameText.Text = result.User.FullName;
            UserRoleText.Text = $"{result.User.Position.Name} - {result.User.Department.Name}";
            StatusText.Text = result.Message;

            // Показываем административную вкладку только для администраторов
            AdminTab.Visibility = _authService.IsAdministrator()
                ? Visibility.Visible
                : Visibility.Collapsed;
        }
        else
        {
            UserNameText.Text = "Не авторизован";
            UserRoleText.Text = result.Message;
            StatusText.Text = "Требуется регистрация";
            AdminTab.Visibility = Visibility.Collapsed;

            MessageBox.Show(result.Message, "Авторизация",
                MessageBoxButton.OK, MessageBoxImage.Warning);
        }
    }

    private async void SwitchUserButton_Click(object sender, RoutedEventArgs e)
    {
        _authService.Logout();
        await AuthenticateUser();
    }

    private void RegistrationRequestsButton_Click(object sender, RoutedEventArgs e)
    {
        // TODO: Открыть окно с запросами на регистрацию
        MessageBox.Show("Функция в разработке", "Информация",
            MessageBoxButton.OK, MessageBoxImage.Information);
    }

    protected override void OnClosed(EventArgs e)
    {
        base.OnClosed(e);
        _context?.Dispose();
    }
}
