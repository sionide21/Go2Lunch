package main

import (
	"container/vector"
)

type LunchPoll struct {
	places       vector.Vector
	indexCounter uint
	votes        map[string]*Place
}

func NewPoll() LunchPoll {
	poll := LunchPoll{
		places:       make(vector.Vector, 5),
		votes:        make(map[string]*Place),
		indexCounter: 1}

	defaultPlace := Place{
		Id:        0,
		Name:      "",
		Votes:     0,
		People:    make(vector.Vector, 5),
		Nominator: nil}

	poll.places.Push(defaultPlace)
	return poll
}

func (p *LunchPoll) addPlace(name, nominator string) uint {
	person := p.getPerson(nominator)
	if person.NominationsLeft > 0 {
		place := &Place{
			Id:        p.indexCounter,
			Nominator: person,
			Name:      name,
			People:    make(vector.Vector, 3)}
		p.places.Push(place)
		defer func() { p.indexCounter++ }()
		person.NominationsLeft--
	}
	return p.indexCounter
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
		if p.votes[who] == nil {
			p.votes[who] = place
			place.Votes++
			place.People.Push(person)
			return true
		}
	}
	return false
}

func (p *LunchPoll) unVote(who string) bool {
	if place, voted := p.votes[who]; voted {
		place.Votes--
		p.votes[who] = nil, false
		person := place.RemovePerson(who)
		if person.Name != "" {
			people := p.places.At(0).(Place).People
			people.Push(person)
			return true
		}
	}
	return false
}

// Helpers
func (p *LunchPoll) getPerson(name string) *Person {
	for _, place := range p.places {
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
	defaultPeopleVector := p.places.At(0).(Place).People
	defaultPeopleVector.Push(person)
	return person
}

func (p *LunchPoll) getPlace(dest uint) (*Place, bool) {
	for _, pl := range p.places {
		place, _ := pl.(*Place)
		if place.Id == dest {
			return place, true
		}
	}
	return nil, false
}

func (p *LunchPoll) remove(sp *Place) bool {
	for i, place := range p.places {
		if place == nil {
			return false
		}

		pl, ok := place.(*Place)
		if ok && pl.Id == sp.Id {
			p.places.Delete(i)
			return true
		}
	}
	return false
}
