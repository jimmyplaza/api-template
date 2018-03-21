package app

// ClearWatch worker

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sprout/api/dbpool"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var (
	Version     = "0.1"
	PackageName = "api"
	LastUpdated = ""
	Authors     = "Jimmy Ko"

	// share in application
	App         *Context
	Timezone, _ = time.LoadLocation("Asia/Taipei")
	debug       bool
	port        = ""
	env         = ""
)

// Context
type Context struct {
	Timezone *time.Location
	DB       *dbpool.DBPool
	Port     string
	Debug    bool
	Env      string
}

// ContextInit for initialize
func ContextInit(rurl string, dburi []string, port, env string, debug bool) *Context {
	App = new(Context)
	// db
	App.DB = dbpool.NewDBPool(dburi)

	App.Port = port
	App.Timezone = Timezone
	App.Debug = debug
	App.Env = env
	return App
}

func init() {
	// Give the default value here
	flag.StringVar(&dbName, "dbname", "inbound_new", `database name to conect`)
	flag.StringVar(&dbHost, "dbhost", "localhost", `database host to conec`)
	flag.StringVar(&dbUser, "dbuser", "root", `database user for connection`)
	flag.StringVar(&port, "port", ":3000", `address for listen default is :3000`)
	flag.BoolVar(&debug, "debug", false, `Flag for DEBUG, Default is: false`)
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	if os.Getenv("DEBUG") != "" {
		debug = true
	}

	if os.Getenv("ENV") != "" {
		env = os.Getenv("ENV")
	}

	fmt.Println(PackageName, Version)
	fmt.Printf("database: %s@%s %s\n", dbUser, dbHost, dbName)

	if debug {
		fmt.Println("Running in DEBUG mode")
	}
}

func printhelp() {
	fmt.Println("Name:", PackageName, Version)
	fmt.Println("Usage:")
	flag.PrintDefaults()
}
func NewContext() *Context {

	log.Println(dbHost)

	dbs := []string{}

	if strings.Contains(dbHost, ",") {
		for _, host := range strings.Split(strings.TrimSpace(dbHost), ",") {
			// dburi := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=1m&parseTime=true&loc=%s", dbUser, dbPass, host, dbPort, dbName, "Asia%2FTaipei")
			dbURI := fmt.Sprintf(" dbname=%s host=%s user=%s sslmode=disable", dbName, host, dbUser)
			log.Println(dbURI)
			dbs = append(dbs, dbURI)
		}

	} else {
		// dburi := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=1m&parseTime=true&loc=%s", dbUser, dbPass, dbHost, dbPort, dbName, "Asia%2FTaipei")
		dbURI := fmt.Sprintf(" dbname=%s host=%s user=%s sslmode=disable", dbName, dbHost, dbUser)
		log.Println(dbURI)
		dbs = append(dbs, dbURI)
	}
	return ContextInit(redisurl, dbs, port, env, debug)
}
