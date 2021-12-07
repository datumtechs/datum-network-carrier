package psutil

import "github.com/shirou/gopsutil/cpu"

// ReadCPUStats retrieves the current CPU stats.
func ReadCPUStats() (int64, int64, int64, error) {
	// passing false to request all cpu times
	timeStats, err := cpu.Times(false)
	if err != nil {
		log.Errorf("Could not read cpu stats err: %s", err)
		return int64(0), int64(0), int64(0), err
	}
	// requesting all cpu times will always return an array with only one time stats entry
	timeStat := timeStats[0]
	globalTime := int64((timeStat.User + timeStat.Nice + timeStat.System) * cpu.ClocksPerSec)
	globalWait := int64((timeStat.Iowait) * cpu.ClocksPerSec)
	localTime := int64(0)
	return globalTime, globalWait, localTime, nil
}