package teedy

import (
	"strconv"
	"time"
)

// Timestamp defines a timestamp encoded as epoch seconds in JSON
type Timestamp struct {
	time.Time
}

// MarshalJSON is used to convert the timestamp to JSON
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(t.Marshal()), nil
}

// UnmarshalJSON is used to convert the timestamp from JSON
func (t *Timestamp) UnmarshalJSON(s []byte) (err error) {
	r := string(s)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(0, q*int64(time.Millisecond))
	t.Time = t.Time.UTC()
	return nil
}

func (t Timestamp) Marshal() string {
	ut := t.UnixNano() / int64(time.Millisecond)
	ret := strconv.FormatInt(ut, 10)
	return ret
}
