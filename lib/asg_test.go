package lib

import "testing"

func TestHowManyServersNeededFor(t *testing.T) {
	tests := []struct {
		MemNeeded   int64
		CPUNeeded   int64
		ServerType  string
		ExpectedNum int64
	}{
		{
			MemNeeded:   3000,
			CPUNeeded:   1024,
			ServerType:  "t2.micro",
			ExpectedNum: 4,
		},
		{
			MemNeeded:   3000,
			CPUNeeded:   1024,
			ServerType:  "t2.small",
			ExpectedNum: 2,
		},
		{
			MemNeeded:   3000,
			CPUNeeded:   1024,
			ServerType:  "t2.large",
			ExpectedNum: 1,
		},
		{
			MemNeeded:   2047,
			CPUNeeded:   1024,
			ServerType:  "t2.micro",
			ExpectedNum: 3,
		},
		{
			MemNeeded:   2047,
			CPUNeeded:   2049,
			ServerType:  "t2.micro",
			ExpectedNum: 3,
		},
	}

	for _, i := range tests {
		results := HowManyServersNeededForAsg(i.ServerType, i.MemNeeded, i.CPUNeeded)
		if results != i.ExpectedNum {
			t.Errorf("Did not get back expected number of %s servers needed for %v mem and %v cpu, expected %v, got %v",
				i.ServerType, i.MemNeeded, i.CPUNeeded, i.ExpectedNum, results)
		}
	}
}
