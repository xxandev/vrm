package vrm

type Installations struct {
	Success bool `json:"success,omitempty"`
	Records []struct {
		IDSite                              int         `json:"idSite,omitempty"`
		AccessLevel                         int         `json:"accessLevel,omitempty"`
		Owner                               bool        `json:"owner,omitempty"`
		IsAdmin                             bool        `json:"is_admin,omitempty"`
		Name                                string      `json:"name,omitempty"`
		Identifier                          string      `json:"identifier,omitempty"`
		IDUser                              int         `json:"idUser,omitempty"`
		PvMax                               int         `json:"pvMax,omitempty"`
		Timezone                            string      `json:"timezone,omitempty"`
		Phonenumber                         interface{} `json:"phonenumber,omitempty"`
		Notes                               interface{} `json:"notes,omitempty"`
		Geofence                            interface{} `json:"geofence,omitempty"`
		GeofenceEnabled                     bool        `json:"geofenceEnabled,omitempty"`
		RealtimeUpdates                     bool        `json:"realtimeUpdates,omitempty"`
		HasMains                            int         `json:"hasMains,omitempty"`
		HasGenerator                        int         `json:"hasGenerator,omitempty"`
		NoDataAlarmTimeout                  int         `json:"noDataAlarmTimeout,omitempty"`
		AlarmMonitoring                     int         `json:"alarmMonitoring,omitempty"`
		InvalidVRMAuthTokenUsedInLogRequest int         `json:"invalidVRMAuthTokenUsedInLogRequest,omitempty"`
		Syscreated                          int         `json:"syscreated,omitempty"`
		GrafanaEnabled                      int         `json:"grafanaEnabled,omitempty"`
		IsPaygo                             int         `json:"isPaygo,omitempty"`
		PaygoCurrency                       interface{} `json:"paygoCurrency,omitempty"`
		PaygoTotalAmount                    interface{} `json:"paygoTotalAmount,omitempty"`
		InverterChargerControl              int         `json:"inverterChargerControl,omitempty"`
		Shared                              bool        `json:"shared,omitempty"`
		DeviceIcon                          string      `json:"device_icon,omitempty"`
		Alarm                               bool        `json:"alarm,omitempty"`
		LastTimestamp                       int         `json:"last_timestamp,omitempty"`
		Tags                                []struct {
			IDTag     int    `json:"idTag,omitempty"`
			Name      string `json:"name,omitempty"`
			Automatic bool   `json:"automatic,omitempty"`
		} `json:"tags,omitempty"`
		CurrentTime     string `json:"current_time,omitempty"`
		TimezoneOffset  int    `json:"timezone_offset,omitempty"`
		Images          bool   `json:"images,omitempty"`
		ViewPermissions struct {
			UpdateSettings bool `json:"update_settings,omitempty"`
			Settings       bool `json:"settings,omitempty"`
			Diagnostics    bool `json:"diagnostics,omitempty"`
			Share          bool `json:"share,omitempty"`
			Vnc            bool `json:"vnc,omitempty"`
			MqttRPC        bool `json:"mqtt_rpc,omitempty"`
			Vebus          bool `json:"vebus,omitempty"`
			Twoway         bool `json:"twoway,omitempty"`
			ExactLocation  bool `json:"exact_location,omitempty"`
			Nodered        bool `json:"nodered,omitempty"`
			NoderedDash    bool `json:"nodered_dash,omitempty"`
			Signalk        bool `json:"signalk,omitempty"`
			Paygo          bool `json:"paygo,omitempty"`
		} `json:"view_permissions,omitempty"`
		Extended []struct {
			IDDataAttribute         int         `json:"idDataAttribute,omitempty"`
			Code                    string      `json:"code,omitempty"`
			Description             string      `json:"description,omitempty"`
			FormatWithUnit          string      `json:"formatWithUnit,omitempty"`
			DataType                string      `json:"dataType,omitempty"`
			IDDeviceType            int         `json:"idDeviceType,omitempty"`
			TextValue               interface{} `json:"textValue,omitempty"`
			Instance                string      `json:"instance,omitempty"`
			Timestamp               string      `json:"timestamp,omitempty"`
			DbusServiceType         interface{} `json:"dbusServiceType,omitempty"`
			DbusPath                interface{} `json:"dbusPath,omitempty"`
			RawValue                interface{} `json:"rawValue,omitempty"`
			FormattedValue          string      `json:"formattedValue,omitempty"`
			Instances               interface{} `json:"instances,omitempty"`
			DataAttributeEnumValues []struct {
				NameEnum  string `json:"nameEnum,omitempty"`
				ValueEnum int    `json:"valueEnum,omitempty"`
			} `json:"dataAttributeEnumValues,omitempty"`
			DataAttributes []struct {
				Instance        int    `json:"instance,omitempty"`
				DbusServiceType string `json:"dbusServiceType,omitempty"`
				DbusPath        string `json:"dbusPath,omitempty"`
			} `json:"dataAttributes,omitempty"`
		} `json:"extended,omitempty"`
		DemoMode      bool        `json:"demo_mode,omitempty"`
		MqttWebhost   string      `json:"mqtt_webhost,omitempty"`
		HighWorkload  bool        `json:"high_workload,omitempty"`
		CurrentAlarms []string    `json:"current_alarms,omitempty"`
		NumAlarms     int         `json:"num_alarms,omitempty"`
		AvatarURL     interface{} `json:"avatar_url,omitempty"`
	} `json:"records,omitempty"`
}
