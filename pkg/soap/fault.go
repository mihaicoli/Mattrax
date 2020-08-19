package soap

import "github.com/mattrax/xml"

// Fault is the body that is returned when an error occurs
type Fault struct {
	XMLName xml.Name                          `xml:"s:Fault"`
	Code    FaultCode                         `xml:"s:Code"`
	Reason  FaultReason                       `xml:"s:Reason>s:Text"`
	Detail  FaultDeviceEnrollmentServiceError `xml:"http://schemas.microsoft.com/windows/pki/2009/01/enrollment s:Detail>DeviceEnrollmentServiceError"`
}

// FaultCode contains the errors causer (Sender or Receiver) and the error code
type FaultCode struct {
	Value   string `xml:"s:Value"`
	Subcode string `xml:"s:Subcode>s:Value"`
}

// FaultReason contains the human readable error message which is shown in the device management logs
type FaultReason struct {
	Lang  string `xml:"xml:lang,attr"`
	Value string `xml:",innerxml"`
}

// FaultDeviceEnrollmentServiceError contains extra error codes (which sometimes have special UI's) and a traceid which can be used trace requests between the client and server logs
type FaultDeviceEnrollmentServiceError struct {
	ErrorType string `xml:"ErrorType"`
	Message   string `xml:"Message"`
	TraceID   string `xml:"TraceId"`
}

// NewFault creates a new fault
func NewFault(causer, code, errortype, reason, traceid string) ResponseEnvelope {
	var deviceEnrollmentServiceError FaultDeviceEnrollmentServiceError
	if errortype != "" {
		deviceEnrollmentServiceError = FaultDeviceEnrollmentServiceError{
			ErrorType: errortype,
			Message:   reason,
			TraceID:   traceid,
		}
	}

	return ResponseEnvelope{
		NamespaceS: "http://www.w3.org/2003/05/soap-envelope",
		NamespaceA: "http://www.w3.org/2005/08/addressing",
		Body: ResponseEnvelopeBody{
			Body: Fault{
				Code: FaultCode{
					Value:   causer,
					Subcode: code,
				},
				Reason: FaultReason{
					Lang:  "en-US",
					Value: reason,
				},
				Detail: deviceEnrollmentServiceError,
			},
		},
	}
}
