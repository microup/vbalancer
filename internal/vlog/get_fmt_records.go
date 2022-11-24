package vlog

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
				st := c.mapLastLogRecords[j]

				if result == "" {
					result = st
				} else {
					if len(st) < 250 {
						result = result + "<BR>" + st
					} else {
						result = result + "<BR><BR>" + st + "<BR>"
					}
				}
			}
		}
	case false:
		{
			for _, st := range c.mapLastLogRecords {
				if result == "" {
					result = st
				} else {
					if len(st) < 250 {
						result = result + "<BR>" + st
					} else {
						result = result + "<BR><BR>" + st + "<BR>"
					}
				}
			}
		}
	}

	return result
}
