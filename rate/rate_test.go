package rate

import (
	"testing"
)

func TestCheckRate(t *testing.T) {

	b := checkRate("10.10.10.5")
	if b == false {
		t.Log("you are over the rate limit ")
	} else {
		t.Log("not over the limit")
	}

}
