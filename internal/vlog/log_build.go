package vlog

import (
	"fmt"
	"time"
	"vbalancer/internal/core"
	"vbalancer/internal/types"
)

// BuildRecord function takes in several input values such as log type, result code,
// remote address, client host, proxy host and string values. It creates a record in the
// specified format with the current date and time by using the "GetDateTimeStr" function
// from the "core" package. The result is returned as a tuple with the log type and the formatted record string.
func BuildRecord(
	typeLog types.TypeLog,
	resultCode types.ResultCode,
	remoteAddr types.RemoteAddr,
	peerAddr types.PeerAddr,
	valuesStr string) (types.TypeLog, string) {
	var recordTime = time.Now()

	var dateStr, timeStr = core.GetDateTimeStr(recordTime)

	var resultStr = core.FmtStringWithDelimiter(";", valuesStr)

	var recordRow = fmt.Sprintf("%s;%s;%s;%d;%s;%s;%s;",
		dateStr,
		timeStr,
		typeLog.GetStr(),
		resultCode,
		remoteAddr,
		peerAddr,
		resultStr)

	return typeLog, recordRow
}

