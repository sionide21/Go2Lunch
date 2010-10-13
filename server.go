package main

import (
	"rpc"
	"rpc/jsonrpc"
	"log"
	"net"
	"os"
	"crypto/hmac"
	"crypto/sha512"
	"crypto/subtle"
	"crypto/rand"
	"io/ioutil"
	"sync"
	"json"
	"flag"
	"strconv"
)

var port = flag.Uint("p", 1234, "Specifies the port to listen on.")
var configFile = flag.String("c", "config.json", "Specify a config file.")
var dataFile = flag.String("d", "data.json", "Specify a data file for persistence.")
var displayHelp = flag.Bool("help", false, "Displays this help message.")

type ServerConfig struct {
	Sekritz map[string]string
}

type LunchTracker struct {
	*LunchPoll
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
	read, err := ioutil.ReadFile(*dataFile)
	if err != nil {
		return nil // Don't error out if the file doesn't exist
	}
	
	err = json.Unmarshal(read, &t.LunchPoll)
	return
}

func saveDataToFile(t *LunchTracker) (err os.Error) {
	data, err := json.MarshalIndent(t.LunchPoll, "", "  ")
	if err != nil {
		return
	}
	
	err = ioutil.WriteFile(*dataFile, data, 0600)
	return
}

func (t *LunchTracker) AddPlace(args *AddPlaceArgs, place *uint) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	*place = t.LunchPoll.addPlace(args.Name, args.Auth.Name)
	return nil
}

func (t *LunchTracker) DelPlace(args *UIntArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	*success = t.LunchPoll.delPlace(args.Num)
	return nil
}

func (t *LunchTracker) Drive(args *UIntArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	*success = t.LunchPoll.drive(args.Auth.Name, args.Num)
	return nil
}

func (t *LunchTracker) UnDrive(args *EmptyArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	*success = t.LunchPoll.unDrive(args.Auth.Name)
	return nil
}

func (t *LunchTracker) Vote(args *UIntArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}

	*success = t.LunchPoll.vote(args.Auth.Name, args.Num)
	return nil
}

func (t *LunchTracker) UnVote(args *EmptyArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	*success = t.LunchPoll.unVote(args.Auth.Name)
	return nil
}

func (t *LunchTracker) DisplayPlaces(args *EmptyArgs, response *[]Place) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	*response = t.LunchPoll.displayPlaces()
	
	return saveDataToFile(t)
}


func (t *LunchTracker) Challenge(name *string, challenge *Bin) os.Error {
	_, valid := userMap[(*name)]
	if !valid {
		valid = checkUser(*name)
		if !valid {
			return nil
		}
	}

	*challenge = make([]byte, 512)
	n, err := rand.Read(*challenge)

	if err != nil || n != 512 {
		panic("Challenge Generation Failed")
	}

	userMap[*name].SChallenge = challenge
	return nil
}

func checkUser(name string) bool {
	loadUsersFromFile()
	cMutex.Lock()
	_, valid := config.Sekritz[name]
	cMutex.Unlock()
	if valid {
		sc := make(Bin, 512)
		userMap[name] = &Auth{Name: name, SChallenge: &sc}
	}

	return valid
}


func verify(a *Auth, d Byter) (bool, os.Error) {
	cMutex.Lock()
	// Is this lock at all necessary?
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
	flag.Parse()

	if *displayHelp {
		flag.PrintDefaults()
		return
	}

	userMap = make(map[string]*Auth)
	err := loadUsersFromFile()
	if err != nil {
		log.Exit("Error reading config file. Have you created it?\nCaused By: ", err)
	}

	t := &LunchTracker{NewPoll()}
	err = loadDataFromFile(t)
	if err != nil {
		log.Exit("Error reading data file.\nCaused by: ", err)
	}
	rpc.Register(t)
	l, e := net.Listen("tcp", ":"+strconv.Uitoa(*port))
	if e != nil {
		log.Exit("listen error:", e)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Stderr(err)
		} else {
			go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}
