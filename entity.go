package main

import (
	"fmt"
	"sync"
	"time"
)

type Fork struct {
	isUsed bool
	lock   sync.Mutex
}
type Philosopher struct {
	name         string
	leftFork     *Fork
	rightFork    *Fork
	eatCount     int
	lastMealTime time.Time
	lock         sync.Mutex
	// might not OO enough to include the following fields here
	gameStartTime time.Time
	stopSignal    chan bool
	killSignal    chan bool
}

// StartRoutine this can accept the timeout config here, time to eat/fallSleep/die
func (p *Philosopher) StartRoutine(timeToEat, timeToSleep int64) {
	for {
		select {
		case <-p.killSignal:
			p.printLog("died")
			return
		case <-p.stopSignal:
			// end game here, either other died or game ended gracefully
			return
		default:
			isStateChanged := p.EatAndSleep(timeToEat, timeToSleep)
			if isStateChanged {
				p.printLog("is thinking")
			}
		}
	}
}

// TryToEat return true if philo can eat, false if fail to eat
func (p *Philosopher) TryToEat() bool {
	p.leftFork.lock.Lock()
	p.rightFork.lock.Lock()
	defer p.leftFork.lock.Unlock()
	defer p.rightFork.lock.Unlock()

	if p.leftFork.isUsed == false && p.rightFork.isUsed == false {
		p.leftFork.isUsed = true
		p.printLog("has taken a fork")
		p.rightFork.isUsed = true
		p.printLog("has taken a fork")
		return true
	}
	return false
}

func (p *Philosopher) releaseForks() {
	p.leftFork.lock.Lock()
	p.rightFork.lock.Lock()
	defer p.leftFork.lock.Unlock()
	defer p.rightFork.lock.Unlock()
	p.leftFork.isUsed = false
	p.rightFork.isUsed = false
}

// EatAndSleep return true if it could eat and it fall asleep later
func (p *Philosopher) EatAndSleep(timeToEat, timeToSleep int64) bool {
	if p.TryToEat() {
		p.lock.Lock()
		p.eatCount++
		p.printLog("is eating")
		p.lastMealTime = time.Now()
		p.lock.Unlock()
		time.Sleep(time.Duration(timeToEat) * time.Millisecond)
		p.releaseForks()
		// after eat will fall asleep immediately
		p.fallSleep(timeToSleep)
		return true
	}
	return false
}

func (p *Philosopher) fallSleep(timeToSleep int64) {
	// todo mutex lock Philo here
	p.printLog("is sleeping")
	time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
}

// if change to channel, would it create latency?
// todo might need to change to a commmon channel to avoid race condition
// common channel might need a new struct, or it could be part of select / loop (hmm)
func (p *Philosopher) printLog(message string) {
	eclipsedTime := time.Since(p.gameStartTime)
	fmt.Printf("%d %s %s", eclipsedTime.Milliseconds(), p.name, message)
}
