package ssh3

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// EXPERIMENTAL_SPEC_VERSION specifies which version of the protocol this software
// is implementing.
// The protocol version string format is:
//
//	major + "." + minor[ + "_" + additional-version-information].
//
// It currently implements a first early version with no specification (alpha).
// Once IETF drafts get published, we plan on having versions such as
// 3.0_draft-michel-ssh3-XX when implementing the IETF specification from
// draft-michel-ssh3-XX.
const PROTOCOL_MAJOR int = 3
const PROTOCOL_MINOR int = 0
const PROTOCOL_EXPERIMENTAL_SPEC_VERSION string = "alpha-00"

const SOFTWARE_IMPLEMENTATION_NAME string = "francoismichel/ssh3"
const SOFTWARE_MAJOR int = 0
const SOFTWARE_MINOR int = 1
const SOFTWARE_PATCH int = 5

const SOFTWARE_RC int = 5

func ThisVersion() Version {
	return Version{
		protocolName: "SSH",
		protocolVersion: ProtocolVersion{
			Major:                   PROTOCOL_MAJOR,
			Minor:                   PROTOCOL_MINOR,
			ExperimentalSpecVersion: PROTOCOL_EXPERIMENTAL_SPEC_VERSION,
		},
		softwareVersion: SoftwareVersion{
			Major: SOFTWARE_MAJOR,
			Minor: SOFTWARE_MINOR,
			Patch: SOFTWARE_PATCH,
		},
	}
}

// Tells if the this version (a.k.a. the version returned by ThisVersion())
// is compatible with `other`.
func IsVersionSupported(other *Version) bool {
	this := ThisVersion()
	// right now, no check for protocol name as it is subject to change

	// strict protocol version checking
	if other.protocolVersion.Major != this.protocolVersion.Major || other.protocolVersion.Minor != this.protocolVersion.Minor {
		return false
	}

	// special case: to our knowledge, experimental spec version older than alpha-00 are only implemented by us (i.e. francoismichel/ssh3)
	if other.protocolVersion.ExperimentalSpecVersion == "" && other.softwareVersion.ImplementationName == SOFTWARE_IMPLEMENTATION_NAME &&
		other.softwareVersion.Major == 0 && other.softwareVersion.Minor == 1 && other.softwareVersion.Patch <= 5 {
		// then, only support software version >= 0.1.4
		return other.softwareVersion.Patch >= 4
	}

	// Starting from here, we have proper experimental spec version signalling.
	// this version is version alpha-00, other versions are not supported.
	// If a server receives a request with an unsupported spec version, the client should
	// start a new request with a compatible version.
	return other.protocolVersion.ExperimentalSpecVersion == "alpha-00"
}

type SoftwareVersion struct {
	ImplementationName string
	Major              int
	Minor              int
	Patch              int
}

type InvalidSoftwareVersion struct {
	softwareVersionString string
}

func (e InvalidSoftwareVersion) Error() string {
	return fmt.Sprintf("invalid protocol version string: %s", e.softwareVersionString)
}

func ParseSoftwareVersion(implementationName string, versionString string) (SoftwareVersion, error) {
	majorDotMinor := strings.Split(versionString, ".")
	if len(majorDotMinor) != 3 {
		log.Error().Msgf("bad SSH version major.minor.patch field")
		return SoftwareVersion{}, InvalidSoftwareVersion{softwareVersionString: versionString}
	}
	major, err := strconv.Atoi(majorDotMinor[0])
	if err != nil {
		log.Error().Msgf("bad software version major value")
		return SoftwareVersion{}, InvalidSoftwareVersion{softwareVersionString: versionString}
	}
	minor, err := strconv.Atoi(majorDotMinor[1])
	if err != nil {
		log.Error().Msgf("bad software version minor value")
		return SoftwareVersion{}, InvalidSoftwareVersion{softwareVersionString: versionString}
	}
	patch, err := strconv.Atoi(majorDotMinor[2])
	if err != nil {
		log.Error().Msgf("bad software version patch value")
		return SoftwareVersion{}, InvalidSoftwareVersion{softwareVersionString: versionString}
	}
	return SoftwareVersion{
		ImplementationName: implementationName,
		Major:              major,
		Minor:              minor,
		Patch:              patch,
	}, nil
}

