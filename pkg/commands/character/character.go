package character

import (
	"fmt"
)

// ============================================================================
// Context
// ============================================================================
// This package implements ESC/POS commands for character formatting and appearance.
// ESC/POS is the command system used by thermal receipt printers to control
// character fonts, sizes, styles, effects, international character sets,
// code pages, and user-defined characters.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Spacing represents character spacing in dots
type Spacing byte

// PrintMode represents the print mode bits for character formatting
type PrintMode byte

const (
	// FontAPm represents Font A in print mode
	FontAPm PrintMode = 0x00
	// FontBPm represents Font B in print mode
	FontBPm PrintMode = 0x01
	// EmphasizedOffPm represents emphasized mode off in print mode
	EmphasizedOffPm PrintMode = 0x00
	// EmphasizedOnPm represents emphasized mode on in print mode
	EmphasizedOnPm PrintMode = 0x08
	// DoubleHeightOffPm represents double height mode off in print mode
	DoubleHeightOffPm PrintMode = 0x00
	// DoubleHeightOnPm represents double height mode on in print mode
	DoubleHeightOnPm PrintMode = 0x10
	// DoubleWidthOffPm represents double width mode off in print mode
	DoubleWidthOffPm PrintMode = 0x00
	// DoubleWidthOnPm represents double width mode on in print mode
	DoubleWidthOnPm PrintMode = 0x20
	// UnderlineOffPm represents underline mode off in print mode
	UnderlineOffPm PrintMode = 0x00
	// UnderlineOnPm represents underline mode on in print mode
	UnderlineOnPm PrintMode = 0x80
)

// UnderlineMode represents the underline thickness mode
type UnderlineMode byte

const (
	// NoDot represents no underline mode
	NoDot UnderlineMode = 0x00
	// OneDot represents single dot underline mode
	OneDot UnderlineMode = 0x01
	// TwoDot represents two dot underline mode
	TwoDot UnderlineMode = 0x02
	// NoDotASCII represents no underline mode (ASCII)
	NoDotASCII UnderlineMode = '0'
	// OneDotASCII represents single dot underline mode (ASCII)
	OneDotASCII UnderlineMode = '1'
	// TwoDotASCII represents two dot underline mode (ASCII)
	TwoDotASCII UnderlineMode = '2'
)

// EmphasizedMode represents the emphasized text mode
type EmphasizedMode byte

const (
	// OffEm represents emphasized mode off (LSB = 0)
	OffEm EmphasizedMode = 0x00
	// OnEm represents emphasized mode on (LSB = 1)
	OnEm EmphasizedMode = 0x01
)

// DoubleStrikeMode represents the double-strike printing mode
type DoubleStrikeMode byte

const (
	// OffDsm represents double-strike mode off (LSB = 0)
	OffDsm DoubleStrikeMode = 0x00
	// OnDsm represents double-strike mode on (LSB = 1)
	OnDsm DoubleStrikeMode = 0x01
)

// FontType represents the character font type
type FontType byte

const (
	// FontA represents Font A type
	FontA FontType = 0x00
	// FontB represents Font B type
	FontB FontType = 0x01
	// FontC represents Font C type
	FontC FontType = 0x02
	// FontD represents Font D type
	FontD FontType = 0x03
	// FontE represents Font E type
	FontE FontType = 0x04
	// FontAAscii represents Font A type (ASCII)
	FontAAscii FontType = '0'
	// FontBAscii represents Font B type (ASCII)
	FontBAscii FontType = '1'
	// FontCAscii represents Font C type (ASCII)
	FontCAscii FontType = '2'
	// FontDAscii represents Font D type (ASCII)
	FontDAscii FontType = '3'
	// FontEAscii represents Font E type (ASCII)
	FontEAscii FontType = '4'
	// SpecialFontA represents Special Font A type
	SpecialFontA FontType = 97
	// SpecialFontB represents Special Font B type
	SpecialFontB FontType = 98
)

// InternationalSet represents the international character set
type InternationalSet byte

