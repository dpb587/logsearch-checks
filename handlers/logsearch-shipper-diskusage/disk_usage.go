package logsearchshipperdiskusage

type DiskUsage struct {
  MissingData bool
  Used float64
  Free float64
}

func (du *DiskUsage) IsMissingData() (bool) {
  return du.MissingData
}

func (du *DiskUsage) GetUsed() (float64) {
  return du.Used
}

func (du *DiskUsage) GetUsedPct() (float64) {
  return du.Used / du.GetTotal() * 100
}

func (du *DiskUsage) GetFree() (float64) {
  return du.Free
}

func (du *DiskUsage) GetFreePct() (float64) {
  return du.Free / du.GetTotal() * 100
}

func (du *DiskUsage) GetTotal() (float64) {
  return du.Used + du.Free
}
