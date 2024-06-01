package csv

import (
	"encoding/csv"
	"os"

	"github.com/gofish2020/MedicalSpider/utils"
)

const bomHeader = "\xEF\xBB\xBF" // UTF8-BOM (Byte Order Mark) 编码头

type File struct {
	file *os.File

	writer *csv.Writer
}

func (f *File) Open(fileName string) error {
	file, err := utils.OpenFile(fileName, "./data/")

	if err != nil {
		return err
	}
	file.WriteString(bomHeader)
	f.file = file
	f.writer = csv.NewWriter(file)

	return nil
}

func (f *File) WriteAll(data [][]string) {
	f.writer.WriteAll(data)
}

func (f *File) Write(data []string) {
	f.writer.Write(data)
}

func (f *File) Close() {

	if f.file != nil {
		f.writer.Flush()
		f.file.Close()
	}

}

func NewFile(fileName string) File {
	f := File{}
	f.Open(fileName)
	return f
}
