package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Estructura para los datos del pull request de GitHub
type PullRequestPayload struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	PullRequest struct {
		Title   string `json:"PULLREQUEST"`
		HTMLURL string `json:"https://github.com/OldROOx/PULLREQUEST"`
		User    struct {
			Login string `json:"OldROOx"`
		} `json:"user"`
		Body string `json:"body"`
	} `json:"pull_request"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

// Estructura para el mensaje de Discord
type DiscordWebhook struct {
	Username  string         `json:"altf421_"`
	AvatarURL string         `json:"avatar_url,omitempty"`
	Content   string         `json:"content,omitempty"`
	Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Color       int    `json:"color"`
	Author      struct {
		Name string `json:"name"`
	} `json:"author"`
}

// URL del webhook de Discord (debes configurar esto)
var discordWebhookURL = os.Getenv("https://discord.com/api/webhooks/1344001053246488688/wCNdfl-uAwABP3R-xiQlOWGngpHtf-9YxhMumZxRVMjTuaxim1tKdtRRJHbAr59SP1GC")

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
		// Mostrar notificación en terminal
		fmt.Printf("\n🔔 Nuevo PR #%d en %s\n", payload.Number, payload.Repository.FullName)
		fmt.Printf("📝 %s\n", payload.PullRequest.Title)
		fmt.Printf("👤 %s\n", payload.PullRequest.User.Login)
		fmt.Printf("🔗 %s\n\n", payload.PullRequest.HTMLURL)

		// Enviar a Discord si la URL del webhook está configurada
		if discordWebhookURL != "" {
			err := sendToDiscord(payload)
			if err != nil {
				log.Printf("Error enviando a Discord: %v", err)
			} else {
				fmt.Println("✅ Notificación enviada a Discord")
			}
		} else {
			fmt.Println("⚠️ No se envió a Discord: DISCORD_WEBHOOK_URL no está configurada")
		}
	}

	fmt.Fprintf(w, "Webhook recibido correctamente")
}

func sendToDiscord(payload *PullRequestPayload) error {
	// Crear mensaje para Discord
	actionText := "creado"
	color := 5814783 // Azul para nuevo PR

	if payload.Action == "reopened" {
		actionText = "reabierto"
		color = 16750899 // Naranja para PR reabierto
	} else if payload.Action == "synchronize" {
		actionText = "actualizado"
		color = 5763719 // Verde para PR actualizado
	}

	message := DiscordWebhook{
		Username: "GitHub PR Bot",
		Embeds: []DiscordEmbed{
			{
				Title:       fmt.Sprintf("Pull Request #%d %s", payload.Number, actionText),
				Description: payload.PullRequest.Title,
				URL:         payload.PullRequest.HTMLURL,
				Color:       color,
				Author: struct {
					Name string `json:"name"`
				}{
					Name: payload.PullRequest.User.Login,
				},
			},
		},
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Enviar a Discord
	req, err := http.NewRequest("POST", discordWebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error de Discord: código %d", resp.StatusCode)
	}

	return nil
}

func main() {
	// Puerto para el servidor local
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Verificar URL de Discord
	if discordWebhookURL == "" {
		fmt.Println("⚠️ DISCORD_WEBHOOK_URL no está configurada. Las notificaciones solo aparecerán en la terminal.")
		fmt.Println("   Configúrala con: export DISCORD_WEBHOOK_URL='https://discord.com/api/webhooks/...'")
	}

	// Configurar rutas
	http.HandleFunc("/webhook", handleWebhook)

	fmt.Printf("🚀 Iniciando servidor API para notificaciones de GitHub en puerto %s\n", port)
	fmt.Printf("📌 URL del webhook: https://6961-189-150-56-23.ngrok-free.app\n", port)
	fmt.Println("🔧 Configura esta URL en la sección de webhooks de tu repositorio GitHub")

	// Iniciar servidor HTTP
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
