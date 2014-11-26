package command_factory

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

type appRunner interface {
	StartDockerApp(name, startCommand, dockerImagePath string) error
	ScaleDockerApp(name string, instances int) error
	StopDockerApp(name string) error
	IsDockerAppUp(name string) (bool, error)
}

type AppRunnerCommandFactory struct {
	appRunnerCommand *appRunnerCommand
}


func NewAppRunnerCommandFactory(appRunner appRunner, output io.Writer, timeout int) *AppRunnerCommandFactory {
	return &AppRunnerCommandFactory{&appRunnerCommand{appRunner, output, timeout}}
}

func (commandFactory *AppRunnerCommandFactory) MakeStartDiegoAppCommand() cli.Command {

	var startFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "docker-image, i",
			Usage: "the docker image to run",
		},
		cli.StringFlag{
			Name:  "start-command, c",
			Usage: "the command to run in the context of the docker image (ie the start command for the app)",
		},
	}

	var startCommand = cli.Command{
		Name:        "start",
		ShortName:   "s",
		Description: "Start a docker app on diego",
		Usage:       "diego-edge-cli start APP_NAME -i DOCKER_IMAGE -c START_COMMAND",
		Action:      commandFactory.appRunnerCommand.startDiegoApp,
		Flags:       startFlags,
	}

	return startCommand
}

func (commandFactory *AppRunnerCommandFactory) MakeScaleDiegoAppCommand() cli.Command {

	var scaleFlags = []cli.Flag{
		cli.IntFlag{
			Name:  "instances, i",
			Usage: "the number of instances to scale to",
		},
	}

	var scaleCommand = cli.Command{
		Name:        "scale",
		Description: "Scale a docker app on diego",
		Usage:       "diego-edge-cli scale APP_NAME --instances NUM_INSTANCES ",
		Action:      commandFactory.appRunnerCommand.scaleDiegoApp,
		Flags:       scaleFlags,
	}

	return scaleCommand
}

func (commandFactory *AppRunnerCommandFactory) MakeStopDiegoAppCommand() cli.Command {

	var stopCommand = cli.Command{
		Name:        "stop",
		Description: "Stop a docker app on diego",
		Usage:       "diego-edge-cli stop APP_NAME",
		Action:      commandFactory.appRunnerCommand.stopDiegoApp,
	}

	return stopCommand
}

type appRunnerCommand struct {
	appRunner appRunner
	output    io.Writer
	timeout   int
}

func (cmd *appRunnerCommand) startDiegoApp(c *cli.Context) {
	startCommand := c.String("start-command")
	dockerImage := c.String("docker-image")
	name := c.Args().First()

	if name == "" || dockerImage == "" || startCommand == "" {
		cmd.incorrectUsage()
		return
	} else if !strings.HasPrefix(dockerImage, "docker:///") {
		cmd.incorrectUsage()
		cmd.say("Docker Image should begin with: docker:///")
		return
	}

	err := cmd.appRunner.StartDockerApp(name, startCommand, dockerImage)

	if err != nil {
		cmd.say(fmt.Sprintf("Error Starting App: %s", err))
		return
	}

	cmd.say("Starting App: " + name)
	for i := 0; i < cmd.timeout; i++ {
		if status, _ := cmd.appRunner.IsDockerAppUp(name); status {
			cmd.newLine()
			cmd.say(name + " is now running. \n")
			return
		} else {
			cmd.dot()
		}
		time.Sleep(time.Second)
	}
	cmd.newLine()
	cmd.say(name + " took too long to start. \n")
}

func (cmd *appRunnerCommand) scaleDiegoApp(c *cli.Context) {
	instances := c.Int("instances")
	appName := c.Args().First()
	if appName == "" {
		cmd.incorrectUsage()
		return
	} else if instances == 0 {
		cmd.say(fmt.Sprintf("Error Scaling to 0 instances - Please stop with: diego-edge-cli stop cool-web-app"))
		return
	}

	err := cmd.appRunner.ScaleDockerApp(appName, instances)

	if err != nil {
		cmd.say(fmt.Sprintf("Error Scaling App: %s", err))
		return
	}

	cmd.say("App Scaled Successfully")
}

func (cmd *appRunnerCommand) stopDiegoApp(c *cli.Context) {
	appName := c.Args().First()
	if appName == "" {
		cmd.incorrectUsage()
		return
	}

	err := cmd.appRunner.StopDockerApp(appName)

	if err != nil {
		cmd.say(fmt.Sprintf("Error Stopping App: %s", err))
		return
	}

	cmd.say("App Stopped Successfully")
}

func (cmd *appRunnerCommand) say(output string) {
	cmd.output.Write([]byte(output))
}

func (cmd *appRunnerCommand) incorrectUsage() {
	cmd.say("Incorrect Usage\n")
}

func (cmd *appRunnerCommand) dot() {
	cmd.output.Write([]byte("."))
}

func (cmd *appRunnerCommand) newLine() {
	cmd.output.Write([]byte("\n"))
}
