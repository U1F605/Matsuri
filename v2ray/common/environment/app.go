package environment

type AppEnvironmentCapabilitySet interface {
	BaseEnvironmentCapabilitySet
	SystemNetworkCapabilitySet
	InstanceNetworkCapabilitySet
	FileSystemCapabilitySet
}

type AppEnvironment interface {
	AppEnvironmentCapabilitySet
	NarrowScope(key string) (AppEnvironment, error)
	doNotImpl()
}
