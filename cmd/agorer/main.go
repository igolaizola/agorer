package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"time"

	"github.com/igolaizola/agorer"
	"github.com/igolaizola/agorer/pkg/agora"
	"github.com/igolaizola/agorer/pkg/example"
	"github.com/igolaizola/agorer/pkg/mail"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

// Build flags
var Version = ""
var Commit = ""
var Date = ""

func main() {
	// Create signal based context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Launch command
	cmd := newCommand()
	if err := cmd.ParseAndRun(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func newCommand() *ffcli.Command {
	fs := flag.NewFlagSet("agorer", flag.ExitOnError)

	return &ffcli.Command{
		ShortUsage: "agorer [flags] <subcommand>",
		FlagSet:    fs,
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
		Subcommands: []*ffcli.Command{
			newVersionCommand(),
			newStockCommand(),
			newSalesCommand(),
			newMockServeCommand(),
			newExampleCommand(),
			newMailCommand(),
		},
	}
}

func newVersionCommand() *ffcli.Command {
	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "igogpt version",
		ShortHelp:  "print version",
		Exec: func(ctx context.Context, args []string) error {
			v := Version
			if v == "" {
				if buildInfo, ok := debug.ReadBuildInfo(); ok {
					v = buildInfo.Main.Version
				}
			}
			if v == "" {
				v = "dev"
			}
			versionFields := []string{v}
			if Commit != "" {
				versionFields = append(versionFields, Commit)
			}
			if Date != "" {
				versionFields = append(versionFields, Date)
			}
			fmt.Println(strings.Join(versionFields, " "))
			return nil
		},
	}
}

func newStockCommand() *ffcli.Command {
	cmd := "stock"
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")

	var cfg agorer.Config
	fs.BoolVar(&cfg.Debug, "debug", false, "debug mode")
	fs.StringVar(&cfg.LogDir, "log-dir", "logs", "output directory")
	fs.StringVar(&cfg.Input, "input", "", "input file or URL")
	fs.StringVar(&cfg.InputType, "input-type", "", "input type (json, agora, agora-json)")
	fs.StringVar(&cfg.Output, "output", "", "output file")
	fs.StringVar(&cfg.OutputType, "output-type", "", "output type (json, sinli)")

	// Agora parameters
	fs.StringVar(&cfg.AgoraToken, "agora-token", "", "agora token")
	// ISBN parameters
	fs.StringVar(&cfg.ISBNDir, "isbn-dir", "data", "isbn directory")

	// Mail parameters
	fs.BoolVar(&cfg.Mail.Dry, "mail-dry", false, "dry run, don't send mail")
	fs.StringVar(&cfg.Mail.Host, "mail-host", "", "mail smtp host")
	fs.IntVar(&cfg.Mail.Port, "mail-port", 0, "mail smtp port")
	fs.StringVar(&cfg.Mail.Username, "mail-user", "", "mail smtp username")
	fs.StringVar(&cfg.Mail.Password, "mail-pass", "", "mail smtp password")

	// SINLI parameters
	fs.StringVar(&cfg.SINLISourceEmail, "sinli-source-email", "", "sinli source email")
	fs.StringVar(&cfg.SINLISourceID, "sinli-source-id", "", "sinli source id")
	fs.StringVar(&cfg.SINLIDestinationEmail, "sinli-destination-email", "", "sinli destination email")
	fs.StringVar(&cfg.SINLIDestinationID, "sinli-destination-id", "", "sinli destination id")
	fs.StringVar(&cfg.SINLIClientName, "sinli-client-name", "", "sinli client name")

	return &ffcli.Command{
		Name:       cmd,
		ShortUsage: fmt.Sprintf("agorer %s [flags] <key> <value data...>", cmd),
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithEnvVarPrefix("AGORER"),
		},
		ShortHelp: fmt.Sprintf("%s agorer command", cmd),
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			return agorer.Stock(ctx, &cfg)
		},
	}
}