const (
	// USA represents the USA international character set
	USA InternationalSet = 0
	// France represents the France international character set
	France InternationalSet = 1
	// Germany represents the Germany international character set
	Germany InternationalSet = 2
	// UK represents the UK international character set
	UK InternationalSet = 3
	// DenmarkI represents the Denmark I international character set
	DenmarkI InternationalSet = 4
	// Sweden represents the Sweden international character set
	Sweden InternationalSet = 5
	// Italy represents the Italy international character set
	Italy InternationalSet = 6
	// SpainI represents the Spain I international character set
	SpainI InternationalSet = 7
	// Japan represents the Japan international character set
	Japan InternationalSet = 8
	// Norway represents the Norway international character set
	Norway InternationalSet = 9
	// DenmarkII represents the Denmark II international character set
	DenmarkII InternationalSet = 10
	// SpainII represents the Spain II international character set
	SpainII InternationalSet = 11
	// LatinAmerica represents the Latin America international character set
	LatinAmerica InternationalSet = 12
	// Korea represents the Korea international character set
	Korea InternationalSet = 13
	// SloveniaCro represents the Slovenia/Croatia international character set
	SloveniaCro InternationalSet = 14
	// China represents the China international character set
	China InternationalSet = 15
	// Vietnam represents the Vietnam international character set
	Vietnam InternationalSet = 16
	// Arabia represents the Arabia international character set
	Arabia InternationalSet = 17

	// IndiaDevanagari represents the Devanagari character set (model-dependent)
	IndiaDevanagari InternationalSet = 66
	// IndiaBengali represents the Bengali character set (model-dependent)
	IndiaBengali InternationalSet = 67
	// IndiaTamil represents the Tamil character set (model-dependent)
	IndiaTamil InternationalSet = 68
	// IndiaTelugu represents the Telugu character set (model-dependent)
	IndiaTelugu InternationalSet = 69
	// IndiaAssamese represents the Assamese character set (model-dependent)
	IndiaAssamese InternationalSet = 70
	// IndiaOriya represents the Oriya character set (model-dependent)
	IndiaOriya InternationalSet = 71
	// IndiaKannada represents the Kannada character set (model-dependent)
	IndiaKannada InternationalSet = 72
	// IndiaMalayalam represents the Malayalam character set (model-dependent)
	IndiaMalayalam InternationalSet = 73
	// IndiaGujarati represents the Gujarati character set (model-dependent)
	IndiaGujarati InternationalSet = 74
	// IndiaPunjabi represents the Punjabi character set (model-dependent)
	IndiaPunjabi InternationalSet = 75
	// IndiaMarathi represents the Marathi character set (model-dependent)
	IndiaMarathi InternationalSet = 82
)

// RotationMode represents the character rotation mode
type RotationMode byte

const (
	// NoRotation represents no rotation mode
	NoRotation RotationMode = 0x00
	// On90Dot1 represents 90-degree rotation with 1-dot spacing
	On90Dot1 RotationMode = 0x01
	// On90Dot15 represents 90-degree rotation with 1.5-dot spacing
	On90Dot15 RotationMode = 0x02
	// NoRotationASCII represents no rotation mode (ASCII mode)
	NoRotationASCII RotationMode = '0'
	// On90Dot1Ascii represents 90-degree rotation with 1-dot spacing (ASCII mode)
	On90Dot1Ascii RotationMode = '1'
	// On90Dot15Ascii represents 90-degree rotation with 1.5-dot spacing (ASCII mode)
	On90Dot15Ascii RotationMode = '2'
)

// PrintColor represents the print color selection
type PrintColor byte

const (
	// Black represents black print color
	Black PrintColor = 0x00
	// Red represents red print color
	Red PrintColor = 0x01
	// BlackASCII represents black print color (ASCII mode)
	BlackASCII PrintColor = '0'
	// RedASCII represents red print color (ASCII mode)
	RedASCII PrintColor = '1'
)

// CodeTable represents the character code table page
type CodeTable byte

