package expprof

import (
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuProfile	 *os.File
var ramProfile  *os.File
var blockProfile *os.File

//内存分析器的取样间隔 单位:字节 每X字节变化取样一次
var ramProfRate = 512 * 1024

//阻塞分析器的取样间隔 单位:个   每X个阻塞取样一次
var blockProfRate = 1

/**启动性能监控*/
func Initialize(cpuFile string,ramFile string,blockFile string) {
	var err error

	//cpu监控
	cpuProfile, err = os.Create(cpuFile)
	if err != nil {
		panic(err)
	}
	err = pprof.StartCPUProfile(cpuProfile)
	if err != nil {
		cpuProfile.Close()
		panic(err)
	}

	//内存监控
	ramProfile, err = os.Create(ramFile)
	if err != nil {
		panic(err)
	}
	runtime.MemProfileRate = ramProfRate

	//阻塞监控
	blockProfile, err = os.Create(blockFile)
	if err != nil {
		panic(err)
	}
	runtime.SetBlockProfileRate(blockProfRate)
}

func UnInitialize() {
	//结束cpu监控
	if cpuProfile != nil {
		pprof.StopCPUProfile()
		defer cpuProfile.Close()
	}

	//结束内存监控
	if ramProfile != nil {
		pprof.WriteHeapProfile(ramProfile)
		defer ramProfile.Close()
	}

	//结束阻塞文件监控
	if blockProfile != nil {
		pprof.Lookup("block").WriteTo(blockProfile, 0)
		defer blockProfile.Close()
	}
}
