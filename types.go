package h264parse

type NALUs struct {
	Units []NAL
}

// Rec. ITU-T H.264 (08/2021) p.43
type NAL struct {
	RefIDC      byte
	UnitType    byte
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

const (
	// Rec. ITU-T H.264 (08/2021) p.65
	Unspecified0                                            = 0  //	Unspecified
	CodedSliceNonIDRPicture                                 = 1  //	Coded slice of a non-IDR picture
	CodedSliceDataPartitionA                                = 2  //	Coded slice data partition A
	CodedSliceDataPartitionB                                = 3  //	Coded slice data partition B
	CodedSliceDataPartitionC                                = 4  //	Coded slice data partition C
	CodedSliceIDRPicture                                    = 5  //	Coded slice of an IDR picture
	SupplementalEnhancementInformation                      = 6  //	Supplemental enhancement information (SEI)
	SequenceParameterSet                                    = 7  //	Sequence parameter set
	PictureParameterSet                                     = 8  //	Picture parameter set
	AccessUnitDelimiter                                     = 9  //	Access unit delimiter
	EndOfSequence                                           = 10 //	End of sequence
	EndOfStream                                             = 11 //	End of stream
	FillerData                                              = 12 //	Filler data
	SequenceParameterSetExtension                           = 13 //	Sequence parameter set extension
	PrefixNALUnit                                           = 14 //	Prefix NAL unit
	SubsetSequenceParameterSet                              = 15 //	Subset sequence parameter set
	DepthParameterSet                                       = 16 //	Depth parameter set
	Reserved17                                              = 17 //	Reserved
	Reserved18                                              = 18 //	Reserved
	CodedSliceAuxiliaryCodedPictureWithoutPartitioning      = 19 //	Coded slice of an auxiliary coded  picture without partitioning
	CodedSliceExtension                                     = 20 //	Coded slice extension
	CodedSliceExtensionDepthViewComponentOr3DAVCTextureView = 21 //	Coded slice extension for a depth view component or a 3D-AVC texture view component
	Reserved22                                              = 22 //	Reserved
	Reserved23                                              = 23 //	Reserved
	Unspecified24                                           = 24 //	Unspecified
	Unspecified25                                           = 25 //	Unspecified
	Unspecified26                                           = 26 //	Unspecified
	Unspecified27                                           = 27 //	Unspecified
	Unspecified28                                           = 28 //	Unspecified
	Unspecified29                                           = 29 //	Unspecified
	Unspecified30                                           = 30 //	Unspecified
	Unspecified31                                           = 31 //	Unspecified

)
