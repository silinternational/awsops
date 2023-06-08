package lib

type InstanceType struct {
	MemoryMb int64
	CPUUnits int64
}

const (
	SingleCPUUnits int64 = 1024
	MbInGb         int64 = 1024
	AgentSize      int64 = 39 // memory consumed by the ECS agent
)

var InstanceTypes = map[string]InstanceType{
	"t2.nano": {
		CPUUnits: 1 * SingleCPUUnits,
		MemoryMb: MbInGb/2 - AgentSize,
	},
	"t2.micro": {
		CPUUnits: 1 * SingleCPUUnits,
		MemoryMb: 1*MbInGb - AgentSize,
	},
	"t2.small": {
		CPUUnits: 1 * SingleCPUUnits,
		MemoryMb: 2*MbInGb - AgentSize,
	},
	"t2.medium": {
		CPUUnits: 2 * SingleCPUUnits,
		MemoryMb: 4*MbInGb - AgentSize,
	},
	"t2.large": {
		CPUUnits: 2 * SingleCPUUnits,
		MemoryMb: 8*MbInGb - AgentSize,
	},
	"t2.xlarge": {
		CPUUnits: 4 * SingleCPUUnits,
		MemoryMb: 16*MbInGb - AgentSize,
	},
	"t2.2xlarge": {
		CPUUnits: 8 * SingleCPUUnits,
		MemoryMb: 32*MbInGb - AgentSize,
	},
	"t3.nano": {
		CPUUnits: 2 * SingleCPUUnits,
		MemoryMb: MbInGb/2 - AgentSize,
	},
	"t3.micro": {
		CPUUnits: 2 * SingleCPUUnits,
		MemoryMb: 1*MbInGb - AgentSize,
	},
	"t3.small": {
		CPUUnits: 2 * SingleCPUUnits,
		MemoryMb: 2*MbInGb - AgentSize,
	},
	"t3.medium": {
		CPUUnits: 2 * SingleCPUUnits,
		MemoryMb: 4*MbInGb - AgentSize,
	},
	"t3.large": {
		CPUUnits: 2 * SingleCPUUnits,
		MemoryMb: 8*MbInGb - AgentSize,
	},
	"t3.xlarge": {
		CPUUnits: 4 * SingleCPUUnits,
		MemoryMb: 16*MbInGb - AgentSize,
	},
	"t3.2xlarge": {
		CPUUnits: 8 * SingleCPUUnits,
		MemoryMb: 32*MbInGb - AgentSize,
	},
}
