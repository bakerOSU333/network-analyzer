package table

import (
	"os"

	"github.com/bakerOSU333/network-analyzer/network"
	"github.com/jedib0t/go-pretty/table"
)

// function to print out beautiful table along with network consumption chart
func PrintNetworkingTable(networkDataMap map[string]network.NetworkData, keyDesc []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// set the header
	t.AppendHeader(table.Row{"Process Name", "Incoming data (MB)", "Outgoing data (MB)", "Time"})

	// key is sorted by MBIn (incoming network)
	for _, processName := range keyDesc {
		dataLength := len(networkDataMap[processName].Time)

		// get the latest network in data
		MBIn := networkDataMap[processName].ReceivedMB[dataLength - 1]

		// get the latest network out data
		MBOut := networkDataMap[processName].SentMB[dataLength - 1]

		// get the latest time recorded
		Time := networkDataMap[processName].Time[dataLength - 1]

		// append them all to the row
		t.AppendRow(table.Row{processName, MBIn, MBOut, Time})
	}
	t.AppendFooter(table.Row{"Table is sorted by Incoming network"})

	// set auto index
	t.SetAutoIndex(true)

	// set style
	t.SetStyle(table.StyleColoredBlackOnMagentaWhite)

	// render the table to the terminal
	t.Render()
}