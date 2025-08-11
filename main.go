package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/beevik/ntp"
)

type ntpServer struct {
	address string
	stratum int
	rtt     time.Duration
	offset  time.Duration
}

// getNTPTime retrieves the current time from an NTP server
func getNTPTime(server string) (string, error) {
	response, err := ntp.Query(server)
	if err != nil {
		return "", err
	}

	ntpTime := time.Now().Add(response.ClockOffset)
	return ntpTime.Format(time.RFC3339), nil
}

// getAccurateNTPTime retrieves time from multiple servers for better accuracy
func getAccurateNTPTime() (string, error) {
	servers := []string{
		"time.nist.gov",             // stratum 1 - US National Institute of Standards
		"time-a-g.nist.gov",         // stratum 1 - NIST server A
		"time-b-g.nist.gov",         // stratum 1 - NIST server B
		"0.beevik-ntp.pool.ntp.org", // stratum 2 - Pool server
		"1.beevik-ntp.pool.ntp.org", // stratum 2 - Pool server
		"2.beevik-ntp.pool.ntp.org", // stratum 2 - Pool server
	}

	var responses []ntpServer
	var validResponses int

	for _, server := range servers {
		response, err := ntp.Query(server)
		if err != nil {
			log.Printf("Warning: %s unavailable: %v", server, err)
			continue
		}

		if response.Stratum > 0 && response.Stratum <= 4 && response.RTT < 2*time.Second {
			responses = append(responses, ntpServer{
				address: server,
				stratum: int(response.Stratum),
				rtt:     response.RTT,
				offset:  response.ClockOffset,
			})
			validResponses++
		}
	}

	if validResponses == 0 {
		return "", fmt.Errorf("no reliable NTP servers available")
	}

	sort.Slice(responses, func(i, j int) bool {
		if responses[i].stratum != responses[j].stratum {
			return responses[i].stratum < responses[j].stratum
		}
		return responses[i].rtt < responses[j].rtt
	})

	var totalWeight float64
	var weightedOffsetNs float64

	for _, resp := range responses {
		stratumWeight := 1.0 / float64(resp.stratum)
		rttMs := resp.rtt.Milliseconds()
		if rttMs == 0 {
			rttMs = 1
		}
		rttWeight := 1.0 / float64(rttMs)
		weight := stratumWeight * rttWeight

		totalWeight += weight
		weightedOffsetNs += float64(resp.offset.Nanoseconds()) * weight
	}

	var finalOffset time.Duration
	if totalWeight == 0 {
		var totalOffset time.Duration
		for _, resp := range responses {
			totalOffset += resp.offset
		}
		finalOffset = totalOffset / time.Duration(len(responses))
	} else {
		finalOffset = time.Duration(weightedOffsetNs / totalWeight)
	}

	accurateTime := time.Now().Add(finalOffset)

	log.Printf("Used %d servers, best: %s (stratum %d, rtt %v)",
		validResponses, responses[0].address, responses[0].stratum, responses[0].rtt)

	return accurateTime.Format(time.RFC3339), nil
}

func main() {
	timeStr, err := getAccurateNTPTime()
	if err != nil {
		log.Printf("Warning: Accurate NTP failed, using fallback: %v", err)
		timeStr, err = getNTPTime("0.beevik-ntp.pool.ntp.org")
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	fmt.Println(timeStr)
}
