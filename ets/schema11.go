package ets

import (
	"encoding/xml"
	"fmt"
	"strings"
)

const schema11Namespace = "http://knx.org/xml/project/11"

type projectInfo11 ProjectInfo

func (pi *projectInfo11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Project struct {
			ID                 string `xml:"Id,attr"`
			ProjectInformation struct {
				Name              string `xml:",attr"`
				Comment           string `xml:",attr"`
				GroupAddressStyle string `xml:",attr"`
			}
		}
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	pi.ID = ProjectID(doc.Project.ID)
	pi.Name = doc.Project.ProjectInformation.Name
	pi.Comment = doc.Project.ProjectInformation.Comment

	switch doc.Project.ProjectInformation.GroupAddressStyle {
	case "ThreeLevel":
		pi.AddressStyle = GroupAddressStyleThree
	case "TwoLevel":
		pi.AddressStyle = GroupAddressStyleTwo
	default:
		pi.AddressStyle = GroupAddressStyleFree
	}

	return nil
}

type deviceInstance11 DeviceInstance

func (di *deviceInstance11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID         string `xml:"Id,attr"`
		ProductID  string `xml:"ProductRefId,attr"`
		ProgramID  string `xml:"Hardware2ProgramRefId,attr"`
		Name       string `xml:",attr"`
		Address    uint16 `xml:",attr"`
		ComObjects []struct {
			RefID         string `xml:"RefId,attr"`
			DatapointType string `xml:",attr"`
			Connectors    struct {
				Elements []struct {
					XMLName xml.Name
					RefID   string `xml:"GroupAddressRefId,attr"`
				} `xml:",any"`
			}
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

		var links = []string{}
		for _, docConnElem := range docComObj.Connectors.Elements {
			ids := strings.Split(docConnElem.RefID, "_")
			if len(ids) == 2 && len(ids[1]) > 0 {
				links = append(links, ids[1])
			}
		}

		comObj := ComObjectInstanceRef{
			ComObjectID:    comObjID,
			ComObjectRefID: comObjRefID,
			DatapointType:  docComObj.DatapointType,
			Links:          links,
		}

		di.ComObjects[n] = comObj

	}

	return nil
}

type line11 Line

func (l *line11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID             string `xml:"Id,attr"`
		Name           string `xml:",attr"`
		Address        uint16 `xml:",attr"`
		DeviceInstance []deviceInstance11
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	l.ID = LineID(doc.ID)
	l.Name = doc.Name
	l.Address = doc.Address
	l.Devices = make([]DeviceInstance, len(doc.DeviceInstance))

	for n, docDeviceInstance := range doc.DeviceInstance {
		l.Devices[n] = DeviceInstance(docDeviceInstance)
	}

	return nil
}

type area11 Area

func (a *area11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID      string `xml:"Id,attr"`
		Name    string `xml:",attr"`
		Address uint16 `xml:",attr"`
		Line    []line11
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	a.ID = AreaID(doc.ID)
	a.Name = doc.Name
	a.Address = doc.Address
	a.Lines = make([]Line, len(doc.Line))

	for n, docLine := range doc.Line {
		a.Lines[n] = Line(docLine)
	}

	return nil
}

type space11 Space

func (sp *space11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID                string    `xml:"Id,attr"`
		Name              string    `xml:",attr"`
		Type              string    `xml:",attr"`
		SubSpaces         []space11 `xml:"Space"`
		DeviceInstanceRef []struct {
			RefID string `xml:"RefId,attr"`
		}
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	ids := strings.Split(doc.ID, "_")
	if len(ids) == 2 {
		sp.ProjectID = ProjectID(ids[0])
		sp.ID = SpaceID(ids[1])
	}

	sp.Name = doc.Name
	sp.Type = doc.Type
	sp.SubSpaces = make([]Space, len(doc.SubSpaces))
	sp.DeviceInstanceIDs = make([]DeviceInstanceID, len(doc.DeviceInstanceRef))

	for n, docSpace := range doc.SubSpaces {
		sp.SubSpaces[n] = Space(docSpace)
	}

	for n, docRef := range doc.DeviceInstanceRef {
		ids := strings.Split(docRef.RefID, "_")
		if len(ids) == 2 {
			sp.DeviceInstanceIDs[n] = DeviceInstanceID(ids[1])
		}
	}

	return nil
}

type groupRange11 GroupRange

