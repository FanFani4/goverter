package house_test

import (
	"database/sql"
	"testing"

	"github.com/FanFani4/goverter/example/house"
	"github.com/FanFani4/goverter/example/house/generated"
	"github.com/stretchr/testify/require"
)

func TestConverter(t *testing.T) {
	var c house.Converter = &generated.ConverterImpl{}

	input := house.DBHouse{
		Address: "Somewhere",
		Apartments: map[int]house.DBApartment{
			5: {
				Position: 1,
				Owner: house.DBPerson{
					ID:         5,
					Name:       "FanFani4",
					MiddleName: sql.NullString{},
					Friends: []house.DBPerson{{
						ID:         5,
						Name:       "my cat (:",
						MiddleName: sql.NullString{String: "sir", Valid: true},
						Friends:    []house.DBPerson{},
					}},
				},
				CoResident: []house.DBPerson{{
					ID:         5,
					Name:       "my cat (:",
					MiddleName: sql.NullString{String: "sir", Valid: true},
					Friends:    []house.DBPerson{},
				}},
			},
		},
	}

	actual := c.ConvertHouse(input)

	expected := house.APIHouse{
		Address: "Somewhere",
		Apartments: map[house.APIRoomNR]house.APIApartment{
			house.APIRoomNR(5): {
				Position:  1,
				OwnerName: "FanFani4",
				Owner: house.APIPerson{
					ID:         5,
					FirstName:  p("FanFani4"),
					MiddleName: nil,
					Friends: []house.APIPerson{{
						ID:         5,
						FirstName:  p("my cat (:"),
						MiddleName: p("sir"),
						Friends:    []house.APIPerson{},
					}},
				},
				CoResident: []house.APIPerson{{
					ID:         5,
					FirstName:  p("my cat (:"),
					MiddleName: p("sir"),
					Friends:    []house.APIPerson{},
				}},
			},
		},
	}

	require.Equal(t, expected, actual)
}

func p(s string) *string {
	return &s
}
