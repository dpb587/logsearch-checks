package check

type CheckData struct {
    Threshold float64 `json:"threshold"`
    Value float64 `json:"value"`
    Units string `json:"units"`
    Extra map[string]float64 `json:"extra,omitempty"`
}

func (check *CheckData) GetThreshold() (float64) {
    return check.Threshold
}

func (check *CheckData) GetValue() (float64) {
    return check.Value
}

func (check *CheckData) GetUnits() (string) {
    return check.Units
}

func (check *CheckData) GetExtraValue(name string) (value float64, ok bool) {
    value, ok = check.Extra[name]

    return
}
