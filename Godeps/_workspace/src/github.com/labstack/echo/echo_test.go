package echo

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type (
	user struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

var u1 = user{
	ID:   "1",
	Name: "Joe",
}

// TODO: Improve me!
func TestEchoMaxParam(t *testing.T) {
	e := New()
	e.MaxParam(8)
	if e.maxParam != 8 {
		t.Errorf("max param should be 8, found %d", e.maxParam)
	}
}

func TestEchoIndex(t *testing.T) {
	e := New()
	e.Index("example/public/index.html")
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(GET, "/", nil)
	e.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Errorf("status code should be 200, found %d", w.Code)
	}
}

func TestEchoStatic(t *testing.T) {
	e := New()
	e.Static("/scripts", "example/public/scripts")
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(GET, "/scripts/main.js", nil)
	e.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Errorf("status code should be 200, found %d", w.Code)
	}
}

func TestEchoMiddleware(t *testing.T) {
	e := New()
	b := new(bytes.Buffer)

	// func(*echo.Context)
	e.Use(func(c *Context) {
		b.WriteString("a")
	})

	// func(echo.HandlerFunc) echo.HandlerFunc
	e.Use(func(h HandlerFunc) HandlerFunc {
		return HandlerFunc(func(c *Context) {
			b.WriteString("b")
			h(c)
		})
	})

	// http.HandlerFunc
	e.Use(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b.WriteString("c")
	}))

	// http.Handler
	e.Use(http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b.WriteString("d")
	})))

	// func(http.Handler) http.Handler
	e.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b.WriteString("e")
			h.ServeHTTP(w, r)
		})
	})

	// func(http.ResponseWriter, *http.Request)
	e.Use(func(w http.ResponseWriter, r *http.Request) {
		b.WriteString("f")
	})

	// Route
	e.Get("/hello", func(c *Context) {
		c.String(200, "world")
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(GET, "/hello", nil)
	e.ServeHTTP(w, r)
	if b.String() != "abcdef" {
		t.Errorf("buffer should be abcdef, found %s", b.String())
	}
	if w.Body.String() != "world" {
		t.Error("body should be world")
	}
}

func TestEchoHandler(t *testing.T) {
	e := New()

	// func(*echo.Context)
	e.Get("/1", func(c *Context) {
		c.String(http.StatusOK, "1")
	})
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(GET, "/1", nil)
	e.ServeHTTP(w, r)
	if w.Body.String() != "1" {
		t.Error("body should be 1")
	}

	// http.Handler/http.HandlerFunc
	e.Get("/2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("2"))
	}))
	w = httptest.NewRecorder()
	r, _ = http.NewRequest(GET, "/2", nil)
	e.ServeHTTP(w, r)
	if w.Body.String() != "2" {
		t.Error("body should be 2")
	}

	// func(http.ResponseWriter, *http.Request)
	e.Get("/3", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("3"))
	})
	w = httptest.NewRecorder()
	r, _ = http.NewRequest(GET, "/3", nil)
	e.ServeHTTP(w, r)
	if w.Body.String() != "3" {
		t.Error("body should be 3")
	}
}

func TestEchoGroup(t *testing.T) {
	b := new(bytes.Buffer)
	e := New()
	e.Use(func(*Context) {
		b.WriteString("1")
	})
	e.Get("/users", func(*Context) {})
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(GET, "/users", nil)
	e.ServeHTTP(w, r)
	if b.String() != "1" {
		t.Errorf("should only execute middleware 1, executed %s", b.String())
	}

	// Group
	g1 := e.Group("/group1")
	g1.Use(func(*Context) {
		b.WriteString("2")
	})
	g1.Get("/home", func(*Context) {})
	b.Reset()
	w = httptest.NewRecorder()
	r, _ = http.NewRequest(GET, "/group1/home", nil)
	e.ServeHTTP(w, r)
	if b.String() != "12" {
		t.Errorf("should execute middleware 1 & 2, executed %s", b.String())
	}

	// Group with no parent middleware
	g2 := e.Group("/group2", func(*Context) {
		b.WriteString("3")
	})
	g2.Get("/home", func(*Context) {})
	b.Reset()
	w = httptest.NewRecorder()
	r, _ = http.NewRequest(GET, "/group2/home", nil)
	e.ServeHTTP(w, r)
	if b.String() != "3" {
		t.Errorf("should execute middleware 3, executed %s", b.String())
	}
}

func TestEchoMethod(t *testing.T) {
	e := New()
	e.Connect("/", func(*Context) {})
	e.Delete("/", func(*Context) {})
	e.Get("/", func(*Context) {})
	e.Head("/", func(*Context) {})
	e.Options("/", func(*Context) {})
	e.Patch("/", func(*Context) {})
	e.Post("/", func(*Context) {})
	e.Put("/", func(*Context) {})
	e.Trace("/", func(*Context) {})
}

func TestEchoNotFound(t *testing.T) {
	e := New()

	// Default NotFound handler
	r, _ := http.NewRequest(GET, "/files", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("status code should be 404, found %d", w.Code)
	}

	// Customized NotFound handler
	e.NotFoundHandler(func(c *Context) {
		c.String(404, "not found")
	})
	w = httptest.NewRecorder()
	e.ServeHTTP(w, r)
	if w.Body.String() != "not found" {
		t.Errorf("body should be `not found`")
	}
}

func verifyUser(u2 *user, t *testing.T) {
	if u2.ID != u1.ID {
		t.Errorf("user id should be %s, found %s", u1.ID, u2.ID)
	}
	if u2.Name != u1.Name {
		t.Errorf("user name should be %s, found %s", u1.Name, u2.Name)
	}
}
