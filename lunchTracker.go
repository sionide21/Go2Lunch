package main

import "os"

type LunchTracker chan LunchPoll

func (t *LunchTracker) AddPlace(args *AddPlaceArgs, place *uint) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	poll := t.getPoll()
	*place = poll.addPlace(args.Name, args.Auth.Name)
	t.putPoll(poll)
	return nil
}

func (t *LunchTracker) DelPlace(args *UIntArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	poll := t.getPoll()
	*success = poll.delPlace(args.Num)
	t.putPoll(poll)
	return nil
}

func (t *LunchTracker) Drive(args *UIntArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	poll := t.getPoll()
	*success = poll.drive(args.Auth.Name, args.Num)
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
	t.putPoll(poll)
	return nil
}

func (t *LunchTracker) Vote(args *UIntArgs, success *bool) os.Error {
	valid, ive := verify(&args.Auth, args)
	if !valid {
		return ive
	}
	poll := t.getPoll()
	*success = poll.vote(args.Auth.Name, args.Num)
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

func (t *LunchTracker) getPoll() LunchPoll {
	return <- *t
}

func (t *LunchTracker) putPoll(l LunchPoll) {
	*t <- l
}