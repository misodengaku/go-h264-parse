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
	VUI
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
	PicInitQPMinus26                      int64
	PicInitQSMinus26                      int64
	ChromaQPIndexOffset                   int64
	DeblockingFilterControlPresentFlag    bool
	ConstrainedIntraPredFlag              bool
	RedundantPicCntPresentFlag            bool
	Transform8x8ModeFlag                  bool
	PicScalingMatrixPresentFlag           bool
	SeqScalingListPresentFlags            []bool
	SecondChromaQPIndexOffset             int64
}

type SEI struct {
	PayloadType  int
	PayloadSize  int
	PayloadBytes []byte
}

// Rec. ITU-T H.264 (08/2021) pp.43-45
type SPS struct {
	ProfileIDC                      byte
	ConstraintSet0Flag              bool
	ConstraintSet1Flag              bool
	ConstraintSet2Flag              bool
	ConstraintSet3Flag              bool
	ConstraintSet4Flag              bool
	ConstraintSet5Flag              bool
	LevelIDC                        byte
	ID                              uint64
	ChromaFormatIDC                 uint64
	SeparateColourPlaneFlag         bool
	BitDepthLumaMinus8              uint64
	BitDepthChromaMinus8            uint64
	QPPrimeYZeroTransformBypassFlag bool
	SeqScalingMatrixPresentFlag     bool
	Log2MaxFrameNumMinus4           uint64
	PicOrderCntType                 uint64
	Log2MaxPicOrderCntLSBMinus4     uint64
	DeltaPicOrderAlwaysZeroFlag     bool
	OffsetForNonRefPic              int64
	OffsetForTopToBottomField       int64
	NumRefFramesInPicOrderCntCycle  uint64
	OffsetForRefFrames              []int64
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

type VUI struct {
	// Rec. ITU-T H.264 (08/2021) p.422-423
	AspectRatioInfoPresentFlag         bool
	AspectRatioIDC                     byte
	SARWidth                           uint16
	SARHeight                          uint16
	OverscanInfoPresentFlag            bool
	OverscanAppropriateFlag            bool
	VideoSignalTypePresentFlag         bool
	VideoFormat                        byte
	VideoFullRangeFlag                 bool
	ColourDescriptionPresentFlag       bool
	ColourPrimaries                    byte
	TransferCharacteristics            byte
	MatrixCoefficients                 byte
	ChromaLocInfoPresentFlag           bool
	ChromaSampleLocTypeTopField        uint64
	ChromaSampleLocTypeBottomField     uint64
	TimingInfoPresentFlag              bool
	NumUnitsInTick                     uint32
	TimeScale                          uint32
	FixedFrameRateFlag                 bool
	NALHRDParametersPresentFlag        bool
	NALHRDParameters                   HRDParameters
	VCLHRDParametersPresentFlag        bool
	VCLHRDParameters                   HRDParameters
	LowDelayHRDFlag                    bool
	PicStructPresentFlag               bool
	BitstreamRestrictionFlag           bool
	MotionVectorsOverPicBoundariesFlag bool
	MaxBytesPerPicDenom                uint64
	MaxBitsPerMBDenom                  uint64
	Log2MaxMVLengthHorizontal          uint64
	Log2MaxMVLengthVertical            uint64
	MaxNumReorderFrames                uint64
	MaxDecFrameBuffering               uint64
}

type HRDParameters struct {
	CPBCntMinus1                       uint64
	BitRateScale                       byte
	CPBSizeScale                       byte
	AlternativeCPBSpecifications       []AlternativeCPBSpecification
	InitialCPBRemovalDelayLengthMinus1 byte
	CPBRemovalDelayLengthMinus1        byte
	DPBOutputDelayLengthMinus1         byte
	TimeOffsetLength                   byte
}

type AlternativeCPBSpecification struct {
	BitRateValueMinus1 uint64
	CPBSizeValueMinus1 uint64
	CBRFlag            bool
}

const (
	AspectRatioIDC_Unspecified = 0
	AspectRatioIDC_1_1         = 1
	AspectRatioIDC_12_11       = 2
	AspectRatioIDC_10_11       = 3
	AspectRatioIDC_16_11       = 4
	AspectRatioIDC_40_33       = 5
	AspectRatioIDC_24_11       = 6
	AspectRatioIDC_20_11       = 7
	AspectRatioIDC_32_11       = 8
	AspectRatioIDC_80_33       = 9
	AspectRatioIDC_18_11       = 10
	AspectRatioIDC_15_11       = 11
	AspectRatioIDC_64_33       = 12
	AspectRatioIDC_160_99      = 13
	AspectRatioIDC_4_3         = 14
	AspectRatioIDC_3_2         = 15
	AspectRatioIDC_2_1         = 16
	AspectRatioIDC_ExtendedSAR = 255
)

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
