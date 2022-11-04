package reminder

import (
	"testing"
	"time"
)

func TestParseSpec(t *testing.T) {
	loc, err := time.LoadLocation("America/Montreal")
	if err != nil {
		t.Fatal(err)
	}
	time, content, _ := parseSpec("12/24 12:41 buy goobers", loc)
	t.Log(time, content)

}
