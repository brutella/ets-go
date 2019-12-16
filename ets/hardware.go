package ets

import (
	"encoding/xml"
	"fmt"
	"io"
)

// TranslationRefID is the ID of a translation.
type TranslationRefID string

type Translation struct {
	ManufacturerID ManufacturerID
	HardwareID     HardwareID
	ProductID      ProductID
	Text           string
}

// LanguageID is the ID of a language.
type LanguageID string

type Language struct {
	ID           LanguageID
	Translations []Translation
}

// ProductID is the ID of a product.
type ProductID string

type Product struct {
	ID             ProductID
	ManufacturerID ManufacturerID
	HardwareID     HardwareID
	Text           string
}

type Hardware2ProgramID string
type Hardware2Program struct {
	ManufacturerID       ManufacturerID
	HardwareID           HardwareID
	ID                   Hardware2ProgramID
	ApplicationProgramID ApplicationProgramID
}

// HardwareID is the ID of a manufacturer.
type HardwareID string

type Hardware struct {
	ID                HardwareID
	Name              string
	Products          []Product
	Hardware2Programs []Hardware2Program
}

// HardwareData contains hardware-specific data.
type HardwareData struct {
	Manufacturer ManufacturerID
	Hardwares    []Hardware
	Languages    []Language
}

// UnmarshalXML implements xml.Unmarshaler.
func (md *HardwareData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Decide which schema to use based on the value of the 'xmlns' attribute.
	ns := getNamespace(start)
	switch ns {
	case schema11Namespace, schema12Namespace, schema13Namespace, schema20Namespace:
		return d.DecodeElement((*hardwareData11)(md), &start)

	default:
		return fmt.Errorf("Unexpected namespace '%s'", ns)
	}
}

// DecodeHardwareData parses the contents of a manufacturer file.
func DecodeHardwareData(r io.Reader) (*HardwareData, error) {
	md := &HardwareData{}
	if err := xml.NewDecoder(r).Decode(md); err != nil {
		return nil, err
	}

	return md, nil
}
