package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
	"fmt"
	"strings"
	"time"
	"crypto/md5"
)

const message = "# THIS FILE WAS GENERATED USING github.com/jackyzha0/hugo-obsidian\n# DO NOT EDIT\n"
func write(links []Link, contentIndex ContentIndex, toIndex bool, out string) error {
	hashedContentIndex := make(ContentIndex);
	for i := range links {
		links[i] = hashLink(links[i])
	}
	for key, content := range contentIndex {
		if strings.HasPrefix(key, "/private/") {
			hashedContentIndex[hashPath(key, "/private/")] = Content{
				Title: "",
				Content: "",
				LastModified: time.Now(),
			}
		} else {
			hashedContentIndex[key] = content
		}
	}

	index := index(links)
	resStruct := struct {
		Index Index
		Links []Link
	}{
		Index: index,
		Links: links,
	}
	marshalledIndex, mErr := yaml.Marshal(&resStruct)
	if mErr != nil {
		return mErr
	}

	writeErr := ioutil.WriteFile(path.Join(out, "linkIndex.yaml"), append([]byte(message), marshalledIndex...), 0644)
	if writeErr != nil {
		return writeErr
	}

	if toIndex {
		marshalledContentIndex, mcErr := yaml.Marshal(&hashedContentIndex)
		if mcErr != nil {
			return mcErr
		}

		writeErr = ioutil.WriteFile(path.Join(out, "contentIndex.yaml"), append([]byte(message), marshalledContentIndex...), 0644)
		if writeErr != nil {
			return writeErr
		}
	}

	return nil
}

func hashPath(path string, prefix string) (string) {
	if strings.HasPrefix(path, prefix) == true {
		path = fmt.Sprintf("%s%x", prefix, md5.Sum([]byte(path)));
	}
	return path
}

func hashLink(l Link) (Link) {
	return Link{
		Source: hashPath(l.Source, "/private/"),
		Target: hashPath(l.Target, "/private/"),
		Text:   "",
	}
}

// constructs index from links
func index(links []Link) (index Index) {
	linkMap := make(map[string][]Link)
	backlinkMap := make(map[string][]Link)
	for _, l := range links {
		l := hashLink(l)
		// backlink (only if internal)
		if _, ok := backlinkMap[l.Target]; ok {
			backlinkMap[l.Target] = append(backlinkMap[l.Target], l)
		} else {
			backlinkMap[l.Target] = []Link{l}
		}

		// regular link
		if _, ok := linkMap[l.Source]; ok {
			linkMap[l.Source] = append(linkMap[l.Source], l)
		} else {
			linkMap[l.Source] = []Link{l}
		}
	}
	index.Links = linkMap
	index.Backlinks = backlinkMap
	return index
}



