package main

import (
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"text/template"

	"gopkg.in/yaml.v2"
)

type TalkType string

const (
	CommunityTalk  TalkType = "COMMUNITY_TYPE"
	CorporateTalk  TalkType = "CORPORATE_TALK"
	ConferenceTalk TalkType = "CONFERENCE_TALK"
	Workshop       TalkType = "WORKSHOP"
)

type Event struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Artifact struct {
	Video string `yaml:"video"`
	Slide string `yaml:"slide"`
}

type Talks []Talk

func (t *Talks) readConf() *Talks {

	yamlFile, err := ioutil.ReadFile("talks.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, t)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return t
}

func (t Talks) Len() int {
	return len(t)
}

func (t Talks) Less(i, j int) bool {
	return t[i].GetTime().After(t[j].GetTime())
}

func (t Talks) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type Talk struct {
	Title    string   `yaml:"title"`
	Type     TalkType `yaml:"type"`
	Event    Event    `yaml:"event"`
	Date     string   `yaml:"dateTime"`
	Artifact Artifact `yaml:"artifact"`
	Tags     []string `yaml:"tags"`
}

func (t Talk) DateStr() string {
	return t.GetTime().Format("02-Jan-2006") // 31-Aug-2017
}

func (t Talk) GetTime() time.Time {
	layout := "2006-01-02T15:04:05-0700"
	tf, err := time.Parse(layout, t.Date)
	if err != nil {
		panic(err)
	}
	return tf
}

func (t Talk) TagsStr() string {
	return strings.Join(t.Tags, ", ")
}

func main() {
	var talks Talks
	talks.readConf()

	sort.Sort(talks)

	f, err := createOrOpenFile("README.md", false)
	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("readmeTemplate").Parse(string(readmeTemplate)))
	t.Execute(f, talks)
}

var readmeTemplate = []byte(`
TECH TALK
====

List of tech talk I've ever gave about Community, Open Source and Cloud Computings

| Date | Event | Title | Slide | Video | Tags |
| -----------  | ----------- | ----------- | ----------- | ----------- | ----------- |
{{- range $index, $element := . }}
| {{ $element.DateStr }} | {{ if $element.Event.URL }} [{{ $element.Event.Name}}]({{$element.Event.URL}}) {{ else }}{{ $element.Event.Name}}{{ end }} | *{{ $element.Title }}* | {{ if $element.Artifact.Slide }} [Slide]({{ $element.Artifact.Slide }} ) {{ end }} | {{ if $element.Artifact.Video }} [Video]({{ $element.Artifact.Video }}) {{ end }}  | {{ $element.TagsStr }} |
{{- end }} 
`)

func createOrOpenFile(filename string, append bool) (*os.File, error) {
	if append {
		return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	}
	return os.Create(filename)
}
