package web

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Pages handles web pages
type Pages struct {
	appVersion string
	root       string
}

// TemplateData is the mapping type we use to pass data to the template engine.
type TemplateData map[string]interface{}

// CreatePageHandler creates a new Pages object.
func CreatePageHandler(appVersion string) *Pages {
	return &Pages{
		appVersion,
		TemplatePathPrefix("templates/layout.html"),
	}
}

func (p *Pages) showIndexPage(w http.ResponseWriter, r *http.Request) {
	p.showTemplate("templates/index.html", w, r, TemplateData{})
}

func (p *Pages) showScorePage(w http.ResponseWriter, r *http.Request) {
	p.showTemplate("templates/score.html", w, r, TemplateData{})
}

// AddRoutes adds routes to serve reporting status requests.
func (p *Pages) AddRoutes(router *mux.Router) {
	router.HandleFunc("/", p.showIndexPage).Methods("GET")
	router.HandleFunc("/score", p.showScorePage).Methods("GET")
	// router.HandleFunc("/as-json", rep.sendStatusReport).Methods("GET")
	// router.HandleFunc("/latest-image", rep.showLatestImagePage).Methods("GET")
	// router.HandleFunc("/worker-action/{worker-id}", rep.workerAction).Methods("POST")

	static := noDirListing(http.StripPrefix("/static/", http.FileServer(http.Dir(p.root+"static"))))
	router.PathPrefix("/static/").Handler(static).Methods("GET")
}

func noDirListing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (p *Pages) showTemplate(templfname string, w http.ResponseWriter, r *http.Request, templateData TemplateData) {
	tmpl := template.New("").Funcs(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
				log.Infof("dict[%q] = %q", key, values[i+1])
			}
			return dict, nil
		},
	})

	tmpl, err := tmpl.ParseFiles(
		p.root+"templates/layout.html",
		p.root+templfname)
	if err != nil {
		log.Errorf("Error parsing HTML template %s: %s", templfname, err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	usedData := TemplateData{
		"Version": p.appVersion,
		"Root":    p.root,
	}
	merge(usedData, templateData)

	err = tmpl.ExecuteTemplate(w, "layout", usedData)
	if err != nil {
		log.Errorf("Error executing HTML template %s: %s", templfname, err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}

// Merges 'two' into 'one'
func merge(one map[string]interface{}, two map[string]interface{}) {
	for key := range two {
		one[key] = two[key]
	}
}
