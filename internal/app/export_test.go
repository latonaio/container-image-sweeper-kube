package app

func NewAppForTest(cmdArgs *CommandArgs) *App {
	return &App{
		cmdArgs: cmdArgs,
	}
}

var ProcessForTest = (*App).process

func NewDockerManagerForTest(client IDockerClient) *DockerManager {
	return &DockerManager{
		client: client,
	}
}
