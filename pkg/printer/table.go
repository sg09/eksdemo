package printer

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

type TablePrinter struct {
	header       []string
	data         [][]string
	rowSeperator bool
}

// TODO: use make to set the size of the slice
func NewTablePrinter() *TablePrinter {
	return &TablePrinter{}
}

func (p *TablePrinter) AppendRow(row []string) {
	p.data = append(p.data, row)
}

func (p *TablePrinter) SetHeader(header []string) {
	p.header = header
}

func (p *TablePrinter) SeparateRows() {
	p.rowSeperator = true
}

func (p *TablePrinter) Print(writer io.Writer) {
	if len(p.data) == 0 {
		fmt.Println("No resources found.")
		return
	}
	table := tablewriter.NewWriter(writer)
	table.SetAutoFormatHeaders(false)

	if p.rowSeperator {
		table.SetRowLine(true)
	}

	table.SetHeader(p.header)
	table.AppendBulk(p.data)
	table.Render()
}
