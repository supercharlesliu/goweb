package goweb

// AppConfig the config structure for a web app
type AppConfig struct {
	Domain   string // Host domain name this app server serves
	BasePath string // The root path of this app
}

type App struct {
	Config *AppConfig
}

func CreateApp(config *AppConfig) *App {
	return &App{
		Config: config,
	}
}

func (app *App) AddRoute(path string, controller *Controller, methodMapper map[string]uintptr) {

}
