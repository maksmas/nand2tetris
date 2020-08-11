package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var eqCounter = 0
var gtCounter = 0
var ltCounter = 0
var callCounter = 0

func main() {
	filenames := make([]string, 0, 0)

	location := os.Args[1]
	info, err := os.Stat(location)

	if err != nil {
		panic(err)
	}

	output := ""
	var asmCommands []string

	if info.Mode().IsDir() {
		asmCommands = bootstrap()
		output = location + "/" + location + ".asm"
		err = filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".vm") {
				filenames = append(filenames, path)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	} else {
		asmCommands = make([]string, 0, 0)
		filenames = append(filenames, location)
		output = strings.Replace(location, ".vm", ".asm", 1)
	}

	for _, filename := range filenames {
		fmt.Println("Processing " + filename)
		lines := readLines(filename)
		filename = strings.Replace(filename, ".vm", "", -1)
		staticFilename := filename
		if strings.Contains(staticFilename, "/") {
			staticFilename = strings.Split(filename, "/")[1]
		}

		if strings.Contains(staticFilename, "\\") {
			staticFilename = strings.Split(filename, "\\")[1]
		}
		asmCommands = append(asmCommands, translate(lines, staticFilename)...)
		fmt.Println("Finished processing " + filename)
	}

	writeLines(asmCommands, output)
}

func bootstrap() []string {
	boot := make([]string, 6, 6)

	boot[0] = "@261"
	boot[1] = "D=A"
	boot[2] = "@SP"
	boot[3] = "M=D"
	boot[4] = "@Sys.init"
	boot[5] = "0;JMP"

	return boot
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
	case "label":
		return []string{fmt.Sprintf("(%s)", cmdParts[1])}
	case "goto":
		return []string{"@" + cmdParts[1], "0;JMP"}
	case "if-goto":
		return []string{"@SP", "AM=M-1", "D=M", "@" + cmdParts[1], "D;JNE"}
	case "function":
		return buildFunction(cmdParts[1], cmdParts[2])
	case "return":
		return buildReturn(filename)
	case "call":
		return buildCall(cmdParts[1], cmdParts[2], filename)
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

func buildPush(segment string, addr string, filename string) []string {
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

func buildFunction(fnName string, localArgsString string) []string {
	localArgs, err := strconv.Atoi(localArgsString)
	if err != nil {
		panic(err)
	}

	fn := make([]string, 1)

	fn[0] = fmt.Sprintf("(%s)", fnName)

	if localArgs > 0 {
		fn = append(fn, "@SP")
		fn = append(fn, "D=M")

		for i := 0; i < localArgs; i++ {
			fn = append(fn, "M=M+1")
		}

		fn = append(fn, "A=D")

		for i := 0; i < localArgs; i++ {
			fn = append(fn, "M=0")
			fn = append(fn, "A=A+1")
		}
	}

	return fn
}

func buildReturn(filename string) []string {
	ret := make([]string, 42, 42)

	callCounter += 1

	ret[0] = "@LCL"
	ret[1] = "D=M-1"
	ret[2] = "D=D-1"
	ret[3] = "D=D-1"
	ret[4] = "D=D-1"
	ret[5] = "D=D-1"
	ret[6] = "A=D"
	ret[7] = "D=M"

	ret[8] = fmt.Sprintf("@ret_addr.%s.%d", filename, callCounter)
	ret[9] = "M=D"

	ret[10] = "@SP"
	ret[11] = "AM=M-1"
	ret[12] = "D=M"
	ret[13] = "@ARG"
	ret[14] = "A=M"
	ret[15] = "M=D"
	ret[16] = "D=A"
	ret[17] = "@SP"
	ret[18] = "M=D+1"

	ret[19] = "@LCL"
	ret[20] = "AM=M-1"
	ret[21] = "D=M"
	ret[22] = "@THAT"
	ret[23] = "M=D"

	ret[24] = "@LCL"
	ret[25] = "AM=M-1"
	ret[26] = "D=M"
	ret[27] = "@THIS"
	ret[28] = "M=D"

	ret[29] = "@LCL"
	ret[30] = "AM=M-1"
	ret[31] = "D=M"
	ret[32] = "@ARG"
	ret[33] = "M=D"

	ret[34] = "@LCL"
	ret[35] = "AM=M-1"
	ret[36] = "D=M"
	ret[37] = "@LCL"
	ret[38] = "M=D"

	ret[39] = fmt.Sprintf("@ret_addr.%s.%d", filename, callCounter)
	ret[40] = "A=M"
	ret[41] = "0;JMP"

	return ret
}

func buildCall(fnName string, argCountString string, filename string) []string {
	argCount, err := strconv.Atoi(argCountString)
	if err != nil {
		panic(err)
	}

	callCounter += 1
	call := make([]string, 46 + argCount, 46 + argCount)

	call[0] = fmt.Sprintf("@ret_addr_lbl.%s.%s.%d", filename, fnName, callCounter)
	call[1] = "D=A"
	call[2] = "@SP"
	call[3] = "M=M+1"
	call[4] = "A=M-1"
	call[5] = "M=D"

	call[6] = "@LCL"
	call[7] = "D=M"
	call[8] = "@SP"
	call[9] = "M=M+1"
	call[10] = "A=M-1"
	call[11] = "M=D"

	call[12] = "@ARG"
	call[13] = "D=M"
	call[14] = "@SP"
	call[15] = "M=M+1"
	call[16] = "A=M-1"
	call[17] = "M=D"

	call[18] = "@THIS"
	call[19] = "D=M"
	call[20] = "@SP"
	call[21] = "M=M+1"
	call[22] = "A=M-1"
	call[23] = "M=D"

	call[24] = "@THAT"
	call[25] = "D=M"
	call[26] = "@SP"
	call[27] = "M=M+1"
	call[28] = "A=M-1"
	call[29] = "M=D"

	call[30] = "@SP"
	call[31] = "D=M"

	for i := 0; i < argCount + 5; i++ {
		call[32 + i] = "D=D-1"
	}

	call[37 + argCount] = "@ARG"
	call[38 + argCount] = "M=D"

	call[39 + argCount] = "@SP"
	call[40 + argCount] = "D=M"
	call[41 + argCount] = "@LCL"
	call[42 + argCount] = "M=D"

	call[43 + argCount] = fmt.Sprintf("@%s", fnName)
	call[44 + argCount] = "0;JMP"

	call[45 + argCount] = fmt.Sprintf("(ret_addr_lbl.%s.%s.%d)", filename, fnName, callCounter)

	return call
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
