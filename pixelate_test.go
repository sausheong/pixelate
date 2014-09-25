package main

import(  
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/julienschmidt/httprouter"
    "strings"
)

// Test the index page
// Route should redirect to the login page
func TestIndex(t *testing.T) {   
  router := httprouter.New()
  router.GET("/", index)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)    
  router.ServeHTTP(writer, request)
  
  if writer.Code != 302 {
		t.Errorf("Response code is %v", writer.Code)
	}  
}

// Test package function that adds a new user
func TestAddUser(t *testing.T) {
  addUser("Sau Sheong", "sausheong@gmail.com", "password")
  if val, ok := users["sausheong@gmail.com"]; !ok {
    t.Errorf("Cannot add user")
  } else {
    if val.Name != "Sau Sheong" {
      t.Errorf("User name is wrong")
    }
  }
}

// Test the sign up form
// Route should sign up a new user and redirect to the login page
func TestSignUp(t *testing.T) {
  router := httprouter.New()
  router.POST("/signup", createUser)
  writer := httptest.NewRecorder()
  body := strings.NewReader("name=Sau Sheong&email=sausheong@gmail.com&password=password")
  request, _ := http.NewRequest("POST", "/signup", body)
  request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  router.ServeHTTP(writer, request)

  if writer.Code != 302 {
    t.Errorf("Response code is %v", writer.Code)
  }
  if writer.Header().Get("Location") != "/login" {
    t.Errorf("Location is %v", writer.Header().Get("Location"))
  }
}

// Test user authentication
// Route should authenticate a user given the email and password
// from a form post 
// If authentication is successful, a cookie that starts with
// pixelate_cookie must be created and added to the browser
func TestAuthenticate(t *testing.T) {
  addUser("Sau Sheong", "sausheong@gmail.com", "password")
  router := httprouter.New()
  router.POST("/login", authenticate)  
  writer := httptest.NewRecorder()
  body := strings.NewReader("email=sausheong@gmail.com&password=password")
  request, _ := http.NewRequest("POST", "/login", body)
  request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  router.ServeHTTP(writer, request)

  if writer.Code != 302 {
    t.Errorf("Response code is %v", writer.Code)
  }
  if !strings.HasPrefix(writer.Header().Get("Set-Cookie"), "pixelate_cookie") {
    t.Errorf("Cookie not set")
  }
}

// Test user authentication
// Make sure that user authentication fails if the wrong password is given
func TestAuthenticateFail(t *testing.T) {
  addUser("Sau Sheong", "sausheong@gmail.com", "password")
  router := httprouter.New()
  router.POST("/login", authenticate)  
  writer := httptest.NewRecorder()
  body := strings.NewReader("email=sausheong@gmail.com&password=wrong_password")
  request, _ := http.NewRequest("POST", "/login", body)
  request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  router.ServeHTTP(writer, request)
  
  if writer.Code != 302 {
    t.Errorf("Response code is %v", writer.Code)
  }
  if writer.Header().Get("Location") != "/login" {
    t.Errorf("Not redirected to login")
  }
  if strings.HasPrefix(writer.Header().Get("Set-Cookie"), "pixelate_cookie") {
    t.Errorf("Cookie is set")
  }
}
