package environment

type TransportEnvironmentCapacitySet interface {
	BaseEnvironmentCapabilitySet
	SystemNetworkCapabilitySet
	InstanceNetworkCapabilitySet
}

type TransportEnvironment interface {
	TransportEnvironmentCapacitySet
	NarrowScope(key string) (TransportEnvironment, error)
	doNotImpl()
}
