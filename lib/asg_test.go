package lib

import "testing"

func TestHowManyServersNeededFor(t *testing.T) {
	tests := []struct {
		MemNeeded   int64
		CPUNeeded   int64
		SmallestMem int64
		SmallestCPU int64
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
		{
			MemNeeded:   886,
			CPUNeeded:   1,
			SmallestMem: 100,
			ServerType:  "t2.micro",
			ExpectedNum: 2, // 1024 - 39 - 100 = 885 MB available per instance
		},
		{
			MemNeeded:   1,
			CPUNeeded:   925,
			SmallestCPU: 100,
			ServerType:  "t2.micro",
			ExpectedNum: 2, // 1024 - 100 = 924 CPU available per instance
		},
	}

	for _, i := range tests {
		resourceSizes := ResourceSizes{
			TotalCPU:       i.CPUNeeded,
			TotalMemory:    i.MemNeeded,
			SmallestCPU:    i.SmallestCPU,
			SmallestMemory: i.SmallestMem,
		}
		results := HowManyServersNeededForAsg(i.ServerType, resourceSizes)
		if results != i.ExpectedNum {
			t.Errorf("Did not get back expected number of %s servers needed for %v mem and %v cpu, expected %v, got %v",
				i.ServerType, i.MemNeeded, i.CPUNeeded, i.ExpectedNum, results)
		}
	}
}
