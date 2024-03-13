package ets

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"
)

func getNamespace(start xml.StartElement) string {
	for _, attr := range start.Attr {
		if attr.Name.Local == "xmlns" {
			return attr.Value
		}
	}

	return ""
}

// ProjectID is a project identifier.
type ProjectID string

type GroupAddressStyle int

const (
	GroupAddressStyleThree GroupAddressStyle = iota
	GroupAddressStyleTwo
	GroupAddressStyleFree
)

// ProjectInfo contains project information. These information are usually stored in
// the P-XXXX/Project.xml file.
type ProjectInfo struct {
	ID           ProjectID
	Name         string
	Comment      string
	AddressStyle GroupAddressStyle
}

// UnmarshalXML implements xml.Unmarshaler.
func (pi *ProjectInfo) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Decide which schema to use based on the value of the 'xmlns' attribute.
	ns := getNamespace(start)
	switch ns {
	case schema11Namespace, schema12Namespace, schema13Namespace, schema14Namespace, schema20Namespace, schema21Namespace, schema22Namespace, schema23Namespace:
		return d.DecodeElement((*projectInfo11)(pi), &start)

	default:
		return fmt.Errorf("Unexpected namespace '%s'", ns)
	}
}

// DecodeProjectInfo parses the contents of project info file.
func DecodeProjectInfo(r io.Reader) (*ProjectInfo, error) {
	info := &ProjectInfo{}
	if err := xml.NewDecoder(r).Decode(info); err != nil {
		return nil, err
	}

	return info, nil
}

// Connector is a connection to a group address.
type Connector struct {
	Receive bool
	RefID   GroupAddressID
}

// ComObjectRefID is the ID of a communication object reference.
type ComObjectRefID string

func (s ComObjectRefID) InstanceID() string {
	com := strings.Split(string(s), "_")
	if len(com) == 4 {
		return com[3]
	}

	return ""
}

// IDs returns the manfucturer, applicaton program, communication object and communcation object reference ids.
func (s ComObjectRefID) IDs() (m string, ap string, o string, r string) {
	com := strings.Split(string(s), "_")
	if len(com) == 4 {
		m = com[0]
		ap = com[1]
		o = com[2]
		r = com[3]
	}

	return
}

// ComObjectInstanceRef connects a communication object reference with zero or more group addresses.
type ComObjectInstanceRef struct {
	ComObjectRefID ComObjectRefID
	ComObjectID    ComObjectID
	DatapointType  string
	Links          []string
}

// DeviceInstanceID is the ID of a device instance.
type DeviceInstanceID string

// DeviceInstance is a device instance.
type DeviceInstance struct {
	ID                 DeviceInstanceID
	ProjectID          ProjectID
	ManufacturerID     ManufacturerID
	HardwareID         HardwareID
	ProductID          ProductID
	Hardware2ProgramID Hardware2ProgramID
	Name               string
	Address            uint16
	ComObjects         []ComObjectInstanceRef
}

// LineID is the ID of a line.
type LineID string

// Line is a line.
type Line struct {
	ID        LineID
	ProjectID ProjectID
	Name      string
	Address   uint16
	Devices   []DeviceInstance
}

// AreaID is the ID of an area.
type AreaID string

// Area is an area.
type Area struct {
	ID        AreaID
	ProjectID ProjectID
	Name      string
	Address   uint16
	Lines     []Line
}

// GroupAddressID is the ID of a group address.
type GroupAddressID string

// GroupAddress is a group address.
type GroupAddress struct {
	ID            GroupAddressID
	ProjectID     ProjectID
	Name          string
	Description   string
	Address       uint16
	DatapointType string
}

// GroupRangeID is the ID of a group range.
type GroupRangeID string

// GroupRange is a range of group addresses.
type GroupRange struct {
	ID         GroupRangeID
	Name       string
	RangeStart uint16
	RangeEnd   uint16
	Addresses  []GroupAddress
	SubRanges  []GroupRange
}

// SpaceID is the ID of a space.
type SpaceID string

const (
	SpaceTypeBuilding          = "Building"
	SpaceTypeBuildingPart      = "BuildingPart"
	SpaceTypeFloor             = "Floor"
	SpaceTypeRoom              = "Room"
	SpaceTypeDistributionBoard = "DistributionBoard"
	SpaceTypeStairway          = "Stairway"
	SpaceTypeCorridor          = "Corridor"
)

// Space is a space for devices and other spaces.
type Space struct {
	ID                SpaceID
	ProjectID         ProjectID
	DeviceInstanceIDs []DeviceInstanceID
	Type              string
	Name              string
	SubSpaces         []Space
}

// Installation is an installation within a project.
type Installation struct {
	Name           string
	Topology       []Area
	Locations      []Space
	GroupAddresses []GroupRange
}

// Project contains an entire project. These information are usually stored within a file located
// at P-XXXX/N.xml (0 <= N <= 16).
type Project struct {
	ID            ProjectID
	Name          string
	Installations []Installation
}

// UnmarshalXML implements xml.Unmarshaler.
func (p *Project) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	log.Println("Unmarshal XML")
	// Decide which schema to use based on the value of the 'xmlns' attribute.
	ns := getNamespace(start)
	switch ns {
	case schema11Namespace, schema12Namespace, schema13Namespace, schema14Namespace:
		return d.DecodeElement((*project11)(p), &start)

	case schema20Namespace:
		return d.DecodeElement((*project20)(p), &start)

	case schema21Namespace, schema22Namespace, schema23Namespace:
		return d.DecodeElement((*project21)(p), &start)

	default:
		return fmt.Errorf("Unexpected namespace '%s'", ns)
	}
}

// DecodeProject parses the contents of a project file.
func DecodeProject(r io.Reader) (*Project, error) {
	proj := &Project{}
	if err := xml.NewDecoder(r).Decode(proj); err != nil {
		return nil, err
	}

	return proj, nil
}
