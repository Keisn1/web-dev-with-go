package main

import (
	"html/template"
	"log"
	"os"
)

type User struct {
	Name string
	Age  int
	Bio  string
}

func main() {
	content, err := os.ReadFile("./helloGo.html")
	if err != nil {
		log.Fatalln("Problem loading file: ", err)
	}

	user := User{Bio: `<script>alert("Haha, you have been h4x0r3d!");</script>`, Name: "Kay", Age: 12}

	t, err := template.New("NewTemplate").Parse(string(content))
	err = t.Execute(os.Stdout, user)
	if err != nil {
		log.Fatalln("Problem Executing template: ", err)
	}
}
