package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/template"

	"github.com/GeertJohan/go.rice"
	"github.com/julienschmidt/httprouter"
)

func Serve(cfg *Config) error {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println(cfg)
	repo := simpleLightRepo{}
	for _, lc := range cfg.Lights {

		ledc, err := NewLEDController(lc)
		if err != nil {
			return err
		}

		err = repo.Add(ledc)
		if err != nil {
			return err
		}
	}
	ctrl := &Controller{
		Lights: &repo,
	}

	router := httprouter.New()
	router.HandlerFunc("GET", "/", ctrl.GetSchedules)

	return http.ListenAndServe(":"+port, router)
}

var Templates struct {
	T     map[string]*template.Template
	Funcs template.FuncMap
	sync.Mutex
}

func T(name string) *template.Template {
	return template.Must(GetTemplate(name))
}

func GetTemplate(name string) (*template.Template, error) {
	Templates.Lock()
	defer Templates.Unlock()

	if Templates.T == nil {
		Templates.T = map[string]*template.Template{}
	}

	inDev := os.Getenv("DEV_RUN") != ""

	if t, ok := Templates.T[name]; ok && !inDev {
		return t, nil
	}

	box, err := rice.FindBox("./templates")
	if err != nil {
		return nil, err
	}
	t := template.New("base").Funcs(Templates.Funcs)

	load := func(name string) {
		var str string
		if err == nil {
			str, err = box.String(name)
		}
		if err == nil {
			t, err = t.Parse(str)
		}
	}

	load("layout.tmpl")
	load(name + ".tmpl")

	if err != nil {
		return nil, err
	}
	Templates.T[name] = t
	return t, nil
}

func renderErr(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(rw, err)
}
func render(rw http.ResponseWriter, req *http.Request, page string, data interface{}) {
	err := T(page).Execute(rw, data)
	if err != nil {
		fmt.Println("bah", err)
	}
}

type Controller struct {
	Lights LightRepo
}

var (
	ErrNoSuchPin = errors.New("no such pin")
)

type simpleLightRepo struct {
	l []*LEDController
	sync.Mutex
}

func (s *simpleLightRepo) All() ([]*LEDController, error) {
	s.Lock()
	defer s.Unlock()
	r := make([]*LEDController, len(s.l))
	copy(r, s.l)
	return r, nil
}

func (s *simpleLightRepo) Add(ledc *LEDController) error {
	s.Lock()
	defer s.Unlock()

	for i := 0; i < len(s.l); i++ {
		if s.l[i].Pin == ledc.Pin {
			return fmt.Errorf("pin %d is already in use by %s", s.l[i].Name)
		}
	}

	s.l = append(s.l, ledc)
	return nil
}

func (s *simpleLightRepo) Delete(pin uint) error {
	s.Lock()
	defer s.Unlock()

	for _, l := range s.l {
		if l.Pin == pin {
			return l.Close()
		}
	}
	return ErrNoSuchPin
}

func (s *simpleLightRepo) Get(pin uint) (*LEDController, error) {
	s.Lock()
	s.Unlock()

	for _, l := range s.l {
		if l.Pin == pin {
			return l, nil
		}
	}
	return nil, ErrNoSuchPin
}

var _ LightRepo = &simpleLightRepo{}

type LightRepo interface {
	All() ([]*LEDController, error)
	Get(pin uint) (*LEDController, error)
	Add(*LEDController) error
	Delete(pin uint) error
}

func (c *Controller) GetSchedules(rw http.ResponseWriter, req *http.Request) {

	l, err := c.Lights.All()
	if err != nil {
		renderErr(rw, err)
		return
	}
	render(rw, req, "index",
		struct{ Lights []*LEDController }{Lights: l})
}
