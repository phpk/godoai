package convert

import (
	"io"
	"strings"

	"godoai/convert/libs"

	"github.com/pbnjay/grate"
	_ "github.com/pbnjay/grate/simple" // tsv and csv support
	_ "github.com/pbnjay/grate/xls"
	_ "github.com/pbnjay/grate/xlsx"
)

func ConvertXlsx(r io.Reader) (string, error) {
	absFileFrom, tmpfromfile, err := libs.GetTempFile(r, "prefix-xlsx-from")
	if err != nil {
		return "", err
	}
	text := ""
	wb, _ := grate.Open(absFileFrom) // open the file
	sheets, _ := wb.List()           // list available sheets
	for _, s := range sheets {       // enumerate each sheet name
		sheet, _ := wb.Get(s) // open the sheet
		for sheet.Next() {    // enumerate each row of data
			row := sheet.Strings() // get the row's content as []string
			//fmt.Println(strings.Join(row, "\t"))
			// 跳过空记录
			if len(row) == 0 {
				continue
			}
			text += strings.Join(row, "\t") + "\n"
		}
	}
	wb.Close()
	libs.CloseTempFile(tmpfromfile)
	return text, nil
}
