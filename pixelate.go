package main

import (
  "net/http"
  "fmt"
  "encoding/base64"
  "strconv"
  "time"
  "bytes"
  "image"
  "image/jpeg"
  "image/color"
  "github.com/julienschmidt/httprouter" // mux  
)


func main() {
  fmt.Println("Pixelate server started.")
  router := httprouter.New()
  
  initData()
  
  router.GET("/", index)
  router.GET("/login", login)
  router.POST("/login", authenticate)
  router.GET("/logout", logout)
    
  router.GET("/signup", signup)
  router.POST("/signup", createUser)
  
  router.POST("/pixelate", pixelate)
  
  router.ServeFiles("/css/*filepath", http.Dir("public/css"))
  router.ServeFiles("/fonts/*filepath", http.Dir("public/fonts"))
  router.ServeFiles("/js/*filepath", http.Dir("public/js"))
  
  server := &http.Server{
  	Addr:           "0.0.0.0:1234",
  	Handler:        router,
  	ReadTimeout:    10 * time.Second,
  	WriteTimeout:   600 * time.Second,
  	MaxHeaderBytes: 1 << 20,
  }
  server.ListenAndServe()
}


// Handlers for the web app

// GET /
func index(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
  isLoggedIn(writer, request)
  t := getTemplate("index")
  err := t.Execute(writer, nil)  
  check(err, "Failed executing index template")
}

// GET /login
// Show the login page
func login(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
  t := getTemplate("login")
  err := t.Execute(writer, nil)  
  check(err, "Failed executing login template")
}

// POST /login
// Authenticate the user given the email and password
func authenticate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {  
  err := request.ParseForm()
  check(err, "Cannot parse form")

  user, err := getUser(request.PostFormValue("email"))
  check(err, "Can't get user")
  
  if user.Password == encrypt(request.PostFormValue("password")) {
    sessionId := addSession(request.PostFormValue("email"))  
    cookie := http.Cookie{
      Name:      "chitchat_cookie", 
      Value:     sessionId,
      Expires:   time.Now().Add(90 * time.Minute),
      HttpOnly:  true,
    }
    http.SetCookie(writer, &cookie)
    http.Redirect(writer, request, "/", 302)
  } else {
    http.Redirect(writer, request, "/login", 302)
  }
  
}

// GET /signup
// Show the signup page
func signup(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
  t := getTemplate("signup")  
  err := t.Execute(writer, nil)  
  check(err, "Failed executing signup template")
}

// POST /signup
// Create the user account
func createUser(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
  err := request.ParseForm()
  check(err, "Cannot parse form")  
  addUser(request.PostFormValue("name"), request.PostFormValue("email"), request.PostFormValue("password"))
  http.Redirect(writer, request, "/login", 302)
}

// GET /logout
// Logs the user out
func logout(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
  cookie, err := request.Cookie("chitchat_cookie")
  if err != http.ErrNoCookie {
    check(err, "Failed to get cookie")
    delete(sessions, cookie.Value)
  }  
  http.Redirect(writer, request, "/", 302)
}

// POST /pixelate
// pixelate a JPEG file
func pixelate(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
  isLoggedIn(writer, request)
  
  // get the content from the POSTed form
  err := request.ParseMultipartForm(10485760) // max body in memory is 10MB
  check(err, "Cannot parse form")  
  file, _, err := request.FormFile("image")  
  check(err, "Failed to get multipart file from form")
  defer file.Close()
  pixelSize, _ := strconv.Atoi(request.FormValue("pixel_size"))
  
  // decode and get original image
  original, _, err := image.Decode(file)
	bounds := original.Bounds()
  // create a new image for the pixelated
  newimage := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.X, bounds.Max.X, bounds.Max.Y))  

  // use the top left most pixel color in each rectangle (size is pixelSize)
	for y := bounds.Min.Y; y < bounds.Max.Y; y = y + pixelSize {
		for x := bounds.Min.X; x < bounds.Max.X; x = x + pixelSize {
      // get the RGBA value
      r, g, b, a := original.At(x, y).RGBA()      
      // set the RGBA value for the similar sized rectangle in the pixelated
      for i := 0; i < pixelSize; i++ {
        for j := 0; j < pixelSize; j++ {
          newimage.SetNRGBA(x+i, y+j, color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
        }
      }
		}
	}
  
  buf1 := new(bytes.Buffer)
  err = jpeg.Encode(buf1, original, nil)  
  check(err, "Cannot convert original to JPEG")
  original_str := base64.StdEncoding.EncodeToString(buf1.Bytes())

  buf2 := new(bytes.Buffer)
  err = jpeg.Encode(buf2, newimage, nil)
  check(err, "Cannot convert pixelated to JPEG")
  pixelated := base64.StdEncoding.EncodeToString(buf2.Bytes())
  
  
  images := map[string]string{
    "original": original_str,
    "pixelated": pixelated,
  }
  t := getTemplate("pixelated")
  err = t.Execute(writer, images)  
  check(err, "Failed executing pixelated template")  
}


