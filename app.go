package visiauth

type App struct {
	id string
}

func NewApp(id string) *App {
	return &App{id}
}

func (a App) ID() string {
	return a.id
}