func newSalesCommand() *ffcli.Command {
	cmd := "sales"
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")

	var day string
	fs.StringVar(&day, "day", time.Now().UTC().Format("2006-01-02"), "day to process")

	var cfg agorer.Config
	fs.BoolVar(&cfg.Debug, "debug", false, "debug mode")
	fs.StringVar(&cfg.LogDir, "log-dir", "logs", "output directory")
	fs.StringVar(&cfg.Input, "input", "", "input file or URL")
	fs.StringVar(&cfg.InputType, "input-type", "", "input type (json, agora, agora-json)")
	fs.StringVar(&cfg.Output, "output", "", "output file")
	fs.StringVar(&cfg.OutputType, "output-type", "", "output type (json, sinli)")

	// Agora parameters
	fs.StringVar(&cfg.AgoraToken, "agora-token", "", "agora token")
	// ISBN parameters
	fs.StringVar(&cfg.ISBNDir, "isbn-dir", "data", "isbn directory")

	// Mail parameters
	fs.BoolVar(&cfg.Mail.Dry, "mail-dry", false, "dry run, don't send mail")
	fs.StringVar(&cfg.Mail.Host, "mail-host", "", "mail smtp host")
	fs.IntVar(&cfg.Mail.Port, "mail-port", 0, "mail smtp port")
	fs.StringVar(&cfg.Mail.Username, "mail-user", "", "mail smtp username")
	fs.StringVar(&cfg.Mail.Password, "mail-pass", "", "mail smtp password")

	// SINLI parameters
	fs.StringVar(&cfg.SINLISourceEmail, "sinli-source-email", "", "sinli source email")
	fs.StringVar(&cfg.SINLISourceID, "sinli-source-id", "", "sinli source id")
	fs.StringVar(&cfg.SINLIDestinationEmail, "sinli-destination-email", "", "sinli destination email")
	fs.StringVar(&cfg.SINLIDestinationID, "sinli-destination-id", "", "sinli destination id")
	fs.StringVar(&cfg.SINLIClientName, "sinli-client-name", "", "sinli client name")

	return &ffcli.Command{
		Name:       cmd,
		ShortUsage: fmt.Sprintf("agorer %s [flags] <key> <value data...>", cmd),
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithEnvVarPrefix("AGORER"),
		},
		ShortHelp: fmt.Sprintf("%s agorer command", cmd),
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			d, err := time.Parse("2006-01-02", day)
			if err != nil {
				return fmt.Errorf("couldn't parse day: %w", err)
			}
			return agorer.Sales(ctx, &cfg, d)
		},
	}
}

func newMockServeCommand() *ffcli.Command {
	cmd := "mock-serve"
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")

	var addr, master string
	fs.StringVar(&addr, "addr", "", "address to listen to")
	fs.StringVar(&master, "master", "", "master file")

	return &ffcli.Command{
		Name:       cmd,
		ShortUsage: fmt.Sprintf("agorer %s [flags] <key> <value data...>", cmd),
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithEnvVarPrefix("AGORER"),
		},
		ShortHelp: fmt.Sprintf("%s agorer command", cmd),
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			if _, err := agora.MockServe(ctx, addr, master); err != nil {
				return err
			}
			<-ctx.Done()
			return nil
		},
	}
}

func newExampleCommand() *ffcli.Command {
	cmd := "example"
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")

	var cfg example.Config
	fs.StringVar(&cfg.SourceEmail, "source-email", "", "source email")
	fs.StringVar(&cfg.DestinationEmail, "destination-email", "", "destination email")
	fs.StringVar(&cfg.SourceID, "source-id", "", "source id")
	fs.StringVar(&cfg.DestinationID, "destination-id", "", "destination id")
	fs.StringVar(&cfg.ProviderName, "provider-name", "", "provider name")
	fs.StringVar(&cfg.ClientName, "client-name", "", "client name")
	fs.StringVar(&cfg.MasterFile, "master-file", "", "master file")
	fs.StringVar(&cfg.DayFile, "day-file", "", "day file")
	fs.StringVar(&cfg.OutputDir, "data", "", "output dir")

	return &ffcli.Command{
		Name:       cmd,
		ShortUsage: fmt.Sprintf("agorer %s [flags] <key> <value data...>", cmd),
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithEnvVarPrefix("AGORER"),
		},
		ShortHelp: fmt.Sprintf("%s agorer command", cmd),
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			return example.Run(ctx, &cfg)
		},
	}
}

func newMailCommand() *ffcli.Command {
	cmd := "mail"
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")

	var cfg mail.Config
	fs.StringVar(&cfg.Host, "host", "", "smtp host")
	fs.IntVar(&cfg.Port, "port", 0, "smtp port")
	fs.StringVar(&cfg.Username, "user", "", "smtp username")
	fs.StringVar(&cfg.Password, "pass", "", "smtp password")
	var from, to, subject, body, file string
	fs.StringVar(&from, "from", "", "from email")
	fs.StringVar(&to, "to", "", "to email")
	fs.StringVar(&subject, "subject", "", "email subject")
	fs.StringVar(&body, "body", "", "email body")
	fs.StringVar(&file, "file", "", "file to attach")

	return &ffcli.Command{
		Name:       cmd,
		ShortUsage: fmt.Sprintf("agorer %s [flags] <key> <value data...>", cmd),
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithEnvVarPrefix("AGORER"),
		},
		ShortHelp: fmt.Sprintf("%s agorer command", cmd),
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			// Read subject from file
			b, err := os.ReadFile(subject)
			if err != nil {
				return fmt.Errorf("couldn't read subject file: %w", err)
			}
			sub := strings.TrimSpace(string(b))
			return mail.Send(ctx, &cfg, from, to, string(sub), body, file)
		},
	}
}
