package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/guoyk93/bastion/sshd"
	"github.com/guoyk93/bastion/types"
	"os"
	"os/signal"
	"syscall"
)

var (
	dev          bool
	overrideKeys bool
	optionsFile  string
	options      types.Options
)

func main() {
	var err error

	// init logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true})

	// load command-line options
	flag.BoolVar(&dev, "dev", false, "dev mode")
	flag.StringVar(&optionsFile, "c", "/etc/bastion/bastion.yml", "bastion config file")
	flag.BoolVar(&overrideKeys, "override-keys", false, "run a one-off operation to override keys for key-managed nodes")
	flag.Parse()

	// load options files
	log.Info().Str("file", optionsFile).Msg("loading options file")
	if options, err = types.LoadOptions(optionsFile); err != nil {
		log.Error().Str("file", optionsFile).Err(err).Msg("failed to load options file")
		os.Exit(1)
		return
	}
	// merge command line options
	if dev {
		options.SSHD.Dev = true
	}
	log.Info().Interface("options", options.SSHD).Msg("options file loaded")

	// adjust logger
	if options.SSHD.Dev {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	// create daemon
	d := sshd.New(options.SSHD)

	// one-off operation - override keys
	if overrideKeys {
		if err = d.OverrideKeys(); err != nil {
			log.Error().Err(err).Msg("failed to override keys")
			os.Exit(1)
		}
	} else {
		// run the signalHandler
		go signalHandler(d)

		// run the sshd
		if err = d.Run(); err != nil {
			log.Error().Err(err).Msg("exited")
			os.Exit(1)
		}
	}
}

func signalHandler(d *sshd.SSHD) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	s := <-shutdown
	log.Info().Str("signal", s.String()).Msg("signal received")
	d.Shutdown()
}
