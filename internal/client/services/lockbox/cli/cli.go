package cli

import (
	"context"
	"fmt"
	"gophKeeper/internal/client/services/lockbox/models"
	"gophKeeper/internal/client/services/lockbox/usecase"

	"github.com/spf13/cobra"
)

type LockBoxCLI struct {
	lockBoxUC usecase.ILockBoxUsecase
}

func NewLockBoxCLI(lockBoxUC usecase.ILockBoxUsecase) *LockBoxCLI {
	return &LockBoxCLI{lockBoxUC: lockBoxUC}
}

func (cli *LockBoxCLI) CreateCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createLock",
		Short: "Create a new lockbox",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			url, _ := cmd.Flags().GetString("url")
			login, _ := cmd.Flags().GetString("login")
			password, _ := cmd.Flags().GetString("password")
			description, _ := cmd.Flags().GetString("description")

			input := models.LockBoxInput{
				Name:        name,
				URL:         url,
				Login:       login,
				Password:    password,
				Description: description,
			}

			_, err := cli.lockBoxUC.CreateLockBox(ctx, &input)
			if err != nil {
				return
			}

			fmt.Println("✅ LockBox успешно создан:")
		},
	}
	cmd.Flags().String("name", "", "Lockbox name")
	cmd.Flags().String("url", "", "URL (необязательно)")
	cmd.Flags().String("login", "", "Login (необязательно)")
	cmd.Flags().String("password", "", "Password (необязательно)")
	cmd.Flags().String("description", "", "Description (необязательно)")

	return cmd
}

func (cli *LockBoxCLI) DeleteCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteLock",
		Short: "Delete a lockbox",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")

			if name == "" {
				fmt.Println("❌ Ошибка: имя LockBox обязательно")
				return
			}

			if err := cli.lockBoxUC.DeleteLockBox(ctx, name); err != nil {
				fmt.Println("❌ Ошибка удаления")
				return
			}
			fmt.Println("✅ Lockbox успешно удалён!")
		},
	}

	cmd.Flags().String("name", "", "Название LockBox (обязательно)")

	return cmd
}
func (cli *LockBoxCLI) GetCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getLock",
		Short: "Get a lockbox by name",
		Run: func(cmd *cobra.Command, args []string) {
			name, err := cmd.Flags().GetString("name")
			if err != nil || name == "" {
				fmt.Println("❌ Ошибка: не указан параметр name")
				return
			}

			lockBox, err := cli.lockBoxUC.GetLockBoxById(ctx, name)
			if err != nil {
				fmt.Println("❌ Ошибка, данное хранилище не найдено")
				return
			}
			if lockBox == nil {
				fmt.Println("❌ Ошибка, пустое хранилище")
				return
			}
			fmt.Println("\n✅ Lockbox найден!")
			fmt.Println("──────────────────────────────────────────────")
			fmt.Printf("🔹 Название:     %s\n", lockBox.Name)
			fmt.Printf("🔗 URL:          %s\n", lockBox.URL)
			fmt.Printf("👤 Логин:        %s\n", lockBox.Login)
			fmt.Printf("🔑 Пароль:       %s\n", lockBox.Password)
			fmt.Printf("📝 Описание:     %s\n", lockBox.Description)
			fmt.Printf("📅 Дата создания:%s\n", lockBox.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("♻️  Обновлено:   %s\n", lockBox.UpdatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println("──────────────────────────────────────────────")
		},
	}

	cmd.Flags().String("name", "", "Lockbox name (обязательно)")

	return cmd
}

