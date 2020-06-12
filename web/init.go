package web

import (
	"context"
	"encoding/gob"
	"github.com/antonlindstrom/pgstore"
	"github.com/coreos/go-oidc"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/juju/loggo"
	"golang.org/x/oauth2"
	"html/template"
	"net/http"
	"phsite/models"
	"time"
)

var (
	ctx            context.Context
	logger         *loggo.Logger
	oauth2Config   oauth2.Config
	oauth2Verifier *oidc.IDTokenVerifier
	store          *pgstore.PGStore
	staticFiles    *packr.Box
	templates      *packr.Box
)

type templateVars interface {
	SetUser(u *models.User)
	SetNavbar(n *[]templateNavbarNode)
}

type templateAlert struct {
	Header string
	Text   string
}

type templateCommon struct {
	AlertError   templateAlert
	AlertSuccess templateAlert
	AlertWarn    templateAlert

	NavBar       *[]templateNavbarNode
	PageTitle    string
	User         *models.User
}

func (t *templateCommon) SetNavbar(n *[]templateNavbarNode) {
	t.NavBar = n
	return
}

func (t *templateCommon) SetUser(u *models.User) {
	t.User = u
	return
}

type templateNavbarNode struct {
	Text     string
	URL      string
	MatchStr string
	FAIcon   string

	Active   bool
	Disabled bool

	Children []*templateNavbarNode
}

func Close() {
	store.Close()
}

func Init(secretKey, providerURL, clientID, clientSecret, callbackURL string) error {
	newLogger := loggo.GetLogger("web")
	logger = &newLogger

	// load templates
	templates = packr.New("htmlTemplates", "./templates")

	// load templates
	staticFiles = packr.New("staticFiles", "./static")

	// Init Sessions
	gs, err := pgstore.NewPGStoreFromPool(models.GetDBConn(), []byte(secretKey))
	if err != nil {
		logger.Errorf("Could not create pgstore: %s", err.Error())
		return err
	}
	store = gs
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	// Register models for GOB
	gob.Register(models.User{})

	// Init OIDC
	ctx = context.Background()
	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		logger.Errorf("Could not create oidc provider: %s", err.Error())
		return err
	}
	oidcConfig := &oidc.Config{ClientID: clientID}
	oauth2Verifier = provider.Verifier(oidcConfig)

	// Configure OAuth2 client
	oauth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// Setup Router
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticFiles)))

	r.HandleFunc("/", HandleIndex).Methods("GET")
	r.HandleFunc("/login", HandleLogin).Methods("GET")
	r.HandleFunc("/logout", HandleLogout).Methods("GET")
	r.HandleFunc("/oauth/callback", HandleOauthCallback).Methods("GET")

	// Protected Pages
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(MiddlewareRequireAuth)
	protected.HandleFunc("/debug", HandleDebug).Methods("GET")
	protected.HandleFunc("/home", HandleHome).Methods("GET")
	protected.HandleFunc("/mastodon", HandleMastodon).Methods("GET")

	go func() {
		err := http.ListenAndServe(":8080", r)
		if err != nil {
			logger.Errorf("Could not start web server %s", err.Error())
		}
	}()
	return nil
}

func CompileTemplate(filename string) (*template.Template, error) {
	tHTML, err := templates.FindString(filename)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("main").Parse(tHTML)
	if err != nil {
		return nil, err
	}

	tHTML, err = templates.FindString("common.html")
	if err != nil {
		return nil, err
	}
	_, err = tmpl.New("common").Parse(tHTML)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// privates
func initSession(response http.ResponseWriter, request *http.Request, tmplVars templateVars) (*sessions.Session, error) {
	// Init Session
	us, err := store.Get(request, "session-key")
	if err != nil {
		return nil, err
	}

	// Get User
	val := us.Values["user"]
	var user models.User
	var ok bool
	if user, ok = val.(models.User); ok {
		tmplVars.SetUser(&user)
	}

	// Get Navbar
	tmplVars.SetNavbar(makeNavbar(request.URL.Path))

	return us, nil
}
