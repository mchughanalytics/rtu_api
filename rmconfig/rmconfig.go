package rmconfig

// rmConfig - struct for remarkable device config file
type rmConfig struct {
	General         map[string]string
	OnboardingTools map[string]string
	Wifinetworks    map[string]string
}

func NewRmConfig() *rmConfig {

	g := map[string]string{
		"ConvertLegacyLinesOnStartup": "",
		"ConvertUuidLinesOnStartup":   "",
		"DeveloperPassword":           "",
		"HomeSortOrder":               "",
		"LatestLockTime":              "",
		"NumberOfPasscodeTries":       "",
		"OnboardingBanner":            "",
		"OnboardingProgress":          "",
		"OnboardingUpdateScreen":      "",
		"Password":                    "",
		"PreviousVersion":             "",
		"RecentTemplates":             "",
		"Setup":                       "",
		"ShareEmailAddresses":         "",
		"WebInterfaceEnabled":         "",
		"deviceid":                    "",
		"devicetoken":                 "",
		"usertoken":                   "",
		"version":                     "",
		"wifion":                      "",
	}

	ot := map[string]string{
		"Brush":         "",
		"BrushColor":    "",
		"Eraser":        "",
		"Highlighter":   "",
		"Paintbrush":    "",
		"PenType":       "",
		"Pencil":        "",
		"PencilType":    "",
		"SelectionTool": "",
		"ToolSize":      "",
		"WritingTool":   "",
	}

	rmc := rmConfig{
		General:         g,
		OnboardingTools: ot,
	}

	return &rmc
}
