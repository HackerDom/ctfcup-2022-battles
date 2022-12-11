package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	readRam    = []int{252, 2, 157, 28, 102, 55, 232, 17, 208, 3}
	compareRam = []int{252, 2, 157, 28, 102, 55, 232, 17, 208, 3, 254, 185, 157, 249, 211}
)

func main() {
	server := http.NewServeMux()
	server.HandleFunc("/prog", HandleProg)

	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Printf("start service failed: %v\n", err)
	}
}

func HandleProg(writer http.ResponseWriter, request *http.Request) {
	rom, err := io.ReadAll(io.LimitReader(request.Body, 100))
	if err != nil {
		response(writer, "read body failed: %v", err)
		return
	}

	ram := make([]int, 15)
	copy(ram, readRam)
	addressReg := 0
	aReg := 0
	bReg := 0

	for pos := 0; pos < len(rom); pos += 2 {
		control := int(rom[pos])
		opcode := getBit(control, 0)

		if opcode == 0 {
			addressReg = int(rom[pos+1])
			continue
		}

		if addressReg < 0 || addressReg > 15 {
			continue
		}

		regSelector1 := getBit(control, 1)
		regSelector2 := getBit(control, 2)

		switch {
		case regSelector1 == 0 && regSelector2 == 1:
			aReg = ram[addressReg]
		case regSelector1 == 1 && regSelector2 == 1:
			bReg = ram[addressReg]
		case regSelector1 == 0 && regSelector2 == 0:
			ram[addressReg] = aReg + bReg
		default:
		}
	}

	for i, j := range ram {
		if j != compareRam[i] {
			response(writer, "wrong flag!")
			return
		}
	}

	response(writer, "Congratulations! Your flag: %v", os.Getenv("FLAG"))
}

func response(writer http.ResponseWriter, format string, a ...interface{}) {
	if _, err := fmt.Fprintf(writer, format, a...); err != nil {
		log.Printf("write response failed: %v\n", err)
	}
}

func getBit(v, index int) int {
	if v&(1<<index) > 0 {
		return 1
	}

	return 0
}
