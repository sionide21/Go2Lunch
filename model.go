package main

import (
	"container/list"
)

type ServerPlace struct {
	Place
	PeopleList *list.List
}

func (t *ServerPlace) remove(p string) bool {
	person := t.PeopleList.Front()

	for {
		if person == nil {
			return false
		}

		pl, _ := person.Value.(*Person)
		if pl.Name == p {
			t.PeopleList.Remove(person)
			return true
		}

		person = person.Next()
	}
	return false
}

type LunchPoll struct {
	places       *list.List
	people       *list.List
	indexCounter uint
	votes        map[string]*ServerPlace
}

func NewPoll() *LunchPoll {
	return &LunchPoll{
		places:       list.New(),
		people:       list.New(),
		votes:        make(map[string]*ServerPlace),
		indexCounter: 1}
}

func (p *LunchPoll) addPlace(name, nominator string) uint {
	person := p.getPerson(nominator)
	if person.NominationsLeft > 0 {
		place := &ServerPlace{
			Place{
				Id:        p.indexCounter,
				Nominator: person,
				Name:      name},
			list.New()}

		p.places.PushBack(place)
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
			place.PeopleList.PushBack(person)
			return true
		}
	}
	return false
}

func (p *LunchPoll) unVote(who string) bool {
	if place, voted := p.votes[who]; voted {
		place.Votes--
		p.votes[who] = nil, false
		place.remove(who)
		return true
	}
	return false
}

func (p *LunchPoll) displayPlaces() []Place {
	ret := make([]Place, p.places.Len())
	var i = 0
	for place := range p.places.Iter() {
		ret[i] = flattenPlace(place.(*ServerPlace))
		i++
	}
	return ret
}

// Helpers
func (p *LunchPoll) getPerson(name string) *Person {
	for p := range p.people.Iter() {
		person, _ := p.(*Person)
		if person.Name == name {
			return person
		}
	}
	person := &Person{
		CanDrive:        false,
		Name:            name,
		NominationsLeft: 2}
	p.people.PushBack(person)
	return person
}

func (p *LunchPoll) getPlace(dest uint) (*ServerPlace, bool) {
	for pl := range p.places.Iter() {
		place, _ := pl.(*ServerPlace)
		if place.Id == dest {
			return place, true
		}
	}
	return nil, false
}

func (p *LunchPoll) getDriver(driverId uint) (*Person, bool) {
	for per := range p.people.Iter() {
		person, ok := per.(*Person)
		if ok {
			return person, true
		}
	}
	return nil, false
}

func (p *LunchPoll) remove(sp *ServerPlace) bool {
	place := p.places.Front()
	for {
		if place == nil {
			return false
		}

		pl, _ := place.Value.(*ServerPlace)
		if pl.Id == sp.Id {
			p.places.Remove(place)
			return true
		}

		place = place.Next()
	}
	return false
}

func flattenPlace(server *ServerPlace) (place Place) {
	peeps := make([]*Person, server.PeopleList.Len())
	place = Place{
		Id:        server.Id,
		Name:      server.Name,
		Votes:     server.Votes,
		Nominator: server.Nominator,
		People:    peeps}
	var i = 0
	for sperson := range server.PeopleList.Iter() {
		place.People[i] = sperson.(*Person)
		i++
	}
	return
}
