package main

import (
    "bufio"
    "fmt"
    "log"
    "math"
    "os"
    "runtime/pprof"
    "slices"
    "strconv"
    "strings"
)

func main() {
    // Start CPU profiling
    if len(os.Args) > 2 && os.Args[2] == "cpu" {
        f, err := os.Create("cpu.pprof")
        if err != nil {
            panic(err)
        } 
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

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

    stationMap := make(map[string]StationAggregate)

    for scanner.Scan() {
        station := ParseTextToStation(scanner.Text())
        if err != nil {
            return  "", err
        }
        if val, exists := stationMap[station.Name]; exists {
            val.Sum += station.Temperature
            val.Count += 1
            if station.Temperature < val.Min {
                val.Min = station.Temperature
            }
            if station.Temperature > val.Max {
                val.Max = station.Temperature
            }
            stationMap[station.Name] = val

            continue
        }
        stationMap[station.Name] = StationAggregate{station.Name, 
        station.Temperature, station.Temperature, station.Temperature, 1}
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
        sb.WriteString(fmt.Sprintf("%s=%.1f/%.1f/%.1f, ", station.Name, 
        station.Min, round(float64(station.Sum / station.Count)), station.Max))
    }
    sb.WriteString("}")
    finalResult := sb.String()[:sb.Len()-3] + "}"
    return finalResult, nil
}

func ParseTextToStation(text string) Station {
    splitted := strings.Split(text, ";")
    temp, err := strconv.ParseFloat(splitted[1], 64)
    if err != nil {
        return Station{}
    }
    return Station{splitted[0], temp}
}


func round(x float64) float64 {
    if x > 0 {
        return math.Ceil(x * 10) / 10
    }
    return math.Floor(x * 10) / 10
}

// Create a struct that contains name station and min max values
type StationAggregate struct {
    Name string
    Min  float64
    Max  float64
    Sum  float64
    Count float64
}

type Station struct {
    Name string
    Temperature float64
}
