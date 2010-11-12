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
	"template"
)

const (
	configFile   = "config"
	templateFile = "template"
)

const (
	destErr       = "Invalid Destination"
	unvoteError   = "UnVoting Failed"
	undriveError  = "UnDriving Failed"
	voteError     = "Vote Failed"
	driveError    = "Drive Failed"
	commentError  = "Comment Failed"
	delPlaceError = "Could not delete. This can happen if there are still votes on the place or if you did not nominate it."

	noSekrit  = "No Auth Token"
	noUser    = "No User Name"
	noHost    = "No Host"
	badRandom = "Random NOT random"
)

const clientVersion = "0.03"

var add = flag.Bool("a", false, "add a place")
var del = flag.Bool("rm", false, "remove a place")
var seats = flag.Uint("d", 0, "driver with _ additional seats")
var unvote = flag.Bool("u", false, "unvote")
var server = flag.String("s", "", "[host]:[port]")
var name = flag.String("n", "", "user name")
var walk = flag.Bool("w", false, "not driving")
var debug = flag.Bool("g", false, "debug")
var up = flag.Bool("p", false, "enable automatic update checks")
var version = flag.Bool("v", false, "show current version")
var printJson = flag.Bool("json", false, "display output in json")
var noVotes = flag.Bool("e", false, "print everyone")
var comment = flag.Bool("c", false, "comment")
var sekrit = ""
var user = ""
var host = ""

type LunchServer struct {
	*rpc.Client
}

func main() {
	defer func() {
		if !*debug {
			if bad := recover(); bad != nil {
				fmt.Println(bad)
			}
		}
	}()

	flag.Parse()

	// Check for new versions of the client application.
	var upChan chan os.Error

	if *up {
		upChan = CheckForUpdates()
	}

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
	remote := &LunchServer{r}

	dest, err := strconv.Atoui(flag.Arg(0))
	var poll *LunchPoll

	switch {
	case *version:
		fmt.Printf("Go2Lunch %s\n", clientVersion)
	case *seats != 0:
		remote.drive(*seats)
	case *walk:
		remote.undrive()
	case *add:
		name := strings.Join(flag.Args(), " ")
		dest = remote.addPlace(name)
	case *del:
		remote.delPlace(dest)
	case *unvote:
		remote.unvote()
	case dest != 0:
		remote.vote(dest)
	case *comment:
		comment := strings.Join(flag.Args(), " ")
		remote.comment(comment)
	default:
		poll = remote.displayPlaces()
	}

	if poll != nil {
		if *printJson {
			out, err := json.Marshal(poll)
			if err != nil {
				panic(err.String())
			}
			fmt.Println(string(out))
		} else {
			for i, p := range poll.Places {
				if i != 0 || *noVotes {
					ppPlace(p)
				}
			}
		}
	}

	if *up {
		upErr := <-upChan
		if upErr != nil {
			panic(upErr)
		}
	}
}

func ppPlace(place *Place) {
	home := os.Getenv("HOME")
	t, err := template.ParseFile(path.Join(home, ".lunch", templateFile), nil)
	if err != nil {
		fmt.Println(place.String())
		return
	}
	t.Execute(place, os.Stdout)
}


func (t *LunchServer) addPlace(name string) (place uint) {
	args := &StringArgs{String: name}
	args.Auth = *(t.calcAuth(args))
	err := t.Call("LunchTracker.AddPlace", &args, &place)
	if err != nil {
		panic(err)
	}
	return
}

func (t *LunchServer) comment(comment string) {
	args := &StringArgs{String: comment}
	args.Auth = *(t.calcAuth(args))
	var suc bool
	err := t.Call("LunchTracker.Comment", &args, &suc)
	if err != nil {
		panic(err)
	}

	if !suc {
		panic(commentError)
	}
	return
}

func (t *LunchServer) delPlace(dest uint) {
	args := &UIntArgs{Num: dest}
	args.Auth = *(t.calcAuth(args))
	var suc bool
	err := t.Call("LunchTracker.DelPlace", args, &suc)
	if err != nil {
		panic(err)
	}
	if !suc {
		panic(delPlaceError)
	}
	return
}

