package repository

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
)

//Persistance must persist using mongodb
type Persistance interface {
	Persist(*mgo.Session)
}

//StartPersistance initialize the loop and make some implementation
type StartPersistance interface {
	Do()
}

//RunPersistance start loop
func RunPersistance(persistInterval int, startPersistance StartPersistance) {
	log.Printf("repository - Starting persistance each %d second(s)\n", persistInterval)
	go func() {
		ticker := time.NewTicker(time.Duration(persistInterval) * time.Second)
		for range ticker.C {
			startPersistance.Do()
		}
	}()
}

// PersistIncrementation persist the incrementations
func PersistIncrementation(p Persistance) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Println(err)
		return
	}
	go func(s *mgo.Session) {
		defer s.Close()
		s.SetMode(mgo.Monotonic, true)
		p.Persist(s)
	}(session)
}
