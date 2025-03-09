package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gophKeeper/internal/client/services/lockbox/models"
	"gophKeeper/pkg/crypt"
	"io"
	"log"
	"net/http"
)

type LockBoxService interface {
	Create(ctx context.Context, data *models.LockBoxInput) (int, error)
	Get(ctx context.Context, name string) (*models.LockBox, error)
	GetAll(ctx context.Context) (*[]models.LockBox, error)
	Update(ctx context.Context, data *models.LockBoxInput) error
	Delete(ctx context.Context, name string) error
	RegisterUser(ctx context.Context, username, password string) error
	AuthUser(ctx context.Context, username, password string) (string, error)
	Authenticated() bool
	UpdateOrCreate(ctx context.Context, data *models.LockBox) error
}

var key string = "superSecretKey19"

type lockBoxService struct {
	baseURL   string
	port      string
	authToken string
	client    *http.Client
	encryptor crypt.Encryptor
}

func NewLockBoxService(baseURL string, port string) LockBoxService {
	return &lockBoxService{
		baseURL:   baseURL,
		port:      port,
		client:    &http.Client{},
		encryptor: crypt.New(key),
	}
}

func (s *lockBoxService) Create(ctx context.Context, data *models.LockBoxInput) (int, error) {
	dataEncrypt, err := crypt.EncryptStruct(data, s.encryptor)
	if err != nil {
		return 0, err
	}
	jsonData, err := json.Marshal(dataEncrypt)
	if err != nil {
		return 0, err
	}

	url := s.baseURL + ":" + s.port + "/api/lock_boxes/create"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.authToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("failed to create lockbox (code %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return response.ID, nil
}

func (s *lockBoxService) Get(ctx context.Context, name string) (*models.LockBox, error) {
	url := fmt.Sprintf("%s:%s/api/lock_boxes/%s", s.baseURL, s.port, name)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", s.authToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get lockbox (code %d): %s", resp.StatusCode, string(body))
	}

	var lockBox models.LockBox
	if err := json.NewDecoder(resp.Body).Decode(&lockBox); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	log.Println(lockBox)
	dataDecrypt, err := crypt.DecryptLockBox(&lockBox, s.encryptor)
	if err != nil {
		return nil, err
	}
	log.Println(dataDecrypt)
	return dataDecrypt, nil
}

func (s *lockBoxService) GetAll(ctx context.Context) (*[]models.LockBox, error) {
	url := fmt.Sprintf("%s:%s/api/lock_boxes/", s.baseURL, s.port)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", s.authToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get lockboxes (code %d): %s", resp.StatusCode, string(body))
	}

	var lockBoxes []models.LockBox
	if err := json.NewDecoder(resp.Body).Decode(&lockBoxes); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	datesDecrypt := make([]models.LockBox, len(lockBoxes))
	for i := range lockBoxes {
		dataDecrypt, err := crypt.DecryptLockBox(&lockBoxes[i], s.encryptor)
		if err != nil {
			return nil, err
		}
		datesDecrypt[i] = *dataDecrypt
	}

	return &datesDecrypt, nil
}

func (s *lockBoxService) Update(ctx context.Context, data *models.LockBoxInput) error {

	url := fmt.Sprintf("%s:%s/api/lock_boxes/", s.baseURL, s.port)
	dataEncrypt, err := crypt.EncryptStruct(data, s.encryptor)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(dataEncrypt)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.authToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update lockbox (code %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *lockBoxService) Delete(ctx context.Context, name string) error {
	url := fmt.Sprintf("%s:%s/api/lock_boxes/%s", s.baseURL, s.port, name)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", s.authToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete lockbox (code %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *lockBoxService) RegisterUser(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password are required")
	}

	data := map[string]string{"username": username, "password": password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := s.baseURL + ":" + s.port + "/api/users/"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	return fmt.Errorf("registration failed (code %d): %s", resp.StatusCode, string(body))
}

func (s *lockBoxService) AuthUser(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username and password are required")
	}

	data := map[string]string{"username": username, "password": password}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := s.baseURL + ":" + s.port + "/api/auth/login"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("authentication failed (code %d), and response body could not be read", resp.StatusCode)
		}
		return "", fmt.Errorf("authentication failed (code %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	s.authToken = response.Token
	return response.Token, nil
}

func (s *lockBoxService) Authenticated() bool {
	return s.authToken != ""
}

func (s *lockBoxService) UpdateOrCreate(ctx context.Context, data *models.LockBox) error {
	dataEncrypt, err := crypt.EncryptLockBox(data, s.encryptor)
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(dataEncrypt)
	if err != nil {
		return err
	}

	url := s.baseURL + ":" + s.port + "/api/lock_boxes/create/update"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.authToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//todo опять же не уверен насчёт статус кодов
	if (resp.StatusCode != http.StatusCreated) && (resp.StatusCode != http.StatusOK) {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create or update lockbox (code %d): %s", resp.StatusCode, string(body))
	}
	//todo тут может не приходить айдишник, проверить. как отреагирует
	var response struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}
