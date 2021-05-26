package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/MattHodge/go-teedy/evernote"

	"github.com/MattHodge/go-teedy/backup"
	"github.com/MattHodge/go-teedy/restore"
	"github.com/MattHodge/go-teedy/teedy"
	"github.com/alexflint/go-arg"
)

type BackupCmd struct {
	DestinationPath string `arg:"-d,required" placeholder:"DST" help:"Path to backup to"`
}

type RestoreCmd struct {
	SourcePath string `arg:"-s,required" placeholder:"SRC" help:"Path to restore from"`
}

type EvernoteCmd struct {
	SourceEnex string `args:"-e,required" placeholder:"ENEXFILE" help:"Path to evernote exported .enex file"`
}

type args struct {
	Backup   *BackupCmd   `arg:"subcommand:backup"`
	Restore  *RestoreCmd  `arg:"subcommand:restore"`
	Evernote *EvernoteCmd `arg:"subcommand:evernote"`
	URL      string       `arg:"-u,required" help:"Teedy Server URL"`
}

func (a *args) Description() string {
	return "teedy-cli allows you to backup from and restore to a https://teedy.io/ server.\n"
}

func main() {
	var args args
	arg.MustParse(&args)

	// validate params
	CheckURL(args.URL)
	username := teedy.GetEnvMustExist("TEEDY_USERNAME")
	password := teedy.GetEnvMustExist("TEEDY_PASSWORD")

	client, err := teedy.NewClient(args.URL, username, password)

	if err != nil {
		log.Fatalf("can't get teedy client: %v", err)
	}

	switch {
	case args.Backup != nil:
		fmt.Printf("Running backup of %s to '%s'\n", args.URL, args.Backup.DestinationPath)
		backupClient := backup.NewBackupClient(client, args.Backup.DestinationPath)

		err = backupClient.Tags()

		if err != nil {
			log.Fatalf("unable to backup tags: %v", err)
		}

		err = backupClient.Documents()

		if err != nil {
			log.Fatalf("unable to backup documents: %v", err)
		}
	case args.Restore != nil:
		fmt.Printf("Running restore of '%s' to '%s'\n", args.Restore.SourcePath, args.URL)

		restoreClient := restore.NewRestoreClient(client, args.Restore.SourcePath)
		err = restoreClient.Tags()

		if err != nil {
			log.Fatalf("unable to restore tags: %v", err)
		}

		_, err = restoreClient.Documents()

		if err != nil {
			log.Fatalf("unable to restore docs: %v", err)
		}
	case args.Evernote != nil:
		fmt.Printf("Running import of Evernote file '%s' to '%s'\n", args.Evernote.SourceEnex, args.URL)

		evClient := evernote.NewImportClient(args.Evernote.SourceEnex, client)

		_, err := evClient.Import()

		if err != nil {
			log.Fatalf("unable to import from evernote enex: %v", err)
		}
	}
}

func CheckURL(uri string) {
	_, err := url.ParseRequestURI(uri)

	if err != nil {
		log.Fatalf("URL '%s' is not valid", uri)
	}
}
