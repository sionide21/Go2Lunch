package main

import (
	"gob"
	"strings"
	"regexp"
)

type LunchPoll struct {
	Places       PlaceVector
	IndexCounter int
	Votes        map[string]*Place
}

var placeRegex = regexp.MustCompile("[^A-Z0-9]")

func init() {
	RegisterTypes()
}

func NewPoll() LunchPoll {
	poll := LunchPoll{
		Places:       make(PlaceVector, 0),
		Votes:        make(map[string]*Place),
		IndexCounter: 1}

	defaultPlace := &Place{
		Id:        0,
		Name:      "No Where",
		Votes:     0,
		People:    make(PersonVector, 0),
		Nominator: nil}

	poll.Places.Push(defaultPlace)
	return poll
}

func (p *LunchPoll) addPlace(name, nominator string) (count int, success bool) {
	count = p.IndexCounter
	if name == "" || p.placeExists(name) {
		return
	}

	person := p.getPerson(nominator)
	if person.NominationsLeft > 0 {
		place := &Place{
			Id:        p.IndexCounter,
			Nominator: person,
			Name:      name,
			People:    make(PersonVector, 0)}
		p.Places.Push(place)
		p.IndexCounter++
		person.NominationsLeft--
		success = true
	}
	return
}

func (p *LunchPoll) delPlace(placeId int, person string) bool {
	place, ok := p.getPlace(placeId)
	if ok {
		if place.Votes == 0 && place.Nominator.Name == person {
			if p.remove(place) {
				place.Nominator.NominationsLeft++
				return true
			}
		}
	}
	return false
}

func (p *LunchPoll) drive(who string, seats int) bool {
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

func (p *LunchPoll) vote(who string, vote int) bool {
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
			people := p.Places.At(0).People
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
		for _, peep := range place.People {
			if peep.Name == name {
				return peep
			}
		}
	}
	person := &Person{
		CanDrive:        false,
		Name:            name,
		NominationsLeft: 2}
	defaultPeopleVector := p.Places.At(0).People
	defaultPeopleVector.Push(person)
	return person
}

func (p *LunchPoll) getPlace(dest int) (*Place, bool) {
	for _, pl := range p.Places {
		if pl.Id == dest {
			return pl, true
		}
	}
	return nil, false
}

func (p *LunchPoll) remove(sp *Place) bool {
	for i, place := range p.Places {
		if place == nil {
			return false
		}

		if place.Id == sp.Id {
			p.Places.Delete(i)
			return true
		}
	}
	return false
}

func sanitizePlace(name string) string {
	return placeRegex.ReplaceAllString(strings.ToUpper(name), "")
}

func (p *LunchPoll) placeExists(name string) bool {
	check := sanitizePlace(name)
	for _, place := range p.Places {
		if check == sanitizePlace(place.Name) {
			return true
		}
	}
	return false
}

func RegisterTypes() {
	gob.Register(make(PlaceVector, 0))
	gob.Register(make(PersonVector, 0))
	gob.Register(
		LunchPoll{
			Places:       make(PlaceVector, 0),
			Votes:        make(map[string]*Place),
			IndexCounter: 1})
	gob.Register(
		Place{
			Id:     0,
			Name:   "No Where",
			Votes:  0,
			People: make(PersonVector, 0),
			Nominator: &Person{
				CanDrive:        false,
				Name:            "",
				NominationsLeft: 2}})
	gob.Register(
		Person{
			CanDrive:        false,
			Name:            "",
			NominationsLeft: 2})
}
