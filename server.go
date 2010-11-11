package main

import (
	"rpc"
	"http"
	"log"
	"net"
	"os"
	"crypto/hmac"
	"crypto/sha512"
	"crypto/subtle"
	"io/ioutil"
	"sync"
	"json"
	"flag"
	"strconv"
	"gob"
)

var port = flag.Uint("p", 1234, "Specifies the port to listen on.")
var configFile = flag.String("c", "config.json", "Specify a config file.")
var dataFile = flag.String("d", "data.gob", "Specify a data file.")
var displayHelp = flag.Bool("help", false, "Displays this help message.")

type ServerConfig struct {
	Sekritz map[string]string
}

var userMap map[string]*Auth
var config *ServerConfig
var cMutex sync.Mutex

func loadUsersFromFile() (err os.Error) {
	tempConfig := &ServerConfig{}
	read, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(read, tempConfig)
	if err != nil {
		return err
	}

	cMutex.Lock()
	config = tempConfig
	cMutex.Unlock()
	return nil
}

func loadDataFromFile(t *LunchTracker) (err os.Error) {
	// read, err := ioutil.ReadFile(*dataFile)
	// if err != nil {
	// 	return nil // Don't error out if the file doesn't exist
	// }
	// 
	// err = json.Unmarshal(read, t.LunchPoll)
	// fmt.Println(t.LunchPoll)

	file, err := os.Open(*dataFile, os.O_RDONLY, 0)
	if err != nil {
		return nil
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	poll := t.getPoll()
	err = dec.Decode(poll)
	t.putPoll(poll)

	return
}

func saveDataToFile(t *LunchTracker) (err os.Error) {
	// data, err := json.MarshalIndent(t.LunchPoll, "", "  ")
	// if err != nil {
	// 	return
	// }
	// 
	// err = ioutil.WriteFile(*dataFile, data, 0600)
	// return

	file, err := os.Open(*dataFile, os.O_WRONLY|os.O_CREAT|os.O_TRUNC, 0600)
	if err != nil {
		return
	}

	enc := gob.NewEncoder(file)

	poll := t.getPoll()
	err = enc.Encode(poll)
	t.putPoll(poll)
	file.Close()

	return
}

func checkUser(name string) bool {
	loadUsersFromFile()
	cMutex.Lock()
	_, valid := config.Sekritz[name]
	cMutex.Unlock()
	if valid {
		bin := make(Bin, 512)
		userMap[name] = &Auth{Name: name, SChallenge: &bin}
	}

	return valid
}


func verify(a *Auth, d Byter) (bool, os.Error) {
	cMutex.Lock()
	key, ok := config.Sekritz[(*a).Name]
	cMutex.Unlock()

	if !ok {
		return false, os.NewError("Unknown User")
	}

	mac := hmac.New(sha512.New, []byte(key))

	mac.Write([]byte((*a).Name))
	mac.Write(*(*a).CChallenge)
	mac.Write(d.Byte())
	mac.Write(*userMap[(*a).Name].SChallenge)
	if subtle.ConstantTimeCompare(mac.Sum(), *(*a).Mac) == 1 {
		return true, nil
	}
	return false, os.NewError("Authentication Failed")
}


func main() {
	log.SetOutput(os.Stderr)
	flag.Parse()

	if *displayHelp {
		flag.PrintDefaults()
		return
	}
	userMap = make(map[string]*Auth)
	err := loadUsersFromFile()
	if err != nil {
		log.Exit("Error reading config file. Have you created it?\nCoused By: ", err)
	}
	t := newPollChan()

	RegisterTypes()
	rpc.Register(t)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+strconv.Uitoa(*port))
	if e != nil {
		log.Exit("listen error:", e)
	}
	http.Serve(l, nil)
	// for {
	// 	conn, err := l.Accept()
	// 	if err != nil {
	// 		log.Println(err)
	// 	} else {
	// 		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	// 	}
	// }
}

func newPollChan() *LunchTracker {
	ch := LunchTracker(make(chan LunchPoll))
	poll := NewPoll()
	go func() {
		for {
			ch <- poll
			poll = <-ch
		}
	}()
	return &ch
}
