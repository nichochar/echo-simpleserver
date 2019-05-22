package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

// User is a simple struct representing a user
type User struct {
	Name string
	Age  int
}

func (u User) String() string {
	return fmt.Sprintf("Name: %s || Age: %d\n", u.Name, u.Age)
}

// DB is an in memory database: map of user names to their structs
var DB = map[string]User{}

func getIndex(c echo.Context) error {
	users := make([]string, len(DB))
	for name := range DB {
		users = append(users, name)
	}
	return c.String(http.StatusOK, "Users\n"+strings.Join(users, "\n"))
}

func getUser(c echo.Context) error {
	name := c.Param("name")
	user, ok := DB[name]
	if !ok {
		return c.String(http.StatusNotFound, "Not found")
	}

	return c.String(http.StatusOK, user.String())
}

func putUser(c echo.Context) error {
	name := c.FormValue("name")
	ageRaw := c.FormValue("age")
	age, err := strconv.Atoi(ageRaw)
	if err != nil {
		return c.String(http.StatusBadRequest, "Age must be an integer")
	}
	_, exists := DB[name]
	if exists {
		return c.String(http.StatusConflict, "User with name "+name+" already exists")
	}
	DB[name] = User{Name: name, Age: age}
	fmt.Printf("Added user:\n%v\n", DB[name])
	return c.String(http.StatusOK, "Added "+name)
}

func main() {
	e := echo.New()
	e.Static("/static", "static")
	e.GET("/", getIndex)
	e.GET("/users/:name", getUser)
	e.PUT("/users", putUser)

	e.Logger.Fatal(e.Start(":1323"))
}
