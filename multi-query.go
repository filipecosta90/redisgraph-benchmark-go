package main

import (
	"log"
	"math/rand"
	"strconv"
)

func sample(cdf []float32) int {
	r := rand.Float32()
	bucket := 0
	for r > cdf[bucket] {
		bucket++
	}
	return bucket
}

func prepareCommandsDistribution(cmds []string, cmdRates []int) (int, []float32) {
	var totalRates int = 0
	var totalDifferentCommands = len(cmds)
	var err error
	for i, rawCmdString := range benchmarkQueries {
		cmds[i] = rawCmdString
		if i >= len(benchmarkQueryRates) {
			cmdRates[i] = 1

		} else {
			cmdRates[i], err = strconv.Atoi(benchmarkQueryRates[i])
			if err != nil {
				log.Fatalf("Error while converting query-rate param %s: %v", benchmarkQueryRates[i], err)
			}
		}
		totalRates += cmdRates[i]
	}
	// probability density function
	if len(benchmarkQueryRates) > 0 && (len(benchmarkQueryRates) != len(benchmarkQueries)) {
		log.Fatalf("When specifiying -query-rate parameter, you need to have the same number of -query and -query-rate parameters. Number of time -query ( %d ) != Number of times -query-params ( %d )", len(benchmarkQueries), len(benchmarkQueryRates))
	}
	pdf := make([]float32, len(benchmarkQueries))
	cdf := make([]float32, len(benchmarkQueries))
	for i := 0; i < len(cmdRates); i++ {
		pdf[i] = float32(cmdRates[i]) / float32(totalRates)
		cdf[i] = 0
	}
	// get cdf
	cdf[0] = pdf[0]
	for i := 1; i < len(cmdRates); i++ {
		cdf[i] = cdf[i-1] + pdf[i]
	}
	return totalDifferentCommands, cdf
}