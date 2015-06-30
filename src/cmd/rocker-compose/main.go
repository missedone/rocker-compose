package main

import (
	"compose"
	"fmt"
	"os"
	"path"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	app := cli.NewApp()
	app.Name = "rocker-compose"
	app.Version = "0.0.1"
	app.Usage = "Tool for docker orchestration"
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Usage:    "execute manifest",
			Action: run,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "log, l",
				},
				cli.StringFlag{
					Name: "manifest, m",
				},
				cli.BoolFlag{
					Name: "verbose, v",
				},
			},
		},
	}
	app.Run(os.Args)
}

func run(ctx *cli.Context) {
	if ctx.Bool("verbose") {
		log.SetLevel(log.DebugLevel)
	}

	if logFilename, err := toAbsolutePath(ctx.String("log")); err == nil {
		logFile, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Warnf("Cannot initialize log file %s due to error %s", logFilename, err)
		}

		if path.Ext(logFilename) == "json" {
			log.SetFormatter(&log.JSONFormatter{})
		}

		log.SetOutput(logFile)
	}

	if configFilename, err := toAbsolutePath(ctx.String("manifest")); err != nil {
		log.Fatal(err)
		//		os.Exit(1) // no config - no pichenka
	} else {
		config, err := compose.ReadConfigFile(configFilename, map[string]interface{}{})
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("Config path: %s\n", configFilename)
		log.Infof("Config: %+q\n", config)
	}

	// if c.GlobalIsSet("tlsverify") {
	//   config.tlsverify = c.GlobalBool("tlsverify")
	//   config.tlscacert = globalString(c, "tlscacert")
	//   config.tlscert = globalString(c, "tlscert")
	//   config.tlskey = globalString(c, "tlskey")
	// }
}

func toAbsolutePath(filePath string) (string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return filePath, fmt.Errorf("No such file or directory: %s", filePath)
	}

	if !path.IsAbs(filePath) {
		wd, err := os.Getwd()
		if err != nil {
			log.Errorf("Cannot get absolute path to %s due to error %s", filePath, err)
			return filePath, err
		}
		return path.Join(wd, filePath), nil
	}
	return filePath, nil
}

// Fix string arguments enclosed with boudle quotes
// 'docker-machine config' gives such arguments
// func globalString(c *cli.Context, name string) string {
// 	str := c.GlobalString(name)
// 	if len(str) >= 2 && str[0] == '\u0022' && str[len(str)-1] == '\u0022' {
// 		str = str[1 : len(str)-1]
// 	}
// 	return str
// }
