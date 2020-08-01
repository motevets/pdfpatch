package manifest

import (
	"io/ioutil"
	"net/url"
	"path"

	"gopkg.in/yaml.v2"
)

type Manifest struct {
	Sources []Source
}

type Source struct {
	URL          string
	Md5Sum       string
	PatchedFiles []string `yaml:"patched_files"`
}

func ParseFile(path string) (theManifest Manifest, err error) {
	manifestData, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(manifestData, &theManifest)
	if err != nil {
		return
	}

	return
}

func (m Manifest) SourceFileNames() (fileNames []string, err error) {
	fileNames = make([]string, len(m.Sources))
	for i, source := range m.Sources {
		fileNames[i], err = source.FileName()
		if err != nil {
			return nil, err
		}
	}
	return
}

func (source Source) FileName() (string, error) {
	parsedURL, err := url.Parse(source.URL)
	if err != nil {
		return "", err
	}
	return path.Base(parsedURL.Path), nil
}
