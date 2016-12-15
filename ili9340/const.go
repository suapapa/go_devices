// Copyright 2016, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ili9340

// ili9340 lcd should be connected in following pins
const (
	PinDC  = "DC"
	PinRST = "RST"
)

const (
	width  = 240
	height = 320
)

const (
	regNOP     = 0x00
	regSWRESET = 0x01
	regRDDID   = 0x04
	regRDDST   = 0x09

	regSLPIN  = 0x10
	regSLPOUT = 0x11
	regPTLON  = 0x12
	regNORON  = 0x13

	regRDMODE     = 0x0A
	regRDMADCTL   = 0x0B
	regRDPIXFMT   = 0x0C
	regRDIMGFMT   = 0x0A
	regRDSELFDIAG = 0x0F

	regINVOFF   = 0x20
	regINVON    = 0x21
	regGAMMASET = 0x26
	regDISPOFF  = 0x28
	regDISPON   = 0x29

	regCASET = 0x2A
	regPASET = 0x2B
	regRAMWR = 0x2C
	regRAMRD = 0x2E

	regPTLAR  = 0x30
	regMADCTL = 0x36

	regMADCTLvalMY  = 0x80
	regMADCTLvalMX  = 0x40
	regMADCTLvalMV  = 0x20
	regMADCTLvalML  = 0x10
	regMADCTLvalRGB = 0x00
	regMADCTLvalBGR = 0x08
	regMADCTLvalMH  = 0x04

	regPIXFMT = 0x3A

	regFRMCTR1 = 0xB1
	regFRMCTR2 = 0xB2
	regFRMCTR3 = 0xB3
	regINVCTR  = 0xB4
	regDFUNCTR = 0xB6

	regPWCTR1 = 0xC0
	regPWCTR2 = 0xC1
	regPWCTR3 = 0xC2
	regPWCTR4 = 0xC3
	regPWCTR5 = 0xC4
	regVMCTR1 = 0xC5
	regVMCTR2 = 0xC7

	regRDID1 = 0xDA
	regRDID2 = 0xDB
	regRDID3 = 0xDC
	regRDID4 = 0xDD

	regGMCTRP1 = 0xE0
	regGMCTRN1 = 0xE1
)