func (v SoftwareVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

type ProtocolVersion struct {
	Major                   int
	Minor                   int
	ExperimentalSpecVersion string
}

func (v ProtocolVersion) String() string {
	return fmt.Sprintf("%d.%d_%s", v.Major, v.Minor, v.ExperimentalSpecVersion)
}

type InvalidProtocolVersion struct {
	protocolVersionString string
}

func (e InvalidProtocolVersion) Error() string {
	return fmt.Sprintf("invalid protocol version string: %s", e.protocolVersionString)
}

func ParseProtocolVersion(versionString string) (ProtocolVersion, error) {
	fields := strings.Split(versionString, "_")
	if len(fields) == 0 || len(fields) > 2 {
		return ProtocolVersion{}, InvalidProtocolVersion{protocolVersionString: versionString}
	}
	majorDotMinor := strings.Split(fields[0], ".")
	if len(majorDotMinor) != 2 {
		log.Error().Msgf("protocol version should be in format x.y, got: %s", fields[0])
		return ProtocolVersion{}, InvalidProtocolVersion{protocolVersionString: versionString}
	}
	major, err := strconv.Atoi(majorDotMinor[0])
	if err != nil {
		log.Error().Msgf("bad protocol version major value %s: %s", majorDotMinor[0], err)
		return ProtocolVersion{}, InvalidSSHVersion{versionString: versionString}
	}
	minor, err := strconv.Atoi(majorDotMinor[1])
	if err != nil {
		log.Error().Msgf("bad protocol version minor value %s: %s", majorDotMinor[1], err)
		return ProtocolVersion{}, InvalidSSHVersion{versionString: versionString}
	}
	experimentalSpecVersion := ""
	if len(fields) == 2 {
		experimentalSpecVersion = fields[1]
	}
	return ProtocolVersion{
		Major:                   major,
		Minor:                   minor,
		ExperimentalSpecVersion: experimentalSpecVersion,
	}, nil
}

type Version struct {
	protocolName    string // having the protocol name here might sound silly but there are discussions about updating the name right now and we want to support a change
	protocolVersion ProtocolVersion
	softwareVersion SoftwareVersion
}

func (v Version) GetProtocolVersion() ProtocolVersion {
	return v.protocolVersion
}

func (v Version) GetSoftwareVersion() SoftwareVersion {
	return v.softwareVersion
}

func NewVersion(protocolName string, protocolVersion ProtocolVersion, softwareVersion SoftwareVersion) *Version {
	return &Version{
		protocolName:    protocolName,
		protocolVersion: protocolVersion,
		softwareVersion: softwareVersion,
	}
}

type InvalidSSHVersion struct {
	versionString string
}

func (e InvalidSSHVersion) Error() string {
	return fmt.Sprintf("invalid ssh version string: %s", e.versionString)
}

type UnsupportedSSHVersion struct {
	versionString string
}

func (e UnsupportedSSHVersion) Error() string {
	return fmt.Sprintf("unsupported ssh version: %s", e.versionString)
}

// GetCurrentVersionString() returns the version string to be exchanged between two
// endpoints for version negotiation
func GetCurrentVersionString() string {
	return fmt.Sprintf("SSH %d.%d %s %d.%d.%d experimental_spec_version=%s", PROTOCOL_MAJOR, PROTOCOL_MINOR, SOFTWARE_IMPLEMENTATION_NAME, SOFTWARE_MAJOR, SOFTWARE_MINOR, SOFTWARE_PATCH, PROTOCOL_EXPERIMENTAL_SPEC_VERSION)
}

func ParseVersionString(versionString string) (version *Version, err error) {
	fields := strings.Fields(versionString)
	if len(fields) < 4 {
		log.Error().Msgf("bad SSH version fields")
		return nil, InvalidSSHVersion{versionString: versionString}
	}
	protocolName := fields[0]

	protocolVersion, err := ParseProtocolVersion(fields[1])
	if err != nil {
		log.Error().Msgf("could not parse protocol version: %s", err)
		return nil, err
	}

	softwareVersion, err := ParseSoftwareVersion(fields[2], fields[3])
	if err != nil {
		log.Error().Msgf("could not parse software version: %s", err)
		return nil, err
	}

	// Temporary tweak to announce a spec version while keeping compatibility with alpha-00 and older versions,
	// as alpha-00 and older versions do strict version checking and error as soon as the protocol version is not "3.0".
	// This will likely disappear once we decide to remove support for alpha-00 and older versions.
	// From that point onwards, the spec version will be announced as part of the version field.
	if len(fields) > 4 {
		for _, field := range fields[4:] {
			subfields := strings.Split(field, "=")
			if len(subfields) == 2 && subfields[0] == "experimental_spec_version" {
				protocolVersion.ExperimentalSpecVersion = subfields[1]
			} else {
				log.Debug().Msgf("skipping custom version field %s", field)
			}
		}

	}
	return NewVersion(protocolName, protocolVersion, softwareVersion), nil
}

// GetCurrentSoftwareVersion() returns the current software version to be displayed to the user
// For version string to be communicated between endpoints, use GetCurrentVersionString() instead.
func GetCurrentSoftwareVersion() string {
	versionStr := ThisVersion().softwareVersion.String()
	if SOFTWARE_RC > 0 {
		versionStr += fmt.Sprintf("-rc%d", SOFTWARE_RC)
	}
	return versionStr
}
