package main

import (
	"container/vector"
)

type LunchPoll struct {
	Places       vector.Vector
	IndexCounter uint
	Votes        map[string]*Place
}

func NewPoll() LunchPoll {
	poll := LunchPoll{
		Places:       make(vector.Vector, 5),
		Votes:        make(map[string]*Place),
		IndexCounter: 1}

	defaultPlace := Place{
		Id:        0,
		Name:      "No Where",
		Votes:     0,
		People:    make(vector.Vector, 5),
		Nominator: nil}

	poll.Places.Push(defaultPlace)
	return poll
}

func (p *LunchPoll) addPlace(name, nominator string) uint {
	person := p.getPerson(nominator)
	if person.NominationsLeft > 0 {
		place := &Place{
			Id:        p.IndexCounter,
			Nominator: person,
			Name:      name,
			People:    make(vector.Vector, 3)}
		p.Places.Push(place)
		defer func() { p.IndexCounter++ }()
		person.NominationsLeft--
	}
	return p.IndexCounter
}

func (p *LunchPoll) delPlace(placeId uint) bool {
	place, ok := p.getPlace(placeId)
	if ok {
		if place.Votes == 0 {
			if p.remove(place) {
				place.Nominator.NominationsLeft++
				return true
			}
		}
	}
	return false
}

func (p *LunchPoll) drive(who string, seats uint) bool {
	person := p.getPerson(who)
	person.CanDrive = true
	person.NumSeats = seats
	return true
}

func (p *LunchPoll) unDrive(who string) bool {
	person := p.getPerson(who)
	person.CanDrive = false
	person.NumSeats = 0
	return true
}

func (p *LunchPoll) vote(who string, vote uint) bool {
	place, ok := p.getPlace(vote)
	if ok {
		person := p.getPerson(who)
		if p.Votes[who] == nil {
			p.Votes[who] = place
			place.Votes++
			place.People.Push(person)
			return true
		}
	}
	return false
}

func (p *LunchPoll) unVote(who string) bool {
	if place, voted := p.Votes[who]; voted {
		place.Votes--
		p.Votes[who] = nil, false
		person := place.RemovePerson(who)
		if person.Name != "" {
			people := p.Places.At(0).(Place).People
			people.Push(person)
			return true
		}
	}
	return false
}

// Helpers
func (p *LunchPoll) getPerson(name string) *Person {
	for _, place := range p.Places {
		if place == nil {
			continue
		}
		for _, peep := range place.(Place).People {
			person, ok := peep.(*Person)
			if ok && person.Name == name {
				return person
			}
		}
	}
	person := &Person{
		CanDrive:        false,
		Name:            name,
		NominationsLeft: 2}
	defaultPeopleVector := p.Places.At(0).(Place).People
	defaultPeopleVector.Push(person)
	return person
}

func (p *LunchPoll) getPlace(dest uint) (*Place, bool) {
	for _, pl := range p.Places {
		place, _ := pl.(*Place)
		if place.Id == dest {
			return place, true
		}
	}
	return nil, false
}

func (p *LunchPoll) remove(sp *Place) bool {
	for i, place := range p.Places {
		if place == nil {
			return false
		}

		pl, ok := place.(*Place)
		if ok && pl.Id == sp.Id {
			p.Places.Delete(i)
			return true
		}
	}
	return false
}
