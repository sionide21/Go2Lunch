package main

import (
	"os"
	"crypto/rand"
	"gob"
)

type LunchTracker chan LunchPoll

func (t *LunchTracker) AddPlace(args *StringArgs, place *int) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	var ok bool
	poll := t.getPoll()
	*place, ok = poll.addPlace(args.String, args.Auth.Name)
	if !ok {
		*place = -1
	}
	t.persist(poll)
	t.putPoll(poll)
	return nil
}

func (t *LunchTracker) Comment(args *StringArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	poll := t.getPoll()
	*success = poll.comment(args.String, args.Auth.Name)
	t.persist(poll)
	t.putPoll(poll)
	return nil
}

func (t *LunchTracker) DelPlace(args *IntArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	poll := t.getPoll()
	*success = poll.delPlace(args.Num, args.Auth.Name)
	t.persist(poll)
	t.putPoll(poll)
	return nil
}

func (t *LunchTracker) Drive(args *IntArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	poll := t.getPoll()
	*success = poll.drive(args.Auth.Name, args.Num)
	t.persist(poll)
	t.putPoll(poll)
	return nil
}

func (t *LunchTracker) UnDrive(args *EmptyArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	poll := t.getPoll()
	*success = poll.unDrive(args.Auth.Name)
	t.persist(poll)
	t.putPoll(poll)
	return nil
}

func (t *LunchTracker) Vote(args *IntArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	poll := t.getPoll()
	*success = poll.vote(args.Auth.Name, args.Num)
	t.persist(poll)
	t.putPoll(poll)
	return nil
}

func (t *LunchTracker) UnVote(args *EmptyArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	poll := t.getPoll()
	*success = poll.unVote(args.Auth.Name)
	t.persist(poll)
	t.putPoll(poll)
	return nil
}

func (t *LunchTracker) DisplayPlaces(args *EmptyArgs, response *LunchPoll) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	*response = t.getPoll()
	t.putPoll(*response)
	return nil
}

func (t *LunchTracker) Challenge(name *string, challenge *Bin) os.Error {
	_, valid := userMap[(*name)]
	if !valid {
		valid = checkUser(*name)
		if !valid {
			return os.NewError("Unknown User")
		}
	}

	*challenge = make(Bin, 512)
	n, err := rand.Read(*challenge)

	if err != nil || n != 512 {
		panic("Challenge Generation Failed")
	}

	userMap[*name].SChallenge = challenge
	return nil
}

func (t *LunchTracker) getPoll() LunchPoll {
	return <-*t
}

func (t *LunchTracker) putPoll(l LunchPoll) {
	*t <- l
}

func (t *LunchTracker) persist(poll LunchPoll) {
	file, err := os.Open(*dataFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encode := gob.NewEncoder(file)
	encode.Encode(poll)
}
