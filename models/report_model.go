package models

type Report struct {
	//Id                   int     `stbl:"Id, PRIMARY_KEY, AUTO_INCREMENT"`
	Date                   string  `stbl:"Date"`
	CastIronMelting        float64 `stbl:"CastIronMelting"`
	ScrapMelting           float64 `stbl:"ScrapMelting"`
	SiInCastIron           float64 `stbl:"SiInCastIron"`
	CastIronTemperature    float64 `stbl:"CastIronTemperature"`
	SContent               float64 `stbl:"SContent"`
	MNLZMelting            float64 `stbl:"MNLZMelting"`
	IngotMelting           float64 `stbl:"IngotMelting"`
	O2Content              float64 `stbl:"O2Content"`
	LimestoneFlow          float64 `stbl:"LimestoneFlow"`
	DolomiteFlow           float64 `stbl:"DolomiteFlow"`
	AluminumPreheating     float64 `stbl:"AluminumPreheating"`
	MixMelting             float64 `stbl:"MixMelting"`
	SiCC                   float64 `stbl:"SiCC"`
	SiModel                float64 `stbl:"SiModel"`
	SiMnCC                 float64 `stbl:"SiMnCC"`
	SiMnModel              float64 `stbl:"SiMnModel"`
	MnCC                   float64 `stbl:"MnCC"`
	MnModel                float64 `stbl:"MnModel"`
	SlagTruncationRatio    float64 `stbl:"SlagTruncationRatio"`
	SlagSkimmingRatio      float64 `stbl:"SlagSkimmingRatio"`
	CCMeltingCycle         float64 `stbl:"CCMeltingCycle"`
	FePercentageInSlag     float64 `stbl:"FePercentageInSlag"`
	SlagSamplingPercentage float64 `stbl:"SlagSamplingPercentage"`
	GoodCCOutput           float64 `stbl:"GoodCCOutput"`
	GoodCCMNLZOutput       float64 `stbl:"GoodCCMNLZOutput"`
	GoodIngotOutput        float64 `stbl:"GoodIngotOutput"`
	ProcessingTime         float64 `stbl:"ProcessingTime"`
	ArcTime                float64 `stbl:"ArcTime"`
	LimestoneConsumption   float64 `stbl:"LimestoneConsumption"`
	FluorsparConsumption   float64 `stbl:"FluorsparConsumption"`
	ArgonOxygenConsumption float64 `stbl:"ArgonOxygenConsumption"`
	ElectricityConsumption float64 `stbl:"ElectricityConsumption"`
	ElectrodeConsumption   float64 `stbl:"ElectrodeConsumption"`
	InletTemperature       float64 `stbl:"InletTemperature"`
	InletOxidation         float64 `stbl:"InletOxidation"`
	UPKSlagAnalysis        float64 `stbl:"UPKSlagAnalysis"`
	CastingCycle           float64 `stbl:"CastingCycle"`
	CastingSpeed           float64 `stbl:"CastingSpeed"`
	CastingStopperSerial   float64 `stbl:"CastingStopperSerial"`
	MNLZ1Streams           float64 `stbl:"MNLZ1Streams"`
	MNLZ2Streams           float64 `stbl:"MNLZ2Streams"`
	MNLZ3Streams           float64 `stbl:"MNLZ3Streams"`
	MNLZ1RepackingDuration float64 `stbl:"MNLZ1RepackingDuration"`
	MNLZ2RepackingDuration float64 `stbl:"MNLZ2RepackingDuration"`
	MNLZ3RepackingDuration float64 `stbl:"MNLZ3RepackingDuration"`
	GoodMNLZOutput         float64 `stbl:"GoodMNLZOutput"`
	MetalRetentionTime     float64 `stbl:"MetalRetentionTime"`
}
