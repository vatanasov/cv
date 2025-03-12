package extractor

import (
	"autobiography/internal/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

var gitHubResponse = `
{
  "data": {
    "viewer": {
      "pinnedItems": {
        "nodes": [
          {
            "name": "elixir-course-materials",
            "description": null,
            "url": "https://github.com/vatanasov/elixir-course-materials",
            "createdAt": "2018-07-01T09:18:38Z",
            "updatedAt": "2023-08-21T17:46:47Z",
            "primaryLanguage": {
              "name": "Elixir"
            }
          },
          {
            "name": "anki_cards_from_images",
            "description": null,
            "url": "https://github.com/vatanasov/anki_cards_from_images",
            "createdAt": "2025-03-11T08:17:35Z",
            "updatedAt": "2025-03-11T08:37:19Z",
            "primaryLanguage": {
              "name": "Go"
            }
          }
        ]
      }
    }
  }
}
`

func TestExtractFromGitHub(t *testing.T) {
	type args struct {
		token Token
	}

	tests := []struct {
		name           string
		args           args
		serverResponse string
		want           []models.GitHubRepo
		wantErr        bool
	}{
		{
			name:           "Success",
			args:           args{token: "test_token"},
			serverResponse: gitHubResponse,
			want: []models.GitHubRepo{
				{
					HtmlUrl:  "https://github.com/vatanasov/elixir-course-materials",
					Language: "Elixir",
					Name:     "elixir-course-materials",
				},
				{
					HtmlUrl:  "https://github.com/vatanasov/anki_cards_from_images",
					Language: "Go",
					Name:     "anki_cards_from_images",
				},
			},
			wantErr: false,
		},
		{
			name:           "Empty Success",
			args:           args{token: "test_token"},
			serverResponse: "{}",
			want:           []models.GitHubRepo{},
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				authorizationChan := make(chan string, 1)
				defer close(authorizationChan)
				server := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							fmt.Println("HERE")
							authorizationChan <- r.Header.Get("Authorization")
							w.Write([]byte(tt.serverResponse))
						},
					),
				)
				defer server.Close()

				r := &RepoApiClient{
					BaseUrl: server.URL,
				}
				got, err := r.ExtractFromGitHub(tt.args.token)

				authorizationHeader := readTimeout(authorizationChan)
				wantAuthorizationHeader := "Bearer " + string(tt.args.token)
				if authorizationHeader != wantAuthorizationHeader {
					t.Errorf(
						"ExtractFromGitHub() Authorization = %v, want %v",
						authorizationHeader,
						wantAuthorizationHeader,
					)
				}

				if (err != nil) != tt.wantErr {
					t.Errorf("ExtractFromGitHub() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ExtractFromGitHub() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func readTimeout(r <-chan string) string {
	select {
	case msg := <-r:
		return msg
	case <-time.After(1 * time.Second):
		panic("timeout reading from channel")
	}
}
