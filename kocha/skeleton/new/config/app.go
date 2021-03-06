package config

import (
	"github.com/naoina/kocha"
	"path/filepath"
	"runtime"
	"time"
)

var (
	AppName   = "{{.appName}}"
	Addr      = "0.0.0.0"
	Port      = 9100
	AppConfig = &kocha.AppConfig{
		AppPath:    rootPath,
		AppName:    AppName,
		RouteTable: kocha.InitRouteTable(kocha.RouteTable(Routes())),
		TemplateSet: kocha.TemplateSetFromPaths(map[string][]string{
			AppName: []string{
				filepath.Join(rootPath, "app", "views"),
			},
		}),

		// Session settings
		Session: kocha.SessionConfig{
			Name:  "{{.appName}}_session",
			Store: &kocha.SessionCookieStore{},

			// Expiration of session cookie, in seconds, from now.
			// Persistent if -1, For not specify, set 0.
			CookieExpires: time.Duration(90) * time.Hour * 24,

			// Expiration of session data, in seconds, from now.
			// Perssitent if -1, For not specify, set 0.
			SessionExpires: time.Duration(90) * time.Hour * 24,
			HttpOnly:       false,

			// AUTO-GENERATED Random keys. DO NOT EDIT.
			SecretKey: "{{.secretKey}}",
			SignedKey: "{{.signedKey}}",
		},

		MaxClientBodySize: 1024 * 1024 * 10, // 10MB
	}

	_, configFileName, _, _ = runtime.Caller(0)
	rootPath                = filepath.Dir(filepath.Join(configFileName, ".."))
)

func init() {
	config := kocha.Config(AppName)
	config.Set("AppName", AppName)
	config.Set("Addr", Addr)
	config.Set("Port", Port)
}
