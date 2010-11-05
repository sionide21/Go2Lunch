package main
import (
	"container/list"
)
type LunchPoll struct {
	places       *vector.Vector
	people       *vector.Vector
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
	return places
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
	p.people.Push(person)
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
