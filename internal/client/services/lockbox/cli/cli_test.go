package cli

import (
	"bytes"
	"context"
	"gophKeeper/internal/client/services/lockbox/usecase"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = old
	return buf.String()
}
func TestCreateCommand(t *testing.T) {
	mockUC := usecase.NewLockBoxUsecaseMock()
	cliObj := NewLockBoxCLI(mockUC)
	ctx := context.Background()
	cmd := cliObj.CreateCommand(ctx)

	// Устанавливаем флаги
	cmd.Flags().Set("name", "TestLock")
	cmd.Flags().Set("url", "http://test.com")
	cmd.Flags().Set("login", "testuser")
	cmd.Flags().Set("password", "testpass")
	cmd.Flags().Set("description", "test description")

	output := captureOutput(func() {
		cmd.Run(cmd, []string{})
	})
	if !strings.Contains(output, "✅ LockBox успешно создан с ID:") {
		t.Errorf("Ожидался успешный вывод, получено: %s", output)
	}
}

func TestDeleteCommand(t *testing.T) {
	mockUC := usecase.NewLockBoxUsecaseMock()
	cliObj := NewLockBoxCLI(mockUC)
	ctx := context.Background()
	cmd := cliObj.DeleteCommand(ctx)

	output := captureOutput(func() {
		cmd.Run(cmd, []string{})
	})
	if !strings.Contains(output, "❌ Ошибка: имя LockBox обязательно") {
		t.Errorf("Ожидалась ошибка для отсутствия имени, получено: %s", output)
	}

	// Тест: успешное удаление
	cmd.Flags().Set("name", "TestLock")
	output = captureOutput(func() {
		cmd.Run(cmd, []string{})
	})
	if !strings.Contains(output, "✅ Lockbox успешно удалён!") {
		t.Errorf("Ожидался вывод успешного удаления, получено: %s", output)
	}
}

func TestGetCommand(t *testing.T) {
	mockUC := usecase.NewLockBoxUsecaseMock()
	cliObj := NewLockBoxCLI(mockUC)
	ctx := context.Background()
	cmd := cliObj.GetCommand(ctx)

	output := captureOutput(func() {
		cmd.Run(cmd, []string{})
	})
	if !strings.Contains(output, "❌ Ошибка: не указан параметр name") {
		t.Errorf("Ожидалась ошибка для отсутствия параметра name, получено: %s", output)
	}

	cmd.Flags().Set("name", "TestLock")
	output = captureOutput(func() {
		cmd.Run(cmd, []string{})
	})
	if !strings.Contains(output, "✅ Lockbox найден!") {
		t.Errorf("Ожидался вывод успешного получения, получено: %s", output)
	}
}

func TestUpdateCommand(t *testing.T) {
	mockUC := usecase.NewLockBoxUsecaseMock()
	cliObj := NewLockBoxCLI(mockUC)
	ctx := context.Background()
	cmd := cliObj.UpdateCommand(ctx)

	cmd.Flags().Set("name", "TestLock")
	cmd.Flags().Set("url", "http://updated.com")
	cmd.Flags().Set("login", "updatedUser")
	cmd.Flags().Set("password", "updatedPass")
	cmd.Flags().Set("description", "updated description")

	output := captureOutput(func() {
		cmd.Run(cmd, []string{})
	})
	if !strings.Contains(output, "✅ Lockbox успешно обновлён!") {
		t.Errorf("Ожидался вывод успешного обновления, получено: %s", output)
	}
}

func TestGetAllCommand(t *testing.T) {
	mockUC := usecase.NewLockBoxUsecaseMock()
	cliObj := NewLockBoxCLI(mockUC)
	ctx := context.Background()
	cmd := cliObj.GetAllCommand(ctx)
	output := captureOutput(func() {
		cmd.Run(cmd, []string{})
	})
	if !strings.Contains(output, "Список Lockbox") {
		t.Errorf("Ожидался вывод списка lockbox, получено: %s", output)
	}
}

func TestNewRegisterCli(t *testing.T) {
	mockUC := usecase.NewLockBoxUsecaseMock()
	cliObj := NewLockBoxCLI(mockUC)
	ctx := context.Background()
	cmd := cliObj.NewRegisterCli(ctx)

	output := captureOutput(func() {
		cmd.Run(cmd, []string{})
	})
	if !strings.Contains(output, "❌ Ошибка: Укажите имя пользователя") {
		t.Errorf("Ожидалась ошибка отсутствия username, получено: %s", output)
	}

	// Тест: успешная регистрация
	cmd.Flags().Set("username", "testuser")
	cmd.Flags().Set("password", "testpass")
	output = captureOutput(func() {
		cmd.Run(cmd, []string{})
	})
	if !strings.Contains(output, "✅ Регистрация успешна!") {
		t.Errorf("Ожидался вывод успешной регистрации, получено: %s", output)
	}
}

func TestNewAuthCli(t *testing.T) {
	mockUC := usecase.NewLockBoxUsecaseMock()
	cliObj := NewLockBoxCLI(mockUC)
	ctx := context.Background()
	cmd := cliObj.NewAuthCli(ctx)

	output := captureOutput(func() {
		cmd.Run(cmd, []string{})
	})
	if !strings.Contains(output, "❌ Ошибка: Укажите имя пользователя") {
		t.Errorf("Ожидалась ошибка отсутствия username в auth, получено: %s", output)
	}

	cmd.Flags().Set("username", "testuser")
	cmd.Flags().Set("password", "testpass")
	output = captureOutput(func() {
		cmd.Run(cmd, []string{})
	})
	if !strings.Contains(output, "✅ Аутентификация успешна!") {
		t.Errorf("Ожидался вывод успешной аутентификации, получено: %s", output)
	}
}

func TestIsAuthenticated(t *testing.T) {
	mockUC := usecase.NewLockBoxUsecaseMock()
	cliObj := NewLockBoxCLI(mockUC)
	if !cliObj.IsAuthenticated() {
		t.Errorf("Ожидалось, что IsAuthenticated вернёт true")
	}
}
