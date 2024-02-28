package oauth

import (
	"app/oauth/oauthutils"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type printURLFunctionType func(url string)

// TODO: change them to env
const (
	oAuthFileName             = "oauth.json"
	oAuthOfflineTokenFileName = "offline-token.json"
	oAuthActivationPort       = ":19941"
)

var (
	savedClient *http.Client = nil
)

func GetClient(printURLFunction printURLFunctionType, scope ...string) (*http.Client, error) {
	if savedClient != nil {
		return savedClient, nil
	}

	// we first fix the credentials.json
	err := oauthutils.FixCredentialsJSON(oAuthFileName, oAuthActivationPort)
	if err != nil {
		return nil, err
	}
	// then we proceed by reading the file
	bytes, err := os.ReadFile(oAuthFileName)
	if err != nil {
		return nil, err
	}
	// we create a google config struct
	config, err := google.ConfigFromJSON(bytes, scope...)
	if err != nil {
		return nil, err
	}
	// lastly, we create the google auth client
	savedClient, err := createClient(config, printURLFunction)
	if err != nil {
		return nil, err
	}

	return savedClient, nil
}

func createClient(config *oauth2.Config, printURLFunction printURLFunctionType) (*http.Client, error) {
	// we first try to create the token from the saved file
	token, err := getOfflineTokenFromFile(oAuthOfflineTokenFileName)
	if err != nil {
		// if it doesn't work, chances are this token is not valid
		// or doesn't exist; we try to get it interactively
		token, err = getOfflineTokenFromWeb(config, printURLFunction)
		if err != nil {
			return nil, err
		}
		// and then we save such token
		err = saveOfflineToken(token)
		if err != nil {
			return nil, err
		}
	}
	return config.Client(context.Background(), token), nil
}

func getOfflineTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// TODO: add confirmation that it's connecting again
func getOfflineTokenFromWeb(config *oauth2.Config, printURLFunction printURLFunctionType) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	if !strings.HasPrefix(authURL, "https://") {
		return nil, errors.New("invalid oauth credentials")
	}

	printURLFunction(authURL)

	var (
		token *oauth2.Token
		err   error
	)
	var (
		wg     sync.WaitGroup
		server *http.Server
	)
	server = &http.Server{Addr: oAuthActivationPort, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the code
		authCode := r.URL.Query().Get("code")

		// save the token
		token, err = config.Exchange(context.TODO(), authCode)
		if err != nil {
			w.Write([]byte("There was an error while authenticating. Check the terminal."))
			wg.Done()
			go server.Shutdown(context.TODO())
			return
		}

		// close the server
		w.Write([]byte("You can now close this tab."))
		wg.Done()
		go server.Shutdown(context.TODO())
	})}
	wg.Add(1)
	server.ListenAndServe()

	wg.Wait()
	if err != nil {
		return nil, err
	}
	return token, nil
}

func saveOfflineToken(token *oauth2.Token) error {
	tokenFile, err := os.OpenFile(oAuthOfflineTokenFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer tokenFile.Close()
	json.NewEncoder(tokenFile).Encode(token)

	return nil
}
