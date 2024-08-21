package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type Candidate struct {
	Name  string
	Votes int
}

var candidates = []Candidate{
	{"Candidate A", 0},
	{"Candidate B", 0},
	{"Candidate C", 0},
}

var mutex sync.Mutex

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read candidate number from the connection
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	candidateNumber := strings.TrimSpace(string(buf[:n]))

	// Cast the vote for the candidate
	castVote(candidateNumber)

	index := parseCandidateNumber(candidateNumber)
	if index != -1 && index < len(candidates) {
		fmt.Printf("Vote cast for candidate %s\n", candidates[index].Name)
	}
}

func castVote(candidateNumber string) {
	mutex.Lock()
	defer mutex.Unlock()

	index := parseCandidateNumber(candidateNumber)
	if index != -1 && index < len(candidates) {
		candidates[index].Votes++
	}
}

func parseCandidateNumber(candidateNumber string) int {
	// Convert the candidate number to an integer
	num := int(candidateNumber[0] - '0')

	// Adjust the number to be zero-based index
	index := num - 1

	if index >= 0 && index < len(candidates) {
		return index
	}

	return -1
}

func main() {
	// Listen for interrupt signal (SIGINT)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	go func() {
		<-sigChan

		// Print total votes per candidate
		fmt.Println("Total votes per candidate:")
		allVotesZero := true
		votesMap := make(map[int]int)
		for _, candidate := range candidates {
			fmt.Printf("%s: %d\n", candidate.Name, candidate.Votes)
			if candidate.Votes > 0 {
				allVotesZero = false
			}
			votesMap[candidate.Votes]++
		}

		if allVotesZero {
			fmt.Println("No votes submitted")
		} else {
			isDraw := false
			for _, count := range votesMap {
				if count > 1 {
					isDraw = true
					break
				}
			}
			if isDraw {
				fmt.Println("It's a draw!")
			} else {
				maxVotes := 0
				var winner string
				for _, candidate := range candidates {
					if candidate.Votes > maxVotes {
						maxVotes = candidate.Votes
						winner = candidate.Name
					}
				}
				fmt.Printf("The winner is: %s\n", winner)
			}
		}
		os.Exit(0)
	}()

	// Listen for incoming connections
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()

	fmt.Println("Server listening on port 8080")

	for {
		// Accept new connections
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}
