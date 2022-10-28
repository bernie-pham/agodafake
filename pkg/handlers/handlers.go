package handlers

import (
	"net/http"

	"github.com/bernie-pham/agodafake/pkg/config"
	"github.com/bernie-pham/agodafake/pkg/model"
	"github.com/bernie-pham/agodafake/pkg/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func Layout(w http.ResponseWriter, req *http.Request) {
	// data := model.TodoPageData{
	// 	PageTitle: "My TODO list",
	// 	Todos: []model.Todo{
	// 		{Title: "Task 1", Done: false},
	// 		{Title: "Task 2", Done: true},
	// 		{Title: "Task 3", Done: true},
	// 	},
	// }
	render.RenderTemplate(w, "layout.html", &model.TemplateData{})

}

func (r *Repository) Home(w http.ResponseWriter, req *http.Request) {
	remoteIP := req.RemoteAddr
	r.App.Session.Put(req.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.html", &model.TemplateData{})
}
func (r *Repository) About(w http.ResponseWriter, req *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."
	remoteIP := r.App.Session.GetString(req.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(w, "about.page.html", &model.TemplateData{
		StringMap: stringMap,
	})
}
