package kocha

import (
	"html/template"
	"regexp"
)

func newTestAppConfig() *AppConfig {
	return &AppConfig{
		AppPath: "apppath/appname",
		AppName: "appname",
		TemplateSet: TemplateSet{
			"appname": map[string]*template.Template{
				"fixture_root_test_ctrl.html":   template.Must(template.New("tmpl1").Parse(`tmpl1`)),
				"fixture_user_test_ctrl.html":   template.Must(template.New("tmpl2").Parse(`tmpl2-{{.id}}`)),
				"fixture_date_test_ctrl.html":   template.Must(template.New("tmpl3").Parse(`tmpl3-{{.name}}-{{.year}}-{{.month}}-{{.day}}`)),
				"fixture_error_test_ctrl.html":  template.Must(template.New("tmpl4").Parse(`tmpl4`)),
				"fixture_json_test_ctrl.json":   template.Must(template.New("tmpl5").Parse(`{"tmpl5":"json"}`)),
				"fixture_teapot_test_ctrl.html": template.Must(template.New("tmpl6").Parse(`teapot`)),
			},
		},
		RouteTable: RouteTable{
			{
				Name:       "root",
				Path:       "/",
				Controller: FixtureRootTestCtrl{},
				MethodTypes: map[string]MethodArgs{
					"Get": MethodArgs{},
				},
				RegexpPath: regexp.MustCompile(`^/$`),
			},
			{
				Name:       "user",
				Path:       "/user/:id",
				Controller: FixtureUserTestCtrl{},
				MethodTypes: map[string]MethodArgs{
					"Get": MethodArgs{
						"id": "int",
					},
				},
				RegexpPath: regexp.MustCompile(`^/user/(?P<id>\d+)$`),
			},
			{
				Name:       "date",
				Path:       "/:year/:month/:day/user/:name",
				Controller: FixtureDateTestCtrl{},
				MethodTypes: map[string]MethodArgs{
					"Get": MethodArgs{
						"year":  "int",
						"month": "int",
						"day":   "int",
						"name":  "string",
					},
				},
				RegexpPath: regexp.MustCompile(`^/(?P<year>\d+)/(?P<month>\d+)/(?P<day>\d+)/user/(?P<name>[\w-]+)$`),
			},
			{
				Name:       "error",
				Path:       "/error",
				Controller: FixtureErrorTestCtrl{},
				MethodTypes: map[string]MethodArgs{
					"Get": MethodArgs{},
				},
				RegexpPath: regexp.MustCompile(`^/error$`),
			},
			{
				Name:       "json",
				Path:       "/json",
				Controller: FixtureJsonTestCtrl{},
				MethodTypes: map[string]MethodArgs{
					"Get": MethodArgs{},
				},
				RegexpPath: regexp.MustCompile(`^/json$`),
			},
			{
				Name:       "teapot",
				Path:       "/teapot",
				Controller: FixtureTeapotTestCtrl{},
				MethodTypes: map[string]MethodArgs{
					"Get": MethodArgs{},
				},
				RegexpPath: regexp.MustCompile(`^/teapot$`),
			},
			{
				Name:       "static",
				Path:       "/static/*path",
				Controller: StaticServe{},
				MethodTypes: map[string]MethodArgs{
					"Get": MethodArgs{
						"path": "url.URL",
					},
				},
				RegexpPath: regexp.MustCompile(`^/static/(?P<path>[\w-./]+)$`),
			},
		},
		Middlewares: append(DefaultMiddlewares, []Middleware{}...),
		Session: SessionConfig{
			Name:      "test_session",
			Store:     &SessionCookieStore{},
			SecretKey: "abcdefghijklmnopqrstuvwxyzABCDEF",
			SignedKey: "abcdefghijklmn",
		},
	}
}
