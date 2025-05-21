package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
)

type APIHandler struct {
	apiBaseURL string
	store      *sessions.CookieStore
	templates  map[string]*template.Template
}

func NewAPIHandler(apiBaseURL string) *APIHandler {
	templates := make(map[string]*template.Template)
	basePath := "./web/templates/base.html"

	templates["login"] = template.Must(template.ParseFiles(
		basePath,
		"./web/templates/login.html",
	))

	templates["register"] = template.Must(template.ParseFiles(
		basePath,
		filepath.Join("web", "templates", "register.html"),
	))

	templates["expressions"] = template.Must(template.ParseFiles(
		basePath,
		filepath.Join("web", "templates", "expressions.html"),
	))

	templates["calculate"] = template.Must(template.ParseFiles(
		basePath,
		filepath.Join("web", "templates", "calculate.html"),
	))

	return &APIHandler{
		apiBaseURL: apiBaseURL,
		store:      sessions.NewCookieStore([]byte("your-secret-key")),
		templates:  templates,
	}
}

func (h *APIHandler) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := h.templates[name]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err := tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *APIHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *APIHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		password := r.FormValue("password")

		reqBody, _ := json.Marshal(map[string]string{
			"login":    login,
			"password": password,
		})
		fmt.Println(reqBody)
		resp, err := http.Post(h.apiBaseURL+"/api/v1/login", "application/json", bytes.NewBuffer(reqBody))
		fmt.Println(resp, err)
		if err != nil {
			http.Error(w, "API connection error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		var authResp struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
			http.Error(w, "Failed to parse token", http.StatusInternalServerError)
			return
		}

		session, _ := h.store.Get(r, "session")
		session.Values["token"] = "Bearer " + authResp.Token // Добавляем префикс "Bearer "
		session.Save(r, w)

		http.Redirect(w, r, "/expressions", http.StatusSeeOther)
		return
	}

	h.renderTemplate(w, "login", nil)
}

func (h *APIHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		password := r.FormValue("password")

		reqBody, _ := json.Marshal(map[string]string{
			"login":    login,
			"password": password,
		})

		resp, err := http.Post(h.apiBaseURL+"/api/v1/register", "application/json", bytes.NewBuffer(reqBody))
		fmt.Println(resp, err)
		if err != nil {
			http.Error(w, "API connection error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Registration failed", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	h.renderTemplate(w, "register", nil)
}

func (h *APIHandler) ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session")
	token, ok := session.Values["token"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	req, _ := http.NewRequest("GET", h.apiBaseURL+"/api/v1/expressions", nil)
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "API connection error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get expressions", http.StatusInternalServerError)
		return
	}

	var response struct {
		List []struct {
			ID     string      `json:"id"`
			Input  string      `json:"input"`
			Status string      `json:"status"`
			Result interface{} `json:"result"`
		} `json:"list"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to parse data", http.StatusInternalServerError)
		return
	}

	h.renderTemplate(w, "expressions", response.List)
}

func (h *APIHandler) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session")
	token, ok := session.Values["token"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		expression := r.FormValue("expression")

		reqBody, _ := json.Marshal(map[string]string{
			"expression": expression,
		})

		req, _ := http.NewRequest("POST", h.apiBaseURL+"/api/v1/calculate", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "API connection error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			http.Error(w, "Calculation failed: "+string(body), resp.StatusCode)
			return
		}

		http.Redirect(w, r, "/expressions", http.StatusSeeOther)
		return
	}

	h.renderTemplate(w, "calculate", nil)
}
