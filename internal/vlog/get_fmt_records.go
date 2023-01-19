package vlog

const maxLenString = 250

//nolint:cyclop
func (v *VLog) GetRecordsText(isReverse bool) string {
	v.mu.Lock()
	defer v.mu.Unlock()

	result := ""

	if v.mapLastLogRecords == nil {
		return result
	}

	if len(v.mapLastLogRecords) == 0 {
		return result
	}

	switch isReverse {
	case true:
		{
			for i, j := 0, len(v.mapLastLogRecords)-1; j > -1; i, j = i+1, j-1 {
				strLastLogRecord := v.mapLastLogRecords[j]

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
			for _, strLastLogRecord := range v.mapLastLogRecords {
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
