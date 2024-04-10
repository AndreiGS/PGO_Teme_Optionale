package pkg

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var CHUNK_SIZE = 512

func GenerateFromFile(ctx context.Context, file string) {
	fileName := strings.Replace(file, " ", "_", -1)
	fileName = strings.Replace(fileName, ".txt", "", -1)
	fileName = fmt.Sprintf("%s_generated.txt", fileName)
	os.Remove(fileName)

	type GenerateEvent struct {
		NoBytesRead int64
		MaxBytes    int64
	}

	f, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		return
	}

	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string
	var noLines int = 0
	var noBytes int64 = 0
	stat, err := os.Stat(file)

	if err != nil {
		return
	}
	var maxBytes int64 = int64(stat.Size())

	fmt.Printf("Generating from file %s...", file)
	fmt.Printf("Generated %d lines", noLines)
	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())

		if len(fileLines) < CHUNK_SIZE {
			continue
		}
		generateFromBatch(fileName, fileLines)
		noLines += CHUNK_SIZE
		noBytes += int64(len(fileScanner.Bytes()))
		fileLines = make([]string, 0)
		fmt.Printf("\rGenerated %d lines", noLines)
		runtime.EventsEmit(ctx, "generate", &GenerateEvent{NoBytesRead: noBytes, MaxBytes: maxBytes})
	}
	if len(fileLines) > 0 {
		generateFromBatch(fileName, fileLines)
	}
	noLines += len(fileLines)
	noBytes += int64(len(fileScanner.Bytes()))
	fmt.Printf("\rGenerated %d lines", noLines)
	runtime.EventsEmit(ctx, "generate", &GenerateEvent{NoBytesRead: noBytes, MaxBytes: maxBytes})
	fmt.Printf("Finished generating from file %s!", file)
}

func generateFromBatch(fileName string, batch []string) {
	input, e := parseInput(batch)
	if e != nil {
		return
	}

	response, err := Generate(input)
	if err != nil {
		return
	}

	lines := make([]string, 0)
	for index, value := range response {
		p := int32(input[index][0])
		q := int32(input[index][1])
		theta := input[index][2]
		phi := input[index][3]
		line := fmt.Sprintf("%d\t%d\t%f\t%f\t%f\t%f\n", p, q, theta, phi, value.Real, value.Imag)

		lines = append(lines, line)
	}

	AppendToFile(fileName, lines)
}

func parseInput(batch []string) ([][]float32, error) {
	input := make([][]float32, 0)
	for i := 0; i < len(batch); i++ {
		values := strings.Fields(batch[i])
		if len(values) < 4 {
			continue
		}
		input = append(input, make([]float32, 0))
		for j := 0; j < 4; j++ {
			value := values[j]
			val, e := strconv.ParseFloat(value, 32)
			if e != nil {
				return nil, fmt.Errorf("error parsing input value \"%s\": %w",
					value, e)
			}
			input[i] = append(input[i], float32(val))
		}
	}
	return input, nil
}
