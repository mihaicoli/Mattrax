package wap

import "github.com/mattrax/xml"

// WapProvisioningDoc contains the manage client configuration
type WapProvisioningDoc struct {
	XMLName        xml.Name            `xml:"wap-provisioningdoc"`
	Version        string              `xml:"version,attr"`
	Characteristic []WapCharacteristic `xml:"characteristic"`
}

// WapCharacteristic is a management client characteristic
type WapCharacteristic struct {
	Type            string `xml:"type,attr,omitempty"`
	Params          []WapParameter
	Characteristics []WapCharacteristic `xml:"characteristic,omitempty"`
}

// WapParameter is a management client paramter (setting) that is set on a characteristic
type WapParameter struct {
	XMLName  xml.Name `xml:"parm"`
	Name     string   `xml:"name,attr,omitempty"`
	Value    string   `xml:"value,attr,omitempty"`
	DataType string   `xml:"datatype,attr,omitempty"`
}
