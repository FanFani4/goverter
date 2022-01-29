// Code generated by github.com/FanFani4/goverter, DO NOT EDIT.

package generated

import errors "github.com/FanFani4/goverter/example/errors"

type ConverterImpl struct{}

func (c *ConverterImpl) ToAPIApartment(source errors.DBApartment) errors.APIApartment {
	var errorsAPIApartment errors.APIApartment
	errorsAPIApartment.Position = source.Position
	errorsAPIApartment.Owner = errors.ToAPIPerson(source.Owner)
	return errorsAPIApartment
}
func (c *ConverterImpl) ToDBApartment(source errors.APIApartment) (errors.DBApartment, error) {
	var errorsDBApartment errors.DBApartment
	errorsDBApartment.Position = source.Position
	errorsDBPerson, err := errors.ToDBPerson(source.Owner)
	if err != nil {
		return errorsDBApartment, err
	}
	errorsDBApartment.Owner = errorsDBPerson
	return errorsDBApartment, nil
}
