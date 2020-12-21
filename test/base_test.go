package test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T){
	tm := struct {
		time.Time
		N int
	}{
		time.Date(2020,12,20,0,0,0,0,time.UTC),
		5,
	}

	tm.String()

	marshal, _ := json.Marshal(tm)
	fmt.Printf("%s\n",marshal)
	fmt.Printf("%s\n",tm)

}


