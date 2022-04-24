package libcore

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"libcore/device"
	"os"
	"path/filepath"
	"strings"
	"time"
	_ "unsafe"

	"github.com/sagernet/libping"
)

//go:linkname systemRoots crypto/x509.systemRoots
var systemRoots *x509.CertPool

func Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func Unsetenv(key string) error {
	return os.Unsetenv(key)
}

func IcmpPing(address string, timeout int32) (int32, error) {
	return libping.IcmpPing(address, timeout)
}

func initCoreDefer() {
	device.AllDefer("InitCore", forceLog)
}

func InitCore(internalAssets string, externalAssets string, prefix string, useOfficial BoolFunc, // extractV2RayAssets
	cachePath string, process string, //InitCore
	enableLog bool, maxKB int32, //SetEnableLog
) {
	defer initCoreDefer()

	isBgProcess := strings.HasSuffix(process, ":bg")

	// Set up log
	SetEnableLog(enableLog, maxKB)
	s := fmt.Sprintln("InitCore called", externalAssets, cachePath, process, isBgProcess)
	err := setupLogger(filepath.Join(cachePath, "neko.log"))

	if err == nil {
		go forceLog(s)
	} else {
		// not fatal
		forceLog(fmt.Sprintln("Log not inited:", s, err.Error()))
	}

	// Set up some component
	go func() {
		defer initCoreDefer()
		device.GoDebug(process)

		externalAssetsPath = externalAssets
		internalAssetsPath = internalAssets
		assetsPrefix = prefix
		Setenv("v2ray.conf.geoloader", "memconservative")

		setupV2rayFileSystem(internalAssets, externalAssets)
		setupResolvers()

		if time.Now().Unix() >= GetExpireTime() {
			outdated = "Your version is too old! Please update!! 版本太旧，请升级！"
		} else if time.Now().Unix() < (GetBuildTime() - 86400) {
			outdated = "Wrong system time! 系统时间错误！"
		}

		// Setup CA Certs
		x509.SystemCertPool()
		roots := x509.NewCertPool()
		systemRoots = roots

		// Extract assets
		if isBgProcess {
			extractV2RayAssets(useOfficial)
		}
	}()

	if !isBgProcess {
		return
	}

	device.AutoGoMaxProcs()
}
