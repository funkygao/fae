package main

import (
	"bytes"
	log "github.com/funkygao/log4go"
	"io/ioutil"
	"strings"
	"text/template"
)

func dumpMaintainConfigPhp(info []string) {
	if config.maintainTargetFile == "" ||
		config.maintainTemplateFile == "" ||
		maintainTemplateContents == "" {
		log.Warn("Invalid maintain conf, disabled")
		return
	}

	templateData := make(map[string]string)
	for _, s := range info {
		// s is like "kingdom_1:30"
		parts := strings.SplitN(s, ":", 2)
		templateData[parts[0]] = parts[1]
	}

	t := template.Must(template.New("maintain").Parse(maintainTemplateContents))
	wr := new(bytes.Buffer)
	t.Execute(wr, templateData)

	err := ioutil.WriteFile(config.maintainTargetFile, wr.Bytes(), 0644)
	if err != nil {
		log.Error("dump[%s]: %s", config.maintainTargetFile, err.Error())
	} else {
		log.Info("dumped[%s]: %+v", config.maintainTargetFile, templateData)
	}

}

func dumpActorConfigPhp(servers []string) {
	if config.actorTargetFile == "" ||
		config.actorTemplateFile == "" ||
		actorTemplateContents == "" {
		log.Warn("Invalid actor conf, disabled")
		return
	}

	log.Info("dumped[%s]: %+v", config.actorTargetFile, servers)
}

func dumpFaeConfigPhp(servers []string) {
	if config.faeTargetFile == "" ||
		config.faeTemplateFile == "" ||
		faeTemplateContents == "" {
		log.Warn("Invalid fae conf, disabled")
		return
	}

	type tempateVar struct {
		Servers []string
		Ports   []string
	}
	templateData := tempateVar{Servers: make([]string, 0), Ports: make([]string, 0)}
	for _, s := range servers {
		// s is like "12.3.11.2:9001"
		parts := strings.SplitN(s, ":", 2)
		templateData.Servers = append(templateData.Servers, parts[0])
		templateData.Ports = append(templateData.Ports, parts[1])
	}

	t := template.Must(template.New("fae").Parse(faeTemplateContents))
	wr := new(bytes.Buffer)
	t.Execute(wr, templateData)

	err := ioutil.WriteFile(config.faeTargetFile, wr.Bytes(), 0644)
	if err != nil {
		log.Error("dump[%s]: %s", config.faeTargetFile, err.Error())
	} else {
		log.Info("dumped[%s]: %+v", config.faeTargetFile, templateData)
	}

}
