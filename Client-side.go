package main

import (
    "fmt"
    "net"
    "os"
    "strconv"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: client [candidate_number]")
        return
    }

    // Parse the candidate number from the command line argument
    candidateNumber, err := strconv.Atoi(os.Args[1])
    if err != nil {
        fmt.Println("Invalid candidate number:", err)
        return
    }

    // Connect to the server
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error connecting:", err.Error())
        return
    }
    defer conn.Close()

    // Send the selected candidate number to the server
    _, err = conn.Write([]byte(fmt.Sprintf("%d", candidateNumber)))
    if err != nil {
        fmt.Println("Error writing:", err.Error())
        return
    }

    fmt.Printf("Vote cast for candidate %d\n", candidateNumber)
}
