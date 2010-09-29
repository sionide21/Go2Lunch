package main

import (
	"flag"
	"strconv"
	"rpc"
	"fmt"
	"strings"
	"os"
	"crypto/hmac"
	"crypto/sha512"
	"crypto/rand"
	"io/ioutil"
	"json"
	"encoding/base64"
	"path"
)

const (
	destErr    = "Invalid Destination"
	configFile = ".lunch_config.json"
)

var add = flag.Bool("a", false, "add a place")
var del = flag.Bool("rm", false, "remove a place")
var seats = flag.Uint("d", 0, "driver with _ additional seats")
var unvote = flag.Bool("u", false, "unvote")
var server = flag.String("s", "", "[host]:[port]")
var name = flag.String("n", "", "user name")
var walk = flag.Bool("w", false, "not driving")
var sekrit = ""
var user = ""
var host = ""

type LunchServer struct {
	*rpc.Client
	Auth
}

func getConfig() (err os.Error) {
	config := make(map[string]string)
	home := os.Getenv("HOME")
	read, err := ioutil.ReadFile(path.Join(home, configFile))
	if err != nil {
		return
	}
	err = json.Unmarshal(read, &config)
	if err != nil {
		return
	}
	var ok bool
	sekrit, ok = config["sekrit"]
	if !ok {
		return os.NewError("No Sekrit")
	}
	host, ok = config["host"]
	if !ok {
		return os.NewError("No Host")
	}
	user, ok = config["user"]
	if !ok {
		return os.NewError("No User")
	}
	return
}

func genConfig() (err os.Error) {
	getConfig()
	config := make(map[string]string)
	if *server != "" {
		config["host"] = *server
	} else {
		config["host"] = host
	}
	if *name != "" {
		config["user"] = *name
	} else {
		config["user"] = user
	}
	if sekrit == "" {
		newSekrit := make([]byte, 512, 512)
		encoder := base64.StdEncoding

		_, err := rand.Read(newSekrit)
		if err != nil {
			panic("Random not random")
		}
		newEncoded := make([]byte, encoder.EncodedLen(len(newSekrit)), encoder.EncodedLen(len(newSekrit)))
		encoder.Encode(newEncoded, newSekrit)
		config["sekrit"] = string(newEncoded)
	}
	if config["sekrit"] == "" {
		err = os.NewError("Missing Sekrit")
		return
	} else {
		sekrit = config["sekrit"]
	}
	if config["user"] == "" {
		err = os.NewError("Missing User")
		return
	} else {
		user = config["user"]
	}
	if config["host"] == "" {
		err = os.NewError("Missing host")
	} else {
		host = config["host"]
	}
	data, err := json.Marshal(config)
	if err != nil {
		return
	}
	home := os.Getenv("HOME")
	err = ioutil.WriteFile(path.Join(home, configFile), data, 0600)
	return
}

func main() {
	flag.Parse()
	err := getConfig()
	if err != nil {
		err = genConfig()
		if err != nil {
			panic(err)
		}
		fmt.Println("\"" + user + "\":\"" + sekrit + "\"")
		return
	}
	r, e := rpc.DialHTTP("tcp", host)
	if e != nil {
		fmt.Println("Cannot connect to server: " + host)
		os.Exit(-1)
	}
	remote := &LunchServer{r, Auth{Name: "User"}}

	if *seats != 0 {
		if !remote.drive(*seats) {
			fmt.Println("Drive Failed")
		}
		return
	}

	if *walk {
		if !remote.undrive() {
			fmt.Println("UnDrive Failed")
		}
		return
	}
	if !*unvote && flag.NArg() == 0 {
		places := remote.displayPlaces()
		if places != nil {
			for _, p := range *places {
				fmt.Println(p.String())
			}
		}
		return
	}

	var dest uint = 0
	if *add && !*del {
		name := strings.Join(flag.Args(), " ")
		dest = remote.addPlace(name)
		return
	}

	if dest == 0 {
		dest, _ = strconv.Atoui(flag.Arg(0))
	}

	if *del {
		if !remote.delPlace(dest) {
			fmt.Println("Could not delete. This can happen if there are still votes on the place or if you did not nominate it.")
		}
		return
	}

	if *unvote {
		if !remote.unvote() {
			fmt.Println("Unvoting Failed")
		}
		return
	} else {
		if !remote.vote(dest) {
			fmt.Println("Vote Failed")
		}
		return
	}
}

func (t *LunchServer) addPlace(name string) (place uint) {
	args := &AddPlaceArgs{Name: name}
	args.Auth = *(t.calcAuth(args))
	t.Call("LunchTracker.AddPlace", &args, &place)
	return
}

func (t *LunchServer) delPlace(dest uint) (suc bool) {
	args := &UIntArgs{Num: dest}
	args.Auth = *(t.calcAuth(args))
	t.Call("LunchTracker.DelPlace", args, &suc)
	return
}

func (t *LunchServer) drive(seats uint) (suc bool) {
	args := &UIntArgs{Num: seats}
	args.Auth = *(t.calcAuth(args))
	t.Call("LunchTracker.Drive", args, &suc)
	return
}

func (t *LunchServer) vote(dest uint) (suc bool) {
	args := &UIntArgs{Num: dest}
	args.Auth = *(t.calcAuth(args))
	t.Call("LunchTracker.Vote", args, &suc)
	return
}

func (t *LunchServer) unvote() (suc bool) {
	args := &EmptyArgs{}
	args.Auth = *(t.calcAuth(args))
	t.Call("LunchTracker.UnVote", args, &suc)
	return
}

func (t *LunchServer) displayPlaces() *[]Place {
	args := &EmptyArgs{}
	args.Auth = *(t.calcAuth(args))
	var places []Place
	t.Call("LunchTracker.DisplayPlaces", args, &places)
	return &places
}

func (t *LunchServer) undrive() (suc bool) {
	args := &EmptyArgs{}
	args.Auth = *(t.calcAuth(args))
	t.Call("LunchTracker.UnDrive", args, &suc)
	return
}

func (t *LunchServer) calcAuth(d Byter) (a *Auth) {
	var challenge []byte
	t.Call("LunchTracker.Challenge", &user, &challenge)
	a = &Auth{Name: user, CChallenge: challenge}
	sum(d, a)
	return
}

func sum(d Byter, a *Auth) {

	challenge := make([]byte, 512)
	_, err := rand.Read(challenge)

	if err != nil {
		panic("Random not random.")
	}
	(*a).CChallenge = challenge

	mac := hmac.New(sha512.New, []byte(sekrit))
	mac.Write([]byte((*a).Name))
	mac.Write((*a).CChallenge)
	mac.Write(d.Byte())
	mac.Write((*a).SChallenge)
	(*a).Mac = mac.Sum()
}
