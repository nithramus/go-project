package main

import (
	"flag"
	"go-project/test/tuto/trace"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

var avatars Avatar = UseFileSystemAvatar

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()
	gomniauth.SetSecurityKey("ezretyzrfjzaeraezraezraezr")
	gomniauth.WithProviders(
		google.New("735975902608-cvtma2109pojkdrojuec540utpi7l3p1.apps.googleusercontent.com",
			"6Pc-qHDWQavqkTYmZlcBXvTJ", "http://localhost:3000/auth/callback/google"))
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/room", r)
	go r.run()
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/upload", MustAuth(&templateHandler{filename: "upload.html"}))
	http.HandleFunc("/uploader", uploadHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	ref := templateHandler{filename: "chat.html"}
	http.HandleFunc("/", ref.ServeHTTP)
	log.Println("start web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServce:", err)
	}
	log.Println("ended")
}
