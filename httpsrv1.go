package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var inch chan string = make(chan string)
var outch chan string = make(chan string)
var password string
var httpPort string

func init() {
	flag.StringVar(&password, "password", "mypassword", "The password for users")
	flag.StringVar(&httpPort, "port", "9999", "The http listen port")
}

func main() {
	flag.Parse()
	http.HandleFunc("/", sayHelloName)
	http.HandleFunc("/login", login)
	http.HandleFunc("/updateIp", updateIp)
	http.HandleFunc("/getIp", getIp)
	go cacheIp(inch, outch)
	fmt.Println("Listen on :", httpPort)
	err := http.ListenAndServe(":"+httpPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe err: ", err)
	}
}

func cacheIp(in chan string, out chan string) {
	ips := make(map[string]string)
	for {
		select {
		case x := <-in:
			data := strings.Split(x, ":")
			if data[0] == "update" {
				ips[data[1]] = data[2]
				fmt.Println("update ip", data[1], data[2])
			}
			if data[0] == "get" {
				ip2, ok := ips[data[1]]
				fmt.Println("get ip", data[1], ip2)
				if ok {
					out <- ip2
				} else {
					out <- "127.0.0.1"
				}
			}
		}
	}
}

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Printf("k %s, v %s", k, strings.Join(v, ","))
	}
	fmt.Fprintf(w, "Hello World")
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		fmt.Println("username: ", r.Form["username"])
		fmt.Println("username.len: ", r.Form["username"][0])
		fmt.Println("password: ", r.Form["password"])
	}
}

func updateIp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user, pass, ok := r.BasicAuth()
	if ok {
		if user[:5] == "robin" && pass == password {
			rip := strings.Split(r.RemoteAddr, ":")[0]
			fmt.Println("remoteAddr", rip)
			inch <- "update:" + user + ":" + rip
		}
	}
}

func getIp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user, pass, ok := r.BasicAuth()
	if ok {
		if user[:5] == "robin" && pass == password {
			rip := strings.Split(r.RemoteAddr, ":")[0]
			fmt.Println("remoteAddr", rip)
			inch <- "get:" + user + ":127.0.0.1"
			w.Write([]byte(<-outch))

		}
	}
}
