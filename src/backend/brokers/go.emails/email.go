// Copyleft (ɔ) 2017 The Caliopen contributors.
// Use of this source code is governed by a GNU AFFERO GENERAL PUBLIC
// license (AGPL) that can be found in the LICENSE file.

// package email_broker handles codec/decodec between external emails and Caliopen message format
package email_broker

import (
	"bytes"
	obj "github.com/CaliOpen/Caliopen/src/backend/defs/go-objects"
	log "github.com/Sirupsen/logrus"
	"github.com/jhillyerd/go.enmime"
	"github.com/satori/go.uuid"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"mime"
	"net/mail"
	"strings"
	"time"
)

func newAddressesFields() (af map[string][]string) {
	af = map[string][]string{
		"From":     []string{},
		"Sender":   []string{},
		"Reply-To": []string{},
		"To":       []string{},
		"Cc":       []string{},
		"Bcc":      []string{},
	}
	return
}

// build a 'ready to send' email from a Caliopen message model
// conforms to
// RFC822 / RFC2822 / RFC5322 (internet message format)
// RFC2045 / RFC2046 / RFC2047 / RFC2048 / RFC2049 (MIME) => TODO
func MarshalEmail(msg *obj.Message, version string, mailhost string) (em *obj.EmailMessage, err error) {

	em = &obj.EmailMessage{
		Email: &obj.Email{
			SmtpMailFrom: []string{},
			SmtpRcpTo:    []string{},
		},
		Message: msg,
	}

	m := gomail.NewMessage()
	addr_fields := newAddressesFields()
	for _, participant := range msg.Participants {
		switch participant.Type {
		case "From":
			addr_fields["From"] = append(addr_fields["From"], m.FormatAddress(participant.Address, participant.Label))
			em.Email.SmtpMailFrom = append(em.Email.SmtpMailFrom, participant.Address) //TODO: handle multisender to conform to RFC5322#3.6.2 (coupled with sender field)
		case "Reply-To":
			addr_fields["Reply-To"] = append(addr_fields["Reply-To"], m.FormatAddress(participant.Address, participant.Label))
		case "To":
			addr_fields["To"] = append(addr_fields["To"], m.FormatAddress(participant.Address, participant.Label))
			em.Email.SmtpRcpTo = append(em.Email.SmtpRcpTo, participant.Address)
		case "Cc":
			addr_fields["Cc"] = append(addr_fields["Cc"], m.FormatAddress(participant.Address, participant.Label))
			em.Email.SmtpRcpTo = append(em.Email.SmtpRcpTo, participant.Address)
		case "Bcc":
			em.Email.SmtpRcpTo = append(em.Email.SmtpRcpTo, participant.Address)
		case "Sender":
			addr_fields["Sender"] = append(addr_fields["Sender"], m.FormatAddress(participant.Address, participant.Label))
			em.Email.SmtpMailFrom = append(em.Email.SmtpMailFrom, participant.Address) //TODO: handle multisender to conform to RFC5322#3.6.2 (coupled with sender field)
		}
	}

	for field, addrs := range addr_fields {
		if len(addrs) > 0 {
			m.SetHeader(field, addrs...)
		}
	}

	if msg.Parent_id != "" {
		m.SetHeader("In-Reply-To", msg.Parent_id)
		//TODO: handle "References" header (RFC5322#3.6.4)
	}

	em.Message.Date = time.Now()
	m.SetHeader("Date", em.Message.Date.Format(time.RFC1123Z))

	messageId := "<" + msg.Message_id.String() + "@" + mailhost + ">" // should be the default domain in case there are multiple 'from' addresses

	m.SetHeader("Message-ID", messageId)
	m.SetHeader("X-Mailer", "Caliopen-"+version)

	//TODO: In-Reply-To header
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/plain", msg.Body)
	//TODO: errors handling

	m.WriteTo(&em.Email.Raw)
	return
}

// executed by natsMsgHandler after an outgoing email has been transmitted to the MTA without error
// it flags the caliopen message to 'sent' in cassandra and elastic
// and stores the raw outbound email
func (b *EmailBroker) SaveIndexSentEmail(ack *DeliveryAck) error {

	// save raw email in db
	json_mail, err := emailToJsonRep(ack.EmailMessage.Email.Raw.String())
	if err == nil {
		json_mail.Envelope.From = ack.EmailMessage.Email.SmtpMailFrom
		json_mail.Envelope.To = ack.EmailMessage.Email.SmtpRcpTo
	}
	raw_email_id, err := b.Store.StoreRaw(ack.EmailMessage.Email.Raw.String())
	if err != nil {
		log.WithError(err).Warn("outbound: storing raw email failed")
		return err
	}
	// update caliopen message status
	fields := make(map[string]interface{})
	fields["raw_msg_id"] = raw_email_id
	fields["is_draft"] = false
	fields["date"] = ack.EmailMessage.Message.Date
	err = b.Store.UpdateMessage(ack.EmailMessage.Message, fields)
	if err != nil {
		log.Warn("Store.UpdateMessage operation failed")
	}
	err = b.Index.UpdateMessage(ack.EmailMessage.Message, fields)
	if err != nil {
		log.Warn("Index.UpdateMessage operation failed")
	}
	return err
}

