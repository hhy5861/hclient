package hclient

type (
	RunnerClient struct {
		opts   []Option
		client *ConfigCache
	}
)

func NewRunnerClient(opts ...Option) *RunnerClient {
	return &RunnerClient{opts: opts}
}

func (runner *RunnerClient) Name() string {
	return "hclient"
}

func (runner *RunnerClient) Start() error {
	runner.client = NewClient(runner.opts...)

	return nil
}

func (runner *RunnerClient) Shutdown() error {
	for k, _ := range runner.client.remotes {
		delete(runner.client.remotes, k)
	}

	return nil
}
