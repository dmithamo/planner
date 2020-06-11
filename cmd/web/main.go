package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dmithamo/planner/pkg/form"
	ilog "github.com/dmithamo/planner/pkg/log" // custom logger
	"github.com/dmithamo/planner/pkg/mysql"
	"github.com/dmithamo/planner/pkg/projects"
	"github.com/golangcollege/sessions"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type application struct {
	infoLogger      *log.Logger
	errLogger       *log.Logger
	port            *string
	staticResServer http.Handler
	mux             *mux.Router
	projects        *projects.Projects
	templates       map[string]*template.Template
	templateData    templateData
	isDevEnv        bool
	session         *sessions.Session
}

type templateData struct {
	Project  *projects.Model
	Projects []*projects.Model
	Form     *form.Form
	Error    error
	FlashMsg interface{}
}

// make this global to both init() and main() to keep main() short
var app *application

// for use in both main && create handler
var initialForm = &form.Form{ValidationErrs: nil, Values: url.Values{}}

// init *injects dependencies into an instance of the app
func init() {
	port := flag.String("p", ":3001", "http address where server will run")
	dsn := flag.String("db", "", "data source name of the db to be used")
	recreateDB := flag.Bool("rdb", false, "drop db tables and recreate them")
	logToFile := flag.Bool("ltf", true, "toggle whether to persist logs to file")
	isDevEnv := flag.Bool("dev", true, "toggle whether app is running in developemnt mode")
	secret := flag.String("secret", "Iamthestonethatthebuilderrefused", "32-byte string used to sign sessions")
	flag.Parse()

	// instantiate loggers
	logservice := &ilog.Log{}
	if *logToFile {
		year, month, date := time.Now().Date()
		logFile := fmt.Sprintf("./log-%v-%v-%v.log", year, month, date)
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

		checkFatalErrorsHelper(err, "app init::logToFile config::fail")
		logservice = &ilog.Log{
			LogFile: f,
		}
		logservice.Initialize()
	}

	// Yes I know this is superflous
	logservice.InfoLogger.Println("::::::::::::::::::::::::::::::::")
	logservice.InfoLogger.Println("{re}starting application")
	logservice.InfoLogger.Println("::::::::::::::::::::::::::::::::")

	// instantiate db
	dbservice := &mysql.Mysql{}
	db, err := dbservice.OpenDB(fmt.Sprintf("%v?parseTime=true", *dsn))
	checkFatalErrorsHelper(err, "app init::db open::fail")
	logservice.InfoLogger.Println("app init::db open::success")
	dbservice.IDB = db
	// defer db.Close() <- This is not necessary. Or is it? TODO. After reading

	// recreate tables in db if need be
	if *recreateDB {
		err := dbservice.DropTables()
		checkFatalErrorsHelper(err, "app init::db drop tables::fail")
		logservice.InfoLogger.Println("app init::db drop tables::success")

		err = dbservice.CreateTables()
		checkFatalErrorsHelper(err, "app init::db create tables::fail")
		logservice.InfoLogger.Println("app init::db create tables::success")
	}

	// instantiate data: projects model
	projects := projects.Projects{
		IDB: db,
	}

	// instantiate template cache
	templateCache, err := buildTemplatesCache("./views/html")
	checkFatalErrorsHelper(err, "app init::templates build cache::fail")
	logservice.InfoLogger.Println("app init::templates build cache::success")

	// instantiate router
	router := mux.NewRouter()

	// instantiate static resources server
	staticDir := "./views/static"
	staticResServer := http.FileServer(http.Dir(staticDir))

	// instantiate session
	session := sessions.New([]byte(*secret))
	session.Lifetime = 30 * time.Minute
	session.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		panic(err)
	}

	// Avengers, Assemble! <assemble all the things>
	app = &application{
		port:            port,
		infoLogger:      logservice.InfoLogger,
		errLogger:       logservice.ErrorLogger,
		mux:             router,
		staticResServer: staticResServer,
		projects:        &projects,
		templateData:    templateData{Form: initialForm},
		templates:       templateCache,
		isDevEnv:        *isDevEnv,
		session:         session,
	}

	logservice.InfoLogger.Println("app init::success")
}

// main runs an instance of the app
func main() {
	secureRouter := app.mux.PathPrefix("").Subrouter() // needs auth
	secureRouter.Use(app.session.Enable)

	app.mux.HandleFunc("/auth", app.landingPage)
	secureRouter.HandleFunc("/", app.listProjects)
	secureRouter.HandleFunc("/projects/create", app.showCreateProjectForm).Methods("GET")
	secureRouter.HandleFunc("/projects/create", app.createproject).Methods("POST")
	secureRouter.HandleFunc("/projects/slug/{projectSlug}", app.viewProject)
	secureRouter.HandleFunc("/settings", app.settings)

	app.mux.NotFoundHandler = http.HandlerFunc(app.notFoundErr)

	app.mux.PathPrefix("/static/").Handler(http.StripPrefix("/static", app.staticResServer))

	standardMiddleware := alice.New(app.panicRecovery, app.requestLogger, app.secureHeaders)
	srv := &http.Server{
		Addr:     *app.port,
		ErrorLog: app.errLogger,
		Handler:  standardMiddleware.Then(app.mux),
	}

	app.infoLogger.Printf("app start::start server [127.0.0.1%v]::success", *app.port)
	app.infoLogger.Println("app start::ready for requests::success")
	app.errLogger.Fatal("app start::fail", srv.ListenAndServe())
}
