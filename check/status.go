package check

const (
    CHECK_OKAY string = "OKAY"
    CHECK_FAIL string = "FAIL"
)

type Status struct {
    Check Check `json:"check"`
    CheckData CheckData `json:"check_data"`
    Annotations map[string]string `json:"annotations,omitempty"`
}

func (s *Status) GetAnnotation(name string) (value string, ok bool) {
    value, ok = s.Annotations[name]

    return
}

func (s *Status) SetAnnotation(name string, value string) {
    if nil == s.Annotations {
        s.Annotations = map[string]string{}
    }

    s.Annotations[name] = value
}
