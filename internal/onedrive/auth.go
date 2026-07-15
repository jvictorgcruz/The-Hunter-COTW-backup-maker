package onedrive

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"backup-maker/internal/browser"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

const clientID = "20ef2407-4157-48b6-af1b-e284d32446a5"
const defaultPort = "8999"
const redirectURL = "http://localhost:" + defaultPort + "/oauth/callback"

var oauthConfig = &oauth2.Config{
	ClientID:    clientID,
	Endpoint:    microsoft.AzureADEndpoint("common"),
	RedirectURL: redirectURL,
	Scopes:      []string{"Files.ReadWrite.AppFolder", "offline_access"},
}

func getTokenPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "backup-maker")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "onedrive_token.json"), nil
}

func loadToken() (*oauth2.Token, error) {
	path, err := getTokenPath()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	return token, err
}

func saveToken(token *oauth2.Token) error {
	path, err := getTokenPath()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(token)
}

func IsAuthenticated() bool {
	token, err := loadToken()
	if err != nil {
		return false
	}
	return token.Valid() || token.RefreshToken != ""
}

func Disconnect() error {
	path, err := getTokenPath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func GetClient(ctx context.Context) (*http.Client, error) {
	token, err := loadToken()
	if err != nil {
		token, err = authenticate(ctx)
		if err != nil {
			return nil, err
		}
	}

	tokenSource := oauthConfig.TokenSource(ctx, token)

	refreshedToken, err := tokenSource.Token()
	if err == nil && refreshedToken.AccessToken != token.AccessToken {
		_ = saveToken(refreshedToken)
	}

	return oauth2.NewClient(ctx, tokenSource), nil
}

func authenticate(ctx context.Context) (*oauth2.Token, error) {
	state := "cotw-backup-state"
	authURL := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	codeChan := make(chan string)
	errChan := make(chan error)

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		queryState := r.URL.Query().Get("state")
		if queryState != state {
			errChan <- fmt.Errorf("invalid OAuth state")
			w.Write([]byte("Falha na autenticação: Estado inválido (State mismatch)."))
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("authorization code not found in request")
			w.Write([]byte("Falha na autenticação: Código de autorização não encontrado."))
			return
		}

		w.Write([]byte("Autenticação bem-sucedida! Você já pode fechar esta aba e voltar para o aplicativo."))
		codeChan <- code
	})

	server := &http.Server{
		Addr:    ":" + defaultPort,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	err := browser.OpenBrowser(authURL)
	if err != nil {
		server.Shutdown(ctx)
		return nil, fmt.Errorf("failed to open browser: %v", err)
	}

	select {
	case code := <-codeChan:
		server.Shutdown(ctx)
		token, err := oauthConfig.Exchange(ctx, code)
		if err != nil {
			return nil, err
		}
		if err := saveToken(token); err != nil {
			return nil, err
		}
		return token, nil
	case err := <-errChan:
		server.Shutdown(ctx)
		return nil, err
	case <-time.After(5 * time.Minute):
		server.Shutdown(ctx)
		return nil, fmt.Errorf("authentication timed out")
	}
}
