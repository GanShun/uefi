package uefi

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	// FlashDescriptorMapSize is the size in bytes of the Intel flash descriptor
	FlashDescriptorMapSize = 0x1000
	// FlashDescriptorMapMaxBase is the maximum base address for a flash descriptor
	// region
	FlashDescriptorMapMaxBase = 0xe0
)

// FlashDescriptorMap represent an Intel flash descriptor. This object provides
// accessors to the various descriptor fields.
type FlashDescriptorMap struct {
	// FLMAP0
	ComponentBase      uint8
	NumberOfFlashChips uint8
	RegionBase         uint8
	NumberOfRegions    uint8
	// FLMAP1
	MasterBase        uint8
	NumberOfMasters   uint8
	PchStrapsBase     uint8
	NumberOfPchStraps uint8
	// FLMAP2
	ProcStrapsBase          uint8
	NumberOfProcStraps      uint8
	IccTableBase            uint8
	NumberOfIccTableEntries uint8
	// FLMAP3
	DmiTableBase            uint8
	NumberOfDmiTableEntries uint8
	Reserved0               uint8
	Reserved1               uint8
}

// NewFlashDescriptorMap initializes a FlashDescriptor from a slice of bytes.
func NewFlashDescriptorMap(buf []byte) (*FlashDescriptorMap, error) {
	if len(buf) < FlashDescriptorMapSize {
		return nil, fmt.Errorf("Flash Descriptor Map size too small: expected %v bytes, got %v",
			FlashDescriptorMapSize,
			len(buf),
		)
	}
	r := bytes.NewReader(buf)
	var descriptor FlashDescriptorMap
	if err := binary.Read(r, binary.LittleEndian, &descriptor); err != nil {
		return nil, err
	}
	return &descriptor, nil
}

func (d FlashDescriptorMap) String() string {
	return fmt.Sprintf("FlashDescriptorMap{NumberOfRegions=%v, NumberOfFlashChips=%v, NumberOfMasters=%v, NumberOfPCHStraps=%v, NumberOfProcessorStraps=%v, NumberOfICCTableEntries=%v, DMITableEntries=%v}",
		d.NumberOfRegions,
		d.NumberOfFlashChips,
		d.NumberOfMasters,
		d.NumberOfPchStraps,
		d.NumberOfProcStraps,
		d.NumberOfIccTableEntries,
		d.NumberOfDmiTableEntries,
	)
}

// Summary prints a multi-line description of the flash descriptor map
func (d FlashDescriptorMap) Summary() string {
	return fmt.Sprintf("FlashDescriptorMap{\n"+
		"    ComponentBase=%v (0x%02x)\n"+
		"    NumberOfFlashChips=%v (0x%02x)\n"+
		"    RegionBase=%v (0x%02x)\n"+
		"    NumberOfRegions=%v (0x%02x)\n"+
		"    MasterBase=%v (0x%02x)\n"+
		"    NumberOfMasters=%v (0x%02x)\n"+
		"    PCHStrapsBase=%v (0x%02x)\n"+
		"    NumberOfPCHStraps=%v (0x%02x)\n"+
		"    ProcessorStrapsBase=%v (0x%02x)\n"+
		"    NumberOfProcessorStraps=%v (0x%02x)\n"+
		"    ICCTableEntriesBase=%v (0x%02x)\n"+
		"    NumberOfICCTableEntries=%v (0x%02x)\n"+
		"    DMITableEntriesBase=%v (0x%02x)\n"+
		"    NumberOfDMITableEntries=%v (0x%02x)\n"+
		"}",
		d.ComponentBase, d.ComponentBase,
		d.NumberOfFlashChips, d.NumberOfFlashChips,
		d.RegionBase, d.RegionBase,
		d.NumberOfRegions, d.NumberOfRegions,
		d.MasterBase, d.MasterBase,
		d.NumberOfMasters, d.NumberOfMasters,
		d.PchStrapsBase, d.PchStrapsBase,
		d.NumberOfPchStraps, d.NumberOfPchStraps,
		d.ProcStrapsBase, d.ProcStrapsBase,
		d.NumberOfProcStraps, d.NumberOfProcStraps,
		d.IccTableBase, d.IccTableBase,
		d.NumberOfIccTableEntries, d.NumberOfIccTableEntries,
		d.DmiTableBase, d.DmiTableBase,
		d.NumberOfDmiTableEntries, d.NumberOfDmiTableEntries,
	)
}

// Validate runs a set of checks on the flash descriptor and returns a list of
// errors specifying what is wrong.
func (d FlashDescriptorMap) Validate() []error {
	errors := make([]error, 0)
	if d.MasterBase > FlashDescriptorMapMaxBase {
		errors = append(errors, fmt.Errorf("MasterBase too large: expected %v bytes, got %v",
			FlashDescriptorMapMaxBase,
			d.MasterBase,
		))
	}
	if d.RegionBase > FlashDescriptorMapMaxBase {
		errors = append(errors, fmt.Errorf("RegionBase too large: expected %v bytes, got %v",
			FlashDescriptorMapMaxBase,
			d.RegionBase,
		))
	}
	if d.MasterBase > FlashDescriptorMapMaxBase {
		errors = append(errors, fmt.Errorf("ComponentBase too large: expected %v bytes, got %v",
			FlashDescriptorMapMaxBase,
			d.MasterBase,
		))
	}
	if d.MasterBase == d.RegionBase {
		errors = append(errors, fmt.Errorf("MasterBase must be different from RegionBase: both are at 0x%x",
			d.MasterBase,
		))
	}
	if d.MasterBase == d.ComponentBase {
		errors = append(errors, fmt.Errorf("MasterBase must be different from ComponentBase: both are at 0x%x",
			d.MasterBase,
		))
	}
	if d.RegionBase == d.ComponentBase {
		errors = append(errors, fmt.Errorf("RegionBase must be different from ComponentBase: both are at 0x%x",
			d.RegionBase,
		))
	}
	return errors
}
