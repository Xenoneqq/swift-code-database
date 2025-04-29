package main_test

import (
	"testing"

	"github.com/Xenoneqq/swift-code-database/handler"
	"github.com/stretchr/testify/assert"
)

func TestCountryCheck_Positive(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(0, handler.CheckCountry("PL", "POLAND"), "Correct ISO2 and country name for Poland should pass")
	assert.Equal(0, handler.CheckCountry("BE", "BELGIUM"), "Correct ISO2 and country name for Belgium should pass")
	assert.Equal(0, handler.CheckCountry("DE", "GERMANY"), "Correct ISO2 and country name for Germany should pass")
	assert.Equal(0, handler.CheckCountry("FR", "FRANCE"), "Correct ISO2 and country name for France should pass")
	assert.Equal(0, handler.CheckCountry("ES", "SPAIN"), "Correct ISO2 and country name for Spain should pass")
	assert.Equal(0, handler.CheckCountry("IT", "ITALY"), "Correct ISO2 and country name for Italy should pass")
	assert.Equal(0, handler.CheckCountry("NL", "NETHERLANDS"), "Correct ISO2 and country name for Netherlands should pass")
	assert.Equal(0, handler.CheckCountry("SE", "SWEDEN"), "Correct ISO2 and country name for Sweden should pass")
	assert.Equal(0, handler.CheckCountry("CH", "SWITZERLAND"), "Correct ISO2 and country name for Switzerland should pass")
	assert.Equal(0, handler.CheckCountry("AT", "AUSTRIA"), "Correct ISO2 and country name for Austria should pass")
	assert.Equal(0, handler.CheckCountry("DK", "DENMARK"), "Correct ISO2 and country name for Denmark should pass")
	assert.Equal(0, handler.CheckCountry("NO", "NORWAY"), "Correct ISO2 and country name for Norway should pass")
	assert.Equal(0, handler.CheckCountry("FI", "FINLAND"), "Correct ISO2 and country name for Finland should pass")
	assert.Equal(0, handler.CheckCountry("IE", "IRELAND"), "Correct ISO2 and country name for Ireland should pass")
	assert.Equal(0, handler.CheckCountry("PT", "PORTUGAL"), "Correct ISO2 and country name for Portugal should pass")
	assert.Equal(0, handler.CheckCountry("GR", "GREECE"), "Correct ISO2 and country name for Greece should pass")
	assert.Equal(0, handler.CheckCountry("HU", "HUNGARY"), "Correct ISO2 and country name for Hungary should pass")
}

func TestCountryCheck_Negative(t *testing.T) {
	assert := assert.New(t)

	assert.NotEqual(0, handler.CheckCountry("PL", "GERMANY"), "Incorrect country name for PL should fail")
	assert.NotEqual(0, handler.CheckCountry("BE", "POLAND"), "Incorrect country name for BE should fail")
	assert.NotEqual(0, handler.CheckCountry("DE", "FRANCE"), "Incorrect country name for DE should fail")
	assert.NotEqual(0, handler.CheckCountry("FR", "GERMANY"), "Incorrect country name for FR should fail")

	assert.NotEqual(0, handler.CheckCountry("GB", "UNITED KINGDOM"), "Correct name, but different ISO2 should fail")
	assert.NotEqual(0, handler.CheckCountry("PL", "POLSKA"), "Correct ISO2, but slightly different name should fail (case-sensitive)")
	assert.NotEqual(0, handler.CheckCountry("POL", "POLAND"), "Incorrect ISO2 format should fail")
	assert.NotEqual(0, handler.CheckCountry("", "POLAND"), "Empty ISO2 should fail")
	assert.NotEqual(0, handler.CheckCountry("PL", ""), "Empty country name should fail")
}
