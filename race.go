package greatspacerace

import (
	"github.com/TSavo/chipmunk"
	"github.com/TSavo/chipmunk/vect"
	"github.com/chourobin/go.firebase"
	"github.com/nu7hatch/gouuid"
	"math"
	"runtime"
	"time"
)

var db = firebase.New("https://greatspacerace.firebaseio.com/")

type Race struct {
	Id      string
	Track   *Track
	Racers  []*PlayerPosition
	Space   *chipmunk.Space
	Started bool
	Ticks   int64
}

func NewRace(track *Track) *Race {
	id, _ := uuid.NewV4()
	race := &Race{id.String(), track, make([]*PlayerPosition, 0), chipmunk.NewSpace(), false, 0}
	staticBody := chipmunk.NewBodyStatic()
	for _, wall := range track.Segments {
		segment := chipmunk.NewSegment(wall.Point1, wall.Point2, 0)
		staticBody.AddShape(segment)
	}
	race.Space.AddBody(staticBody)
	return race
}

type Location struct {
	Position vect.Vect
	Angle    vect.Float
	Tick     int64
}

type PlayerPosition struct {
	Name                 string
	Dimensions, Position vect.Vect
	Angle                vect.Float
	Checkpoint           int
	LapTimes             []int64
	Finished             bool
	player               *Player
	messages             []*Message
	locations            []*Location
}

func (this *Race) RunRace() {
	this.StartRace()
outer:
	for {
		raceUpdate := Message{"RaceUpdate", this.Racers, this.Ticks}
		for _, racer := range this.Racers {
			racer.player.SendChan <- &raceUpdate
		}
		runtime.Gosched()
		finished := make(chan interface{}, len(this.Racers))
		for _, racer := range this.Racers {
			if racer.Finished {
				finished <- nil
				continue
			}
			go func(racer *PlayerPosition) {
				timeOut := time.After(time.Second)
			control:
				for {
					select {
					case x := <-racer.player.ReceiveChan:
						if x.Tick > 0 && x.Tick != this.Ticks {
							continue control
						}
						racer.messages = append(racer.messages, x)
						if x.Type == "Control" {
							controlMessage := x.Data.(map[string]interface{})
							racer.player.Ship.Controller.Thrust = vect.Float(controlMessage["Thrust"].(float64))
							racer.player.Ship.Controller.Turning = vect.Float(controlMessage["Rotation"].(float64))
						}
						break control
					case <-timeOut:
						break control
					}
				}
				racer.player.Ship.ApplyThrust(racer.player.Ship.Controller.Thrust)
				racer.player.Ship.ApplyRotation(racer.player.Ship.Controller.Turning)
				finished <- true
			}(racer)
		}
		for _ = range this.Racers {
			<-finished
		}
		this.MoveShips()
		for _, racer := range this.Racers {
			if !racer.Finished {
				continue outer
			}
		}
		break
	}
	results := make(map[string]interface{})
	racers := make([]map[string]interface{}, len(this.Racers))
	results["Id"] = this.Id
	results["Date"] = time.Now().String()
	results["Track"] = this.Track
	for x, racer := range this.Racers {
		racers[x]["LapTimes"] = racer.LapTimes
		racers[x]["Messages"] = racer.messages
		racers[x]["Name"] = racer.Name
		racers[x]["Locations"] = racer.locations
		racer.player.SendChan <- &Message{"RaceFinished", racer.LapTimes, -1}
	}
	results["Racers"] = racers
	db.Set("/"+this.Id, results)
	return
}

func (this *Race) MoveShips() {
	this.StepRace(1.0 / 60.0)
	for _, racer := range this.Racers {
		if racer.Finished {
			continue
		}
		motion := Segment{racer.Position, racer.player.Ship.Body.Position()}
		racer.Position = racer.player.Ship.Body.Position()
		racer.Angle = racer.player.Ship.Body.Angle()
		racer.locations = append(racer.locations, &Location{racer.Position, racer.Angle, this.Ticks})
		if racer.Checkpoint == len(this.Track.Checkpoints) {
			if motion.Intersects(&this.Track.Goal) {
				racer.Checkpoint = 0
				elapsed := int64(0)
				for _, x := range racer.LapTimes {
					elapsed += x
				}
				racer.LapTimes = append(racer.LapTimes, this.Ticks-elapsed)
				racer.player.SendChan <- &Message{"LapFinished", this.Ticks - elapsed, -1}
				if len(racer.LapTimes) == this.Track.Laps {
					this.Space.RemoveBody(racer.player.Ship.Body)
					racer.Finished = true
					racer.player.SendChan <- &Message{"YourRaceFinished", racer.LapTimes, -1}
				}
			}
		} else {
			if motion.Intersects(&this.Track.Checkpoints[racer.Checkpoint]) {
				racer.Checkpoint++
			}
		}
	}
	this.Ticks++
}

func (this *Race) StartRace() {
	this.Started = true
	startPoints := this.Track.Goal.GetStartingPositions(len(this.Racers))
	for x, racer := range this.Racers {
		racer.player.Ship.Body.SetPosition(startPoints[x])
		racer.Position = startPoints[x]
		racer.player.Ship.Body.SetAngle(vect.Float(this.Track.StartingAngle * 2 * math.Pi))
		racer.Angle = vect.Float(this.Track.StartingAngle * 2 * math.Pi)
		this.Space.AddBody(racer.player.Ship.Body)
	}
}

func (this *Race) RegisterRacer(player *Player) {
	box := chipmunk.NewBox(vect.Vector_Zero, player.Ship.Prototype.Width, player.Ship.Prototype.Height)
	box.SetElasticity(0.9)
	body := chipmunk.NewBody(player.Ship.Prototype.Mass, box.Moment(player.Ship.Prototype.Moment))
	body.AddShape(box)
	player.Ship.Body = body
	this.Racers = append(this.Racers, &PlayerPosition{player.Name, vect.Vect{player.Ship.Prototype.Width, player.Ship.Prototype.Height}, vect.Vector_Zero, 0, 0, make([]int64, 0), false, player, make([]*Message, 0), make([]*Location, 0)})
}

func (this *Race) StepRace(dt vect.Float) {
	this.Space.Step(dt)
}
