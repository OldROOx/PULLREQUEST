package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Estructura para los datos del pull request
type PullRequestPayload struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	PullRequest struct {
		Title   string `json:"title"`
		HTMLURL string `json:"html_url"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
		Body string `json:"body"`
	} `json:"pull_request"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	payload := &PullRequestPayload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(payload)
	if err != nil {
		http.Error(w, "Error al decodificar el payload", http.StatusBadRequest)
		log.Printf("Error decodificando payload: %v", err)
		return
	}

	// Solo procesar eventos de pull request
	if payload.Action == "opened" || payload.Action == "reopened" || payload.Action == "synchronize" {
		// Mostrar notificación en terminal con colores ANSI
		fmt.Printf("\n\033[1;36m🔔 Nuevo PR #%d en %s\033[0m\n", payload.Number, payload.Repository.FullName)
		fmt.Printf("\033[1;33m📝 %s\033[0m\n", payload.PullRequest.Title)
		fmt.Printf("\033[1;32m👤 %s\033[0m\n", payload.PullRequest.User.Login)
		fmt.Printf("\033[1;34m🔗 %s\033[0m\n\n", payload.PullRequest.HTMLURL)
	}

	fmt.Fprintf(w, "Webhook recibido correctamente")
}

func main() {
	// Puerto para el servidor local (usando variable de entorno o puerto por defecto)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configurar rutas
	http.HandleFunc("/webhook", handleWebhook)

	fmt.Printf("🚀 Iniciando servidor API para notificaciones de GitHub en puerto %s\n", port)
	fmt.Printf("📌 URL del webhook: http://tu-ip-o-dominio:%s/webhook\n", port)
	fmt.Println("🔧 Configura esta URL en la sección de webhooks de tu repositorio GitHub")

	// Iniciar servidor HTTP
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
