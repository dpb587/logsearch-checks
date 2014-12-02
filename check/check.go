package check

type Check struct {
    Owner string `json:"owner"`
    Name string `json:"name"`
    Status string `json:"status"`
}

func (check *Check) GetOwner() (string) {
    return check.Owner
}

func (check *Check) GetName() (string) {
    return check.Name
}

func (check *Check) GetStatus() (string) {
    return check.Status
}
