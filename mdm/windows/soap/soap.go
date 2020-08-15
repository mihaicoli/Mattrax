package soap

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/mattrax/xml"
	"github.com/rs/zerolog/log"
)

// MaxRequestBodySize is the maximum amount of data that is allowed in a single request
const MaxRequestBodySize = 10000

// MustUnderstand is a easily way to create SOAP tags with s:mustUnderstand
type MustUnderstand struct {
	MustUnderstand string `xml:"s:mustUnderstand,attr,omitempty"`
	Value          string `xml:",innerxml"`
}

// RequestHeader is a generic SOAP body for requests
type RequestHeader struct {
	Action    string `xml:"a:Action"`
	MessageID string `xml:"a:MessageID"`
	ReplyTo   struct {
		Address string `xml:"a:Address"`
	} `xml:"a:ReplyTo"`
	To           string `xml:"a:To"`
	WSSESecurity struct {
		Username            string `xml:"wsse:UsernameToken>wsse:Username"`
		Password            string `xml:"wsse:UsernameToken>wsse:Password"`
		BinarySecurityToken string `xml:"wsse:BinarySecurityToken"`
	} `xml:"wsse:Security"`
}

// ResponseEnvelope is a generic Envelope used for the servers responses
type ResponseEnvelope struct {
	XMLName    xml.Name             `xml:"s:Envelope"`
	NamespaceS string               `xml:"xmlns:s,attr"`
	NamespaceA string               `xml:"xmlns:a,attr"`
	Header     ResponseHeader       `xml:"s:Header"`
	Body       ResponseEnvelopeBody `xml:"s:Body"`
}

// Populate sets reasonable default values on the response envelope
func (e *ResponseEnvelope) Populate(action string) {
	e.NamespaceS = "http://www.w3.org/2003/05/soap-envelope"
	e.NamespaceA = "http://www.w3.org/2005/08/addressing"
	e.Header.Action = MustUnderstand{
		MustUnderstand: "1",
		Value:          action,
	}
	e.Header.ActivityID = uuid.New().String()
	e.Body.NamespaceXSI = "http://www.w3.org/2001/XMLSchema-instance"
	e.Body.NamespaceXSD = "http://www.w3.org/2001/XMLSchema"
}

// ResponseHeader is a generic SOAP body for responses
type ResponseHeader struct {
	Action     MustUnderstand `xml:"a:Action,omitempty"`
	ActivityID string         `xml:"a:ActivityID,omitempty"`
	RelatesTo  string         `xml:"a:RelatesTo,omitempty"`
}

// ResponseEnvelopeBody is a generic s:Body which contains the endpoint specific response
type ResponseEnvelopeBody struct {
	NamespaceXSI string `xml:"xmlns:xsi,attr,omitempty"`
	NamespaceXSD string `xml:"xmlns:xsd,attr,omitempty"`
	Body         interface{}
}

// Read safely decodes a SOAP request from the HTTP body into a struct
func Read(v interface{}, r *http.Request, w http.ResponseWriter) bool {
	if r.ContentLength > MaxRequestBodySize {
		log.Debug().Int64("length", r.ContentLength).Int("max-length", MaxRequestBodySize).Msg("Request body larger than supported value")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return true
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)
	if err := xml.NewDecoder(r.Body).Decode(v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Str("caller", fmt.Sprintf("%T", v)).Err(err).Msg("Error decoding request")
		return true
	}

	return false
}

// Respond encodes a SOAP response from a struct into the HTTP response
func Respond(v ResponseEnvelope, w http.ResponseWriter) {
	body, err := xml.Marshal(v)
	if err != nil {
		log.Error().Str("caller", fmt.Sprintf("%T", v.Body.Body)).Err(err).Msg("error: marshalling soap body")
		if fmt.Sprintf("%T", v.Body.Body) != "SOAPFault" {
			w.WriteHeader(http.StatusInternalServerError)
			var res = NewFault("s:Receiver", "s:InternalServiceFault", "", "Mattrax encountered an error. Please check the server logs for more info", "")
			Respond(res, w)
		}
		return
	}

	w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	if _, err := w.Write(body); err != nil {
		log.Error().Str("caller", fmt.Sprintf("%T", v)).Err(err).Msg("error: writing body to client")
	}
}
