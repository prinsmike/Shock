package conf

import (
	"flag"
	"fmt"
	"github.com/kless/goconfig/config"
	"os"
	"strings"
)

type idxOpts struct {
	unique   bool
	dropDups bool
	sparse   bool
}

// Setup conf variables
var (
	// Reload
	RELOAD = ""

	// Config File
	CONFIG_FILE = ""

	// Ports
	SITE_PORT = 0
	API_PORT  = 0

	// SSL
	SSL_ENABLED   = false
	SSL_KEY_FILE  = ""
	SSL_CERT_FILE = ""

	// Anonymous-Access-Control 
	ANON_WRITE      = false
	ANON_READ       = true
	ANON_CREATEUSER = false

	// Auth
	AUTH_TYPE               = "" //globus, oauth, basic
	GLOBUS_TOKEN_URL        = ""
	GLOBUS_PROFILE_URL      = ""
	OAUTH_REQUEST_TOKEN_URL = ""
	OAUTH_AUTH_TOKEN_URL    = ""
	OAUTH_ACCESS_TOKEN_URL  = ""

	// Admin
	ADMIN_EMAIL = ""
	SECRET_KEY  = ""

	// Directories
	DATA_PATH = ""
	SITE_PATH = ""
	LOGS_PATH = ""

	// Mongodb 
	MONGODB = ""

	// Node Indices
	NODE_IDXS map[string]idxOpts = nil
)

func init() {
	flag.StringVar(&CONFIG_FILE, "conf", "/usr/local/shock/conf/shock.cfg", "path to config file")
	flag.StringVar(&RELOAD, "reload", "", "path or url to shock data. WARNING this will drop all current data.")
	flag.Parse()
	c, err := config.ReadDefault(CONFIG_FILE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: error reading conf file: %v\n", err)
		os.Exit(1)
	}

	// Ports
	SITE_PORT, _ = c.Int("Ports", "site-port")
	API_PORT, _ = c.Int("Ports", "api-port")

	// SSL
	SSL_ENABLED, _ = c.Bool("SSL", "enable")
	if SSL_ENABLED {
		SSL_KEY_FILE, _ = c.String("SSL", "key")
		SSL_CERT_FILE, _ = c.String("SSL", "cert")
	}

	// Access-Control 
	ANON_WRITE, _ = c.Bool("Anonymous", "write")
	ANON_READ, _ = c.Bool("Anonymous", "read")
	ANON_CREATEUSER, _ = c.Bool("Anonymous", "create-user")

	// Auth
	AUTH_TYPE, _ = c.String("Auth", "type")
	switch AUTH_TYPE {
	case "globus":
		GLOBUS_TOKEN_URL, _ = c.String("Auth", "globus_token_url")
		GLOBUS_PROFILE_URL, _ = c.String("Auth", "globus_profile_url")
	case "oauth":
		OAUTH_REQUEST_TOKEN_URL, _ = c.String("Auth", "oauth_request_token_url")
		OAUTH_AUTH_TOKEN_URL, _ = c.String("Auth", "oauth_auth_token_url")
		OAUTH_ACCESS_TOKEN_URL, _ = c.String("Auth", "oauth_access_token_url")
	case "basic":
		// nothing yet
	}

	// Admin
	ADMIN_EMAIL, _ = c.String("Admin", "email")
	SECRET_KEY, _ = c.String("Admin", "secretkey")

	// Directories
	SITE_PATH, _ = c.String("Directories", "site")
	DATA_PATH, _ = c.String("Directories", "data")
	LOGS_PATH, _ = c.String("Directories", "logs")

	// Mongodb
	MONGODB, _ = c.String("Mongodb", "hosts")

	// parse Node-Indices
	NODE_IDXS = map[string]idxOpts{}
	nodeIdx, _ := c.Options("Node-Indices")
	for _, opt := range nodeIdx {
		val, _ := c.String("Node-Indices", opt)
		opts := idxOpts{}
		for _, parts := range strings.Split(val, ",") {
			p := strings.Split(parts, ":")
			if p[0] == "unique" {
				if p[1] == "true" {
					opts.unique = true
				} else {
					opts.unique = false
				}
			} else if p[0] == "dropDups" {
				if p[1] == "true" {
					opts.dropDups = true
				} else {
					opts.dropDups = false
				}
			} else if p[0] == "sparse" {
				if p[1] == "true" {
					opts.sparse = true
				} else {
					opts.sparse = false
				}
			}
		}
		NODE_IDXS[opt] = opts
	}

}

func Print() {
	fmt.Printf("##### Admin #####\nemail:\t%s\nsecretkey:\t%s\n\n", ADMIN_EMAIL, SECRET_KEY)
	fmt.Printf("####### Anonymous ######\nread:\t%t\nwrite:\t%t\ncreate-user:\t%t\n\n", ANON_READ, ANON_WRITE, ANON_CREATEUSER)
	if AUTH_TYPE == "basic" {
		fmt.Printf("##### Auth #####\ntype:\tbasic\n\n")
	} else if AUTH_TYPE == "globus" {
		fmt.Printf("##### Auth #####\ntype:\tglobus\ntoken_url:\t%s\nprofile_url:\t%s\n\n", GLOBUS_TOKEN_URL, GLOBUS_PROFILE_URL)
	}
	fmt.Printf("##### Directories #####\nsite:\t%s\ndata:\t%s\nlogs:\t%s\n\n", SITE_PATH, DATA_PATH, LOGS_PATH)
	if SSL_ENABLED {
		fmt.Printf("##### SSL #####\nenabled:\t%t\nkey:\t%s\ncert:\t%s\n\n", SSL_ENABLED, SSL_KEY_FILE, SSL_CERT_FILE)
	} else {
		fmt.Printf("##### SSL #####\nenabled:\t%t\n\n", SSL_ENABLED)
	}
	fmt.Printf("##### Mongodb #####\nhost(s):\t%s\n\n", MONGODB)
	fmt.Printf("##### Ports #####\nsite:\t%d\napi:\t%d\n\n", SITE_PORT, API_PORT)
}
