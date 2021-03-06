// Copyleft (ɔ) 2017 The Caliopen contributors.
// Use of this source code is governed by a GNU AFFERO GENERAL PUBLIC
// license (AGPL) that can be found in the LICENSE file.

package REST

import (
	"context"
	"errors"
	. "github.com/CaliOpen/Caliopen/src/backend/defs/go-objects"
	pi "github.com/CaliOpen/Caliopen/src/backend/main/go.main/pi"
	"github.com/satori/go.uuid"
)

func (rest *RESTfacility) LocalsIdentities(user_id string) ([]LocalIdentity, error) {
	return rest.store.GetLocalsIdentities(user_id)
}

// get contact from db
// aggregate contact's identities
// then update PI for each identity
func (rest *RESTfacility) ContactIdentities(user_id, contact_id string) (identities []ContactIdentity, err error) {
	_, e := uuid.FromString(contact_id)
	if user_id != "" && contact_id != "" && e == nil {
		contact, err := rest.store.GetContact(user_id, contact_id)
		if err != nil || contact == nil {
			err = errors.New("[RESTfacility.ContactIdentities] error when retrieving contact : " + err.Error())
			return []ContactIdentity{}, err
		}

		for _, email := range contact.Emails {
			identities = append(identities, ContactIdentity{
				Identifier:   email.Address,
				Label:        email.Label,
				PrivacyIndex: PrivacyIndex{},
				Protocol:     EmailProtocol,
			})
		}

		for _, identity := range contact.Identities {
			identities = append(identities, ContactIdentity{
				Identifier:   identity.Name,
				Label:        identity.Name,
				PrivacyIndex: PrivacyIndex{},
				Protocol:     identity.Type,
			})

		}

		for _, im := range contact.Ims {
			identities = append(identities, ContactIdentity{
				Identifier:   im.Address,
				Label:        im.Label,
				PrivacyIndex: PrivacyIndex{},
				Protocol:     im.Protocol,
			})
		}

		for _, phone := range contact.Phones {
			identities = append(identities, ContactIdentity{
				Identifier:   phone.Number,
				Label:        phone.Number,
				PrivacyIndex: PrivacyIndex{},
				Protocol:     SmsProtocol,
			})
		}

		for i := range identities {
			// (for now, this func is a monkey)
			pi.UpdatePIContactIdentity(context.TODO(), contact, &identities[i]) // TODO: do not use context pkg, it is not for this king of context
		}
	} else {
		err = errors.New("[RESTfacility.ContactIdentities] unprocessable parameters")
		return
	}
	return
}
