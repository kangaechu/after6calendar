package after6calendar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	// スクリプトのディレクトリを取得
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	tokFile := "token.json"
	tokFilePath := filepath.Join(dir, tokFile)
	tok, err := tokenFromFile(tokFilePath)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFilePath, tok)
	}
	return config.Client(ctx, tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatalf("Unable to write client secret file: %v", err)
	}
}

func authenticate(ctx context.Context) *http.Client {
	// スクリプトのディレクトリを取得
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile(filepath.Join(dir, "credentials.json"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	client := getClient(ctx, config)
	return client

}

func getAfter6Programs(client *http.Client) *calendar.Events {
	opt := option.WithHTTPClient(client)
	srv, err := calendar.NewService(context.Background(), opt)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	timeMin := time.Now().AddDate(0, 0, -14).Format(time.RFC3339)
	timeMax := time.Now().AddDate(0, 0, 7).Format(time.RFC3339)
	events, err := srv.Events.
		List("after6junction905954@gmail.com").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(timeMin).
		TimeMax(timeMax).
		OrderBy("startTime").
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	return events
}

func getAfter6Program(client *http.Client, start time.Time) *calendar.Event {
	opt := option.WithHTTPClient(client)
	srv, err := calendar.NewService(context.Background(), opt)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	timeMin := start.Format(time.RFC3339)
	timeMax := start.Add(1 * time.Hour).Format(time.RFC3339)
	events, err := srv.Events.
		List("after6junction905954@gmail.com").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(timeMin).
		TimeMax(timeMax).
		OrderBy("startTime").
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	for _, event := range events.Items {
		if event.Start.DateTime != "" { // 時間も指定されている
			return event
		}
	}
	return nil
}

func GetEventsJson() {
	ctx := context.Background()
	srv := authenticate(ctx)
	events := getAfter6Programs(srv)

	// convert to json
	jsonBytes, err := json.Marshal(events.Items)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}
	out := new(bytes.Buffer)
	err = json.Indent(out, jsonBytes, "", "    ")
	if err != nil {
		log.Fatal("cannot indent json", err)
	}

	// output to file
	file, err := os.Create(`after6.json`)
	if err != nil {
		log.Fatal("error when opening json", err)
	}
	defer file.Close()

	_, err = file.Write(out.Bytes())
	if err != nil {
		log.Fatal("cannot write json", err)
	}
}

func GetProgramSummary(start time.Time) *string {
	ctx := context.Background()
	srv := authenticate(ctx)
	event := getAfter6Program(srv, start)
	return &event.Summary
}
