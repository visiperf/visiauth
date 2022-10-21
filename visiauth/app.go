package visiauth

import "encoding/json"

type App struct {
	id string
}

func NewApp(id string) *App {
	return &App{id}
}

func (a App) ID() string {
	return a.id
}

func (a App) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}{
		ID:   a.ID(),
		Type: "app",
	})
}
