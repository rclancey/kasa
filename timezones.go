package kasa

import (
	"log"
	"time"
)

type Timezone int

var timezoneIndex = map[Timezone]string{
	0: "Etc/GMT+12", //"UTC-12:00 - International Date Line West",
	1: "Etc/GMT+11", //"UTC-11:00 - Coordinated Universal Time",
	2: "US/Hawaii", //"UTC-10:00 - Hawaii",
	3: "US/Alaska", //"UTC-09:00 - Alaska"
	4: "America/Tijuana", //"UTC-08:00 - Baji California",
	5: "America/Los_Angeles", //"UTC-08:00 - Pacific Standard Time",
	6: "America/Los_Angeles", //"UTC-08:00 - Pacific Daylight Time",
	7: "America/Phoenix", //"UTC-07:00 - Arizona",
	8: "America/Chihuahua", // "UTC-07:00 - Chihuahua, La Paz, Mazatlan",
	9: "America/Denver", //"UTC-07:00 - Mountain Standard Time",
	10: "America/Denver", //"UTC-07:00 - Mountain Daylight Time",
	11: "UTC-06:00 - Central America",
	12: "America/Chicago", //"UTC-06:00 - Central Standard Time",
	13: "America/Chicago", //"UTC-06:00 - Central Daylight Time",
	14: "America/Mexico_City", //"UTC-06:00 - Guadalajara, Mexico City",
	15: "Canada/Saskatchewan", //"UTC-06:00 - Saskatchewan",
	16: "America/Bogota", //"UTC-05:00 - Bogota, Lima, Quito",
	17: "America/New_York", //"UTC-05:00 - Eastern Standard Time",
	18: "America/New_York", //"UTC-05:00 - Eastern Daylight Time",
	19: "America/Indiana/Indianapolis", //"UTC-05:00 - Indiana (East)",
	20: "UTC-04:30 - Caracas",
	21: "UTC-04:00 - Asunicion",
	22: "Canada/Atlantic", //"UTC-04:00 - Atlantic Standard Time",
	23: "Canada/Atlantic", //"UTC-04:00 - Atlantic Daylight Time",
	24: "UTC-04:00 - Cuiaba",
	25: "UTC-04:00 - Georgetown",
	26: "UTC-04:00 - Santiago",
	27: "Canada/Newfoundland", //"UTC-03:30 - Newfoundland",
	28: "UTC-03:00 - Brasilia",
	29: "UTC-03:00 - Buenos Aires",
	30: "UTC-03:00 - Cayenne, Fortaleza",
	31: "UTC-03:00 - Greenland",
	32: "UTC-03:00 - Montevideo",
	33: "UTC-03:00 - Salvador",
	34: "UTC-02:00 - Coordindated Universal Time",
	35: "UTC-01:00 - Azores",
	36: "UTC-01:00 - Cabo Verde Is.",
	37: "UTC", //"UTC - Casablanca",
	38: "UTC", //"UTC - Coordindated Universal Time",
	39: "Europe/London", //"UTC - Dublin, Edinburgh, Lisbon, London",
	40: "UTC", //"UTC - Monrovia, Reykjavik",
	41: "Europe/Berlin", //"UTC+01:00 - Amsterdam, Berlin, Bern",
	42: "UTC+01:00 - Belgrade, Bratislava",
	43: "UTC+01:00 - Brussels, Copenhagen",
	44: "UTC+01:00 - Sarajevo, Skopje, Warsaw",
	45: "UTC+01:00 - West Central Africa",
	46: "UTC+01:00 - Windkoek",
	47: "UTC+02:00 - Amman",
	48: "UTC+02:00 - Athens, Bucharest",
	49: "UTC+02:00 - Beirut",
	50: "UTC+02:00 - Cairo",
	51: "UTC+02:00 - Damascus",
	52: "UTC+02:00 - E. Europe",
	53: "UTC+02:00 - Harare, Pretoria",
	54: "UTC+02:00 - Helsinki, Kyiv, Riga, Sofia",
	55: "UTC+02:00 - Istanbul",
	56: "UTC+02:00 - Jerusalem",
	57: "UTC+02:00 - Kalinigrad (RTZ 1)",
	58: "UTC+02:00 - Tripoli",
	59: "UTC+02:00 - Baghdad",
	60: "UTC+03:00 - Kuwait, Riyadh",
	61: "UTC+03:00 - Minsk",
	62: "Europe/Moscow", //"UTC+03:00 - Moscow, St. Petersburg",
	63: "UTC+03:00 - Nairobi",
	64: "UTC+03:30 - Tehran",
	65: "UTC+04:00 - Abu Dhabi, Muscat",
	66: "UTC+04:00 - Baku",
	67: "UTC+04:00 - Izhevsk, Samara (RTZ 3)",
	68: "UTC+04:00 - Port Louis",
	69: "UTC+04:00 - Tbilisi",
	70: "UTC+04:00 - Yerevan",
	71: "UTC+04:30 - Kabal",
	72: "UTC+05:00 - Ashgabat, Tashkent",
	73: "UTC+05:00 - Ekaterinburg (RTZ 4)",
	74: "UTC+05:00 - Islamabad, Karachi",
	75: "UTC+05:30 - Chennai, Kolkata, Mumbai",
	76: "UTC+05:30 - Sri Jayawardenepura",
	77: "UTC+05:45 - Kathmandu",
	78: "UTC+06:00 - Astana",
	79: "UTC+06:00 - Dhaka",
	80: "UTC+06:00 - Novosibirsk (RTZ 5)",
	81: "UTC+06:30 - Yangon (Rangoon)",
	82: "UTC+07:00 - Bankok, Hanoi, Jakarta",
	83: "UTC+07:00 - Kransnoyarsk (RTZ 6)",
	84: "UTC+08:00 - Beijing, Chongqing, Hong Kong",
	85: "UTC+08:00 - Irkutsk (RTZ 7)",
	86: "UTC+08:00 - Kuala Lumpur, Singapore",
	87: "UTC+08:00 - Perth",
	88: "UTC+08:00 - Taipei",
	89: "UTC+08:00 - Ulaanbaatar",
	90: "UTC+09:00 - Osaka, Sapporo, Tokyo",
	91: "UTC+09:00 - Seoul",
	92: "UTC+09:00 - Yakutsk (RTZ 8)",
	93: "UTC+09:30 - Adelaide",
	94: "UTC+09:30 - Darwin",
	95: "UTC+10:00 - Brisbane",
	96: "UTC+10:00 - Canberra, Melbourne, Sydney",
	97: "UTC+10:00 - Guam, Port Moresby",
	98: "UTC+10:00 - Hobart",
	99: "UTC+10:00 - Magadan",
	100: "UTC+10:00 - Vladivostok, Magadan (RTZ 9)",
	101: "UTC+11:00 - Chokurdakh (RTZ 10)",
	102: "UTC+11:00 - Solomon Is., New Caledonia",
	103: "UTC+12:00 - Anadyr, Petropavlovsk",
	104: "Pacific/Auckland", //"UTC+12:00 - Auckland, Wellington",
	105: "UTC+12:00 - Corrdinated Universal Time",
	106: "UTC+12:00 - Fiji",
	107: "UTC+13:00 - Nuku'alofa",
	108: "UTC+13:00 - Samoa",
	109: "UTC+14:00 - Kiritimati Island",
}

func (tz Timezone) String() string {
	return timezoneIndex[tz]
}

func (tz Timezone) Location() *time.Location {
	loc, err := time.LoadLocation(tz.String())
	if err != nil {
		log.Printf("no timezone for %d/%s: %s", tz, tz.String(), err)
		return nil
	}
	return loc
}
