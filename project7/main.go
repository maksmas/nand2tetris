package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const BASIC = "BasicTest"
const POINTER = "PointerTest"
const SIMPLE_ADD = "SimpleAdd"
const STACK = "StackTest"
const STATIC = "StaticTest"

var eqCounter = 0
var gtCounter = 0
var ltCounter = 0

var filenames = []string{BASIC, POINTER, SIMPLE_ADD, STACK, STATIC}

func main() {
	filename := os.Args[1]
	fmt.Println("Processing " + filename)
	lines := readLines(filename)
	filename = strings.Replace(filename, ".vm", "", -1)
	staticFilename := filename
	if strings.Contains(staticFilename, "/") {
		staticFilename = strings.Split(filename, "/")[1]
	}
	asmCommands := translate(lines, staticFilename)
	writeLines(asmCommands, filename+".asm")
	fmt.Println("Finished processing " + filename)
}

func readLines(filename string) []string {
	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lines := make([]string, 0, 0)
	for i := 0; ; i++ {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		lines = append(lines, string(line))
	}

	return lines
}

func translate(lines []string, filename string) []string {
	processed := make([]string, 0, len(lines))

	for _, line := range lines {
		vmCommand := removeNoise(line)

		if len(vmCommand) == 0 {
			continue
		}

		processed = append(processed, "//"+line)
		asmCommands := translateToAsm(vmCommand, filename)

		processed = append(processed, asmCommands...)
	}

	return processed
}

func removeNoise(line string) string {
	parts := strings.Split(line, "//")
	return strings.TrimSpace(parts[0])
}

func translateToAsm(vmCommand string, filename string) []string {
	cmdParts := strings.Split(vmCommand, " ")

	switch cmdParts[0] {
	case "add":
		return []string{"@SP", "AM=M-1", "D=M", "@SP", "A=M-1", "M=D+M"}
	case "sub":
		return []string{"@SP", "AM=M-1", "D=M", "@SP", "A=M-1", "M=M-D"}
	case "neg":
		return []string{"@SP", "A=M-1", "M=-M"}
	case "eq":
		eqCounter += 1
		return buildEq(eqCounter)
	case "gt":
		gtCounter += 1
		return buildGt(gtCounter)
	case "lt":
		ltCounter += 1
		return buildLt(ltCounter)
	case "and":
		return []string{"@SP", "AM=M-1", "D=M", "@SP", "A=M-1", "M=D&M"}
	case "or":
		return []string{"@SP", "AM=M-1", "D=M", "@SP", "A=M-1", "M=D|M"}
	case "not":
		return []string{"@SP", "A=M-1", "M=!M"}
	case "pop":
		return buildPop(cmdParts[1], cmdParts[2], filename)
	case "push":
		return buildPush(cmdParts[1], cmdParts[2], filename)
	default:
		panic("Wrong vmCommand " + vmCommand)
	}
}

func buildEq(eqCounter int) []string {
	eq := make([]string, 14, 14)
	gl := "TRUE" + strconv.Itoa(eqCounter)

	eq[0] = "@SP"
	eq[1] = "AM=M-1"
	eq[2] = "D=M"

	eq[3] = "@SP"
	eq[4] = "A=M-1"
	eq[5] = "D=D-M"
	eq[6] = "@" + gl
	eq[7] = "D;JEQ"
	eq[8] = "D=-1"
	eq[9] = "(" + gl + ")"
	eq[10] = "D=!D"

	eq[11] = "@SP"
	eq[12] = "A=M-1"
	eq[13] = "M=D"

	return eq
}

func buildGt(gtCounter int) []string {
	gt := make([]string, 17, 17)
	counter := strconv.Itoa(gtCounter)
	gl := "GREATER" + counter
	gl2 := "GT_RES" + counter

	gt[0] = "@SP"
	gt[1] = "AM=M-1"
	gt[2] = "D=M"

	gt[3] = "@SP"
	gt[4] = "A=M-1"
	gt[5] = "D=M-D"
	gt[6] = "@" + gl
	gt[7] = "D;JGT"
	gt[8] = "D=0"
	gt[9] = "@" + gl2
	gt[10] = "0;JMP"
	gt[11] = "(" + gl + ")"
	gt[12] = "D=-1"
	gt[13] = "(" + gl2 + ")"

	gt[14] = "@SP"
	gt[15] = "A=M-1"
	gt[16] = "M=D"

	return gt
}

