package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func WriteToCsv(speedSort []Speed) error {
	speedResultFile, err := os.Create(OutFilePath)
	if err != nil {
		return err
	}
	defer func(speedResultFile *os.File) {
		err := speedResultFile.Close()
		if err != nil {
			return
		}
	}(speedResultFile)
	writer := csv.NewWriter(speedResultFile)
	defer writer.Flush()
	var data [][]string
	for _, v := range speedSort {
		tmp := []string{v.Server, fmt.Sprintf("%d ms", v.TimeOut)}
		data = append(data, tmp)
	}
	for _, row := range data {
		err := writer.Write(row)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadFromCsv(filePath string) ([]Speed, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	reader := csv.NewReader(file)
	var speedSort []Speed
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		server := row[0]
		timeOut, err := strconv.Atoi(strings.Replace(row[1], " ms", "", -1))
		if err != nil {
			return nil, err
		}
		speedSort = append(speedSort, Speed{
			Server:  server,
			TimeOut: int64(timeOut),
		})
	}
	return speedSort, nil
}