func (t *LunchServer) drive(seats uint) {
	args := &UIntArgs{Num: seats}
	args.Auth = *(t.calcAuth(args))
	var suc bool
	err := t.Call("LunchTracker.Drive", args, &suc)
	if err != nil {
		panic(err)
	}
	if !suc {
		panic(driveError)
	}
	return
}

func (t *LunchServer) vote(dest uint) {
	args := &UIntArgs{Num: dest}
	args.Auth = *(t.calcAuth(args))
	var suc bool
	err := t.Call("LunchTracker.Vote", args, &suc)
	if err != nil {
		panic(err)
	}
	if !suc {
		panic(voteError)
	}
	return
}

func (t *LunchServer) unvote() {
	args := &EmptyArgs{}
	args.Auth = *(t.calcAuth(args))
	var suc bool
	err := t.Call("LunchTracker.UnVote", args, &suc)
	if err != nil {
		panic(err)
	}
	if !suc {
		panic(unvoteError)
	}
	return
}

func (t *LunchServer) displayPlaces() *LunchPoll {
	args := &EmptyArgs{}
	args.Auth = *(t.calcAuth(args))
	var poll LunchPoll
	err := t.Call("LunchTracker.DisplayPlaces", args, &poll)
	if err != nil && err != os.EOF {
		panic(err)
	}
	return &poll
}

func (t *LunchServer) undrive() {
	args := &EmptyArgs{}
	args.Auth = *(t.calcAuth(args))
	var suc bool
	err := t.Call("LunchTracker.UnDrive", args, &suc)
	if err != nil {
		panic(err)
	}
	if !suc {
		panic(undriveError)
	}
	return
}

func (t *LunchServer) calcAuth(d Byter) (a *Auth) {
	var challenge *Bin
	err := t.Call("LunchTracker.Challenge", &user, &challenge)
	if err != nil {
		panic(err)
	}
	a = &Auth{Name: user, SChallenge: challenge}
	sum(d, a)
	return

}

func getConfig() (err os.Error) {
	config := make(map[string]string)
	home := os.Getenv("HOME")

	var read []byte
	if read, err = ioutil.ReadFile(path.Join(home, ".lunch", configFile)); err != nil {
		return
	}

	if err = json.Unmarshal(read, &config); err != nil {
		return
	}

	var ok bool
	if user, ok = config["user"]; !ok {
		return os.NewError(noUser)
	}
	if sekrit, ok = config["sekrit"]; !ok {
		return os.NewError(noSekrit)
	}

	if host, ok = config["host"]; !ok {
		return os.NewError(noHost)
	}

	return
}

func genConfig() (err os.Error) {
	config := make(map[string]string)
	config["host"] = *server
	config["user"] = *name
	config["sekrit"] = string(*makeSekrit())

	switch {
	case config["sekrit"] == "":
		err = os.NewError(noSekrit)
	case config["user"] == "":
		err = os.NewError(noUser)
	case config["host"] == "":
		err = os.NewError(noHost)
	}

	if err != nil {
		return
	}

	data, err := json.Marshal(config)
	if err != nil {
		return
	}

	user = config["user"]
	sekrit = config["sekrit"]
	host = config["host"]
	home := os.Getenv("HOME")

	confPath := path.Join(home, ".lunch")
	if err = os.MkdirAll(confPath, 0700); err != nil {
		return
	}
	err = ioutil.WriteFile(path.Join(confPath, configFile), data, 0600)
	return
}

func makeSekrit() *[]byte {
	newSekrit := make([]byte, 512, 512)
	encoder := base64.StdEncoding

	_, err := rand.Read(newSekrit)
	if err != nil {
		panic(badRandom)
	}
	newEncoded := make([]byte, encoder.EncodedLen(len(newSekrit)), encoder.EncodedLen(len(newSekrit)))
	encoder.Encode(newEncoded, newSekrit)
	return &newEncoded
}

func sum(d Byter, a *Auth) {

	challenge := make(Bin, 512)
	_, err := rand.Read(challenge)

	if err != nil {
		panic(badRandom)
	}
	(*a).CChallenge = &challenge

	mac := hmac.New(sha512.New, []byte(sekrit))
	mac.Write([]byte((*a).Name))
	mac.Write(*(*a).CChallenge)
	mac.Write(d.Byte())
	mac.Write(*(*a).SChallenge)
	bin := Bin(mac.Sum())
	(*a).Mac = &bin
}
