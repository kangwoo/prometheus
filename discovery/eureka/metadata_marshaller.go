// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eureka

import (
	"encoding/xml"
	"regexp"
)

type MetaData struct {
	Map   map[string]string
	Class string
}

type Vraw struct {
	Content []byte `xml:",innerxml"`
	Class   string `xml:"class,attr"`
}

func (s *MetaData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var attributes = make([]xml.Attr, 0)
	if s.Class != "" {
		attributes = append(attributes, xml.Attr{
			Name: xml.Name{
				Local: "class",
			},
			Value: s.Class,
		})
	}
	start.Attr = attributes
	tokens := []xml.Token{start}

	for key, value := range s.Map {
		t := xml.StartElement{Name: xml.Name{Space: "", Local: key}}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{Name: t.Name})
	}

	tokens = append(tokens, xml.EndElement{
		Name: start.Name,
	})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	err := e.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (s *MetaData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	s.Map = make(map[string]string)
	vraw := &Vraw{}
	err := d.DecodeElement(vraw, &start)
	if err != nil {
		return err
	}
	dataInString := string(vraw.Content)
	regex, err := regexp.Compile(`\s*<([^<>]+)>([^<>]+)</[^<>]+>\s*`)
	if err != nil {
		return err
	}
	subMatches := regex.FindAllStringSubmatch(dataInString, -1)
	for _, subMatch := range subMatches {
		s.Map[subMatch[1]] = subMatch[2]
	}
	s.Class = vraw.Class
	return nil
}