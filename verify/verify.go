package main

import (
	"fmt"

	"killyosaur.github.io/barriers"
)

var num_steps int = 100000000
var step float64 = 1.0 / float64(num_steps)
var pi_1 float64 = 0.0
var pi_2 float64 = 0.0

func calculate_pi(node *barriers.NodeCollection, rank int, sums [10]float64, done chan bool) {
	sums[rank] = 0.0
	fmt.Printf("%v Rank: current sums = %v, pre-barrier 1\n", rank, sums[rank])

	for i := 1; i < num_steps; i += 10 {
		x := (float64(i) - .5) * step
		sums[rank] += 4.0 / (1.0 + x*x)
	}

	fmt.Printf("%v Rank: current sums = %v, pre-barrier 2\n", rank, sums[rank])
	node.Barrier(rank)

	if rank == 0 {
		for i := 0; i < 10; i++ {
			pi_1 += sums[i]
		}
	}
	fmt.Printf("%v Rank: pi_1 = %v, first barrier\n", rank, pi_1)

	node.Barrier(rank)

	sums[rank] = 0.0

	for i := 1; i < num_steps; i += 10 {
		x := (float64(i) - .5) * step
		sums[rank] += 4.0 / (1.0 + x*x)
	}
	fmt.Printf("%v Rank: current sums = %v, second barrier\n", rank, sums[rank])

	node.Barrier(rank)

	if rank == 0 {
		for i := 0; i < 10; i++ {
			pi_2 += sums[i]
		}
	}
	fmt.Printf("%v Rank: pi_2 = %v, third barrier\n", rank, pi_2)

	done <- true
}

func main() {
	var sums [10]float64

	nodes := barriers.Init(10)

	channels := make([]chan bool, 10)

	for i := 0; i < 10; i++ {
		fmt.Printf("goroutine %v\n", i)
		channels[i] = make(chan bool)
		go calculate_pi(nodes, i, sums, channels[i])
	}

	pi_2 /= float64(num_steps)
	pi_1 /= float64(num_steps)

	<-channels[0]
	<-channels[1]
	<-channels[2]
	<-channels[3]
	<-channels[4]
	<-channels[5]
	<-channels[6]
	<-channels[7]
	<-channels[8]
	<-channels[9]

	fmt.Printf("%v == %v", pi_1, pi_2)
}
