package proxyblock

import (
	"sync"
	"time"
)

// https://blog.cloudflare.com/counting-things-a-lot-of-different-things/

// OneMinute is
var OneMinute = 1 * time.Minute

// CountData is
type CountData struct {
	sync.RWMutex
	PrevMinuteTime time.Time
	MaxPerMin      int
	CurMinCount    int
	PrevMinCount   int
	WindowSize     int
}

// NewCountData is
func NewCountData(MaxPerMinCount int) *CountData {
	c := CountData{MaxPerMin: MaxPerMinCount, WindowSize: 60}
	c.PrevMinuteTime = c.GetTimeNow()
	return &c
}

// Increment updates the current minute count
func (c *CountData) Increment() {
	c.Lock()
	defer c.Unlock()
	c.CurMinCount++
}

// ResetCount sets the previous time window variable
func (c *CountData) ResetCount(t time.Time) {
	// determine whether we have data for previous minute or not
	if !t.Equal(c.PrevMinuteTime.Add(OneMinute)) {
		c.Lock()
		defer c.Unlock()
		// first create a zero value for last minute
		c.PrevMinuteTime = t.Add(OneMinute * -1)
		c.PrevMinCount = 0
		// then add our count to the current minute
		c.CurMinCount = 1
	} else {
		c.Lock()
		defer c.Unlock()
		c.CurMinCount = 1
		c.PrevMinuteTime = t
	}

}

// GetTimeNow returns the current time properly truncated for use with our app -- it zeroes out everything past seconds
func (c *CountData) GetTimeNow() time.Time {
	return time.Now().Truncate(time.Second)
}

// ShouldBlock is function to determine whether we should block a request based on hit volume over the last
// interval
func (c *CountData) ShouldBlock() bool {
	// make sure we have a timestamp ONLY up to seconds for when we compare equality
	nowTime := c.GetTimeNow()

	// check if our stored prevTime is within a minute of nowTime.  If not, reset the count and other values
	if nowTime.Sub(c.PrevMinuteTime) >= OneMinute {
		c.ResetCount(nowTime)
	}

	// sliding window calculation
	// previousCount * ratio of (how many seconds into current minute/60 seconds) + currentCount
	slidingWindowRate := c.PrevMinCount*((60-nowTime.Second())/60) + c.CurMinCount
	if slidingWindowRate > c.MaxPerMin {
		return true
	}
	c.Increment()
	return false
}
