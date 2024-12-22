package blnkgo

import "errors"

// validate fields in Idenity based on the type of identity selected
func ValidateCreateIdentity(identity Identity) error {
	if identity.IdentityType == Individual {
		if identity.FirstName == "" {
			return errors.New("FirstName is required for Individual")
		}
		if identity.LastName == "" {
			return errors.New("LastName is required for Individual")
		}
		if identity.DOB == nil {
			return errors.New("DateOfBirth is required for Individual")
		}
		if identity.Gender == "" {
			return errors.New("gender is required for Individual")
		}
		if identity.Nationality == "" {
			return errors.New("nationality is required for Individual")
		}
	} else if identity.IdentityType == Organization {
		if identity.OrganizationName == "" {
			return errors.New("organizationName is required for Organization")
		}
	} else {
		return errors.New("invalid IdentityType")
	}
	return nil
}
