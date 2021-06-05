package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/MattHodge/go-teedy/evernote"

	"github.com/MattHodge/go-teedy/backup"
	"github.com/MattHodge/go-teedy/restore"
	"github.com/MattHodge/go-teedy/teedy"
	"github.com/alexflint/go-arg"
)

// values are replaced by go releaser
var version = "development"
var commit = "none"
var date = "none"

type BackupCmd struct {
	DestinationPath string `arg:"-d,required" placeholder:"DST" help:"Path to backup to"`
}

type RestoreCmd struct {
	SourcePath string `arg:"-s,required" placeholder:"SRC" help:"Path to restore from"`
}

type EvernoteCmd struct {
	SourceEnex string   `arg:"-e,--source-enex,required" placeholder:"ENEXFILE" help:"Path to evernote exported .enex file"`
	TagId      []string `arg:"-t" placeholder:"TAGID" help:"A tag ID from Teedy to add to the imported document"`
	Language   string   `arg:"-l" placeholder:"LANGUAGE" default:"eng"`
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

func (a *args) Version() string {
	return fmt.Sprintf("Version: %s\nCommit: %s\nDate: %s\n	", version, commit, date)
}

func main() {
	var args args
	arg.MustParse(&args)

	// validate params
	checkUrl(args.URL)
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

		err := validateLanguage(args.Evernote.Language)

		if err != nil {
			log.Fatalf("%v", err)
		}

		evClientOpts := []evernote.ImportClientOption{
			evernote.WithLanguage(args.Evernote.Language),
		}

		for _, tag := range args.Evernote.TagId {
			evClientOpts = append(evClientOpts, evernote.WithTagID(tag))
		}

		evClient := evernote.NewImportClient(args.Evernote.SourceEnex, client, evClientOpts...)

		_, err = evClient.Import()

		if err != nil {
			log.Fatalf("unable to import from evernote enex: %v", err)
		}
	}
}

func checkUrl(uri string) {
	_, err := url.ParseRequestURI(uri)

	if err != nil {
		log.Fatalf("URL '%s' is not valid", uri)
	}
}

func validateLanguage(language string) error {
	// sourced from: https://github.com/sismics/docs/blob/fd4c627c61ce40b5431d9a7e091043faa8a1dd30/README.md#available-environment-variables
	validLanguages := []string{
		"eng",
		"fra",
		"ita",
		"deu",
		"spa",
		"por",
		"pol",
		"rus",
		"ukr",
		"ara",
		"hin",
		"chi_sim",
		"chi_tra",
		"jpn",
		"tha",
		"kor",
		"nld",
		"tur",
		"heb",
		"hun",
		"fin",
		"swe",
		"lav",
		"dan",
	}

	for _, vl := range validLanguages {
		if vl == language {
			return nil
		}
	}

	return fmt.Errorf("provided language '%s' is not valid for Teedy. Choose from: %s", language, strings.Join(validLanguages, ", "))
}
