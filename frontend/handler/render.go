package handler

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type templateData struct {
	StringMap            map[string]string
	IntMap               map[string]int
	FloatMap             map[string]float32
	Data                 map[string]interface{}
	CSRFToken            string
	Flash                string
	Warning              string
	Error                string
	IsAuthenticated      int
	UserID               int
	API                  string
	FrontendWsAddr       string
	CSSVersion           string
	StripeSecretKey      string
	StripePublishableKey string
}

func formatCurrency(n int) string {
	f := float32(n) / float32(100)
	return fmt.Sprintf("$%.2f", f)
}

var functions = template.FuncMap{
	"formatCurrency": formatCurrency,
}

//go:embed templates
var templateFS embed.FS

func (server *Server) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData, partials ...string) error {
	var t *template.Template
	var err error
	templateToRender := fmt.Sprintf("templates/%s.page.gohtml", page)

	_, templateInMap := server.templateCache[templateToRender]

	if templateInMap {
		t = server.templateCache[templateToRender]
	} else {
		t, err = server.parseTemplate(partials, page, templateToRender)
		if err != nil {
			log.Error().Err(err).Msg("renderTemplate")
			return err
		}
	}

	if td == nil {
		td = &templateData{}
	}

	td = server.addDefaultData(td, r)

	err = t.Execute(w, td)
	if err != nil {
		log.Error().Err(err).Msg("renderTemplate")
		return err
	}

	return nil
}

func (server *Server) parseTemplate(partials []string, page, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error

	// build partials
	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("templates/%s.partial.gohtml", x)
		}
	}

	if len(partials) > 0 {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.gohtml", strings.Join(partials, ","), templateToRender)
	} else {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.gohtml", templateToRender)
	}
	if err != nil {
		log.Error().Err(err).Msg("parseTemplate")
		return nil, err
	}

	server.templateCache[templateToRender] = t
	return t, nil
}

func (server *Server) addDefaultData(td *templateData, r *http.Request) *templateData {
	td.API = server.config.MainServerAddr
	td.FrontendWsAddr = server.config.FrontendWsAddr
	td.StripeSecretKey = server.config.StripeSecret
	td.StripePublishableKey = server.config.StripeKey

	if server.Session.Exists(r.Context(), "userID") {
		td.IsAuthenticated = 1
		td.UserID = server.Session.GetInt(r.Context(), "userID")
	} else {
		td.IsAuthenticated = 0
		td.UserID = 0
	}

	return td
}
