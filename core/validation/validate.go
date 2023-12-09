package validation

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		// You can register custom validations here, if needed
		err := validate.RegisterValidation("min_ids", minIDsRule)
		if err != nil {
			log.Error().Err(err)
		}
		err = validate.RegisterValidation("validRelationID", validRelationID)
		if err != nil {
			log.Error().Err(err)
		}
	})
	return validate
}

func validRelationID(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	if field.Kind() != reflect.Struct || !field.FieldByName("ID").IsValid() {
		return false
	}

	idField := field.FieldByName("ID")
	if idField.Uint() > 0 {
		return true
	}
	return false
}

type DefaultValidator struct{}

func (v *DefaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == "struct" {
		validatorMain := GetValidator() // from our utils package
		return validatorMain.Struct(obj)
	}
	return nil
}

func (v *DefaultValidator) Engine() interface{} {
	return GetValidator()
}

func kindOfData(data interface{}) string {
	return fmt.Sprintf("%T", data)
}

//goland:noinspection ALL,SpellCheckingInspection
var TagMessages = map[string]string{
	// Custom Tag
	"validRelationID": "Invalid ID for the relationship",
	"LoginViaEnum":    "Invalid Login Via type, allowed types are STAFF_ID, EMAIL, USERNAME",
	"StatusEnum":      "Invalid Status, allowed statuses are ENABLE, DISABLE, DRAFT",

	// Fields Related Tag
	"eqcsfield":     "Field must be equal to %s (case-sensitive).",
	"eqfield":       "Field must be equal to %s.",
	"fieldcontains": "Field must contain the value of %s.",
	"fieldexcludes": "Field must not contain the value of %s.",
	"gtcsfield":     "Field must be greater than %s (case-sensitive).",
	"gtecsfield":    "Field must be greater than or equal to %s (case-sensitive).",
	"gtefield":      "Field must be greater than or equal to %s.",
	"gtfield":       "Field must be greater than %s.",
	"ltcsfield":     "Field must be less than %s (case-sensitive).",
	"ltecsfield":    "Field must be less than or equal to %s (case-sensitive).",
	"ltefield":      "Field must be less than or equal to %s.",
	"ltfield":       "Field must be less than %s.",
	"necsfield":     "Field must not be equal to %s (case-sensitive).",
	"nefield":       "Field must not be equal to %s.",

	// Network Related Tag
	"cidr":             "Field must be a valid CIDR address.",
	"cidrv4":           "Field must be a valid IPv4 CIDR address.",
	"cidrv6":           "Field must be a valid IPv6 CIDR address.",
	"datauri":          "Field must be a valid data URI.",
	"fqdn":             "Field must be a valid fully qualified domain name.",
	"hostname":         "Field must be a valid hostname.",
	"hostname_port":    "Field must be a valid hostname and port.",
	"hostname_rfc1123": "Field must be a valid RFC-1123 hostname.",
	"ip":               "Field must be a valid IP address.",
	"ip4_addr":         "Field must be a valid IPv4 address.",
	"ip6_addr":         "Field must be a valid IPv6 address.",
	"ip_addr":          "Field must be a valid IP address (either IPv4 or IPv6).",
	"ipv4":             "Field must be a valid IPv4 address.",
	"ipv6":             "Field must be a valid IPv6 address.",
	"mac":              "Field must be a valid MAC address.",
	"tcp4_addr":        "Field must be a valid TCP IPv4 address.",
	"tcp6_addr":        "Field must be a valid TCP IPv6 address.",
	"tcp_addr":         "Field must be a valid TCP address.",
	"udp4_addr":        "Field must be a valid UDP IPv4 address.",
	"udp6_addr":        "Field must be a valid UDP IPv6 address.",
	"udp_addr":         "Field must be a valid UDP address.",
	"unix_addr":        "Field must be a valid Unix address.",
	"uri":              "Field must be a valid URI.",
	"url":              "Field must be a valid URL.",
	"http_url":         "Field must be a valid HTTP or HTTPS URL.",
	"url_encoded":      "Field must be URL-encoded.",
	"urn_rfc2141":      "Field must be a valid URN as per RFC 2141.",

	// Strings Related Tag
	"alpha":           "Field must contain only alphabetic characters.",
	"alphanum":        "Field must contain only alphanumeric characters.",
	"alphanumunicode": "Field must contain only Unicode alphanumeric characters.",
	"alphaunicode":    "Field must contain only Unicode alphabetic characters.",
	"ascii":           "Field must contain only ASCII characters.",
	"boolean":         "Field must be either true or false.",
	"contains":        "Field must contain the specified substring.",
	"containsany":     "Field must contain any of the specified characters.",
	"containsrune":    "Field must contain the specified rune.",
	"endsnotwith":     "Field must not end with the specified substring.",
	"endswith":        "Field must end with the specified substring.",
	"excludes":        "Field must not contain the specified substring.",
	"excludesall":     "Field must not contain any of the specified characters.",
	"excludesrune":    "Field must not contain the specified rune.",
	"lowercase":       "Field must contain only lowercase characters.",
	"multibyte":       "Field must contain one or more multibyte characters.",
	"number":          "Field must be a number.",
	"numeric":         "Field must be numeric.",
	"printascii":      "Field must contain only printable ASCII characters.",
	"startsnotwith":   "Field must not start with the specified substring.",
	"startswith":      "Field must start with the specified substring.",
	"uppercase":       "Field must contain only uppercase characters.",

	// Format Related Tag
	"base64":                        "Field must be a valid Base64 string.",
	"base64url":                     "Field must be a valid Base64 URL string.",
	"base64rawurl":                  "Field must be a valid raw Base64 URL string.",
	"bic":                           "Field must be a valid Business Identifier Code (BIC).",
	"bcp47_language_tag":            "Field must be a valid BCP 47 language tag.",
	"btc_addr":                      "Field must be a valid Bitcoin address.",
	"btc_addr_bech32":               "Field must be a valid Bech32-encoded Bitcoin address.",
	"credit_card":                   "Field must be a valid credit card number.",
	"mongodb":                       "Field must be a valid MongoDB identifier.",
	"cron":                          "Field must be a valid cron expression.",
	"spicedb":                       "Field must be a valid Spicedb expression.",
	"datetime":                      "Field must be a valid date-time expression.",
	"e164":                          "Field must be a valid E.164 formatted phone number.",
	"email":                         "Field must be a valid email address.",
	"eth_addr":                      "Field must be a valid Ethereum address.",
	"hexadecimal":                   "Field must be a valid hexadecimal number.",
	"hexcolor":                      "Field must be a valid hex color code.",
	"hsl":                           "Field must be a valid HSL color.",
	"hsla":                          "Field must be a valid HSLA color.",
	"html":                          "Field must be a valid HTML text.",
	"html_encoded":                  "Field must be a valid HTML-encoded text.",
	"isbn":                          "Field must be a valid ISBN.",
	"isbn10":                        "Field must be a valid ISBN-10.",
	"isbn13":                        "Field must be a valid ISBN-13.",
	"iso3166_1_alpha2":              "Field must be a valid ISO 3166-1 alpha-2 country code.",
	"iso3166_1_alpha3":              "Field must be a valid ISO 3166-1 alpha-3 country code.",
	"iso3166_1_alpha_numeric":       "Field must be a valid ISO 3166-1 alpha-numeric country code.",
	"iso3166_2":                     "Field must be a valid ISO 3166-2 region code.",
	"iso4217":                       "Field must be a valid ISO 4217 currency code.",
	"json":                          "Field must be a valid JSON.",
	"jwt":                           "Field must be a valid JWT token.",
	"latitude":                      "Field must be a valid latitude.",
	"longitude":                     "Field must be a valid longitude.",
	"luhn_checksum":                 "Field must satisfy the Luhn checksum.",
	"postcode_iso3166_alpha2":       "Field must be a valid ISO 3166-1 alpha-2 postcode.",
	"postcode_iso3166_alpha2_field": "Field must be a valid ISO 3166-1 alpha-2 postcode for the specific field.",
	"rgb":                           "Field must be a valid RGB color.",
	"rgba":                          "Field must be a valid RGBA color.",
	"ssn":                           "Field must be a valid Social Security Number.",
	"timezone":                      "Field must be a valid time zone.",
	"uuid":                          "Field must be a valid UUID.",
	"uuid3":                         "Field must be a valid UUID v3.",
	"uuid3_rfc4122":                 "Field must be a valid RFC 4122 UUID v3.",
	"uuid4":                         "Field must be a valid UUID v4.",
	"uuid4_rfc4122":                 "Field must be a valid RFC 4122 UUID v4.",
	"uuid5":                         "Field must be a valid UUID v5.",
	"uuid5_rfc4122":                 "Field must be a valid RFC 4122 UUID v5.",
	"uuid_rfc4122":                  "Field must be a valid RFC 4122 UUID.",
	"md4":                           "Field must be a valid MD4 hash.",
	"md5":                           "Field must be a valid MD5 hash.",
	"sha256":                        "Field must be a valid SHA-256 hash.",
	"sha384":                        "Field must be a valid SHA-384 hash.",
	"sha512":                        "Field must be a valid SHA-512 hash.",
	"ripemd128":                     "Field must be a valid RIPEMD-128 hash.",
	"tiger128":                      "Field must be a valid Tiger-128 hash.",
	"tiger160":                      "Field must be a valid Tiger-160 hash.",
	"tiger192":                      "Field must be a valid Tiger-192 hash.",
	"semver":                        "Field must be a valid Semantic Version.",
	"ulid":                          "Field must be a valid ULID.",
	"cve":                           "Field must be a valid Common Vulnerabilities and Exposures (CVE) identifier.",

	// Comparision Related Tag
	"eq":             "Field must be equal to %v.",
	"eq_ignore_case": "Field must be equal to %v, case insensitive.",
	"gt":             "Field must be greater than %v.",
	"gte":            "Field must be greater than or equal to %v.",
	"lt":             "Field must be less than %v.",
	"lte":            "Field must be less than or equal to %v.",
	"ne":             "Field must not be equal to %v.",
	"ne_ignore_case": "Field must not be equal to %v, case insensitive.",

	// File and Directory related tags
	"dir":       "Field must be a valid directory.",
	"dirpath":   "Field must be a valid directory path.",
	"file":      "Field must be a valid file.",
	"filepath":  "Field must be a valid file path.",
	"image":     "Field must be a valid image file.",
	"isdefault": "Field must have the default value.",

	// Length and QuantityLine related tags
	"len":   "Field must have a length of %v.",
	"max":   "Field must have a maximum of %v.",
	"min":   "Field must have a minimum of %v.",
	"oneof": "Field must be one of [%v].",

	// Requirement related tags
	"required":             "Field is required.",
	"required_if":          "Field is required if %v.",
	"required_unless":      "Field is required unless %v.",
	"required_with":        "Field is required when %v is provided.",
	"required_with_all":    "Field is required when all of %v are provided.",
	"required_without":     "Field is required when %v is not provided.",
	"required_without_all": "Field is required when none of %v are provided.",

	// Exclusion related tags
	"excluded_if":          "Field must be empty if %v.",
	"excluded_unless":      "Field must be empty unless %v.",
	"excluded_with":        "Field must be empty when %v is provided.",
	"excluded_with_all":    "Field must be empty when all of %v are provided.",
	"excluded_without":     "Field must be empty when %v is not provided.",
	"excluded_without_all": "Field must be empty when none of %v are provided.",

	// Unique value
	"unique": "Field must be unique.",
}