const (
	// PC437 represents USA, Standard Europe code table
	PC437 CodeTable = 0
	// Katakana represents Katakana code table
	Katakana CodeTable = 1
	// PC850 represents Multilingual code table
	PC850 CodeTable = 2
	// PC860 represents Portuguese code table
	PC860 CodeTable = 3
	// PC863 represents Canadian-French code table
	PC863 CodeTable = 4
	// PC865 represents Nordic code table
	PC865 CodeTable = 5
	// Hiragana represents Hiragana code table
	Hiragana CodeTable = 6
	// OnePassKanji1 represents One-Pass Kanji Type 1 code table
	OnePassKanji1 CodeTable = 7
	// OnePassKanji2 represents One-Pass Kanji Type 2 code table
	OnePassKanji2 CodeTable = 8
	// PC851 represents Greek code table
	PC851 CodeTable = 11
	// PC853 represents Turkish code table
	PC853 CodeTable = 12
	// PC857 represents Turkish code table
	PC857 CodeTable = 13
	// PC737 represents Greek code table
	PC737 CodeTable = 14
	// ISO88597 represents Greek code table
	ISO88597 CodeTable = 15
	// WPC1252 represents WPC1252 code table
	WPC1252 CodeTable = 16
	// PC866 represents Cyrillic #2 code table
	PC866 CodeTable = 17
	// PC852 represents Latin 2 code table
	PC852 CodeTable = 18
	// PC858 represents Euro code table
	PC858 CodeTable = 19
	// ThaiCode42 represents Thai Code 42 code table
	ThaiCode42 CodeTable = 20
	// ThaiCode11 represents Thai Code 11 code table
	ThaiCode11 CodeTable = 21
	// ThaiCode13 represents Thai Code 13 code table
	ThaiCode13 CodeTable = 22
	// ThaiCode14 represents Thai Code 14 code table
	ThaiCode14 CodeTable = 23
	// ThaiCode16 represents Thai Code 16 code table
	ThaiCode16 CodeTable = 24
	// ThaiCode17 represents Thai Code 17 code table
	ThaiCode17 CodeTable = 25
	// ThaiCode18 represents Thai Code 18 code table
	ThaiCode18 CodeTable = 26
	// TCVN31 represents TCVN-3 Type 1 Vietnamese code table
	TCVN31 CodeTable = 30
	// TCVN32 represents TCVN-3 Type 2 Vietnamese code table
	TCVN32 CodeTable = 31
	// PC720 represents Arabic code table
	PC720 CodeTable = 32
	// WPC775 represents Baltic Rim code table
	WPC775 CodeTable = 33
	// PC855 represents Cyrillic code table
	PC855 CodeTable = 34
	// PC861 represents Icelandic code table
	PC861 CodeTable = 35
	// PC862 represents Hebrew code table
	PC862 CodeTable = 36
	// PC864 represents Arabic code table
	PC864 CodeTable = 37
	// PC869 represents Greek code table
	PC869 CodeTable = 38
	// ISO88592 represents Latin 2 code table
	ISO88592 CodeTable = 39
	// ISO885915 represents Latin 9 code table
	ISO885915 CodeTable = 40
	// PC1098 represents Farsi code table
	PC1098 CodeTable = 41
	// PC1118 represents Lithuanian code table
	PC1118 CodeTable = 42
	// PC1119 represents Lithuanian code table
	PC1119 CodeTable = 43
	// PC1125 represents Ukrainian code table
	PC1125 CodeTable = 44
	// WPC1250 represents Latin 2 code table
	WPC1250 CodeTable = 45
	// WPC1251 represents Cyrillic code table
	WPC1251 CodeTable = 46
	// WPC1253 represents Greek code table
	WPC1253 CodeTable = 47
	// WPC1254 represents Turkish code table
	WPC1254 CodeTable = 48
	// WPC1255 represents Hebrew code table
	WPC1255 CodeTable = 49
	// WPC1256 represents Arabic code table
	WPC1256 CodeTable = 50
	// WPC1257 represents Baltic Rim code table
	WPC1257 CodeTable = 51
	// WPC1258 represents Vietnamese code table
	WPC1258 CodeTable = 52
	// KZ1048 represents Kazakhstan code table
	KZ1048 CodeTable = 53

	// Devanagari represents Devanagari code table (model-dependent)
	Devanagari CodeTable = 66
	// Bengali represents Bengali code table (model-dependent)
	Bengali CodeTable = 67
	// Tamil represents Tamil code table (model-dependent)
	Tamil CodeTable = 68
	// Telugu represents Telugu code table (model-dependent)
	Telugu CodeTable = 69
	// Assamese represents Assamese code table (model-dependent)
	Assamese CodeTable = 70
	// Oriya represents Oriya code table (model-dependent)
	Oriya CodeTable = 71
	// Kannada represents Kannada code table (model-dependent)
	Kannada CodeTable = 72
	// Malayalam represents Malayalam code table (model-dependent)
	Malayalam CodeTable = 73
	// Gujarati represents Gujarati code table (model-dependent)
	Gujarati CodeTable = 74
	// Punjabi represents Punjabi code table (model-dependent)
	Punjabi CodeTable = 75
	// Marathi represents Marathi code table (model-dependent)
	Marathi CodeTable = 82

	// Special254 represents special code table 254 (model-dependent)
	Special254 CodeTable = 254
	// Special255 represents special code table 255 (model-dependent)
	Special255 CodeTable = 255
)

// UpsideDownMode represents the upside-down printing mode
type UpsideDownMode byte