func (cli *LockBoxCLI) UpdateCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateLock",
		Short: "Update a lockbox",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			url, _ := cmd.Flags().GetString("url")
			login, _ := cmd.Flags().GetString("login")
			password, _ := cmd.Flags().GetString("password")
			description, _ := cmd.Flags().GetString("description")

			input := models.LockBoxInput{
				Name:        name,
				URL:         url,
				Login:       login,
				Password:    password,
				Description: description,
			}

			if err := cli.lockBoxUC.UpdateLockBox(ctx, &input); err != nil {
				fmt.Println("❌ Ошибка в обновлении данных")
				return
			}

			fmt.Println("✅ Lockbox успешно обновлён!")
		},
	}

	cmd.Flags().String("name", "", "Название LockBox (обязательно)")
	cmd.Flags().String("url", "", "URL")
	cmd.Flags().String("login", "", "Login")
	cmd.Flags().String("password", "", "Password")
	cmd.Flags().String("description", "", "Description")

	return cmd
}
func (cli *LockBoxCLI) GetAllCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getLocks",
		Short: "Get all lockboxes",
		Run: func(cmd *cobra.Command, args []string) {
			lockBoxes, err := cli.lockBoxUC.GetLockBoxAll(ctx)
			if err != nil {
				fmt.Println("❌ Ошибка в получении хранилищ:")
				return
			}

			if len(*lockBoxes) == 0 {
				fmt.Println("🔍 Нет сохранённых Lockbox.")
				return
			}

			fmt.Println("\n📦 Список Lockbox:")
			fmt.Println("──────────────────────────────────────────────────────────────────────")
			for i, lockBox := range *lockBoxes {
				fmt.Printf("[%d] 🔹 Название:  %s\n", i+1, lockBox.Name)
				fmt.Printf("    🔗 URL:       %s\n", lockBox.URL)
				fmt.Printf("    👤 Логин:     %s\n", lockBox.Login)
				fmt.Printf("    📝 Описание:  %s\n", lockBox.Description)
				fmt.Printf("    📅 Создан:    %s\n", lockBox.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Printf("    ♻️  Обновлено: %s\n", lockBox.UpdatedAt.Format("2006-01-02 15:04:05"))
				fmt.Println("──────────────────────────────────────────────────────────────────────")
			}
		},
	}

	return cmd
}

func (cli *LockBoxCLI) RegisterCommands(ctx context.Context, root *cobra.Command) {
	root.AddCommand(
		cli.CreateCommand(ctx),
		cli.DeleteCommand(ctx),
		cli.GetCommand(ctx),
		cli.UpdateCommand(ctx),
		cli.GetAllCommand(ctx),
	)
}
func (cli *LockBoxCLI) NewRegisterCli(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new user",
		Run: func(cmd *cobra.Command, args []string) {
			username, err := cmd.Flags().GetString("username")
			if err != nil || username == "" {
				fmt.Println("❌ Ошибка: Укажите имя пользователя")
				return
			}

			password, err := cmd.Flags().GetString("password")
			if err != nil || password == "" {
				fmt.Println("❌ Ошибка: Укажите пароль")
				return
			}

			if err := cli.lockBoxUC.Register(ctx, username, password); err != nil {
				fmt.Println("❌ Ошибка:", err)
				return
			}

			fmt.Println("✅ Регистрация успешна!")
		},
	}

	cmd.Flags().String("username", "", "Username (обязательно)")
	cmd.Flags().String("password", "", "Password (обязательно)")
	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")

	return cmd
}

func (cli *LockBoxCLI) NewAuthCli(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Login to account",
		Run: func(cmd *cobra.Command, args []string) {
			username, err := cmd.Flags().GetString("username")
			if err != nil || username == "" {
				fmt.Println("❌ Ошибка: Укажите имя пользователя")
				return
			}

			password, err := cmd.Flags().GetString("password")
			if err != nil || password == "" {
				fmt.Println("❌ Ошибка: Укажите пароль")
				return
			}

			if err := cli.lockBoxUC.Authenticate(ctx, username, password); err != nil {
				fmt.Println("❌ Ошибка:", err)
				return
			}

			fmt.Println("✅ Аутентификация успешна!")
		},
	}

	cmd.Flags().String("username", "", "Username (обязательно)")
	cmd.Flags().String("password", "", "Password (обязательно)")
	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")

	return cmd
}

func (cli *LockBoxCLI) IsAuthenticated() bool {
	return cli.lockBoxUC.IsAuthenticated()
}
