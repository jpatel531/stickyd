package pm

import (
	"github.com/jpatel531/stickyd/stats"
	"math"
	"sort"
	"strconv"
	"time"
)

type ProcessedMetrics struct {
	StickyDMetrics   map[string]int64              `json:"stickyd_metrics,omitempty"`
	Counters         map[string]float64            `json:"counters,omitempty"`
	CounterRates     map[string]float64            `json:"counter_rates,omitempty"`
	Sets             map[string][]interface{}      `json:"sets,omitempty"`
	Gauges           map[string]float64            `json:"gauges,omitempty"`
	TimerData        map[string]map[string]float64 `json:"timer_data,omitempty"`
	PercentThreshold []int                         `json:"pctThreshold,omitempty"`
}

func ProcessMetrics(appStats *stats.AppStats, flushInterval int, percentThreshold []int) *ProcessedMetrics {
	startTime := time.Now().Unix()

	counters := appStats.Counters.Map()
	counterRates := map[string]float64{}

	for key, counter := range counters {
		counterRates[key] = counter / float64(flushInterval/1000)
	}

	timers := appStats.Timers.Map()
	timerData := make(map[string]map[string]float64)
	timerCounters := appStats.TimerCounters
	for key, values := range timers {
		currentTimerData := make(map[string]float64)
		if count := len(values); count > 0 {
			sort.Float64s(values)

			min := values[0]
			max := values[count-1]

			cumulativeValues := []float64{min}
			cumulSumSquaresValues := []float64{min * min}

			for i, value := range values[1:] {
				cumulativeValues = append(
					cumulativeValues,
					value+cumulativeValues[i],
				)
				cumulSumSquaresValues = append(
					cumulSumSquaresValues,
					(value*value)+cumulSumSquaresValues[i],
				)
			}

			sum := min
			sumSquares := min * min
			mean := min
			thresholdBoundary := max

			for _, pct := range percentThreshold {
				numInThreshold := count

				if count > 1 {
					numInThreshold = int((math.Abs(float64(pct)) / 100) * float64(count))
					if numInThreshold == 0 {
						continue
					}

					if pct > 0 {
						thresholdBoundary = values[numInThreshold-1]
						sum = cumulativeValues[numInThreshold-1]
						sumSquares = cumulSumSquaresValues[numInThreshold-1]
					} else {
						thresholdBoundary = values[count-numInThreshold]
						sum = cumulativeValues[count-1] - cumulativeValues[count-numInThreshold-1]
						sumSquares = cumulSumSquaresValues[count-1] - cumulSumSquaresValues[count-numInThreshold-1]
					}
					mean = sum / float64(numInThreshold)
				}

				cleanPct := strconv.Itoa(pct)

				currentTimerData["count_"+cleanPct] = float64(numInThreshold)
				currentTimerData["mean_"+cleanPct] = mean

				if pct > 0 {
					currentTimerData["upper_"+cleanPct] = thresholdBoundary
				} else {
					currentTimerData["lower_"+cleanPct] = thresholdBoundary
				}

				currentTimerData["sum_"+cleanPct] = sum
				currentTimerData["sum_squares_"+cleanPct] = sumSquares
			}
			sum = cumulativeValues[count-1]
			sumSquares = cumulSumSquaresValues[count-1]
			mean = float64(sum) / float64(count)

			var sumOfDiffs float64
			for _, value := range values {
				sumOfDiffs += (value - mean) * (value - mean)
			}

			mid := int(math.Floor(float64(count) / float64(2)))
			var median float64
			if count%2 == 0 {
				median = (values[mid-1] + values[mid]) / 2
			} else {
				median = values[mid]
			}

			stddev := math.Sqrt(sumOfDiffs / float64(count))
			currentTimerData["std"] = stddev
			currentTimerData["upper"] = max
			currentTimerData["lower"] = min
			currentTimerData["count"] = timerCounters.Get(key)
			currentTimerData["count_ps"] = timerCounters.Get(key) / float64(flushInterval/100)
			currentTimerData["sum"] = sum
			currentTimerData["sum_squares"] = sumSquares
			currentTimerData["mean"] = mean
			currentTimerData["median"] = median

			// TODO manage histogram
		} else {
			currentTimerData["count"] = 0
			currentTimerData["count_ps"] = 0
		}

		timerData[key] = currentTimerData
	}

	stickydMetrics := map[string]int64{
		"processingTime": startTime - time.Now().Unix(),
	}

	return &ProcessedMetrics{
		StickyDMetrics:   stickydMetrics,
		Counters:         counters,
		CounterRates:     counterRates,
		Sets:             appStats.Sets.Map(),
		Gauges:           appStats.Gauges.Map(),
		TimerData:        timerData,
		PercentThreshold: percentThreshold,
	}
}