func buildLt(ltCounter int) []string {
	lt := make([]string, 17, 17)
	counter := strconv.Itoa(ltCounter)
	gl := "LESSER" + counter
	gl2 := "lt_RES" + counter

	lt[0] = "@SP"
	lt[1] = "AM=M-1"
	lt[2] = "D=M"

	lt[3] = "@SP"
	lt[4] = "A=M-1"
	lt[5] = "D=M-D"
	lt[6] = "@" + gl
	lt[7] = "D;JLT"
	lt[8] = "D=0"
	lt[9] = "@" + gl2
	lt[10] = "0;JMP"
	lt[11] = "(" + gl + ")"
	lt[12] = "D=-1"
	lt[13] = "(" + gl2 + ")"

	lt[14] = "@SP"
	lt[15] = "A=M-1"
	lt[16] = "M=D"

	return lt
}

func buildPop(segment string, addr string, filename string) []string {
	intAddr, err := strconv.Atoi(addr)
	if err != nil {
		panic(err)
	}

	if segment == "static" {
		return []string{"@SP", "AM=M-1", "D=M", "@" + filename + "." + addr, "M=D"}
	} else if segment == "temp" {
		return []string{"@SP", "AM=M-1", "D=M", "@" + strconv.Itoa(intAddr+5), "M=D"}
	} else if segment == "pointer" {
		ch := "THIS"

		if addr == "1" {
			ch = "THAT"
		}

		return []string{"@SP", "AM=M-1", "D=M", "@" + ch, "M=D"}
	}

	commands := make([]string, intAddr+6, intAddr+6)

	commands[0] = "@SP"
	commands[1] = "AM=M-1"
	commands[2] = "D=M"
	commands[3] = "@" + mapSegment(segment)
	commands[4] = "A=M"
	for i := 0; i < intAddr; i++ {
		commands[i+5] = "A=A+1"
	}
	commands[intAddr+5] = "M=D"
	return commands
}

func buildPush(segment string, addr string, filename string ) []string {
	intAddr, err := strconv.Atoi(addr)
	if err != nil {
		panic(err)
	}

	if segment == "constant" {
		return []string{"@" + addr, "D=A", "@SP", "AM=M+1", "A=A-1", "M=D"}
	} else if segment == "static" {
		return []string{"@" + filename + "." + addr, "D=M", "@SP", "AM=M+1", "A=A-1", "M=D"}
	} else if segment == "temp" {
		return []string{"@" + strconv.Itoa(intAddr+5), "D=M", "@SP", "AM=M+1", "A=A-1", "M=D"}
	} else if segment == "pointer" {
		ch := "THIS"

		if addr == "1" {
			ch = "THAT"
		}

		return []string{"@" + ch, "D=M", "@SP", "AM=M+1", "A=A-1", "M=D"}
	}

	commands := make([]string, intAddr+7, intAddr+7)
	commands[0] = "@" + mapSegment(segment)
	commands[1] = "A=M"
	for i := 0; i < intAddr; i++ {
		commands[i+2] = "A=A+1"
	}
	commands[intAddr+2] = "D=M"
	commands[intAddr+3] = "@SP"
	commands[intAddr+4] = "AM=M+1"
	commands[intAddr+5] = "A=A-1"
	commands[intAddr+6] = "M=D"
	return commands
}

func mapSegment(segment string) string {
	switch segment {
	case "local":
		return "LCL"
	case "argument":
		return "ARG"
	case "this":
		return "THIS"
	case "that":
		return "THAT"
	default:
		panic("invalid segment " + segment)
	}
}

func writeLines(lines []string, targetFile string) {
	fout, err := os.Create(targetFile)
	if err != nil {
		panic("Failed to create file")
	}
	defer fout.Close()

	writer := bufio.NewWriter(fout)
	defer writer.Flush()
	for _, line := range lines {
		if _, err := fmt.Fprintln(writer, line); err != nil {
			panic(err)
		}
	}
}
