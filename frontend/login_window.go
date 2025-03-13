package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"net/http"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	a := app.New()
	w := a.NewWindow("Авторизация")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Введите email")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	loginButton := widget.NewButton("Войти", func() {
		// Формируем запрос
		data := AuthRequest{
			Email:    emailEntry.Text,
			Password: passwordEntry.Text,
		}

		jsonData, _ := json.Marshal(data)
		resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body)) // Здесь можно парсить JSON и обрабатывать ответ

		// Если авторизация успешна - открыть главное окно
		if resp.StatusCode == http.StatusOK {
			w.Hide()
			mainWindow()
		}
	})

	w.SetContent(container.NewVBox(
		emailEntry,
		passwordEntry,
		loginButton,
	))

	w.ShowAndRun()
}

func mainWindow() {
	w := app.New().NewWindow("Главное окно")
	w.SetContent(widget.NewLabel("Вы вошли в систему!"))
	w.Show()
}
