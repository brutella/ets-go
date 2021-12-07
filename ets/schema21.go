package ets

import (
	"encoding/xml"
	"strings"
)

const schema21Namespace = "http://knx.org/xml/project/21"

type area21 Area

func (a *area21) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID      string `xml:"Id,attr"`
		Name    string `xml:",attr"`
		Address uint16 `xml:",attr"`
		Line    []struct {
			ID      string `xml:"Id,attr"`
			Name    string `xml:",attr"`
			Address uint16 `xml:",attr"`
			Segment struct {
				ID             string `xml:"Id,attr"`
				Name           string `xml:",attr"`
				Number         int    `xml:",attr"`
				DeviceInstance []deviceInstance20
			}
		}
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	ids := strings.Split(doc.ID, "_")
	if len(ids) == 2 {
		a.ProjectID = ProjectID(ids[0])
		a.ID = AreaID(ids[1])
	}
	a.Name = doc.Name
	a.Address = doc.Address
	a.Lines = make([]Line, len(doc.Line))

	for n, docLine := range doc.Line {
		line := Line{
			Name:    docLine.Name,
			Address: docLine.Address,
			Devices: make([]DeviceInstance, len(docLine.Segment.DeviceInstance)),
		}

		ids := strings.Split(docLine.ID, "_")
		if len(ids) == 2 {
			line.ProjectID = ProjectID(ids[0])
			line.ID = LineID(ids[1])
		}

		for n, segmentDevice := range docLine.Segment.DeviceInstance {
			line.Devices[n] = DeviceInstance(segmentDevice)
		}
		a.Lines[n] = line
	}

	return nil
}

type installation21 Installation

func (i *installation21) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Name        string         `xml:",attr"`
		Areas       []area21       `xml:"Topology>Area"`
		GroupRanges []groupRange11 `xml:"GroupAddresses>GroupRanges>GroupRange"`
		Locations   []space11      `xml:"Locations>Space"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	i.Name = doc.Name
	i.Topology = make([]Area, len(doc.Areas))
	i.GroupAddresses = make([]GroupRange, len(doc.GroupRanges))
	i.Locations = make([]Space, len(doc.Locations))

	for n, docArea := range doc.Areas {
		i.Topology[n] = Area(docArea)
	}

	for n, docGrpRange := range doc.GroupRanges {
		i.GroupAddresses[n] = GroupRange(docGrpRange)
	}

	for n, docSpace := range doc.Locations {
		i.Locations[n] = Space(docSpace)
	}

	return nil
}

type project21 Project

func (p *project21) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Project struct {
			ID            string           `xml:"Id,attr"`
			Installations []installation21 `xml:"Installations>Installation"`
		}
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	p.ID = ProjectID(doc.Project.ID)
	p.Installations = make([]Installation, len(doc.Project.Installations))

	for i, docInst := range doc.Project.Installations {
		p.Installations[i] = Installation(docInst)
	}

	return nil
}