func (gar *groupRange11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID           string `xml:"Id,attr"`
		Name         string `xml:",attr"`
		RangeStart   uint16 `xml:",attr"`
		RangeEnd     uint16 `xml:",attr"`
		GroupAddress []struct {
			ID            string `xml:"Id,attr"`
			Name          string `xml:",attr"`
			Address       uint16 `xml:",attr"`
			Description   string `xml:",attr"`
			DatapointType string `xml:",attr"`
		}
		GroupRange []groupRange11
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	gar.ID = GroupRangeID(doc.ID)
	gar.Name = doc.Name
	gar.RangeStart = doc.RangeStart
	gar.RangeEnd = doc.RangeEnd
	gar.Addresses = make([]GroupAddress, len(doc.GroupAddress))
	gar.SubRanges = make([]GroupRange, len(doc.GroupRange))

	for n, ga := range doc.GroupAddress {
		ids := strings.Split(ga.ID, "_")
		if len(ids) != 2 {
			return fmt.Errorf("Invalid GroupAddress Id %s", ga.ID)
		}

		gar.Addresses[n] = GroupAddress{
			ProjectID:     ProjectID(ids[0]),
			ID:            GroupAddressID(ids[1]),
			Name:          ga.Name,
			Description:   ga.Description,
			Address:       ga.Address,
			DatapointType: ga.DatapointType,
		}
	}

	for n, docGrpRange := range doc.GroupRange {
		gar.SubRanges[n] = GroupRange(docGrpRange)
	}

	return nil
}

type installation11 Installation

