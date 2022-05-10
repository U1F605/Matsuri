package libcore

import (
	"context"
	"errors"
	"fmt"
	"libcore/protect"
	"log"
	gonet "net"
	"strconv"
	"strings"
	"sync"
	"time"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/app/dispatcher"
	"github.com/v2fly/v2ray-core/v5/common/buf"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/features/dns"
	dns_feature "github.com/v2fly/v2ray-core/v5/features/dns"
	v2rayDns "github.com/v2fly/v2ray-core/v5/features/dns"
	"github.com/v2fly/v2ray-core/v5/features/dns/localdns"
	"github.com/v2fly/v2ray-core/v5/features/routing"
	"github.com/v2fly/v2ray-core/v5/features/stats"
	"github.com/v2fly/v2ray-core/v5/infra/conf/serial"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
)

func GetV2RayVersion() string {
	return core.Version() + "-喵"
}

type V2RayInstance struct {
	access       sync.Mutex
	started      bool
	core         *core.Instance
	statsManager stats.Manager
	dispatcher   *dispatcher.DefaultDispatcher
	dnsClient    dns.Client
}

func NewV2rayInstance() *V2RayInstance {
	return &V2RayInstance{}
}

func (instance *V2RayInstance) LoadConfig(content string) error {
	if outdated != "" {
		return errors.New(outdated)
	}

	instance.access.Lock()
	defer instance.access.Unlock()

	config, err := serial.LoadJSONConfig(strings.NewReader(content))
	if err != nil {
		log.Println(content, err.Error())
		return err
	}

	c, err := core.New(config)
	if err != nil {
		return err
	}

	instance.core = c
	instance.statsManager = c.GetFeature(stats.ManagerType()).(stats.Manager)
	instance.dispatcher = c.GetFeature(routing.DispatcherType()).(routing.Dispatcher).(*dispatcher.DefaultDispatcher)
	instance.dnsClient = c.GetFeature(dns.ClientType()).(dns.Client)

	instance.setupDialer()
	return nil
}

func (instance *V2RayInstance) Start() error {
	instance.access.Lock()
	defer instance.access.Unlock()
	if instance.started {
		return newError("already started")
	}
	if instance.core == nil {
		return newError("not initialized")
	}
	err := instance.core.Start()
	if err != nil {
		return err
	}
	instance.started = true
	return nil
}

func (instance *V2RayInstance) QueryStats(tag string, direct string) int64 {
	if instance.statsManager == nil {
		return 0
	}
	counter := instance.statsManager.GetCounter(fmt.Sprintf("outbound>>>%s>>>traffic>>>%s", tag, direct))
	if counter == nil {
		return 0
	}
	return counter.Set(0)
}

func (instance *V2RayInstance) Close() error {
	instance.access.Lock()
	defer instance.access.Unlock()
	if instance.started {
		return instance.core.Close()
	}
	return nil
}

func (instance *V2RayInstance) dialContext(ctx context.Context, destination net.Destination) (net.Conn, error) {
	ctx = core.WithContext(ctx, instance.core)
	r, err := instance.dispatcher.Dispatch(ctx, destination)
	if err != nil {
		return nil, err
	}
	var readerOpt buf.ConnectionOption
	if destination.Network == net.Network_TCP {
		readerOpt = buf.ConnectionOutputMulti(r.Reader)
	} else {
		readerOpt = buf.ConnectionOutputMultiUDP(r.Reader)
	}
	return buf.NewConnection(buf.ConnectionInputMulti(r.Writer), readerOpt), nil
}

// Nekomura

var staticHosts = make(map[string][]net.IP)
var tryDomains = make([]string, 0)                                                    // server's domain, set when enhanced domain mode
var androidResolver = &net.Resolver{PreferGo: false}                                  // Using Android API, lookup from current network.
var androidUnderlyingResolver = &simpleSekaiWrapper{androidResolver: androidResolver} // Using Android API, lookup from non-VPN network.
var dc dns.Client

type simpleSekaiWrapper struct {
	androidResolver *net.Resolver
	sekaiResolver   LocalResolver // passed from java (only when VPNService)
}

func (p *simpleSekaiWrapper) LookupIP(network, host string) (ret []net.IP, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	ok := make(chan interface{})
	defer cancel()

	go func() {
		defer func() {
			select {
			case <-ctx.Done():
			default:
				ok <- nil
			}
			close(ok)
		}()
		ret, err = p.androidResolver.LookupIP(context.Background(), network, host)
	}()

	select {
	case <-ok:
		return
	}
}

// setup dialer and resolver for v2ray (v2ray options)
func (v2ray *V2RayInstance) setupDialer() {
	setupResolvers()
	dc = v2ray.dnsClient

	// All lookup except dnsClient -> dc.LookupIP()
	// and also set protectedDialer
	if _, ok := dc.(v2rayDns.ClientWithIPOption); ok {
		internet.UseAlternativeSystemDialer(&protect.ProtectedDialer{
			Resolver: func(domain string) ([]net.IP, error) {
				if ips, ok := staticHosts[domain]; ok && ips != nil {
					return ips, nil
				}

				return dc.LookupIP(&dns.MatsuriDomainStringEx{
					Domain:     domain,
					OptNetwork: "ip",
				})
			},
		})
	}
}

func setupResolvers() {
	// golang lookup -> androidResolver
	gonet.DefaultResolver = androidResolver

	// dnsClient lookup -> androidUnderlyingResolver.LookupIP()
	internet.UseAlternativeSystemDNSDialer(&protect.ProtectedDialer{
		Resolver: func(domain string) ([]net.IP, error) {
			return androidUnderlyingResolver.LookupIP("ip", domain)
		},
	})

	// "localhost" localDns lookup -> androidUnderlyingResolver.LookupIP()
	localdns.SetLookupFunc(androidUnderlyingResolver.LookupIP)
}
