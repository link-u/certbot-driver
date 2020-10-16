package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
)

const ImageName = "docker.io/certbot/dns-route53:latest"

type CliConfig struct {
	StandardLog bool
}

type AwsConfig struct {
	IamPath string
}

type CertConfig struct {
	Directory string
	Domains   []string
}

type Config struct {
	EmailAddress string
	Cli          CliConfig
	Cert         CertConfig
	Aws          AwsConfig
}

func usage(custom func()) {
	_, _ = fmt.Fprintf(os.Stderr, "certbot-bot create [DOMAINS]...\n")
	_, _ = fmt.Fprintf(os.Stderr, "certbot-bot renew\n")
	flag.PrintDefaults()
	if custom != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\n")
		custom()
	}
	os.Exit(-1)
}

func main() {
	var config Config
	var log *zap.Logger
	var err error
	var app *kingpin.Application
	var createCommand *kingpin.CmdClause
	var renewCommand *kingpin.CmdClause
	var command string
	var domains []string
	{ // Flags
		app = kingpin.New("certbot-bot", "Control certbot automatically").Version("v1.0.0")
		{
			createCommand = app.Command("create", "create new certs")
			createCommand.Flag("cli.standard-log", "Print logs in standard format, not in json").
				BoolVar(&config.Cli.StandardLog)
			createCommand.Flag("cert.directory", "Directory to store the certs").
				PlaceHolder("(path/to/cert)").
				Required().
				StringVar(&config.Cert.Directory)
			createCommand.Flag("email-address", "your email address").
				PlaceHolder("(aoba@example.com)").
				Required().
				StringVar(&config.EmailAddress)
			createCommand.Flag("aws.iam", "Path to IAM").
				PlaceHolder("(iam.conf)").
				Required().
				ExistingFileVar(&config.Aws.IamPath)
			createCommand.Arg("domains", "target domains").Required().StringsVar(&domains)
		}
		{
			renewCommand = app.Command("renew", "renew existing certs")
			renewCommand.Flag("cli.standard-log", "Print logs in standard format, not in json").
				BoolVar(&config.Cli.StandardLog)
			renewCommand.Flag("cert.directory", "Directory to store the certs").
				PlaceHolder("(path/to/cert)").
				Required().
				ExistingDirVar(&config.Cert.Directory)
			renewCommand.Flag("aws.iam", "Path to IAM").
				PlaceHolder("(iam.conf)").
				Required().
				ExistingFileVar(&config.Aws.IamPath)
		}
		command, err = app.Parse(os.Args[1:])
		if err != nil {
			app.Usage(os.Args[1:])
			panic(fmt.Sprintf("Failed to parse args: %v", err))
		}
	}
	// Check weather terminal or not
	if config.Cli.StandardLog || isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		log, err = zap.NewDevelopment()
	} else {
		log, err = zap.NewProduction()
	}
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	undo := zap.ReplaceGlobals(log)
	defer undo()
	log.Info("Log System Initialized.")

	switch command {
	case createCommand.FullCommand():
		log.Info("create new domains", zap.Strings("domains", domains))
		createCerts(&config, domains)
	case renewCommand.FullCommand():
		log.Info("renew existing domains", zap.Strings("domains", domains))
		renewCerts(&config)
	default:
		log.Fatal("Unknown command", zap.String("command", command))
	}

}
