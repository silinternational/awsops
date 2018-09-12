package lib

type InstanceType struct {
	MemoryMb int64
	CPUUnits int64
}

var SingleCPUUnits int64 = 1024
var MbInGb int64 = 1024

var InstanceTypes = map[string]InstanceType{
	"t2.nano": {
		CPUUnits: 1 * SingleCPUUnits,
		MemoryMb: 512,
	},
	"t2.micro": {
		CPUUnits: 1 * SingleCPUUnits,
		MemoryMb: 1 * MbInGb,
	},
	"t2.small": {
		CPUUnits: 1 * SingleCPUUnits,
		MemoryMb: 2 * MbInGb,
	},
	"t2.medium": {
		CPUUnits: 2 * SingleCPUUnits,
		MemoryMb: 4 * MbInGb,
	},
	"t2.large": {
		CPUUnits: 2 * SingleCPUUnits,
		MemoryMb: 8 * MbInGb,
	},
	"t2.xlarge": {
		CPUUnits: 4 * SingleCPUUnits,
		MemoryMb: 16 * MbInGb,
	},
	"t2.2xlarge": {
		CPUUnits: 8 * SingleCPUUnits,
		MemoryMb: 32 * MbInGb,
	},
}
