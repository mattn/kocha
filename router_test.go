package kocha

import (
	"net/http"
	"reflect"
	"testing"
)

func TestInitRouteTable(t *testing.T) {
	actual := InitRouteTable(RouteTable{
		{
			Name:       "root",
			Path:       "/",
			Controller: FixtureRootTestCtrl{},
		},
		{
			Name:       "root_indirect",
			Path:       "/indirect",
			Controller: &FixtureRootTestCtrl{},
		},
		{
			Name:       "user",
			Path:       "/user/:id",
			Controller: FixtureUserTestCtrl{},
		},
		{
			Name:       "date",
			Path:       "/:year/:month/:day/user/:name",
			Controller: FixtureDateTestCtrl{},
		},
		{
			Name:       "static",
			Path:       "/static/*path",
			Controller: StaticServe{},
		},
	})
	expected := RouteTable{
		{
			Name:       "root",
			Path:       "/",
			Controller: FixtureRootTestCtrl{},
			MethodTypes: map[string]MethodArgs{
				"Get": MethodArgs{},
			},
		},
		{
			Name:       "root_indirect",
			Path:       "/indirect",
			Controller: FixtureRootTestCtrl{},
			MethodTypes: map[string]MethodArgs{
				"Get": MethodArgs{},
			},
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
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %q, but %q", expected, actual)
	}

	for _, v := range []interface{}{
		nil,
		"",
		"hoge",
		struct{ Controller interface{} }{},
		struct{ Controller interface{} }{Controller: ""},
		struct{ Controller interface{} }{Controller: "hoge"},
		struct{ Controller interface{} }{Controller: 1},
	} {
		func() {
			defer func() {
				if err := recover(); err == nil {
					t.Errorf("panic doesn't happened by %v", v)
				}
			}()
			InitRouteTable(RouteTable{
				{
					Name:       "testroute",
					Path:       "/",
					Controller: v,
				},
			})
		}()
	}

	// test for validate the single argument mismatch between controller method and route parameter.
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("panic doesn't happened")
			}
		}()
		InitRouteTable(RouteTable{
			{
				Name:       "testroute",
				Path:       "/",
				Controller: FixtureUserTestCtrl{},
			},
		})
	}()

	// test for validate the argument names mismatch between controller method and route parameter.
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("panic doesn't happened")
			}
		}()
		InitRouteTable(RouteTable{
			{
				Name:       "testroute",
				Path:       "/:name",
				Controller: FixtureUserTestCtrl{},
			},
		})
	}()

	// test for validate the duplicated route parameters.
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("panic doesn't happened")
			}
		}()
		InitRouteTable(RouteTable{
			{
				Name:       "testroute",
				Path:       "/:id/:id",
				Controller: FixtureUserTestCtrl{},
			},
		})
	}()

	// test for validate the multiple arguments mismatch between controller method and route parameters.
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("panic doesn't happened")
			}
		}()
		InitRouteTable(RouteTable{
			{
				Name:       "testroute",
				Path:       "/",
				Controller: FixtureDateTestCtrl{},
			},
		})
	}()

	// test for validate the invalid return value.
	for _, v := range []interface{}{
		FixtureInvalidReturnValueTypeTestCtrl{},
		FixtureInvalidNumberOfReturnValueTestCtrl{},
	} {
		func() {
			defer func() {
				if err := recover(); err == nil {
					t.Errorf("panic doesn't happened")
				}
			}()
			InitRouteTable(RouteTable{
				{
					Name:       "testroute",
					Path:       "/",
					Controller: v,
				},
			})
		}()
	}
}

func TestReverse(t *testing.T) {
	oldAppConfig := appConfig
	appConfig = newTestAppConfig()
	defer func() {
		appConfig = oldAppConfig
	}()

	actual := Reverse("root")
	expected := "/"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}

	actual = Reverse("user", 77)
	expected = "/user/77"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}

	actual = Reverse("date", 2013, 10, 26, "naoina")
	expected = "/2013/10/26/user/naoina"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}

	for _, v := range []string{"/hoge.png", "hoge.png"} {
		actual = Reverse("static", v)
		expected = "/static/hoge.png"
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}
}

