package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type farm struct {
	ants_number int
	rooms       map[string][]int
	start       map[string][]int
	end         map[string][]int
	links       map[string][]string
}

func main() {
	var myFarm farm
	myFarm.Read("test.txt")
	BiBFS(&myFarm) // Passing pointer to farm to avoid copying it
	fmt.Println("number of ants is : ", myFarm.ants_number)
	fmt.Println("rooms are : ", myFarm.rooms)
	fmt.Println("start is : ", myFarm.start)
	fmt.Println("end is : ", myFarm.end)
	fmt.Println("links are : ", myFarm.links)
	fmt.Println("adjacent is : ", Graph(myFarm))
}

func (myFarm *farm) Read(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Println("error reading", err)
	}
	content := strings.Split(string(bytes), "\n")

	myFarm.rooms = make(map[string][]int)
	myFarm.start = make(map[string][]int)
	myFarm.end = make(map[string][]int)
	myFarm.links = make(map[string][]string)

	var st, en int
	number, err := strconv.Atoi(content[0])
	if err != nil {
		log.Println("couldn't convert", err)
	}
	myFarm.ants_number = number

	for index := range content {
		if strings.TrimSpace(content[index]) == "##start" {
			st++
			if index+1 <= len(content)-1 {
				split := strings.Split(strings.TrimSpace(content[index+1]), " ")
				x, err := strconv.Atoi(split[1])
				y, err2 := strconv.Atoi(split[2])
				if err == nil && err2 == nil {
					myFarm.start[split[0]] = []int{x, y}
				}
			}
		} else if strings.TrimSpace(content[index]) == "##end" {
			en++
			if index+1 <= len(content)-1 {
				split := strings.Split(strings.TrimSpace(content[index+1]), " ")
				x, err := strconv.Atoi(split[1])
				y, err2 := strconv.Atoi(split[2])
				if err == nil && err2 == nil {
					myFarm.end[split[0]] = []int{x, y}
				}
			}
		} else if strings.Contains(content[index], "-") {
			split := strings.Split(strings.TrimSpace(content[index]), "-")
			if len(split) == 2 {
				myFarm.links[split[0]] = append(myFarm.links[split[0]], split[1])
			}
		} else if strings.Count(content[index], " ") == 2 {
			split := strings.Split(strings.TrimSpace(content[index]), " ")
			if len(split) == 3 {
				x, err := strconv.Atoi(split[1])
				y, err2 := strconv.Atoi(split[2])
				if err == nil || err2 == nil {
					myFarm.rooms[split[0]] = []int{x, y}
				}
			}
		} else if (strings.HasPrefix(strings.TrimSpace(content[index]), "#") || strings.HasPrefix(strings.TrimSpace(content[index]), "L")) && (strings.TrimSpace(content[index]) != "##start" && strings.TrimSpace(content[index]) != "##end") {
			continue
		}
	}
	if en != 1 || st != 1 {
		log.Println("rooms setup is incorrect", err)
	}
}

func Graph(farm farm) map[string][]string {
	adjacent := make(map[string][]string)
	for room := range farm.rooms {
		adjacent[room] = []string{}
	}
	for room, links := range farm.links {
		for _, link := range links {
			adjacent[room] = append(adjacent[room], link)
			adjacent[link] = append(adjacent[link], room)
		}
	}

	return adjacent
}

func BiBFS(myFarm *farm) {
	adjacent := Graph(*myFarm)
	var startRoom, endRoom string
	for k := range myFarm.start {
		startRoom = k
	}
	for k := range myFarm.end {
		endRoom = k
	}

	fmt.Printf("Starting room: %s, End room: %s\n", startRoom, endRoom)

	// For each adjacent room to start, do a bidirectional search
	for _, adj := range adjacent[startRoom] {
		fmt.Printf("\nSearching path through adjacent room: %s\n", adj)

		// Initialize fresh maps for each path search
		VisitedStart := make(map[string]bool)
		VisitedEnd := make(map[string]bool)
		ParentsStart := make(map[string]string)
		ParentsEnd := make(map[string]string)

		// Mark start and end rooms as visited in their respective directions
		VisitedStart[startRoom] = true
		VisitedEnd[endRoom] = true

		// Initialize start side with the adjacent room
		QueueStart := []string{adj}
		VisitedStart[adj] = true
		ParentsStart[adj] = startRoom

		// Initialize end side
		QueueEnd := []string{endRoom}

		// Search from this adjacent room
		for len(QueueStart) > 0 && len(QueueEnd) > 0 {

			if meetingRoom := bfsStep(adjacent, &QueueStart, VisitedStart, VisitedEnd, ParentsStart); meetingRoom != "" {
				printPath(meetingRoom, ParentsStart, ParentsEnd)
				break
			}

			if meetingRoom := bfsStep(adjacent, &QueueEnd, VisitedEnd, VisitedStart, ParentsEnd); meetingRoom != "" {
				printPath(meetingRoom, ParentsEnd, ParentsStart)
				break
			}
		}

		if len(QueueStart) == 0 || len(QueueEnd) == 0 {
			fmt.Printf("No path found through room %s\n", adj)
		}
	}
}

func bfsStep(adjacent map[string][]string, Queue *[]string, Visited, OppositeVisited map[string]bool, Parents map[string]string) string {
	if len(*Queue) == 0 {
		return ""
	}

	current := (*Queue)[0]
	*Queue = (*Queue)[1:]

	for _, link := range adjacent[current] {
		if !Visited[link] {
			Visited[link] = true
			Parents[link] = current
			*Queue = append(*Queue, link)

			if OppositeVisited[link] {
				return link
			}
		}
	}

	return ""
}

func printPath(meetingRoom string, ParentsStart, ParentsEnd map[string]string) {
	// Reconstruct path from start to meeting room
	pathStart := []string{meetingRoom}
	current := meetingRoom
	for ParentsStart[current] != "" {
		current = ParentsStart[current]
		pathStart = append([]string{current}, pathStart...)
	}

	// Reconstruct path from end to meeting room
	pathEnd := []string{}
	current = meetingRoom
	for ParentsEnd[current] != "" {
		current = ParentsEnd[current]
		pathEnd = append([]string{current}, pathEnd...)
	}

	pathEnd = append(pathEnd, current)

	// Combine both paths
	fullPath := append(pathStart, pathEnd[1:]...)
	var paths []string
	paths = append(paths, fullPath...)

	fmt.Printf("\nFull path from start to end: %v\n", fullPath)
	fmt.Println(pathEnd, pathStart, meetingRoom)
	fmt.Printf("\npaths: %v\n", paths)
}
