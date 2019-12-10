// Code generated by Flow Go SDK. DO NOT EDIT.

package generated

import (
	"bytes"
	"fmt"
	values1 "github.com/dapperlabs/flow-go/sdk/abi/encoding/values"
	"github.com/dapperlabs/flow-go/sdk/abi/types"
	"github.com/dapperlabs/flow-go/sdk/abi/values"
)

type AlbumView interface {
	Artist() ArtistView
	Name() string
	Rating() *uint8
	Year() uint16
}
type albumView struct {
	_artist ArtistView
	_name   string
	_rating *uint8
	_year   uint16
}

func (t *albumView) Artist() ArtistView {
	return t._artist
}
func (t *albumView) Name() string {
	return t._name
}
func (t *albumView) Rating() *uint8 {
	return t._rating
}
func (t *albumView) Year() uint16 {
	return t._year
}
func AlbumViewfromValue(value values.Value) (AlbumView, error) {
	composite, err := values.CastToComposite(value)
	if err != nil {
		return nil, err
	}

	_artist, err := ArtistViewfromValue(composite.Fields[uint(0x0)])
	if err != nil {
		return nil, err
	}

	_name, err := values.CastToString(composite.Fields[uint(0x1)])
	if err != nil {
		return nil, err
	}

	_rating, err := __converter0(composite.Fields[uint(0x2)])
	if err != nil {
		return nil, err
	}

	_year, err := values.CastToUInt16(composite.Fields[uint(0x3)])
	if err != nil {
		return nil, err
	}

	return &albumView{
		_artist: _artist,
		_name:   _name,
		_rating: _rating,
		_year:   _year,
	}, nil
}
func DecodeAlbumView(b []byte) (AlbumView, error) {
	r := bytes.NewReader(b)
	dec := values1.NewDecoder(r)
	v, err := dec.DecodeComposite(albumType)
	if err != nil {
		return nil, err
	}

	return AlbumViewfromValue(v)
}
func DecodeAlbumViewVariableSizedArray(b []byte) ([]AlbumView, error) {
	r := bytes.NewReader(b)
	dec := values1.NewDecoder(r)
	v, err := dec.DecodeVariableSizedArray(types.VariableSizedArray{ElementType: albumType})
	if err != nil {
		return nil, err
	}

	array := make([]AlbumView, len(v))
	for i, t := range v {
		array[i], err = AlbumViewfromValue(t.(values.Composite))
		if err != nil {
			return nil, err
		}

	}
	return array, nil
}

var albumType = types.Composite{
	Fields: map[string]*types.Field{
		"artist": {
			Identifier: "artist",
			Type:       artistType,
		},
		"name": {
			Identifier: "name",
			Type:       types.String{},
		},
		"rating": {
			Identifier: "rating",
			Type:       types.Optional{Of: types.UInt8{}},
		},
		"year": {
			Identifier: "year",
			Type:       types.UInt16{},
		},
	},
	Identifier: "Album",
	Initializers: [][]*types.Parameter{{&types.Parameter{
		Field: types.Field{
			Identifier: "artist",
			Type:       artistType,
		},
		Label: "",
	}, &types.Parameter{
		Field: types.Field{
			Identifier: "name",
			Type:       types.String{},
		},
		Label: "",
	}, &types.Parameter{
		Field: types.Field{
			Identifier: "year",
			Type:       types.UInt16{},
		},
		Label: "",
	}, &types.Parameter{
		Field: types.Field{
			Identifier: "rating",
			Type:       types.Optional{Of: types.UInt8{}},
		},
		Label: "",
	}}},
}

type AlbumConstructor interface {
	Encode() ([]byte, error)
}
type albumConstructor struct {
	artist ArtistView
	name   string
	year   uint16
	rating *uint8
}

