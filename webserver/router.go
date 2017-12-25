package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/vitalyisaev2/buildgraph/vcs/gitlab"
)

const (
	gitlabEventHeader = "X-Gitlab-Event"
	gitlabPushHook    = "Push Hook"
)

var (
	defaultCtx = context.Background()
)

// newRouter builds new router instance
func newRouter(s Webserver) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/vcs/gitlab/events/push", s.GitlabPushEvent)
	return router
}

func (s *server) GitlabPushEvent(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get(gitlabEventHeader) != gitlabPushHook {
		http.Error(w, fmt.Sprintf("wrong %s header value", gitlabEventHeader), 400)
		return
	}

	if r.Body == nil {
		http.Error(w, "please send request body", 400)
		return
	}

	var event gitlab.PushEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = s.services.Storage.SavePushEvent(defaultCtx, &event)
	if err != nil {
		s.services.Logger.WithError(err).Error("failed to save event")
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
}
