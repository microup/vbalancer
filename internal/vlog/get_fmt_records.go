package vlog

const maxLenString = 250

//nolint:cyclop
func (v *VLog) GetRecordsText(isReverse bool) string {
	v.Mu.Lock()
	defer v.Mu.Unlock()

	result := ""

	if v.MapLastLogRecords == nil {
		return result
	}

	if len(v.MapLastLogRecords) == 0 {
		return result
	}

	switch isReverse {
	case true:
		{
			for i, j := 0, len(v.MapLastLogRecords)-1; j > -1; i, j = i+1, j-1 {
				strLastLogRecord := v.MapLastLogRecords[j]

				if result == "" {
					result = strLastLogRecord
				} else {
					if len(strLastLogRecord) < maxLenString {
						result = result + "<BR>" + strLastLogRecord
					} else {
						result = result + "<BR><BR>" + strLastLogRecord + "<BR>"
					}
				}
			}
		}
	case false:
		{
			for _, strLastLogRecord := range v.MapLastLogRecords {
				if result == "" {
					result = strLastLogRecord
				} else {
					if len(strLastLogRecord) < maxLenString {
						result = result + "<BR>" + strLastLogRecord
					} else {
						result = result + "<BR><BR>" + strLastLogRecord + "<BR>"
					}
				}
			}
		}
	}

	return result
}