func (p albumConstructor) toValue() values.ConstantSizedArray {
	return values.ConstantSizedArray{values.NewValueOrPanic(p.artist), values.NewValueOrPanic(p.name), values.NewValueOrPanic(p.year), values.NewValueOrPanic(p.rating)}
}
func (p albumConstructor) Encode() ([]byte, error) {
	var w bytes.Buffer
	encoder := values1.NewEncoder(&w)
	err := encoder.EncodeConstantSizedArray(p.toValue())
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
func NewAlbumConstructor(artist ArtistView, name string, year uint16, rating *uint8) (AlbumConstructor, error) {
	return albumConstructor{
		artist: artist,
		name:   name,
		rating: rating,
		year:   year,
	}, nil
}

type ArtistView interface {
	Country() string
	Members() *[]string
	Name() string
}
type artistView struct {
	_country string
	_members *[]string
	_name    string
}

func (t *artistView) Country() string {
	return t._country
}
func (t *artistView) Members() *[]string {
	return t._members
}
func (t *artistView) Name() string {
	return t._name
}
func ArtistViewfromValue(value values.Value) (ArtistView, error) {
	composite, err := values.CastToComposite(value)
	if err != nil {
		return nil, err
	}

	_country, err := values.CastToString(composite.Fields[uint(0x0)])
	if err != nil {
		return nil, err
	}

	_members, err := __converter1(composite.Fields[uint(0x1)])
	if err != nil {
		return nil, err
	}

	_name, err := values.CastToString(composite.Fields[uint(0x2)])
	if err != nil {
		return nil, err
	}

	return &artistView{
		_country: _country,
		_members: _members,
		_name:    _name,
	}, nil
}
func DecodeArtistView(b []byte) (ArtistView, error) {
	r := bytes.NewReader(b)
	dec := values1.NewDecoder(r)
	v, err := dec.DecodeComposite(artistType)
	if err != nil {
		return nil, err
	}

	return ArtistViewfromValue(v)
}
func DecodeArtistViewVariableSizedArray(b []byte) ([]ArtistView, error) {
	r := bytes.NewReader(b)
	dec := values1.NewDecoder(r)
	v, err := dec.DecodeVariableSizedArray(types.VariableSizedArray{ElementType: artistType})
	if err != nil {
		return nil, err
	}

	array := make([]ArtistView, len(v))
	for i, t := range v {
		array[i], err = ArtistViewfromValue(t.(values.Composite))
		if err != nil {
			return nil, err
		}

	}
	return array, nil
}

var artistType = types.Composite{
	Fields: map[string]*types.Field{
		"country": {
			Identifier: "country",
			Type:       types.String{},
		},
		"members": {
			Identifier: "members",
			Type:       types.Optional{Of: types.VariableSizedArray{ElementType: types.String{}}},
		},
		"name": {
			Identifier: "name",
			Type:       types.String{},
		},
	},
	Identifier: "Artist",
	Initializers: [][]*types.Parameter{{&types.Parameter{
		Field: types.Field{
			Identifier: "name",
			Type:       types.String{},
		},
		Label: "",
	}, &types.Parameter{
		Field: types.Field{
			Identifier: "members",
			Type:       types.Optional{Of: types.VariableSizedArray{ElementType: types.String{}}},
		},
		Label: "",
	}, &types.Parameter{
		Field: types.Field{
			Identifier: "country",
			Type:       types.String{},
		},
		Label: "",
	}}},
}

type ArtistConstructor interface {
	Encode() ([]byte, error)
}
type artistConstructor struct {
	name    string
	members *[]string
	country string
}

func (p artistConstructor) toValue() values.ConstantSizedArray {
	return values.ConstantSizedArray{values.NewValueOrPanic(p.name), values.NewValueOrPanic(p.members), values.NewValueOrPanic(p.country)}
}
func (p artistConstructor) Encode() ([]byte, error) {
	var w bytes.Buffer
	encoder := values1.NewEncoder(&w)
	err := encoder.EncodeConstantSizedArray(p.toValue())
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
func NewArtistConstructor(name string, members *[]string, country string) (ArtistConstructor, error) {
	return artistConstructor{
		country: country,
		members: members,
		name:    name,
	}, nil
}
func __converter0(p values.Value) (*uint8, error) {
	var ret0 uint8
	var go0 interface{}
	cast0, ok := p.(values.Optional)
	if !ok {
		return nil, fmt.Errorf("cannot cast %T", p)

	}
	go0 = cast0.ToGoValue()

	var err error
	if go0 == nil {
		return nil, nil
	} else {
		cast1, ok := cast0.Value.(values.UInt8)
		if !ok {
			return nil, fmt.Errorf("cannot cast %T", cast0.Value)

		}

		ret0, err = values.CastToUInt8(cast1)
		if err != nil {
			return nil, err
		}

	}
	return &ret0, nil

}
func __converter1(p values.Value) (*[]string, error) {
	var ret0 []string
	var go0 interface{}
	cast0, ok := p.(values.Optional)
	if !ok {
		return nil, fmt.Errorf("cannot cast %T", p)

	}
	go0 = cast0.ToGoValue()

	var err error
	if go0 == nil {
		return nil, nil
	} else {
		var ret1 []string
		cast1, ok := cast0.Value.(values.VariableSizedArray)
		if !ok {
			return nil, fmt.Errorf("cannot cast %T", cast0.Value)

		}

		if err != nil {
			return nil, err
		}
		ret1 = make([]string, len(cast1))
		for i1, elem1 := range cast1 {
			cast2, ok := elem1.(values.String)
			if !ok {
				return nil, fmt.Errorf("cannot cast %T", elem1)

			}

			ret1[i1], err = values.CastToString(cast2)
			if err != nil {
				return nil, err
			}

		}
		ret0 = ret1

	}
	return &ret0, nil

}
