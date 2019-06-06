package printer

func ServerBandwidth(bandwidth []map[string]string) {
	col := columns{"DATE", "INCOMING BYTES", "OUTGOING BYTES"}
	display(col)
	for _, b := range bandwidth {
		display(columns{b["date"], b["incoming"], b["outgoing"]})
	}
	flush()
}
