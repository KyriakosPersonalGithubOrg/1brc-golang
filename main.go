package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
    var measurementsPath string
    if len(os.Args) > 1 {
        measurementsPath = os.Args[1]
    } else {
        measurementsPath = "measurements.txt"
    }
    result, err := calculate(measurementsPath)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)
}

func calculate(filePath string) (string, error) {
    os, err := os.Open(filePath)
    if err != nil {
        return "", err 
    }
    defer os.Close()


    scanner := bufio.NewScanner(os)

    stationMap := make(map[string]Station)

    for scanner.Scan() {
        stationName, temperature, err := ParseTextToStation(scanner.Text())
        if err != nil {
            return  "", err
        }
        if val, exists := stationMap[stationName]; exists {
            val.Sum += temperature
            val.Count += 1
            if temperature < val.Min {
                val.Min = temperature
            }
            if temperature > val.Max {
                val.Max = temperature
            }
            stationMap[stationName] = val

            continue
        }
        stationMap[stationName] = Station{stationName, temperature, temperature, temperature, 1}
    }

    keys := make([]string, 0, len(stationMap))
    for k := range stationMap {
        keys = append(keys, k)
    }

    slices.Sort(keys)

    var sb strings.Builder
    sb.WriteString("{")
    for _, k := range keys {
        station := stationMap[k]
        sb.WriteString(fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", station.Name, station.Min, round(float64(station.Sum / station.Count)), station.Max))
    }
    sb.WriteString("}")
    finalResult := sb.String()[:sb.Len()-3] + "}"
    return finalResult, nil
}

func ParseTextToStation(text string) (station string, temperature float64, err error) {
    splitted := strings.Split(text, ";")
    temp, err := strconv.ParseFloat(splitted[1], 64)
    if err != nil {
        return "", 0, nil
    }
    return splitted[0], float64(temp), nil
}


func round(x float64) float64 {
    if x > 0 {
        return math.Ceil(x * 10) / 10
    }
    return math.Floor(x * 10) / 10
}

// Create a struct that contains name station and min max values
type Station struct {
    Name string
    Min  float64
    Max  float64
    Sum  float64
    Count float64
}
//
