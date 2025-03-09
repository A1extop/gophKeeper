package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	"gophKeeper/internal/client/config"
	db "gophKeeper/internal/client/db"
	cli2 "gophKeeper/internal/client/services/lockbox/cli"
	clients2 "gophKeeper/internal/client/services/lockbox/clients"
	repos2 "gophKeeper/internal/client/services/lockbox/repository"
	usecase2 "gophKeeper/internal/client/services/lockbox/usecase"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// todo так, синхронизацию с серверной надо доделать
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.New()
	db, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	lockBoxCli := initCLI(ctx, cfg, db)
	if !authFlow(lockBoxCli, ctx) {
		fmt.Println("Ошибка аутентификации. Завершение работы.")
		return
	}

	commandLoop(lockBoxCli, ctx)
}

func initCLI(ctx context.Context, cfg *config.Config, db *sql.DB) *cli2.LockBoxCLI {
	lockBoxRepository := repos2.NewSQLiteRepository(db)
	lockBoxService := clients2.NewLockBoxService(cfg.PgHost, cfg.Port)
	lockBoxUsecase := usecase2.NewLockBoxUsecase(lockBoxService, lockBoxRepository)
	lockBoxCli := cli2.NewLockBoxCLI(lockBoxUsecase)

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = lockBoxUsecase.SyncUpdatesToServer(ctx)
				//if err != nil {
				//	fmt.Println(err)
				//}
				_ = lockBoxUsecase.SyncUpdatesToLocal(ctx)
				//if err != nil {
				//	fmt.Println(err)
				//}
			case <-ctx.Done():
				return
			}

		}
	}()
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = lockBoxRepository.PurgeExpiredLocks()
			case <-ctx.Done():
				return
			}

		}
	}()

	return lockBoxCli
}

func authFlow(lockBoxCli *cli2.LockBoxCLI, ctx context.Context) bool {
	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1. Регистрация")
		fmt.Println("2. Авторизация")
		fmt.Println("3. Выход")
		fmt.Print("Введите номер: ")

		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			continue
		}

		if choice == 3 {
			return false
		}

		fmt.Print("Имя пользователя: ")
		var username string
		_, err = fmt.Scanln(&username)
		if err != nil {
			fmt.Println("Ошибка ввода имени:", err)
			continue
		}

		fmt.Print("Пароль: ")
		var password string
		_, err = fmt.Scanln(&password)
		if err != nil {
			fmt.Println("Ошибка ввода пароля:", err)
			continue
		}

		var cmd *cobra.Command
		switch choice {
		case 1:
			cmd = lockBoxCli.NewRegisterCli(ctx)
		case 2:
			cmd = lockBoxCli.NewAuthCli(ctx)
		default:
			fmt.Println("Некорректный выбор, попробуйте снова.")
			continue
		}

		cmd.SetArgs([]string{"--username", username, "--password", password})
		if err := cmd.Execute(); err != nil {
			fmt.Println("Ошибка аутентификации:", err)
			continue
		}

		if lockBoxCli.IsAuthenticated() {
			return true
		}

		fmt.Println("Не удалось авторизоваться. Попробуйте снова.")
	}
}

func commandLoop(lockBoxCli *cli2.LockBoxCLI, ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nВыберите команду:")
		fmt.Println("1. Создать запись")
		fmt.Println("2. Удалить запись")
		fmt.Println("3. Посмотреть все")
		fmt.Println("4. Получить запись")
		fmt.Println("5. Обновить запись")
		fmt.Println("6. Выход")
		fmt.Print("Введите номер команды: ")

		var choice int
		_, err := fmt.Scanf("%d\n", &choice)
		if err != nil {
			fmt.Println("Ошибка ввода: ожидается число")
			reader.ReadString('\n')
			continue
		}

		if choice == 6 {
			fmt.Println("Завершение работы.")
			return
		}

		switch choice {
		case 1:
			createLockBoxFlow(lockBoxCli, ctx)
		case 2:
			deleteLockBoxFlow3(lockBoxCli, ctx)
		case 3:
			getAllLockBoxFlow2(lockBoxCli, ctx)
		case 4:
			getLockBoxFlow1(lockBoxCli, ctx)
		case 5:
			updateLockBoxFlow4(lockBoxCli, ctx)
		default:
			fmt.Println("Некорректный выбор, попробуйте снова.")
		}
	}
}

