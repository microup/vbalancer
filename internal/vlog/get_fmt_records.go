package vlog

const maxLenString = 250

//nolint:cyclop
func (c *VLog) GetRecordsText(isReverse bool) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	result := ""

	if c.mapLastLogRecords == nil {
		return result
	}

	if len(c.mapLastLogRecords) == 0 {
		return result
	}

	switch isReverse {
	case true:
		{
			for i, j := 0, len(c.mapLastLogRecords)-1; j > -1; i, j = i+1, j-1 {
				strLastLogRecord := c.mapLastLogRecords[j]

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
			for _, strLastLogRecord := range c.mapLastLogRecords {
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
