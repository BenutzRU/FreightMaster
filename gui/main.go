package main

import (
	"FreightMaster/backend/database"
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io"
	"net/http"
	"runtime/debug"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Shipment struct {
	ID               uint       `json:"id"`
	UserID           uint       `json:"user_id"`
	ClientID         uint       `json:"client_id"`
	Description      string     `json:"description"`
	Status           string     `json:"status"`
	Cost             float64    `json:"cost"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
	DeliveryMethodID uint       `json:"delivery_method_id"`
	BranchID         uint       `json:"branch_id"`
}

var username string
var password string
var userRole string

func main() {
	a := app.New()
	showLoginWindow(a)
	a.Run()
}

func showLoginWindow(a fyne.App) {
	w := a.NewWindow("FreightMaster Login")
	w.Resize(fyne.NewSize(400, 300))

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	errorLabel := widget.NewLabel("")
	errorLabel.TextStyle = fyne.TextStyle{Italic: true}

	loginButton := widget.NewButton("Login", func() {
		request := AuthRequest{
			Username: usernameEntry.Text,
			Password: passwordEntry.Text,
		}

		body, _ := json.Marshal(request)
		req, err := http.NewRequest("POST", "http://localhost:8080/login", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error creating login request:", err)
			errorLabel.SetText("Failed to create request: " + err.Error())
			return
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			errorLabel.SetText("Failed to connect to server: " + err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Println("Login failed with status:", resp.StatusCode, "Body:", string(bodyBytes))
			errorLabel.SetText("Invalid username or password")
			return
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Println("Error decoding response:", err)
			errorLabel.SetText("Error processing response")
			return
		}

		role, ok := result["role"].(string)
		if !ok {
			fmt.Println("Role not found in response:", result)
			errorLabel.SetText("Invalid server response")
			return
		}

		username = usernameEntry.Text
		password = passwordEntry.Text
		userRole = role
		fmt.Println("Login successful, username:", username, "role:", userRole)
		w.Close()
		showMainWindow(a)
	})

	registerButton := widget.NewButton("Register", func() {
		request := AuthRequest{
			Username: usernameEntry.Text,
			Password: passwordEntry.Text,
		}

		body, _ := json.Marshal(request)
		req, err := http.NewRequest("POST", "http://localhost:8080/register", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error creating register request:", err)
			errorLabel.SetText("Failed to create request: " + err.Error())
			return
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			errorLabel.SetText("Failed to connect to server: " + err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Println("Registration failed with status:", resp.StatusCode, "Body:", string(bodyBytes))
			errorLabel.SetText("Failed to register")
			return
		}

		errorLabel.SetText("Registered successfully. Please login.")
	})

	content := container.NewVBox(
		widget.NewLabel("FreightMaster Login"),
		usernameEntry,
		passwordEntry,
		errorLabel,
		loginButton,
		registerButton,
	)

	w.SetContent(content)
	w.Show()
}

func showMainWindow(a fyne.App) {
	fmt.Println("Entering showMainWindow, role:", userRole)
	w := a.NewWindow("FreightMaster")
	w.Resize(fyne.NewSize(800, 600))

	tabs := container.NewAppTabs(
		container.NewTabItem("Shipments", createShipmentsTab()),
		container.NewTabItem("Users", createUsersTab()),
	)
	fmt.Println("Tabs created")

	w.SetContent(tabs)
	fmt.Println("Content set")

	w.Show()
	fmt.Println("Main window shown")
}

func createShipmentsTab() fyne.CanvasObject {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic in createShipmentsTab:", r)
			debug.PrintStack()
		}
	}()
	fmt.Println("Creating Shipments tab")
	shipments, err := getShipments()
	if err != nil {
		fmt.Println("Error getting shipments:", err)
		return widget.NewLabel("Error: " + err.Error())
	}
	fmt.Println("Shipments fetched successfully:", len(shipments))
	label := widget.NewLabel(fmt.Sprintf("Shipments: %d", len(shipments)))
	fmt.Println("Label created with text:", label.Text)
	return label
}

func createUsersTab() fyne.CanvasObject {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic in createUsersTab:", r)
			debug.PrintStack()
		}
	}()
	fmt.Println("Creating Users tab")

	// Создаём канал для получения данных
	usersChan := make(chan []database.User)
	errChan := make(chan error)

	// Запускаем запрос в отдельной горутине
	go func() {
		users, err := getUsers()
		if err != nil {
			errChan <- err
			return
		}
		usersChan <- users
	}()

	// Создаём временный лейбл
	label := widget.NewLabel("Loading users...")

	// Ожидаем результат
	go func() {
		select {
		case users := <-usersChan:
			fmt.Println("Users fetched successfully:", len(users))
			label.SetText(fmt.Sprintf("Users: %d", len(users)))
		case err := <-errChan:
			fmt.Println("Error getting users:", err)
			label.SetText("Error: " + err.Error())
		}
	}()

	return label
}

// gui/main.go (обновлённая функция getShipments)
func getShipments() ([]Shipment, error) {
	fmt.Println("Sending GET request to /api/shipments with username:", username, "password:", password)
	req, err := http.NewRequest("GET", "http://localhost:8080/api/shipments", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		debug.PrintStack()
		return nil, err
	}
	req.Header.Set("X-Username", username)
	req.Header.Set("X-Password", password)
	fmt.Println("Request headers set:", req.Header)

	fmt.Println("About to execute httpClient.Do...")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		debug.PrintStack()
		return nil, err
	}
	fmt.Println("httpClient.Do executed successfully")

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic during response handling:", r)
			debug.PrintStack()
		}
		if resp != nil {
			resp.Body.Close()
		}
	}()

	fmt.Println("Received response with status:", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		debug.PrintStack()
		return nil, err
	}
	fmt.Println("Response body:", string(body))
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode, "Response body:", string(body))
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var shipments []Shipment
	fmt.Println("Attempting to decode JSON...")
	if err := json.NewDecoder(resp.Body).Decode(&shipments); err != nil {
		fmt.Println("Error decoding response:", err)
		debug.PrintStack()
		return nil, err
	}
	fmt.Println("Successfully fetched shipments:", len(shipments))
	return shipments, nil
}

func getUsers() ([]database.User, error) {
	fmt.Println("Sending GET request to /api/users with username:", username, "password:", password)
	req, err := http.NewRequest("GET", "http://localhost:8080/api/users", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		debug.PrintStack()
		return nil, err
	}
	req.Header.Set("X-Username", username)
	req.Header.Set("X-Password", password)
	fmt.Println("Request headers set:", req.Header)

	fmt.Println("About to execute httpClient.Do...")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		debug.PrintStack()
		return nil, err
	}
	fmt.Println("httpClient.Do executed successfully")

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic during response handling:", r)
			debug.PrintStack()
		}
		if resp != nil {
			resp.Body.Close()
		}
	}()

	fmt.Println("Received response with status:", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		debug.PrintStack()
		return nil, err
	}
	fmt.Println("Response body:", string(body))
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode, "Response body:", string(body))
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var users []database.User
	fmt.Println("Attempting to decode JSON...")
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		fmt.Println("Error decoding response:", err)
		debug.PrintStack()
		return nil, err
	}
	fmt.Println("Successfully fetched users:", len(users))
	return users, nil
}

func addShipment(shipment Shipment) error {
	fmt.Println("Adding shipment:", shipment)
	body, err := json.Marshal(shipment)
	if err != nil {
		fmt.Println("Error marshaling shipment:", err)
		return err
	}
	req, err := http.NewRequest("POST", "http://localhost:8080/admin/shipments", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return err
	}
	req.Header.Set("X-Username", username)
	req.Header.Set("X-Password", password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return err
	}
	defer resp.Body.Close()
	fmt.Println("Add shipment response status:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Unexpected status code:", resp.StatusCode, "Response body:", string(body))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func addUser(user database.User) error {
	fmt.Println("Adding user:", user.Username)
	body, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error marshaling user:", err)
		return err
	}
	req, err := http.NewRequest("POST", "http://localhost:8080/admin/users", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return err
	}
	req.Header.Set("X-Username", username)
	req.Header.Set("X-Password", password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return err
	}
	defer resp.Body.Close()
	fmt.Println("Add user response status:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Unexpected status code:", resp.StatusCode, "Response body:", string(body))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
