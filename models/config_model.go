package models

type ConStringMS struct {
	Server   string `yaml:"server"`
	UserID   string `yaml:"user_id"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type ConStringPG struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type Querries struct {
	ConsumptionMnlz              string `yaml:"consumption_mnlz"`
	ConsumptionIngot             string `yaml:"consumption_ingot"`
	MeltdownsCasting             string `yaml:"meltdowns_casting"`
	MeltdownsOnrs                string `yaml:"meltdowns_onrs"`
	ScrapConsumptionMnlz         string `yaml:"scrap_consumption_mnlz"`
	ScrapConsumptionIngot        string `yaml:"scrap_consumption_ingot"`
	GetData                      string `yaml:"get_data"`
	GetSlagTruncation            string `yaml:"get_slag_truncation"`
	GetGreaterThanZeroData       string `yaml:"get_data_greater_than_zero"`
	GetFerroalloysForMNLZ        string `yaml:"get_pherroalloy_mnlz"`
	GetFerroalloysForIngots      string `yaml:"get_pherroalloy_ingot"`
	GetUpk                       string `yaml:"get_upk"`
	GetGoodCastIron              string `yaml:"get_good_cast_iron"`
	GetInletTemperatureOxidation string `yaml:"get_inlet_temperature_oxidation"`
	GetMnlz                      string `yaml:"get_mnlz"`
	GetMineralData               string `yaml:"get_mineral_data"`
}

type Measurings struct {
	FlHeatLf         int `yaml:"fl_heat_lf"`
	WHeat            int `yaml:"w_heat"`
	WIron            int `yaml:"w_iron"`
	WLom             int `yaml:"w_lom"`
	WHeatMnlz        int `yaml:"w_heat_mnlz"`
	SiIron           int `yaml:"si_iron"`
	TIron            int `yaml:"t_iron"`
	PIron            int `yaml:"p_iron"`
	OxiBof           int `yaml:"oxi_bof"`
	WCaov            int `yaml:"w_caov"`
	WDolo            int `yaml:"w_dolo"`
	WAlrz            int `yaml:"w_alrz"`
	WSmesrz          int `yaml:"w_smesrz"`
	SIron            int `yaml:"s_iron"`
	WFesi            int `yaml:"w_fesi"`
	Fesi65           int `yaml:"fesi65"`
	WSimn            int `yaml:"w_simn"`
	Simn             int `yaml:"simn"`
	WFemn78          int `yaml:"w_femn78"`
	WFemn88          int `yaml:"w_femn88"`
	WFemn95          int `yaml:"w_femn95"`
	Femn78           int `yaml:"femn78"`
	Femn88           int `yaml:"femn88"`
	FlHeatCutoffSlag int `yaml:"fl_heat_cutoff_slag"`
	FlSlop           int `yaml:"fl_slop"`
	ShlMgo           int `yaml:"shl_mgo"`
	WFecr            int `yaml:"w_fecr"`
	WFecr025         int `yaml:"w_fecr025"`
	WFecr800         int `yaml:"w_fecr800"`
	TmTreatment      int `yaml:"tm_treatment"`
	NLf              int `yaml:"n_lf"`
	M1e              int `yaml:"m1e"`
	M2e              int `yaml:"m2e"`
	M3e              int `yaml:"m3e"`
	M4e              int `yaml:"m4e"`
	M5e              int `yaml:"m5e"`
	M6e              int `yaml:"m6e"`
	M7e              int `yaml:"m7e"`
	M8e              int `yaml:"m8e"`
	WCaoLf           int `yaml:"w_cao_lf"`
	WCao90Lf         int `yaml:"w_cao90_lf"`
	WCaf2Lf          int `yaml:"w_caf2_lf"`
	WApkLf           int `yaml:"w_apk_lf"`
	T1               int `yaml:"t1"`
	Treatment        int `yaml:"treatment"`
	O21              int `yaml:"o2_1"`
	AvgSpeed         int `yaml:"avgspeed"`
	NCasting         int `yaml:"ncasting"`
	CntHeat          int `yaml:"cnt_heat"`
	PCount           int `yaml:"pcount"`
	WeightLs         int `yaml:"weight_ls"`
	WeightGrossStart int `yaml:"weight_grossstart"`
	WeightGrossEnd   int `yaml:"weight_grossend"`
	TmBetween        int `yaml:"tm_between"`
	CCMBegin         int `yaml:"ccm_begin"`
	Dtk              int `yaml:"dtk"`
	Dtn              int `yaml:"dtn"`
	SampleTime       int `yaml:"sample_time"`
	HeatStart        int `yaml:"heat_start"`
	DtBegin          int `yaml:"dt_begin"`
	DtEnd            int `yaml:"dt_end"`
	S1Oil            int `yaml:"s1_oil"`
	S2Oil            int `yaml:"s2_oil"`
	S3Oil            int `yaml:"s3_oil"`
	S4Oil            int `yaml:"s4_oil"`
	S5Oil            int `yaml:"s5_oil"`
	S6Oil            int `yaml:"s6_oil"`
	DkTap            int `yaml:"dk_tap"`
	EndSeries        int `yaml:"end_series"`
	Serialization    int `yaml:"serialization"`
	Nccm             int `yaml:"n_ccm"`
	Energy1          int `yaml:"energy1"`
	Energy2          int `yaml:"energy2"`
	Energy3          int `yaml:"energy3"`
	Energy4          int `yaml:"energy4"`
	Energy5          int `yaml:"energy5"`
	Energy6          int `yaml:"energy6"`
	Energy7          int `yaml:"energy7"`
	Energy8          int `yaml:"energy8"`
	Temperature1     int `yaml:"temperature1"`
	Temperature2     int `yaml:"temperature2"`
	Temperature3     int `yaml:"temperature3"`
	Temperature4     int `yaml:"temperature4"`
	Temperature5     int `yaml:"temperature5"`
	Temperature6     int `yaml:"temperature6"`
	Temperature7     int `yaml:"temperature7"`
	Temperature8     int `yaml:"temperature8"`
	SteelGrade       int `yaml:"steel_grade"`
}

type SteelMark struct {
	SteelType string `yaml:"steel_type"`
	Min       int    `yaml:"min"`
	Max       int    `yaml:"max"`
}

type Config struct {
	ConStringPgDb      ConStringPG `yaml:"postgres_data"`
	Querries           Querries    `yaml:"querries"`
	Measurings         Measurings  `yaml:"measurings"`
	CachePath          string      `yaml:"cache"`
	LoggerLevel        string      `yaml:"logger"`
	TypeMS             string      `yaml:"type_ms"`
	TypePG             string      `yaml:"type_pg"`
	SteelMarks         []SteelMark `yaml:"steel_temperatures"`
}
