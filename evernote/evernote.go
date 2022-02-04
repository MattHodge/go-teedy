package evernote

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/MattHodge/go-teedy/teedy"
	"github.com/asaskevich/govalidator"
	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/file"
)

var reFileAndExt = regexp.MustCompile(`(.*)(\.[\w\d]+)`)

type ImportClient struct {
	client   *teedy.Client
	source   string
	tagId    []string
	language string
}

type ImportClientOption func(*ImportClient)

func WithTagID(tagId string) ImportClientOption {
	return func(ic *ImportClient) {
		ic.tagId = append(ic.tagId, tagId)
	}
}

func WithLanguage(language string) ImportClientOption {
	return func(ic *ImportClient) {
		ic.language = language
	}
}

func NewImportClient(enexFile string, client *teedy.Client, opts ...ImportClientOption) *ImportClient {
	ic := &ImportClient{
		client:   client,
		source:   enexFile,
		language: "eng",
	}

	for _, opt := range opts {
		opt(ic)
	}

	return ic
}

func (ic *ImportClient) Import() ([]*teedy.Document, error) {
	export, err := load(ic.source)

	if err != nil {
		return nil, fmt.Errorf("unable to import evernote export: %v", err)
	}

	var retDocuments []*teedy.Document

	for _, note := range export.Notes {
		doc := teedy.NewDocument(note.Title, ic.language)
		doc.Description = string(note.Content)

		if note.Tags != nil { //if evernote tags exist, iterate through them and find if they already exist; if not, add them
			for _, word := range note.Tags {
				newtag, err := ic.client.Tag.GetByName(word)

				if err != nil {
					fmt.Printf("Tag GetByName error: %s\n", err)
					continue
				}

				if newtag == nil { //Tag not found; make a new tag
					addedtag, err := teedy.NewTag(word, "#4caf50", "") //fill out the TAG structure, tag colour green for Evernote, no parent

					if err != nil {
						return nil, fmt.Errorf("error in newtag structure for %s: %v", word, err)
					} else {
						resulttag, err := ic.client.Tag.Add(addedtag) //add the tag to Teedy

						if err != nil {
							return nil, fmt.Errorf("unable to create tag %s: %v", word, err)
						} else {
							doc.Tags = append(doc.Tags, &teedy.Tag{Id: resulttag.Id}) //add the new tag ID to the new document.
						}
					}
				} else { //Tag was found,add the id to the document.
					doc.Tags = append(doc.Tags, &teedy.Tag{Id: newtag.Id})
				}
			}
		}

		// add any command line specified tags
		for _, tagId := range ic.tagId {
			doc.Tags = append(doc.Tags, &teedy.Tag{Id: tagId})
		}

		convertedTime, err := convertDate(note.Created)

		if err != nil {
			return nil, fmt.Errorf("unable to convert date '%s' from note '%s'", note.Created, note.Title)
		}

		doc.CreateDate = convertedTime

		//Truncate input with limits taken from Teedy's html note input page.
		titlelimit := 100
		desclimit := 4000

		titletruncated := false
		desctruncated := false
		originaltitle := doc.Title

		for i := titlelimit; len(doc.Title) > titlelimit; i-- { //Truncate truncates to nearest word, occasionally exceeding the max length
			doc.Title = govalidator.Truncate(doc.Title, i, "") //Reduce the truncated length until it fits
			titletruncated = true
		}

		if titletruncated {
			fmt.Printf("Title truncated for note '%s'\n", originaltitle)
		}

		for i := desclimit; len(doc.Description) > desclimit; i-- { //Do the same for description
			doc.Description = govalidator.Truncate(doc.Description, i, "")
			desctruncated = true
		}

		if desctruncated {
			fmt.Printf("Description truncated for note '%s'\n", originaltitle)
		}

		addedDoc, err := ic.client.Document.Add(doc)

		retDocuments = append(retDocuments, addedDoc)

		if err != nil {
			return nil, fmt.Errorf("unable to add document to teedy for note '%s': %v", note.Title, err)
		}

		// resources are the attachments to a note
		for _, res := range note.Resources {
			fullFileName := fileName(res)

			_, err := ic.client.File.AddReader(addedDoc.Id, "", fullFileName, res.Mime, decoder(res.Data))

			if err != nil {
				return nil, fmt.Errorf("unable to upload file '%s' teedy from from note '%s': %v", fullFileName, note.Title, err)
			}
		}
	}
	return retDocuments, nil
}

func load(enexFile string) (*enex.Export, error) {
	enexContent, err := os.Open(enexFile)
	if err != nil {
		return nil, fmt.Errorf("cannot load '%s': %v", enexFile, err)
	}

	return enex.Decode(enexContent)
}

// decoder converts evernote attachment data from base64 into an io.Reader
func decoder(d enex.Data) io.Reader {
	if d.Encoding == "base64" {
		return base64.NewDecoder(base64.StdEncoding, bytes.NewReader(bytes.TrimSpace(d.Content)))
	}

	return bytes.NewReader(d.Content)
}

// guessName of the res with the following priority:
// 1. Filename attribute
// 2. SourceUrl attribute
// 3. ID of the res
// 4. File type as name
func guessName(r enex.Resource) string {
	switch {
	case r.Attributes.Filename != "":
		return r.Attributes.Filename
	case r.Attributes.SourceUrl != "":
		return strings.TrimSpace(path.Base(r.Attributes.SourceUrl))
	case r.ID != "":
		return r.ID
	default:
		return r.Type
	}
}

func guessExt(mimeType string) string {
	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(ext) == 0 {
		return ""
	}
	return ext[0]
}

// fileName takes an evernote resources and returns a file name with extension
func fileName(r enex.Resource) string {
	name := guessName(r)
	// Try to split a file into name and extension
	ff := reFileAndExt.FindStringSubmatch(name)
	if len(ff) < 2 {
		// Guess the extension by the mime type
		fileName := file.BaseName(name)
		extension := guessExt(r.Mime)
		return fmt.Sprintf("%s%s", fileName, extension)
	}

	// Return sanitized filename
	fileName := file.BaseName(ff[len(ff)-2])

	// further sanitize
	fileName = govalidator.SafeFileName(fileName)

	extension := ff[len(ff)-1]
	return fmt.Sprintf("%s%s", fileName, extension)
}

// convertDate takes an evernote note date and converts it to a teedy timestamp
func convertDate(date string) (*teedy.Timestamp, error) {
	// 20060102T150405Z
	layout := "20060102T150405Z"
	t, err := time.Parse(layout, date)

	if err != nil {
		return nil, err
	}

	res := teedy.Timestamp{
		Time: t,
	}

	return &res, nil
}
