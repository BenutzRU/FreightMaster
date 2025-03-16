// gui/main.go
package main

import (
	"FreightMaster/backend/database" // Импортируем для типов
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io"
	"net/http"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Shipment struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	ClientID    uint   `json:"client_id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Cost        string `json:"cost"`
}

var token string
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
		resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			errorLabel.SetText("Failed to connect to server: " + err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Login failed with status:", resp.StatusCode)
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Println("Response body:", string(bodyBytes)) // Отладка тела ответа
			errorLabel.SetText("Invalid username or password")
			return
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Println("Error decoding response:", err)
			errorLabel.SetText("Error processing response")
			return
		}

		token, ok := result["token"].(string)
		if !ok {
			fmt.Println("Token not found in response:", result)
			errorLabel.SetText("Invalid server response")
			return
		}
		userRole, ok = result["role"].(string)
		if !ok {
			fmt.Println("Role not found in response:", result)
			errorLabel.SetText("Invalid server response")
			return
		}

		fmt.Println("Login successful, token:", token, "role:", userRole) // Полный токен
		w.Close()
		showMainWindow(a)
	})

	registerButton := widget.NewButton("Register", func() {
		request := AuthRequest{
			Username: usernameEntry.Text,
			Password: passwordEntry.Text,
		}

		body, _ := json.Marshal(request)
		resp, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error connecting to server:", err) // Отладка в терминале
			errorLabel.SetText("Failed to connect to server: " + err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			fmt.Println("Registration failed with status:", resp.StatusCode) // Отладка в терминале
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
	fmt.Println("Main window created")
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
	fmt.Println("Creating Shipments tab") // Отладка
	shipmentsList := widget.NewList(
		func() int {
			shipments, err := getShipments()
			if err != nil {
				fmt.Println("Error getting shipments:", err) // Отладка
				return 0
			}
			fmt.Println("Fetched shipments:", len(shipments)) // Отладка
			return len(shipments)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			shipments, err := getShipments()
			if err != nil {
				fmt.Println("Error getting shipments in list update:", err) // Отладка
				return
			}
			o.(*widget.Label).SetText(fmt.Sprintf("ID: %d, Desc: %s, Status: %s, Cost: %s", shipments[i].ID, shipments[i].Description, shipments[i].Status, shipments[i].Cost))
		},
	)
	fmt.Println("Shipments list created") // Отладка

	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetPlaceHolder("Description")

	statusEntry := widget.NewEntry()
	statusEntry.SetPlaceHolder("Status")

	costEntry := widget.NewEntry()
	costEntry.SetPlaceHolder("Cost")

	addButton := widget.NewButton("Add Shipment", func() {
		shipment := Shipment{
			Description: descriptionEntry.Text,
			Status:      statusEntry.Text,
			Cost:        costEntry.Text,
		}
		addShipment(shipment)
		shipmentsList.Refresh()
	})

	if userRole != "admin" {
		addButton.Disable()
	}

	content := container.NewVBox(
		widget.NewLabel("Shipments"),
		shipmentsList,
		descriptionEntry,
		statusEntry,
		costEntry,
		addButton,
	)
	fmt.Println("Shipments tab content created") // Отладка
	return content
}

func createUsersTab() fyne.CanvasObject {
	fmt.Println("Creating Users tab") // Отладка
	usersList := widget.NewList(
		func() int {
			users, err := getUsers()
			if err != nil {
				fmt.Println("Error getting users:", err) // Отладка
				return 0
			}
			fmt.Println("Fetched users:", len(users)) // Отладка
			return len(users)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			users, err := getUsers()
			if err != nil {
				fmt.Println("Error getting users in list update:", err) // Отладка
				return
			}
			o.(*widget.Label).SetText(fmt.Sprintf("ID: %d, Username: %s, Role: %s", users[i].ID, users[i].Username, users[i].Role))
		},
	)
	fmt.Println("Users list created") // Отладка

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	addButton := widget.NewButton("Add User", func() {
		user := database.User{
			Username: usernameEntry.Text,
			Password: passwordEntry.Text,
			Role:     "user",
		}
		addUser(user)
		usersList.Refresh()
	})

	if userRole != "admin" {
		addButton.Disable()
	}

	content := container.NewVBox(
		widget.NewLabel("Users"),
		usersList,
		usernameEntry,
		passwordEntry,
		addButton,
	)
	fmt.Println("Users tab content created") // Отладка
	return content
}

func getShipments() ([]Shipment, error) {
	fmt.Println("Sending GET request to /api/shipments with token:", token) // Полный токен
	if token == "" {
		fmt.Println("Token is empty!")
		return nil, fmt.Errorf("token is empty")
	}
	req, err := http.NewRequest("GET", "http://localhost:8080/api/shipments", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Unexpected status code:", resp.StatusCode, "Response body:", string(body)) // Отладка
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var shipments []Shipment
	if err := json.NewDecoder(resp.Body).Decode(&shipments); err != nil {
		fmt.Println("Error decoding response:", err)
		return nil, err
	}
	return shipments, nil
}

func getUsers() ([]database.User, error) {
	fmt.Println("Sending GET request to /api/users with token:", token) // Полный токен
	if token == "" {
		fmt.Println("Token is empty!")
		return nil, fmt.Errorf("token is empty")
	}
	req, err := http.NewRequest("GET", "http://localhost:8080/api/users", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Unexpected status code:", resp.StatusCode, "Response body:", string(body)) // Отладка
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var users []database.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		fmt.Println("Error decoding response:", err)
		return nil, err
	}
	return users, nil
}

func readBody(resp *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(buf.Bytes())) // Сбрасываем тело для повторного чтения
	return buf.String()
}

func addShipment(shipment Shipment) {
	fmt.Println("Adding shipment:", shipment) // Отладка
	body, _ := json.Marshal(shipment)
	req, err := http.NewRequest("POST", "http://localhost:8080/admin/shipments", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating POST request:", err) // Отладка
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err) // Отладка
		return
	}
	defer resp.Body.Close()
	fmt.Println("Add shipment response status:", resp.StatusCode) // Отладка
}

func addUser(user database.User) {
	fmt.Println("Adding user:", user.Username)
	body, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "http://localhost:8080/admin/users", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating POST request:", err) // Отладка
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err) // Отладка
		return
	}
	defer resp.Body.Close()
	fmt.Println("Add user response status:", resp.StatusCode) // Отладка
}
