module libcore

go 1.18

require (
	github.com/Dreamacro/clash v1.10.0
	github.com/miekg/dns v1.1.48
	github.com/sagernet/gomobile v0.0.0-20220214172500-89df302623c8
	github.com/sagernet/libping v0.1.1
	github.com/sirupsen/logrus v1.8.1
	github.com/ulikunitz/xz v0.5.10
	github.com/v2fly/v2ray-core/v5 v5.0.3
	go.uber.org/automaxprocs v1.5.1
	golang.org/x/sys v0.0.0-20220429233432-b5fbb4746d32
)

replace gvisor.dev/gvisor => github.com/sagernet/gvisor v0.0.0-20220213143053-df431bee78e3

replace github.com/v2fly/v2ray-core/v5 v5.0.0 => ../../v2ray-core

require (
	github.com/cheekybits/genny v1.0.0 // indirect
	github.com/dgryski/go-metro v0.0.0-20211217172704-adc40b04c140 // indirect
	github.com/ebfe/bcrypt_pbkdf v0.0.0-20140212075826-3c8d2dcb253a // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-chi/chi/v5 v5.0.7 // indirect
	github.com/go-chi/render v1.0.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.1 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jhump/protoreflect v1.12.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lucas-clemente/quic-go v0.27.0 // indirect
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40 // indirect
	github.com/marten-seemann/qtls-go1-16 v0.1.5 // indirect
	github.com/marten-seemann/qtls-go1-17 v0.1.1 // indirect
	github.com/marten-seemann/qtls-go1-18 v0.1.1 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pires/go-proxyproto v0.6.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/riobard/go-bloom v0.0.0-20200614022211-cdc8013cb5b3 // indirect
	github.com/seiflotfy/cuckoofilter v0.0.0-20220411075957-e3b120b3f5fb // indirect
	github.com/v2fly/BrowserBridge v0.0.0-20210430233438-0570fc1d7d08 // indirect
	github.com/v2fly/VSign v0.0.0-20201108000810-e2adc24bf848 // indirect
	github.com/v2fly/ss-bloomring v0.0.0-20210312155135-28617310f63e // indirect
	github.com/xtaci/smux v1.5.16 // indirect
	go.starlark.net v0.0.0-20220328144851-d1966c6b9fcd // indirect
	go4.org/intern v0.0.0-20220301175310-a089fc204883 // indirect
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20211027215541-db492cf91b37 // indirect
	golang.org/x/crypto v0.0.0-20220427172511-eb4f295cb31f // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.10 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	google.golang.org/genproto v0.0.0-20220429170224-98d788798c3e // indirect
	google.golang.org/grpc v1.46.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	inet.af/netaddr v0.0.0-20211027220019-c74959edd3b6 // indirect
)
