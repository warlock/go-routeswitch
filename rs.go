package main

import(
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"time"
	"log"
	"os/exec"
)

const (
	CONF_FILE="conf.json"
	LOCK_FILE="rs.lock"
)

type Config struct {
 //   Mail string
    Routers Router
}

type Router []string

func main() {
	fmt.Println("Route Switch - 0.9go+4linux")
	path, err := os.Getwd()
	checkErr(err)
  data, err := ioutil.ReadFile(path + "/" + CONF_FILE)
  checkErr(err)
  var c Config
	err = json.Unmarshal(data, &c)
	checkErr(err)
	var routers Router
	routers = c.Routers

  if(checkLocker(path)) {
    fmt.Println("Another RS process is working. Or delete file \"rs.lock\".")
  } else {
    makeLocker(path)
    for !checkRoute() {
    	fmt.Println("Connection failed. Want a new route.")
			for _, router := range(routers) {
				fmt.Print("Route changed to: ", router)
				_, err := exec.Command("route","del","default","eth0").Output()
				checkErr(err)
				time.Sleep(2 * time.Second)
				_, err = exec.Command("route","add","default","gw", router,"eth0").Output()
				checkErr(err)
				fmt.Println(" CHANGED!")
				time.Sleep(2 * time.Second)
				fmt.Print("Checking connection... ")
				if !checkRoute() {
					fmt.Println(" FAILED!")
				} else {
					fmt.Println(" OK!")
					break
				}
			}
		}
    deleteLocker(path)
    fmt.Println("THE END")
  }
}


func checkRoute() bool {
	if _, err := http.Get("http://google.com/"); err == nil  {
		return true
	} else {
		return false
	}
}

func checkLocker(path string) bool {
	if _, err := os.Stat(path + "/" + LOCK_FILE); err == nil {
		return true
	} else {
		return false
	}
}

func makeLocker(path string) {
	f, err := os.Create(path + "/" + LOCK_FILE)
    checkErr(err)
	defer f.Close()
}

func deleteLocker(path string) {
	err := os.Remove(path + "/" + LOCK_FILE);
	checkErr(err)
}

func checkErr(e error) {
	if e != nil {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Panic(e)
	}
}