func (i *installation11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Name        string         `xml:",attr"`
		Areas       []area11       `xml:"Topology>Area"`
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

type project11 Project

func (p *project11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Project struct {
			ID            string           `xml:"Id,attr"`
			Installations []installation11 `xml:"Installations>Installation"`
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

type Ids struct {
	Manufacturer ManufacturerID
	AppProgram   ApplicationProgramID
	Module       ModuleID
	ComObject    ComObjectID
	ComObjectRef ComObjectRefID
}

func parseIds(s string) Ids {
	var ids Ids
	for _, id := range strings.Split(s, "_") {
		if len(id) == 0 {
			continue
		}

		parts := strings.Split(id, "-")
		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "M":
			ids.Manufacturer = ManufacturerID(id)
		case "A":
			ids.AppProgram = ApplicationProgramID(id)
		case "MD":
			ids.Module = ModuleID(id)
		case "O":
			ids.ComObject = ComObjectID(id)
		case "R":
			ids.ComObjectRef = ComObjectRefID(id)
		}
	}
	return ids
}

type comObject11 ComObject

func (co *comObject11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID                string `xml:"Id,attr"`
		Name              string `xml:",attr"`
		Text              string `xml:",attr"`
		FunctionText      string `xml:",attr"`
		ObjectSize        string `xml:",attr"`
		DatapointType     string `xml:",attr"`
		Priority          string `xml:",attr"`
		ReadFlag          string `xml:",attr"`
		WriteFlag         string `xml:",attr"`
		CommunicationFlag string `xml:",attr"`
		TransmitFlag      string `xml:",attr"`
		UpdateFlag        string `xml:",attr"`
		ReadOnInitFlag    string `xml:",attr"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}
	ids := parseIds(doc.ID)

	if len(ids.Manufacturer) == 0 || len(ids.AppProgram) == 0 || len(ids.ComObject) == 0 {
		return fmt.Errorf("Invalid ComObjectId %s", doc.ID)
	}

	co.ManufacturerID = ids.Manufacturer
	co.ApplicationProgramID = ids.AppProgram
	co.ID = ids.ComObject
	co.ModuleID = ids.Module
	co.Name = doc.Name
	co.Text = doc.Text
	co.FunctionText = doc.FunctionText
	co.ObjectSize = doc.ObjectSize
	co.DatapointType = doc.DatapointType
	co.Priority = doc.Priority
	co.ReadFlag = doc.ReadFlag == "Enabled"
	co.WriteFlag = doc.WriteFlag == "Enabled"
	co.CommunicationFlag = doc.CommunicationFlag == "Enabled"
	co.TransmitFlag = doc.TransmitFlag == "Enabled"
	co.UpdateFlag = doc.UpdateFlag == "Enabled"
	co.ReadOnInitFlag = doc.ReadOnInitFlag == "Enabled"

	return nil
}

type comObjectRef11 ComObjectRef

// Id="M-0080_A-1012-10-5227-O00C5_O-0_R-1" RefId="M-0080_A-1012-10-5227-O00C5_O-0"
func (cor *comObjectRef11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID                string `xml:"Id,attr"`
		RefID             string `xml:"RefId,attr"`
		Name              string `xml:",attr"`
		Text              string `xml:",attr"`
		Description       string `xml:",attr"`
		FunctionText      string `xml:",attr"`
		ObjectSize        string `xml:",attr"`
		DatapointType     string `xml:",attr"`
		Priority          string `xml:",attr"`
		ReadFlag          string `xml:",attr"`
		WriteFlag         string `xml:",attr"`
		CommunicationFlag string `xml:",attr"`
		TransmitFlag      string `xml:",attr"`
		UpdateFlag        string `xml:",attr"`
		ReadOnInitFlag    string `xml:",attr"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	ids := parseIds(doc.ID)

	if len(ids.Manufacturer) == 0 || len(ids.AppProgram) == 0 || len(ids.ComObject) == 0 || len(ids.ComObjectRef) == 0 {
		return fmt.Errorf("Invalid ComObjectRefId %s", doc.ID)
	}

	cor.ManufacturerID = ids.Manufacturer
	cor.ApplicationProgramID = ids.AppProgram
	cor.ComObjectID = ids.ComObject
	cor.ID = ids.ComObjectRef
	cor.Name = doc.Name
	cor.Text = doc.Text
	cor.FunctionText = doc.FunctionText
	cor.ObjectSize = doc.ObjectSize
	cor.DatapointType = doc.DatapointType
	cor.Priority = doc.Priority
	cor.ReadFlag = doc.ReadFlag == "Enabled"
	cor.WriteFlag = doc.WriteFlag == "Enabled"
	cor.CommunicationFlag = doc.CommunicationFlag == "Enabled"
	cor.TransmitFlag = doc.TransmitFlag == "Enabled"
	cor.UpdateFlag = doc.UpdateFlag == "Enabled"
	cor.ReadOnInitFlag = doc.ReadOnInitFlag == "Enabled"

	return nil
}

type applicationProgram11 ApplicationProgram

func (ap *applicationProgram11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var doc struct {
		ID      string `xml:"Id,attr"` // Id="M-0080_A-1012-10-5227-O00C5"
		Name    string `xml:",attr"`
		Version uint   `xml:"ApplicationVersion,attr"`
		Static  struct {
			Objects    []comObject11    `xml:"ComObjectTable>ComObject"`
			ObjectRefs []comObjectRef11 `xml:"ComObjectRefs>ComObjectRef"`
		}
		Modules []struct {
			ID         string           `xml:"Id,attr"` // Id="M-0080_A-1012-10-5227-O00C5_MD-1"
			Objects    []comObject11    `xml:"Static>ComObjects>ComObject"`
			ObjectRefs []comObjectRef11 `xml:"Static>ComObjectRefs>ComObjectRef"`
		} `xml:"ModuleDefs>ModuleDef"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	ids := parseIds(doc.ID)
	if len(ids.Manufacturer) == 0 || len(ids.AppProgram) == 0 {
		return fmt.Errorf("Invalid ApplicationProgram RefID %s", doc.ID)
	}

	ap.ManufacturerID = ids.Manufacturer
	ap.ID = ids.AppProgram
	ap.Name = doc.Name
	ap.Version = doc.Version

	comObjects := doc.Static.Objects
	comObjectRefs := doc.Static.ObjectRefs
	for _, module := range doc.Modules {
		comObjects = append(comObjects, module.Objects[:]...)
		comObjectRefs = append(comObjectRefs, module.ObjectRefs[:]...)
	}

	ap.Objects = make([]ComObject, len(comObjects))
	ap.ObjectRefs = make([]ComObjectRef, len(comObjectRefs))

	for n, docComObj := range comObjects {
		ap.Objects[n] = ComObject(docComObj)
	}

	for n, docComObjRef := range comObjectRefs {
		ap.ObjectRefs[n] = ComObjectRef(docComObjRef)
	}

	return nil
}

type manufacturerData11 ManufacturerData

func (md *manufacturerData11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Manufacturer struct {
			ID        string                 `xml:"RefId,attr"`
			Programs  []applicationProgram11 `xml:"ApplicationPrograms>ApplicationProgram"`
			Languages []language11           `xml:"Languages>Language"`
		} `xml:"ManufacturerData>Manufacturer"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	md.ID = ManufacturerID(doc.Manufacturer.ID)
	md.Programs = make([]ApplicationProgram, len(doc.Manufacturer.Programs))
	for n, docProg := range doc.Manufacturer.Programs {
		md.Programs[n] = ApplicationProgram(docProg)
	}

	md.Languages = make([]Language, len(doc.Manufacturer.Languages))
	for n, docLang := range doc.Manufacturer.Languages {
		md.Languages[n] = Language(docLang)
	}

	return nil
}

type product11 Product

func (pr *product11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID   string `xml:"Id,attr"`
		Text string `xml:",attr"`
	}

	// <Product Id="M-0080_H-2014.5F10.5F14-1_P-EB10430442"
	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}
	ids := strings.Split(doc.ID, "_")
	if len(ids) != 3 {
		return fmt.Errorf("Invalid Product Id %s", doc.ID)
	}

	pr.ManufacturerID = ManufacturerID(ids[0])
	pr.HardwareID = HardwareID(ids[1])
	pr.ID = ProductID(ids[2])
	pr.Text = doc.Text

	return nil
}

type hardware2Program11 Hardware2Program

func (hp *hardware2Program11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID          string `xml:"Id,attr"`
		ProgramRefs []struct {
			RefID string `xml:"RefId,attr"`
		} `xml:"ApplicationProgramRef"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	ids := strings.Split(doc.ID, "_")
	if len(ids) == 3 {
		hp.ManufacturerID = ManufacturerID(ids[0])
		hp.HardwareID = HardwareID(ids[1])
		hp.ID = Hardware2ProgramID(ids[2])
	}

	hp.ApplicationProgramIDs = make([]ApplicationProgramID, len(doc.ProgramRefs))
	for i, ref := range doc.ProgramRefs {
		refIds := strings.Split(ref.RefID, "_")
		if len(refIds) == 2 {
			hp.ApplicationProgramIDs[i] = ApplicationProgramID(refIds[1])
		}
	}

	return nil
}

type hardware11 Hardware

func (hw *hardware11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Hardware struct {
			ID                string               `xml:"RefId,attr"`
			Name              string               `xml:",attr`
			Products          []product11          `xml:"Products>Product"`
			Hardware2Programs []hardware2Program11 `xml:"Hardware2Programs>Hardware2Program"`
		} `xml:"Hardware"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	hw.ID = HardwareID(doc.Hardware.ID)
	hw.Name = doc.Hardware.Name
	hw.Products = make([]Product, len(doc.Hardware.Products))
	hw.Hardware2Programs = make([]Hardware2Program, len(doc.Hardware.Hardware2Programs))

	for n, docProd := range doc.Hardware.Products {
		hw.Products[n] = Product(docProd)
	}

	for n, docProg := range doc.Hardware.Hardware2Programs {
		hw.Hardware2Programs[n] = Hardware2Program(docProg)
	}

	return nil
}

type language11 Language

func (lang *language11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID   string `xml:"Identifier,attr"`
		Unit []struct {
			Elements []struct {
				RefID        string `xml:"RefId,attr"`
				Translations []struct {
					Name string `xml:"AttributeName,attr"`
					Text string `xml:"Text,attr"`
				} `xml:"Translation"`
			} `xml:"TranslationElement"`
		} `xml:"TranslationUnit"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	lang.ID = LanguageID(doc.ID)

	var ts []Translation
	for _, u := range doc.Unit {
		for _, e := range u.Elements {
			texts := map[string]string{}
			for _, t := range e.Translations {
				texts[t.Name] = t.Text
			}
			t := Translation{
				RefID: TranslationRefID(e.RefID),
				Texts: texts,
			}
			ts = append(ts, t)
		}
	}

	lang.Translations = ts

	return nil
}

type hardwareData11 HardwareData

func (md *hardwareData11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Manufacturer struct {
			ID        string       `xml:"RefId,attr"`
			Hardwares []hardware11 `xml:"Hardware"`
			Languages []language11 `xml:"Languages>Language"`
		} `xml:"ManufacturerData>Manufacturer"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	md.Manufacturer = ManufacturerID(doc.Manufacturer.ID)
	md.Hardwares = make([]Hardware, len(doc.Manufacturer.Hardwares))
	md.Languages = make([]Language, len(doc.Manufacturer.Languages))

	for n, hw := range doc.Manufacturer.Hardwares {
		md.Hardwares[n] = Hardware(hw)
	}

	for n, lang := range doc.Manufacturer.Languages {
		md.Languages[n] = Language(lang)
	}

	return nil
}
