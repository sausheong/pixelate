package main

// 
// Stores data needed by the application 
// 

import (
  "fmt"
  "time"
  "crypto/sha1"
  "crypto/rand"
)

type User struct {
  Email     string
  Name      string
  Password  string
}

type Session struct {
  Email     string
  Timestamp time.Time
}

var sessions map[string]Session
var users map[string]User

func initData() {
  sessions = make(map[string]Session)
  users    = make(map[string]User)
}

// add a new user
func addUser(name string, email string, password string) {
  user := User{
    Email: email,
    Name: name,
    Password: encrypt(password),
  }
  users[email] = user
}

// add a session into the session store
func addSession(email string) (sessionId string) {
  sessionId = createUUID()
  session := Session{
    Email: email,
    Timestamp: time.Now(),
  }  
  sessions[sessionId] = session
  return
}

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func createUUID() (uuid string) {
  u := new([16]byte)
  _, err := rand.Read(u[:])
  check(err, "Cannot generate UUID")
  // 0x40 is reserved variant from RFC 4122  
  u[8] = (u[8] | 0x40) & 0x7F
  // Set the four most significant bits (bits 12 through 15) of the
  // time_hi_and_version field to the 4-bit version number.  
  u[6] = (u[6] & 0xF) | (0x4 << 4)
  uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
  return
}

func encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))  
  return
}