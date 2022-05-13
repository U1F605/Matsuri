package libcore

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"libcore/device"
	"os"
	"path/filepath"
	"strings"
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

func InitCore(internalAssets string, externalAssets string, prefix string, useOfficial BoolFunc, // extractV2RayAssets
	cachePath string, process string, //InitCore
	enableLog bool, maxKB int32, //SetEnableLog
) {
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
		externalAssetsPath = externalAssets
		internalAssetsPath = internalAssets
		assetsPrefix = prefix
		Setenv("v2ray.conf.geoloader", "memconservative")

		setupV2rayFileSystem(internalAssets, externalAssets)
		setupResolvers()

		// Setup CA Certs
		x509.SystemCertPool()
		roots := x509.NewCertPool()
		roots.AppendCertsFromPEM([]byte(mozillaCA))
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

	// CA for other programs
	go func() {
		f, err := os.OpenFile(filepath.Join(internalAssets, "ca.pem"), os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			forceLog("open ca.pem: " + err.Error())
		} else {
			if b, _ := ioutil.ReadAll(f); b == nil || string(b) != mozillaCA {
				f.Truncate(0)
				f.Seek(0, 0)
				f.Write([]byte(mozillaCA))
			}
			f.Close()
		}
	}()
}
