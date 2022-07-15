package main

import (
	"net/http"
	"regexp"
	"encoding/json"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "log"
    "strings"
)

var (
	listPluginRe = regexp.MustCompile(`^\/api\/v1\/wordpress\/plugins[\/]*$`)
	getPluginRe = regexp.MustCompile(`^\/api\/v1\/wordpress\/plugins\/(\w+)$`)
	createPluginRe = regexp.MustCompile(`^\/api\/v1\/wordpress\/plugins[\/]*$`)
	editPluginRe = regexp.MustCompile(`^\/api\/v1\/wordpress\/plugins[\/]*$`)
)

// plugin represents the required REST resource
type plugin struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Version string `json:"version"`
	Author string `json:"author"`
	Description string `json:"description"`
}

// REFERENCE #3: https://stackoverflow.com/questions/69401999/golang-proper-way-to-send-json-response-with-status
type GetAllResponseBody struct {
    StatusCode int  `json:"http_status_code"`
    Data []plugin `json:"plugins"`
}

type GetOneResponseBody struct {
    StatusCode int  `json:"http_status_code"`
    Data plugin `json:"plugin"`
}

type OtherResponseBody struct {
    Status string `json:"status"`
    StatusCode int  `json:"http_status_code"`
}

// REFERENCE #1: https://golang.cafe/blog/golang-rest-api-example.html
// if you want to use in-memory data store, initial API version was using this (code has been ported to MySQL)
// in the above reference, mutex was used to guard access to the following map (for the purpose of concurrent access) 
//
// type datastore struct {
// 		m map[string]plugin
// 		*sync.RWMutex
// }

type pluginHandler struct {}

// GET <base_url>/plugins – Fetch entries from database
func (h *pluginHandler) List(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	results, err := db.Query("SELECT * FROM onedotcom_plugins")
    if err != nil {
        log.Fatal(err.Error())
    }
	plugins := make([]plugin, 0)
    for results.Next() {
        var p plugin
        err = results.Scan(&p.ID, &p.Name, &p.Version, &p.Author, &p.Description)
        if err != nil {
            log.Fatal(err.Error())
        }
        plugins = append(plugins, p)
    }

	jsonBytes, err := json.Marshal(GetAllResponseBody{StatusCode: http.StatusOK, Data: plugins})
	if err != nil {
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// EXTRA: INDIVIDUAL RECORDS RETRIEVAL BY ID FUNCTIONALITY ALSO ADDED TO API
// GET Handler: retrieve plugin by ID: parse the PluginID from request.URL.Path
func (h *pluginHandler) Get(w http.ResponseWriter, r *http.Request) {
	matches := getPluginRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w, r)
		return
	}
	db := dbConn()
	results, err := db.Query("SELECT * FROM onedotcom_plugins WHERE id = ?", matches[1])
    if err != nil {
        log.Fatal(err.Error())
    }

    var p plugin
    noResultsFound := true
	for results.Next() {
        err = results.Scan(&p.ID, &p.Name, &p.Version, &p.Author, &p.Description)
        if err != nil {
            log.Fatal(err.Error())
        }
        noResultsFound = false
    }
	
	if noResultsFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Plugin not found"))
		return
	}

	jsonBytes, err := json.Marshal(GetOneResponseBody{StatusCode: http.StatusOK, Data: p})
	if err != nil {
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// POST <base_url>/plugins – Insert a new entry in database
func (h *pluginHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p plugin
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		internalServerError(w, r)
		return
	}

	db := dbConn()
	newRecordQuery, err := db.Prepare("INSERT INTO onedotcom_plugins (id, name, version, author, description) VALUES(?,?,?,?,?)")
    if err != nil {
        log.Fatal(err.Error())
    }
    newRecordQuery.Exec(p.ID, p.Name, p.Version, p.Author, p.Description)

	jsonBytes, err := json.Marshal(OtherResponseBody{Status: "success", StatusCode: http.StatusCreated})
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

// PATCH <base_url>/plugins – Update version field in database for CDN plugin
func (h *pluginHandler) Update(w http.ResponseWriter, r *http.Request) {
	var p plugin
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Fatal(err.Error())
		internalServerError(w, r)
		return
	}

	if p.Version == "" || p.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing ID or Version information."))
		return
	}

	db := dbConn()
	newRecordQuery, err := db.Prepare("UPDATE onedotcom_plugins SET version = ? WHERE id = ?")
    if err != nil {
        log.Fatal(err.Error())
    }
    newRecordQuery.Exec(p.Version, p.ID)

	jsonBytes, err := json.Marshal(OtherResponseBody{Status: "success", StatusCode: http.StatusAccepted})
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonBytes)	
}

func (h *pluginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// REQUIREMENT 8. API should be blocked for IP address “172.16.0.2”
	ipInfo := strings.Split(r.RemoteAddr, ":")
	for _, ip := range ipInfo {
    	if ip == "172.16.0.2" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Access not allowed!"))
			return    		
    	}
	}

	// REQUIREMENT 7. API should reject requests if request header “Content-Type: application/json” is not present
	// REFERENCE #4: https://golangbyexample.com/headers-http-request-golang
	requestHeadersContentType := r.Header.Values("Content-type")
	if len(requestHeadersContentType) == 0 {
		fmt.Println("No content-type found")
		notAllowed(w, r)
		return
	}
	if requestHeadersContentType[0] != "application/json" {
		fmt.Println("No appropriate content-type found")
		notAllowed(w, r)
		return
	}

	// all requests are going to be routed here
	w.Header().Set("Content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && listPluginRe.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodGet && getPluginRe.MatchString(r.URL.Path):
		h.Get(w, r)
		return
	case r.Method == http.MethodPost && createPluginRe.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	case r.Method == http.MethodPatch && editPluginRe.MatchString(r.URL.Path):
		h.Update(w, r)
		return
	default:
		// REQUIREMENT 6. API should allow only GET, POST, PATCH operations
		notAllowed(w, r)
		return
	}
}

func notAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Access not allowed!"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Route not found!"))
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error"))
}

// REFERENCE #2: https://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:1234567890@tcp(127.0.0.1:3306)/interview_tasks")
    if err != nil {
        log.Fatal(err.Error())
    }
    return db
}

func main() {
	port := ":4545"
	mux := http.NewServeMux()
	pluginH := &pluginHandler{}
	mux.Handle("/api/v1/wordpress/plugins", pluginH)
	mux.Handle("/api/v1/wordpress/plugins/", pluginH)
	fmt.Println("Server started listening on http://localhost" + port)
	http.ListenAndServe(port, mux)
}