func TestReverse_with_unknown_route_name(t *testing.T) {
	oldAppConfig := appConfig
	appConfig = newTestAppConfig()
	defer func() {
		appConfig = oldAppConfig
	}()

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("panic doesn't happened")
		}
	}()
	Reverse("unknown")
}

func TestReverse_with_few_arguments(t *testing.T) {
	oldAppConfig := appConfig
	appConfig = newTestAppConfig()
	defer func() {
		appConfig = oldAppConfig
	}()

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("panic doesn't happened")
		}
	}()
	Reverse("user")
}

func TestReverse_with_many_arguments(t *testing.T) {
	oldAppConfig := appConfig
	appConfig = newTestAppConfig()
	defer func() {
		appConfig = oldAppConfig
	}()

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("panic doesn't happened")
		}
	}()
	Reverse("user", 77, 100)
}

func TestReverse_with_type_mismatch(t *testing.T) {
	oldAppConfig := appConfig
	appConfig = newTestAppConfig()
	defer func() {
		appConfig = oldAppConfig
	}()

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("panic doesn't happened")
		}
	}()
	Reverse("user", "naoina")
}

func Test_dispatch_with_route_missing(t *testing.T) {
	oldAppConfig := appConfig
	appConfig = newTestAppConfig()
	defer func() {
		appConfig = oldAppConfig
	}()
	req, err := http.NewRequest("GET", "/missing", nil)
	if err != nil {
		t.Fatal(err)
	}
	controller, method, args := dispatch(req)
	if controller != nil {
		t.Errorf("Expect %v, but %v", nil, controller)
	}
	if method != nil {
		t.Errorf("Expect %v, but %v", nil, method)
	}
	if args != nil {
		t.Errorf("Expect %v, but %v", nil, args)
	}
}

func Test_dispatch(t *testing.T) {
	oldAppConfig := appConfig
	appConfig = newTestAppConfig()
	defer func() {
		appConfig = oldAppConfig
	}()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	controller, method, args := dispatch(req)
	if _, ok := controller.Interface().(*FixtureRootTestCtrl); !ok {
		t.Errorf("Expect %v, but %v", reflect.ValueOf(&FixtureRootTestCtrl{}), controller)
	}
	actual := method.Type().String()
	methodExpected := "func() kocha.Result"
	if !reflect.DeepEqual(actual, methodExpected) {
		t.Errorf("Expect %v, but %v", methodExpected, actual)
	}
	if len(args) != 0 {
		t.Errorf("Expect length is %v, but %v", 0, len(args))
	}

	req, err = http.NewRequest("GET", "/user/777", nil)
	if err != nil {
		t.Fatal(err)
	}
	controller, method, args = dispatch(req)
	if _, ok := controller.Interface().(*FixtureUserTestCtrl); !ok {
		t.Errorf("Expect %v, but %v", reflect.ValueOf(&FixtureUserTestCtrl{}), controller)
	}
	actual = method.Type().String()
	methodExpected = "func(int) kocha.Result"
	if !reflect.DeepEqual(actual, methodExpected) {
		t.Errorf("Expect %v, but %v", methodExpected, actual)
	}
	argsExpected := []interface{}{777}
	for i, arg := range args {
		if !reflect.DeepEqual(arg.Interface(), argsExpected[i]) {
			t.Errorf("Expect %v, but %v", argsExpected[i], arg)
		}
	}

	req, err = http.NewRequest("GET", "/2013/10/19/user/naoina", nil)
	if err != nil {
		t.Fatal(err)
	}
	controller, method, args = dispatch(req)
	if _, ok := controller.Interface().(*FixtureDateTestCtrl); !ok {
		t.Errorf("Expect %v, but %v", reflect.ValueOf(&FixtureDateTestCtrl{}), controller)
	}
	actual = method.Type().String()
	methodExpected = "func(int, int, int, string) kocha.Result"
	if !reflect.DeepEqual(actual, methodExpected) {
		t.Errorf("Expect %v, but %v", methodExpected, actual)
	}
	argsExpected = []interface{}{2013, 10, 19, "naoina"}
	for i, arg := range args {
		if !reflect.DeepEqual(arg.Interface(), argsExpected[i]) {
			t.Errorf("Expect %v, but %v", argsExpected[i], arg)
		}
	}
}
