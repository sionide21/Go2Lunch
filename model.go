package main

import (
	"container/vector"
	"os"
	"fmt"
	"json"
)

type ServerPlace struct {
	Place
	PeopleList *vector.Vector
}

func (t *ServerPlace) remove(p string) bool {	
	for i, person := range *t.PeopleList {
		pl, _ := person.(*Person)
		if pl.Name == p {
			t.PeopleList.Delete(i)
			return true
		}
	}
	
	return false
}

type LunchPoll struct {
	places       *vector.Vector
	people       *vector.Vector
	indexCounter uint
	votes        map[string]*ServerPlace
}

func NewPoll() *LunchPoll {
	places, people := make(vector.Vector, 0), make(vector.Vector, 0)
	return &LunchPoll{
		places:       &people,
		people:       &places,
		votes:        make(map[string]*ServerPlace),
		indexCounter: 1}
}

func (p *LunchPoll) addPlace(name, nominator string) uint {
	person := p.getPerson(nominator)
	if person.NominationsLeft > 0 {
		people := make(vector.Vector, 0)
		place := &ServerPlace{
			Place{
				Id:        p.indexCounter,
				Nominator: person,
				Name:      name},
			&people}

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
			place.PeopleList.Push(person)
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
	for i, place := range *p.places {
		ret[i] = flattenPlace(place.(*ServerPlace))
	}
	return ret
}

// Helpers
func (p *LunchPoll) getPerson(name string) *Person {
	for _, p := range *p.people {
		person, _ := p.(*Person)
		if person.Name == name {
			return person
		}
	}
	person := &Person{
		CanDrive:        false,
		Name:            name,
		NominationsLeft: 2}
	p.people.Push(person)
	return person
}

func (p *LunchPoll) getPlace(dest uint) (*ServerPlace, bool) {
	for _, pl := range *p.places {
		place, _ := pl.(*ServerPlace)
		if place.Id == dest {
			return place, true
		}
	}
	return nil, false
}

func (p *LunchPoll) remove(sp *ServerPlace) bool {
	for i, place := range *p.places {
		pl, _ := place.(*ServerPlace)
		if pl.Id == sp.Id {
			p.places.Delete(i)
			return true
		}
	}
	
	return false
}

// func (p *LunchPoll) MarshalJSON() ([]byte, os.Error) {
// 	return json.Marshal(map[string]interface{}{
// 		"places":p.displayPlaces(),
// 		"indexCounter":p.indexCounter})
// }
func (p *LunchPoll) UnmarshalJSON(data []byte) os.Error {
	poll := make(map[string]interface{})
	err := json.Unmarshal(data, &poll)
	fmt.Println("JSON:", poll)
	return err
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
	for _, sperson := range *server.PeopleList {
		place.People[i] = sperson.(*Person)
		i++
	}
	return
}
