package ets

import (
	"encoding/xml"
	"fmt"
	"strings"
)

const schema20Namespace = "http://knx.org/xml/project/20"

type deviceInstance20 DeviceInstance

func (di *deviceInstance20) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID         string `xml:"Id,attr"`
		ProductID  string `xml:"ProductRefId,attr"`
		ProgramID  string `xml:"Hardware2ProgramRefId,attr"`
		Name       string `xml:",attr"`
		Address    uint16 `xml:",attr"`
		ComObjects []struct {
			RefID             string `xml:"RefId,attr"`
			DatapointType     string `xml:",attr"`
			Links             string `xml:",attr"`
			ReadFlag          string `xml:",attr"`
			WriteFlag         string `xml:",attr"`
			CommunicationFlag string `xml:",attr"`
			TransmitFlag      string `xml:",attr"`
			UpdateFlag        string `xml:",attr"`
			ReadOnInitFlag    string `xml:",attr"`
		} `xml:"ComObjectInstanceRefs>ComObjectInstanceRef"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	ids := strings.Split(doc.ID, "_")
	if len(ids) == 2 {
		di.ProjectID = ProjectID(ids[0])
		di.ID = DeviceInstanceID(ids[1])
	}

	prodIds := strings.Split(doc.ProductID, "_")
	if len(prodIds) == 3 {
		di.ManufacturerID = ManufacturerID(prodIds[0])
		di.HardwareID = HardwareID(prodIds[1])
		di.ProductID = ProductID(prodIds[2])
	}

	progIds := strings.Split(doc.ProgramID, "_")
	if len(progIds) == 3 {
		di.Hardware2ProgramID = Hardware2ProgramID(progIds[2])
	}

	di.Name = doc.Name
	di.Address = doc.Address
	di.ComObjects = make([]ComObjectInstanceRef, len(doc.ComObjects))

	for n, docComObj := range doc.ComObjects {
		ids := strings.Split(docComObj.RefID, "_")
		var comObjID ComObjectID
		var comObjRefID ComObjectRefID
		for _, id := range ids {
			if strings.HasPrefix(id, "O-") {
				comObjID = ComObjectID(id)
			} else if strings.HasPrefix(id, "R-") {
				comObjRefID = ComObjectRefID(id)
			}
		}

		if len(comObjID) == 0 && len(comObjRefID) == 0 {
			return fmt.Errorf("Invalid ComObjectInstanceRefId %s", docComObj.RefID)
		}

		comObj := ComObjectInstanceRef{
			ComObjectID:       comObjID,
			ComObjectRefID:    comObjRefID,
			DatapointType:     docComObj.DatapointType,
			Links:             make([]string, 0),
			ReadFlag:          docComObj.ReadFlag == "Enabled",
			WriteFlag:         docComObj.WriteFlag == "Enabled",
			CommunicationFlag: docComObj.CommunicationFlag == "Enabled",
			TransmitFlag:      docComObj.TransmitFlag == "Enabled",
			UpdateFlag:        docComObj.UpdateFlag == "Enabled",
			ReadOnInitFlag:    docComObj.ReadOnInitFlag == "Enabled",
		}

		links := strings.Split(docComObj.Links, " ")
		for _, link := range links {
			if len(link) > 0 {
				comObj.Links = append(comObj.Links, link)
			}
		}

		di.ComObjects[n] = comObj
	}

	return nil
}

type line20 Line

func (l *line20) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID             string `xml:"Id,attr"`
		Name           string `xml:",attr"`
		Address        uint16 `xml:",attr"`
		DeviceInstance []deviceInstance20
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	ids := strings.Split(doc.ID, "_")
	if len(ids) == 2 {
		l.ProjectID = ProjectID(ids[0])
		l.ID = LineID(ids[1])
	}

	l.Name = doc.Name
	l.Address = doc.Address
	l.Devices = make([]DeviceInstance, len(doc.DeviceInstance))

	for n, docDeviceInstance := range doc.DeviceInstance {
		l.Devices[n] = DeviceInstance(docDeviceInstance)
	}

	return nil
}

type area20 Area

func (a *area20) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID      string `xml:"Id,attr"`
		Name    string `xml:",attr"`
		Address uint16 `xml:",attr"`
		Line    []line20
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
		a.Lines[n] = Line(docLine)
	}

	return nil
}

type installation20 Installation

func (i *installation20) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Name        string         `xml:",attr"`
		Areas       []area20       `xml:"Topology>Area"`
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

type project20 Project

func (p *project20) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Project struct {
			ID            string           `xml:"Id,attr"`
			Installations []installation20 `xml:"Installations>Installation"`
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
