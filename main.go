package main

import (
	"errors"
	"flag"
	"fmt"
)

func main() {
	numPhilosophers, timeToEat, timeToSleep, timeTimeToDie, numberOfTimesToEat, error := getFlags()

	if error != nil {
		fmt.Println(error)
		return
	}

}

func createPhilosophers(numPhilosophers int) []Philosopher {
	philosophers := make([]Philosopher, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = Philosopher{
			name: fmt.Sprintf("Philosopher %d", i),
			leftFork: i,
			rightFork: (i + 1) % numPhilosophers,
			state: "thinking"
		}
	}
	return philosophers
}

func getFlags() (int, uint, uint, uint, uint, error) {
	numPhilosophers := flag.Int("number_of_philosophers", 5, "number of philosophers")
	timeToEat := flag.Uint("time_to_eat", 1000, "time to eat in milliseconds")
	timeToSleep := flag.Uint("time_to_sleep", 1000, "time to sleep in milliseconds")
	timeTimeToDie := flag.Uint("time_to_die", 1000, "time to die in milliseconds")
	numberOfTimesToEat := flag.Uint("number_of_times_to_eat", 0, "number of times to eat")
	flag.Parse()

	if *numberOfTimesToEat == 0 {
		return 0, 0, 0, 0, 0, errors.New("number of times to eat must be greater than 0")
	}

	if *timeToEat == 0 {
		return 0, 0, 0, 0, 0, errors.New("time to eat must be greater than 0")
	}

	if *timeToSleep == 0 {
		return 0, 0, 0, 0, 0, errors.New("time to sleep must be greater than 0")
	}

	if *timeTimeToDie == 0 {
		return 0, 0, 0, 0, 0, errors.New("time to die must be greater than 0")
	}

	return *numPhilosophers, *timeToEat, *timeToSleep, *timeTimeToDie, *numberOfTimesToEat, nil
}