const (
	// OffUdm represents upside-down mode off (LSB = 0)
	OffUdm UpsideDownMode = 0x00
	// OnUdm represents upside-down mode on (LSB = 1)
	OnUdm UpsideDownMode = 0x01
)

// Size represents the character size configuration
type Size byte

const (
	// Size1x1 represents normal character size (width=1 height=1)
	Size1x1 Size = 0x00
	// Size2x1 represents double width character size (width=2 height=1)
	Size2x1 Size = 0x10
	// Size3x1 represents triple width character size (width=3 height=1)
	Size3x1 Size = 0x20
	// Size4x1 represents quadruple width character size (width=4 height=1)
	Size4x1 Size = 0x30
	// Size1x2 represents double height character size (width=1 height=2)
	Size1x2 Size = 0x01
	// Size2x2 represents double width and height character size (width=2 height=2)
	Size2x2 Size = 0x11

	// HeightMask is used to extract height bits from n
	HeightMask Size = 0x07
	// WidthMask is used to extract width bits from n
	WidthMask Size = 0x70
	// WidthShift is used to shift width bits to LSB position
	WidthShift Size = 4
)

// ReverseMode represents the white/black reverse printing mode
type ReverseMode byte

const (
	// OffRm represents white/black reverse mode off (LSB = 0)
	OffRm ReverseMode = 0x00
	// OnRm represents white/black reverse mode on (LSB = 1)
	OnRm ReverseMode = 0x01
	// OffRmASCII represents white/black reverse mode off (ASCII)
	OffRmASCII ReverseMode = '0'
	// OnRmASCII represents white/black reverse mode on (ASCII)
	OnRmASCII ReverseMode = '1'
)

// SmoothingMode represents the character smoothing mode
type SmoothingMode byte

const (
	// OffSm represents smoothing mode off (LSB = 0)
	OffSm SmoothingMode = 0x00
	// OnSm represents smoothing mode on (LSB = 1)
	OnSm SmoothingMode = 0x01
	// OffSmASCII represents smoothing mode off (ASCII)
	OffSmASCII SmoothingMode = '0'
	// OnSmASCII represents smoothing mode on (ASCII)
	OnSmASCII SmoothingMode = '1'
)

// ============================================================================
// Error Definitions
// ============================================================================

