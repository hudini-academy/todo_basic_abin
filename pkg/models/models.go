package models 
 
import ( 
    "errors" 
    "time" 
)
 
var (
    ErrNoRecord = errors.New("models: no matching record found") 
    //If user tries to login with an incorrect email address
    ErrInvalidCredentials = errors.New("models:invalid credentials")
    //If user tries to use existing email address
    ErrDupicateEmail =errors.New("models:duplicate email")
)
type ToDo struct { 
    ID      int 
    Title   string 
    Created time.Time 
    Expires time.Time 
}
//define new user type
type User struct{
    ID              int
    Name            string
    Email           string
    HashedPassword  []byte
    Created         time.Time
}