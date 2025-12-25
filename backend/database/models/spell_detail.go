package models

type SpellDetail struct {
	*SpellTemplateFull
	Icon     string `json:"icon"`
	ToolTip  string `json:"toolTip"`
	CastTime string `json:"castTime"`
	Range    string `json:"range"`
	Duration string `json:"duration"`
	Power    string `json:"power"`
}
