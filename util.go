package main

import (
  "net/http"
  "html/template"
  "errors"
  "log"
)

// check if the user is currently logged in
func isLoggedIn(writer http.ResponseWriter, request *http.Request)(authenticated bool){
  cookie, err := request.Cookie("chitchat_cookie")
  if err == http.ErrNoCookie {
    http.Redirect(writer, request, "/login", 302)
  } else {
    check(err, "Failed to get cookie")
    if _, ok := sessions[cookie.Value]; !ok {
      http.Redirect(writer, request, "/login", 302)
    }
  }  
  return
}

// Get the user
func getUser(email string) (user User, err error) {
  if u, ok := users[email]; ok {
    user = u
  } else {
    err = errors.New("User not found")
  }  
  return
}

// get the template
func getTemplate(name string)(t *template.Template) {
	t = template.New(name)
	t = template.Must(t.ParseGlob("templates/*.html"))  
  return
}

// checks errors
func check(err error, msg string) {
  if err != nil {
    log.Println(msg, err)
  }
}

