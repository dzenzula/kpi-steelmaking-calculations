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
	ConsumptionMnlz       string `yaml:"consumption_mnlz"`
	ConsumptionIngot      string `yaml:"consumption_ingot"`
	MeltdownsCasting      string `yaml:"meltdowns_casting"`
	MeltdownsOnrs         string `yaml:"meltdowns_onrs"`
	ScrapConsumptionMnlz  string `yaml:"scrap_consumption_mnlz"`
	ScrapConsumptionIngot string `yaml:"scrap_consumption_ingot"`
	SiInCastIron          string `yaml:"si_in_cast_iron"`
	CastIronTemperature   string `yaml:"cast_iron_temperature"`
	SContentPercentage    string `yaml:"s_content_percentage"`
	MnlzProduction        string `yaml:"mnlz_production"`
	IngotProduction       string `yaml:"ingot_production"`
	O2ContentAtTheOutlet  string `yaml:"o2_content_at_the_outlet"`
	GetLime               string `yaml:"get_lime"`
	GetDolomite           string `yaml:"get_dolomite"`
	GetAluminum           string `yaml:"get_aluminum"`
	GetMixture            string `yaml:"get_mixture"`
}

type Config struct {
	ConStringMsDev     ConStringMS `yaml:"mssql_data_bof"`
	ConStringPgRawData ConStringPG `yaml:"postgres_raw_data"`
	ConStringPgDev     ConStringPG `yaml:"postgres_dev"`
	Querries           Querries    `yaml:"querries"`
	TypeMS             string      `yaml:"type_ms"`
	TypePG             string      `yaml:"type_pg"`
}