func createLockBoxFlow(lockBoxCli *cli2.LockBoxCLI, ctx context.Context) {
	for {
		fmt.Println("\nСоздание нового LockBox")

		fmt.Print("Название: ")
		var name string
		_, err := fmt.Scanln(&name)
		if err != nil || name == "" {
			fmt.Println("❌ Ошибка ввода Названия")
			continue
		}

		fmt.Print("URL: ")
		var url string
		_, _ = fmt.Scanln(&url)

		fmt.Print("Логин: ")
		var login string
		_, _ = fmt.Scanln(&login)

		fmt.Print("Пароль: ")
		var password string
		_, _ = fmt.Scanln(&password)

		fmt.Print("Описание (необязательно): ")
		var description string
		_, _ = fmt.Scanln(&description)

		cmd := lockBoxCli.CreateCommand(ctx)
		cmd.SetArgs([]string{
			"--name", name,
			"--url", url,
			"--login", login,
			"--password", password,
			"--description", description,
		})

		if err := cmd.Execute(); err != nil {
			fmt.Println("❌ Ошибка создания LockBox:", err)
			continue
		}

		break
	}
}
func getLockBoxFlow1(lockBoxCli *cli2.LockBoxCLI, ctx context.Context) {
	for {
		fmt.Println("\nполучение LockBox")

		fmt.Print("Название: ")
		var name string
		_, err := fmt.Scanln(&name)
		if err != nil || name == "" {
			fmt.Println("❌ Ошибка ввода Названия")
			continue
		}

		cmd := lockBoxCli.GetCommand(ctx)
		cmd.Flags().Set("name", name)

		if err := cmd.Execute(); err != nil {
			fmt.Println("❌ Ошибка получения LockBox:", err)
			continue
		}

		break
	}
}

func getAllLockBoxFlow2(lockBoxCli *cli2.LockBoxCLI, ctx context.Context) {
	for {
		fmt.Println("\nполучение LockBoxes")

		cmd := lockBoxCli.GetAllCommand(ctx)

		if err := cmd.Execute(); err != nil {
			fmt.Println("❌ Ошибка получения LockBoxes:", err)
			continue
		}

		break
	}
}
func deleteLockBoxFlow3(lockBoxCli *cli2.LockBoxCLI, ctx context.Context) {
	for {
		fmt.Println("\nУдаление LockBox")

		fmt.Print("Название: ")
		var name string
		_, err := fmt.Scanln(&name)
		if err != nil || name == "" {
			fmt.Println("❌ Ошибка ввода названия")
			continue
		}

		cmd := lockBoxCli.DeleteCommand(ctx)
		cmd.Flags().Set("name", name)

		if err := cmd.Execute(); err != nil {
			fmt.Println("❌ Ошибка удаления LockBox:", err)
			continue
		}

		break
	}
}
func updateLockBoxFlow4(lockBoxCli *cli2.LockBoxCLI, ctx context.Context) {
	for {
		fmt.Println("\nАпдейт LockBox")

		fmt.Print("Название: ")
		var name string
		_, err := fmt.Scanln(&name)
		if err != nil || name == "" {
			fmt.Println("❌ Ошибка ввода Названия")
			continue
		}

		fmt.Print("URL: ")
		var url string
		_, _ = fmt.Scanln(&url)

		fmt.Print("Логин: ")
		var login string
		_, _ = fmt.Scanln(&login)

		fmt.Print("Пароль: ")
		var password string
		_, _ = fmt.Scanln(&password)

		fmt.Print("Описание (необязательно): ")
		var description string
		_, _ = fmt.Scanln(&description)

		cmd := lockBoxCli.UpdateCommand(ctx)
		cmd.SetArgs([]string{
			"--name", name,
			"--url", url,
			"--login", login,
			"--password", password,
			"--description", description,
		})

		if err := cmd.Execute(); err != nil {
			fmt.Println("❌ Ошибка обновления LockBox:", err)
			continue
		}

		break
	}
}
