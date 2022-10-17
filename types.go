package h264parse

//go:generate stringer -type=NALUnitType

type NALUs struct {
	Units []NAL
}

// Rec. ITU-T H.264 (08/2021) p.43
type NAL struct {
	RefIDC      byte
	UnitType    NALUnitType
	RBSPByte    []byte
	HeaderBytes []byte

	SPS
	PPS
	SEI
}

// Rec. ITU-T H.264 (08/2021) pp.47-48
type PPS struct {
	ID                                    uint64
	SPS_ID                                uint64
	EntropyCodingModeFlag                 bool
	BottomFieldPicOrderInFramePresentFlag bool
	NumSliceGroupsMinus1                  uint64
	SliceGroupMapType                     uint64
	RunLengthMinus1                       []uint64
	TopLeft                               []uint64
	BottomRight                           []uint64
	NumRefIndexL0DefaultActiveMinus1      uint64
	NumRefIndexL1DefaultActiveMinus1      uint64
	WeightedPredFlag                      bool
	WeightedBipredIdc                     byte
	PicInitQPMinus26                      uint64 // TODO: signed
	PicInitQSMinus26                      uint64 // TODO: signed
	ChromaQPIndexOffset                   uint64 // TODO: signed
	DeblockingFilterControlPresentFlag    bool
	ConstrainedIntraPredFlag              bool
	RedundantPicCntPresentFlag            bool
	Transform8x8ModeFlag                  bool
	PicScalingMatrixPresentFlag           bool
	SeqScalingListPresentFlags            []bool
	SecondChromaQPIndexOffset             uint64 // TODO: signed
}

type SEI struct {
	PayloadType  int
	PayloadSize  int
	PayloadBytes []byte
}

// Rec. ITU-T H.264 (08/2021) pp.43-45
type SPS struct {
	ProfileIDC                      byte
	ConstraintSet0Flag              byte
	ConstraintSet1Flag              byte
	ConstraintSet2Flag              byte
	ConstraintSet3Flag              byte
	ConstraintSet4Flag              byte
	ConstraintSet5Flag              byte
	LevelIDC                        byte
	ID                              uint64
	ChromaFormatIDC                 uint64
	BitDepthLumaMinus8              uint64
	BitDepthChromaMinus8            uint64
	QPPrimeYZeroTransformBypassFlag bool
	SeqScalingMatrixPresentFlag     bool
	Log2MaxFrameNumMinus4           uint64
	PicOrderCntType                 uint64
	MaxNumRefFrames                 uint64
	GapsInFrameNumValueAllowedFlag  bool
	PicWidthInMBSMinus1             uint64
	PicHeightInMapUnitsMinus1       uint64
	FrameMBSOnlyFlag                bool
	MBAdaptiveFrameFieldFlag        bool
	Direct8x8InferenceFlag          bool
	FrameCroppingFlag               bool
	FrameCropLeftOffset             uint64
	FrameCropRightOffset            uint64
	FrameCropTopOffset              uint64
	FrameCropBottomOffset           uint64
	VUIParametersPresentFlag        bool
}

type NALUnitType byte

const (
	// Rec. ITU-T H.264 (08/2021) p.65
	Unspecified0                                            = NALUnitType(0)  //	Unspecified
	CodedSliceNonIDRPicture                                 = NALUnitType(1)  //	Coded slice of a non-IDR picture
	CodedSliceDataPartitionA                                = NALUnitType(2)  //	Coded slice data partition A
	CodedSliceDataPartitionB                                = NALUnitType(3)  //	Coded slice data partition B
	CodedSliceDataPartitionC                                = NALUnitType(4)  //	Coded slice data partition C
	CodedSliceIDRPicture                                    = NALUnitType(5)  //	Coded slice of an IDR picture
	SupplementalEnhancementInformation                      = NALUnitType(6)  //	Supplemental enhancement information (SEI)
	SequenceParameterSet                                    = NALUnitType(7)  //	Sequence parameter set
	PictureParameterSet                                     = NALUnitType(8)  //	Picture parameter set
	AccessUnitDelimiter                                     = NALUnitType(9)  //	Access unit delimiter
	EndOfSequence                                           = NALUnitType(10) //	End of sequence
	EndOfStream                                             = NALUnitType(11) //	End of stream
	FillerData                                              = NALUnitType(12) //	Filler data
	SequenceParameterSetExtension                           = NALUnitType(13) //	Sequence parameter set extension
	PrefixNALUnit                                           = NALUnitType(14) //	Prefix NAL unit
	SubsetSequenceParameterSet                              = NALUnitType(15) //	Subset sequence parameter set
	DepthParameterSet                                       = NALUnitType(16) //	Depth parameter set
	Reserved17                                              = NALUnitType(17) //	Reserved
	Reserved18                                              = NALUnitType(18) //	Reserved
	CodedSliceAuxiliaryCodedPictureWithoutPartitioning      = NALUnitType(19) //	Coded slice of an auxiliary coded  picture without partitioning
	CodedSliceExtension                                     = NALUnitType(20) //	Coded slice extension
	CodedSliceExtensionDepthViewComponentOr3DAVCTextureView = NALUnitType(21) //	Coded slice extension for a depth view component or a 3D-AVC texture view component
	Reserved22                                              = NALUnitType(22) //	Reserved
	Reserved23                                              = NALUnitType(23) //	Reserved
	Unspecified24                                           = NALUnitType(24) //	Unspecified
	Unspecified25                                           = NALUnitType(25) //	Unspecified
	Unspecified26                                           = NALUnitType(26) //	Unspecified
	Unspecified27                                           = NALUnitType(27) //	Unspecified
	Unspecified28                                           = NALUnitType(28) //	Unspecified
	Unspecified29                                           = NALUnitType(29) //	Unspecified
	Unspecified30                                           = NALUnitType(30) //	Unspecified
	Unspecified31                                           = NALUnitType(31) //	Unspecified

)
