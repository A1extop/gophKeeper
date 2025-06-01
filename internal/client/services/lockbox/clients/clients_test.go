package clients

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gophKeeper/internal/client/services/lockbox/models"
	"gophKeeper/pkg/crypt"
)

// extractHostPort получает базовый URL и порт из адреса тестового сервера.
func extractHostPort(url string) (baseURL, port string) {
	trimmed := strings.TrimPrefix(url, "http://")
	host, port, err := net.SplitHostPort(trimmed)
	if err != nil {
		panic(err)
	}
	return "http://" + host, port
}

func TestCreate(t *testing.T) {
	expectedID := 123

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Ожидался метод POST, получили %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/lock_boxes/create") {
			t.Errorf("Неверный путь запроса: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": expectedID})
	}))
	defer ts.Close()

	baseURL, port := extractHostPort(ts.URL)
	svc := NewLockBoxService(baseURL, port)
	svc.(*lockBoxService).authToken = "dummy"

	input := &models.LockBoxInput{
		Name: "test",
		URL:  "example",
	}

	id, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Неожиданная ошибка: %v", err)
	}
	if id != expectedID {
		t.Errorf("Ожидался id %d, получили %d", expectedID, id)
	}
}

func TestGet(t *testing.T) {
	expectedLockBox := models.LockBox{
		Name: "test",
		URL:  "example",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Ожидался GET, получили %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/api/lock_boxes/") {
			t.Errorf("Неверный путь запроса: %s", r.URL.Path)
		}
		encryptor := crypt.New("superSecretKey19")
		encrypted := expectedLockBox
		dataEncrypted, err := crypt.EncryptLockBox(&encrypted, encryptor)
		if err != nil {
			t.Fatalf("Ошибка шифрования в тесте: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(*dataEncrypted)
	}))
	defer ts.Close()

	baseURL, port := extractHostPort(ts.URL)
	svc := NewLockBoxService(baseURL, port)
	svc.(*lockBoxService).authToken = "dummy"

	lb, err := svc.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("Неожиданная ошибка: %v", err)
	}
	if lb.Name != expectedLockBox.Name || lb.URL != expectedLockBox.URL {
		t.Errorf("Ожидался lockbox %+v, получили %+v", expectedLockBox, lb)
	}
}

func TestGetAll(t *testing.T) {
	expectedLockBoxes := []models.LockBox{
		{Name: "test1", URL: "example1"},
		{Name: "test2", URL: "example2"},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Ожидался GET, получили %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/lock_boxes/") {
			t.Errorf("Неверный путь запроса: %s", r.URL.Path)
		}
		encryptor := crypt.New("superSecretKey19")
		encryptedLockBoxes := make([]models.LockBox, len(expectedLockBoxes))
		copy(encryptedLockBoxes, expectedLockBoxes)
		var datesEncrypted []models.LockBox
		for i := range encryptedLockBoxes {
			dataEncrypted, err := crypt.EncryptLockBox(&encryptedLockBoxes[i], encryptor)
			if err != nil {

				t.Fatalf("Ошибка шифрования lockbox[%d]: %v", i, err)
			}
			datesEncrypted = append(datesEncrypted, *dataEncrypted)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(encryptedLockBoxes)
	}))
	defer ts.Close()

	baseURL, port := extractHostPort(ts.URL)
	svc := NewLockBoxService(baseURL, port)
	svc.(*lockBoxService).authToken = "dummy"

	lockBoxes, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Неожиданная ошибка: %v", err)
	}
	if len(*lockBoxes) != len(expectedLockBoxes) {
		t.Errorf("Ожидалось %d lockbox, получено %d", len(expectedLockBoxes), len(*lockBoxes))
	}
}

func TestUpdate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Ожидался PUT, получили %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/lock_boxes/") {
			t.Errorf("Неверный путь запроса: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	baseURL, port := extractHostPort(ts.URL)
	svc := NewLockBoxService(baseURL, port)
	svc.(*lockBoxService).authToken = "dummy"

	input := &models.LockBoxInput{
		Name: "test",
		URL:  "updated",
	}
	if err := svc.Update(context.Background(), input); err != nil {
		t.Fatalf("Неожиданная ошибка: %v", err)
	}
}

func TestDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Ожидался DELETE, получили %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/api/lock_boxes/") {
			t.Errorf("Неверный путь запроса: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	baseURL, port := extractHostPort(ts.URL)
	svc := NewLockBoxService(baseURL, port)
	svc.(*lockBoxService).authToken = "dummy"

	if err := svc.Delete(context.Background(), "test"); err != nil {
		t.Fatalf("Неожиданная ошибка: %v", err)
	}
}

func TestRegisterUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Ожидался POST, получили %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/users/") {
			t.Errorf("Неверный путь запроса: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	baseURL, port := extractHostPort(ts.URL)
	svc := NewLockBoxService(baseURL, port)

	if err := svc.RegisterUser(context.Background(), "user", "pass"); err != nil {
		t.Fatalf("Неожиданная ошибка: %v", err)
	}
}

func TestAuthUser(t *testing.T) {
	expectedToken := "testtoken"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Ожидался POST, получили %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/auth/login") {
			t.Errorf("Неверный путь запроса: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": expectedToken})
	}))
	defer ts.Close()

	baseURL, port := extractHostPort(ts.URL)
	svc := NewLockBoxService(baseURL, port)

	token, err := svc.AuthUser(context.Background(), "user", "pass")
	if err != nil {
		t.Fatalf("Неожиданная ошибка: %v", err)
	}
	if token != expectedToken {
		t.Errorf("Ожидался токен %s, получили %s", expectedToken, token)
	}
	if !svc.Authenticated() {
		t.Errorf("Ожидалось, что сервис будет аутентифицирован")
	}
}

func TestUpdateOrCreate(t *testing.T) {
	expectedID := 123
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Ожидался POST, получили %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/lock_boxes/create/update") {
			t.Errorf("Неверный путь запроса: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": expectedID})
	}))
	defer ts.Close()

	baseURL, port := extractHostPort(ts.URL)
	svc := NewLockBoxService(baseURL, port)
	svc.(*lockBoxService).authToken = "dummy"

	lockBox := &models.LockBox{
		Name: "test",
		URL:  "example",
	}
	if err := svc.UpdateOrCreate(context.Background(), lockBox); err != nil {
		t.Fatalf("Неожиданная ошибка: %v", err)
	}
}