// gets a raw email and transforms into a Caliopen 'message' object
// belonging to an user
func (b *EmailBroker) UnmarshalEmail(em *obj.EmailMessage, user_id obj.UUID) (msg *obj.Message, err error) {

	parsed_mail, err := mail.ReadMessage(&em.Email.Raw)
	if err != nil {
		log.Warn("unable to parse email with raw_id : %s", em.Message.Raw_msg_id)
		return nil, err
	}

	var m_id obj.UUID
	m_id.UnmarshalBinary(uuid.NewV4().Bytes())
	mail_date, err := parsed_mail.Header.Date()
	if err != nil {
		log.WithError(err).Warn("unable to parse email's date")
	}

	mail_body, err := ioutil.ReadAll(parsed_mail.Body)
	if err != nil {
		log.WithError(err).Warn("unable to parse email's body")
	}
	//TODO: Attachments, Externals_references, identities, parent_id
	msg = &obj.Message{
		Body:             string(mail_body),
		Date:             mail_date,
		Date_insert:      time.Now(),
		Is_unread:        true,
		Message_id:       m_id,
		Participants:     []obj.Participant{},
		Privacy_features: obj.PrivacyFeatures{},
		Raw_msg_id:       em.Message.Raw_msg_id,
		Subject:          parsed_mail.Header.Get("subject"),
		Type:             obj.EmailProtocol,
		User_id:          user_id,
	}

	for field, _ := range newAddressesFields() {
		p, err := b.unmarshalParticipants(parsed_mail.Header, field, user_id)
		if err == nil {
			msg.Participants = append(msg.Participants, p...)
		}
	}

	return
}

// if an user_id is provided, the func will try to find a matching contact for each recipient within user's contactbook in db
// otherwise, contact_id will be nil for recipient
func (b *EmailBroker) unmarshalParticipants(h mail.Header, address_type string, user_id ...obj.UUID) (participants []obj.Participant, err error) {
	participants = []obj.Participant{}
	addresses, err := h.AddressList(address_type)
	if err != nil {
		return participants, err
	}
	for _, a := range addresses {
		rcpt := obj.Participant{
			Type:        address_type,
			Protocol:    obj.EmailProtocol,
			Address:     a.Address,
			Label:       a.Name,
			Contact_ids: []obj.UUID{},
		}
		if len(user_id) == 1 {
			contact_ids, err := b.Store.LookupContactsByIdentifier(user_id[0].String(), a.Address)
			if err == nil {
				for _, id := range contact_ids {
					var contact_id obj.UUID
					uuid, _ := uuid.FromString(id)
					contact_id.UnmarshalBinary(uuid.Bytes())
					rcpt.Contact_ids = append(rcpt.Contact_ids, contact_id)
				}
			}
		}
		participants = append(participants, rcpt)
	}

	return
}

// returns an EmailJson object which is an object ready to
// output our json representation of the raw email
func emailToJsonRep(email string) (json_email obj.EmailJson, err error) {
	reader := bytes.NewReader([]byte(email))
	msg, err := mail.ReadMessage(reader)   // Read email using Go's net/mail
	mime, err := enmime.ParseMIMEBody(msg) // Parse message body with enmime
	if err != nil {
		return
	}

	json_email = obj.EmailJson{
		Addresses: []obj.EmailAddress{},
		Headers:   map[string][]string{},
	}

	for k, v := range msg.Header {
		switch strings.ToLower(k) {
		case "from", "to", "cc", "bcc", "reply-to", "sender":
			addr_ptrs, err := msg.Header.AddressList(k)
			if err == nil {
				var addr []obj.EmailAddress
				for _, addr_ptr := range addr_ptrs {
					addr = append(addr, obj.EmailAddress{
						*addr_ptr,
						k,
					})
				}
				json_email.Addresses = append(json_email.Addresses, addr...)
			}
		case "date":
			json_email.Date = v[0] //TODO : manage multiple date fields
		case "subject":
			json_email.Subject = v[0] //TODO : manage multiple subject fields
		default:
			json_email.Headers[k] = v
		}
	}

	json_email.Html = mime.HTML
	json_email.Plain = mime.Text
	json_email.IsTextFromHTML = mime.IsTextFromHTML

	if mime.Root != nil {
		//message was MIME encoded, build the mime tree

		root_boundary, _ := boundary(mime.GetHeader("Content-Type"))
		json_email.MimeRoot = obj.MimeRoot{
			Attachments_count: len(mime.Attachments),
			Root_boundary:     root_boundary,
			Inline_count:      len(mime.Inlines),
			Parts:             []obj.Part{},
		}

		child := mime.Root.FirstChild()
		if child != nil {
			json_email.MimeRoot.Parts = addChildPart(json_email.MimeRoot.Parts, child)
		}

	}
	return json_email, nil
}

// build part objects recursively
func addChildPart(parent []obj.Part, part enmime.MIMEPart) []obj.Part {

	child := obj.Part{
		ChildParts: []obj.Part{},
	}
	content_type := part.Header().Get("Content-Type")
	child.Boundary, _ = boundary(content_type)
	child.Content = part.Content()
	child.MediaType, child.Params, _ = mime.ParseMediaType(content_type)

	mediatype, _, _ := mime.ParseMediaType(part.Header().Get("Content-Disposition"))
	if strings.ToLower(mediatype) == "attachment" ||
		strings.ToLower(mediatype) == "inline" {
		child.Is_attachment = true
	}
	if strings.ToLower(mediatype) == "inline" {
		child.Is_inline = true
	}
	if strings.ToLower(child.MediaType) == "attachment" {
		child.Is_attachment = true
	}

	sub_child := part.FirstChild()
	if sub_child != nil {
		child.ChildParts = addChildPart(child.ChildParts, sub_child)
	}

	sibling := part.NextSibling()
	if sibling != nil {
		parent = addChildPart(parent, sibling)
	}

	return append(parent, child)
}

// returns the boundary string from a part header
func boundary(s string) (boundary string, err error) {
	mediatype, params, err := mime.ParseMediaType(s)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(mediatype, "multipart/") {
		boundary = params["boundary"]
		if err != nil {
			return "", err
		}
	}
	return boundary, nil
}
