package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"
)

var start = time.Now()

func main() {
	numPhilosophers, timeToEat, timeToSleep, timeToDie, totalEatCountToStopGame, error := getFlags()

	if error != nil {
		fmt.Println(error)
		return
	}

	// init forks
	forks := make([]Fork, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		forks[i] = Fork{
			isUsed: false,
		}
	}
	// init philosophers
	philosophers := make([]Philosopher, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = Philosopher{
			name:          strconv.Itoa(i),
			leftFork:      &forks[i],
			rightFork:     &forks[(i+1)%numPhilosophers],
			gameStartTime: start,
		}
	}

	// use index instead of range value, because it will create a copy of value here which wanted to avoid.
	for i := range philosophers {
		go philosophers[i].StartRoutine(timeToEat, timeToSleep)
	}

	isGameStopped := false
	for !isGameStopped {
		totalEatCount := 0
		// loop to check status
		for i := range philosophers {
			// check if any philosophers died
			isGameStopped = checkPhiloDead(philosophers, i, timeToDie)
			if isGameStopped {
				break
			}
			philosophers[i].lock.Lock()
			totalEatCount += philosophers[i].eatCount
			philosophers[i].lock.Unlock()
		}

		if totalEatCount >= totalEatCountToStopGame {
			isGameStopped = true
		}
	}
}

// checkPhiloDead return true
func checkPhiloDead(philosophers []Philosopher, i int, timeToDie int64) bool {
	philosophers[i].lock.Lock()
	lastMealDuration := time.Since(philosophers[i].lastMealTime)
	if timeToDie >= lastMealDuration.Milliseconds() {
		philosophers[i].lock.Unlock()
		philosophers[i].killSignal <- true
		// stop all other thread
		for k := range philosophers {
			if k != i {
				philosophers[i].stopSignal <- true
			}
		}
		return true
	}
	philosophers[i].lock.Unlock()
	return false
}

func getFlags() (int, int64, int64, int64, int, error) {
	numPhilosophers := flag.Int("number_of_philosophers", 5, "number of philosophers")
	timeToEat := flag.Int64("time_to_eat", 1000, "time to eat in milliseconds")
	timeToSleep := flag.Int64("time_to_sleep", 1000, "time to fallSleep in milliseconds")
	timeTimeToDie := flag.Int64("time_to_die", 1000, "time to die in milliseconds")
	numberOfTimesToEat := flag.Int("number_of_times_to_eat", 0, "number of times to eat")
	flag.Parse()

	if *numberOfTimesToEat == 0 {
		return 0, 0, 0, 0, 0, errors.New("number of times to eat must be greater than 0")
	}

	if *timeToEat == 0 {
		return 0, 0, 0, 0, 0, errors.New("time to eat must be greater than 0")
	}

	if *timeToSleep == 0 {
		return 0, 0, 0, 0, 0, errors.New("time to fallSleep must be greater than 0")
	}

	if *timeTimeToDie == 0 {
		return 0, 0, 0, 0, 0, errors.New("time to die must be greater than 0")
	}

	return *numPhilosophers, *timeToEat, *timeToSleep, *timeTimeToDie, *numberOfTimesToEat, nil
}