// ErrUnderlineMode indicates an invalid underline mode
var (
	ErrUnderlineMode   = fmt.Errorf("invalid underline mode(try 0-2 or '0'..'2')")
	ErrRotationMode    = fmt.Errorf("invalid rotation mode(try 0-2 or '0'..'2')")
	ErrPrintColor      = fmt.Errorf("invalid print color(try 0-1 or '0'..'1')")
	ErrCharacterFont   = fmt.Errorf("invalid character font(try 0-4 or '0'..'4')")
	ErrCodeTablePage   = fmt.Errorf("invalid code table page(try check model support!)")
	ErrCharacterSet    = fmt.Errorf("invalid international character set(try check model support!)")
	ErrCharacterWidth  = fmt.Errorf("invalid character width(try 1-8)")
	ErrCharacterHeight = fmt.Errorf("invalid character height(try 1-8)")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Compile-time check that Commands implements Capability
var _ Capability = (*Commands)(nil)

// Capability defines the main interface for character-related commands
type Capability interface {
	SetRightSideCharacterSpacing(spacing Spacing) []byte
	SelectPrintModes(modeBits PrintMode) []byte
	SetUnderlineMode(thickness UnderlineMode) ([]byte, error)
	SetEmphasizedMode(mode EmphasizedMode) []byte
	SetDoubleStrikeMode(mode DoubleStrikeMode) []byte
	SelectCharacterFont(fontType FontType) ([]byte, error)
	SelectInternationalCharacterSet(charset InternationalSet) ([]byte, error)
	Set90DegreeClockwiseRotationMode(rotationMode RotationMode) ([]byte, error)
	SelectPrintColor(color PrintColor) ([]byte, error)
	SelectCharacterCodeTable(page CodeTable) ([]byte, error)
	SetUpsideDownMode(mode UpsideDownMode) []byte
	SelectCharacterSize(sizeConfig Size) []byte
	SetWhiteBlackReverseMode(mode ReverseMode) []byte
	SetSmoothingMode(mode SmoothingMode) []byte
	EffectsCmds() EffectsCapability
	CodeConversionCmds() CodeConversionCapability
	UserDefinedCmds() UserDefinedCapability
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands groups all character-related capabilities
type Commands struct {
	Effects        EffectsCapability
	CodeConversion CodeConversionCapability
	UserDefined    UserDefinedCapability
}

// NewCommands creates a new Commands instance with initialized sub-commands
func NewCommands() *Commands {
	// Constructor that initializes sub-commands; returns pointer to avoid
	// copies when passing the structure.
	return &Commands{
		Effects:        &EffectsCommands{},
		CodeConversion: &CodeConversionCommands{},
		UserDefined:    &UserDefinedCommands{},
	}
}

// EffectsCmds returns the effects capability.
func (c *Commands) EffectsCmds() EffectsCapability {
	return c.Effects
}

// CodeConversionCmds returns the code conversion capability.
func (c *Commands) CodeConversionCmds() CodeConversionCapability {
	return c.CodeConversion
}

// UserDefinedCmds returns the user-defined character capability.
func (c *Commands) UserDefinedCmds() UserDefinedCapability {
	return c.UserDefined
}

// ============================================================================
// Helper Functions
// ============================================================================

// TODO: Check if is better to have them defined in composer package

// NewSize creates a Size value for character width and height magnification.
//
// Format:
//
//	Not applicable (helper function)
//
// Range:
//
//	width = 1–8
//	height = 1–8
//
// Default:
//
//	Not applicable
//
// Parameters:
//
//	width: Character width magnification (1-8)
//	height: Character height magnification (1-8)
//
// Notes:
//   - This is a helper function to build Size values for SelectCharacterSize
//   - The returned Size encodes both width and height in a single byte
//
// Errors:
//
//	Returns ErrCharacterWidth if width is outside the range 1-8.
//	Returns ErrCharacterHeight if height is outside the range 1-8.
func NewSize(width, height byte) (Size, error) {
	if width < 1 || width > 8 {
		return 0, ErrCharacterWidth
	}
	if height < 1 || height > 8 {
		return 0, ErrCharacterHeight
	}
	w := (width - 1) << 4
	h := height - 1
	return Size(w | h), nil
}

// ============================================================================
// Validation Functions
// ============================================================================

// ValidateUnderlineMode validates if underline mode is valid
func ValidateUnderlineMode(mode UnderlineMode) error {
	switch mode {
	case NoDot, OneDot, TwoDot, NoDotASCII, OneDotASCII, TwoDotASCII:
		return nil
	default:
		return ErrUnderlineMode
	}
}

// ValidateFontType validates if font type is valid
func ValidateFontType(font FontType) error {
	switch font {
	case FontA, FontB, FontC, FontD, FontE,
		FontAAscii, FontBAscii, FontCAscii, FontDAscii, FontEAscii,
		SpecialFontA, SpecialFontB:
		return nil
	default:
		return ErrCharacterFont
	}
}

// ValidateInternationalSet validates if international character set is valid
func ValidateInternationalSet(charset InternationalSet) error {
	// Standard range
	if charset <= Arabia {
		return nil
	}
	// India-specific range
	if charset >= IndiaDevanagari && charset <= IndiaPunjabi {
		return nil
	}
	// Marathi
	if charset == IndiaMarathi {
		return nil
	}
	return ErrCharacterSet
}

// ValidateRotationMode validates if rotation mode is valid
func ValidateRotationMode(mode RotationMode) error {
	switch mode {
	case NoRotation, On90Dot1, On90Dot15,
		NoRotationASCII, On90Dot1Ascii, On90Dot15Ascii:
		return nil
	default:
		return ErrRotationMode
	}
}

// ValidatePrintColor validates if print color is valid
func ValidatePrintColor(color PrintColor) error {
	switch color {
	case Black, Red, BlackASCII, RedASCII:
		return nil
	default:
		return ErrPrintColor
	}
}

// ValidateCodeTable validates if code table page is valid
func ValidateCodeTable(page CodeTable) error {
	// Common pages 0-8
	if page <= OnePassKanji2 {
		return nil
	}
	// Pages 11-19
	if page >= PC851 && page <= PC858 {
		return nil
	}
	// Thai codes 20-26
	if page >= ThaiCode42 && page <= ThaiCode18 {
		return nil
	}
	// Vietnamese 30-31
	if page >= TCVN31 && page <= TCVN32 {
		return nil
	}
	// Pages 32-53
	if page >= PC720 && page <= KZ1048 {
		return nil
	}
	// India-specific pages 66-75
	if page >= Devanagari && page <= Punjabi {
		return nil
	}
	// Marathi
	if page == Marathi {
		return nil
	}
	// Special pages
	if page == Special254 || page == Special255 {
		return nil
	}
	return ErrCodeTablePage
}
