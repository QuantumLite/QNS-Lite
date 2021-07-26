package config

// The global configuration variable
var config Config

// The Config struct gathers all of the configuration parameters.
type Config struct {
	size        int
	memory      int
	lifetime    int
	numRequests int
	// The number of paths each request tries to reserve for itself.
	aggressiveness         int
	recoveryAggressiveness int
	opportunismDegree      int
	recoverySpan           int
	p_gen                  float64
	p_swap                 float64
	hasRecovery            bool
	recoveryHasContention  bool
	hasContention          bool
	isOpportunistic        bool
	isMultiPath            bool
}

// init() initializes the config variable once the package is imported.
func init() {
	config.size = 5
	config.memory = 6
	config.lifetime = 30
	config.numRequests = 20
	config.opportunismDegree = 1
	config.recoverySpan = 2
	config.p_gen = 0.5
	config.p_swap = 1
	config.hasRecovery = false
	config.hasContention = true
	config.recoveryHasContention = true
	config.isOpportunistic = false
	config.isMultiPath = true
	config.recoveryAggressiveness = 3
	if !config.isMultiPath {
		config.aggressiveness = 1
	} else {
		config.aggressiveness = 3
	}
}

// GetConfig returns the configuration.
func GetConfig() Config {
	return config
}

// SetConfig allows to manually set the configuration.
func SetConfig(size, memory, lifetime int, p_gen, p_swap float64, hasRecovery bool) {
	config.size = size
	config.memory = memory
	config.lifetime = lifetime
	config.p_gen = p_gen
	config.p_swap = p_swap
	config.hasRecovery = hasRecovery
}

func SetOpportunism(isOpportunistic bool) {
	config.isOpportunistic = isOpportunistic
}

func SetOpportunismDegree(opportunismDegree int) {
	config.opportunismDegree = opportunismDegree
}

func SetPGen(p_gen float64) {
	config.p_gen = p_gen
}

func SetPSwap(p_swap float64) {
	config.p_swap = p_swap
}

func SetSize(size int) {
	config.size = size
}

func SetLifetime(lifetime int) {
	config.lifetime = lifetime
}

func SetNumRequests(numRequests int) {
	config.numRequests = numRequests
}

// GetSize returns the size of the network.
func (conf Config) GetSize() int {
	return conf.size
}

// GetMemory returns the memory of every node.
func (conf Config) GetMemory() int {
	return conf.memory
}

// GetLifetime returns the life time of every entangled pair.
func (conf Config) GetLifetime() int {
	return conf.lifetime
}

// GetNumRequests returns the number of requests.
func (conf Config) GetNumRequests() int {
	return conf.numRequests
}

// GetPGen returns the probability of a successful entanglement generation.
func (conf Config) GetPGen() float64 {
	return conf.p_gen
}

// GetPSwap returns the probability of a successful entanglement swapping.
func (conf Config) GetPSwap() float64 {
	return conf.p_swap
}

// GetAggressiveness returns the aggressiveness of the requests.
func (conf Config) GetAggressiveness() int {
	return conf.aggressiveness
}

func (conf Config) GetRecoveryAggressiveness() int {
	return conf.recoveryAggressiveness
}

func (conf Config) GetOpportunismDegree() int {
	return conf.opportunismDegree
}

func (conf Config) GetRecoverySpan() int {
	return conf.recoverySpan
}

func (conf Config) GetHasRecovery() bool {
	return conf.hasRecovery
}

func (conf Config) GetHasContention() bool {
	return conf.hasContention
}

func (conf Config) GetRecoveryHasContention() bool {
	return conf.recoveryHasContention
}

func (conf Config) GetIsOpportunistic() bool {
	return conf.isOpportunistic
}

func (conf Config) GetIsMultiPath() bool {
	return conf.isMultiPath
}

// TODO: CHECK THIS!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func (conf *Config) SetAggressiveness(aggressiveness int) {
	conf.aggressiveness = aggressiveness
}

func SetAggressiveness(aggressiveness int) {
	config.aggressiveness = aggressiveness
}
