package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	KB        int64  = 1000
	MB        int64  = 1000 * KB
	GB        int64  = 1000 * MB
	addedFile string = ".x_"
)

func main() {
	fmt.Println("Start")
	path := flag.String("f", ".", "Путь до файла")
	sizeFile := flag.String("s", "50", "Размер файла")
	typeSize := flag.String("t", "byte", "Тип размерности: byte, KB, MB, GB")

	flag.Parse()

	if path == nil || sizeFile == nil || typeSize == nil {
		fmt.Println("Введены не все параметры. Ознакомиться с ними можно по команде help")
		os.Exit(1)
	}
	fmt.Printf("%s %s %s\n", *path, *sizeFile, *typeSize)

	fmt.Printf("Открываю файл %s. . .\n", *path)
	fileIn, err := os.Open(*path)
	if err != nil {
		fmt.Printf("Произошла ошибка при чтении файла %s: %s", *path, err.Error())
		os.Exit(1)
	}
	defer fileIn.Close()
	fmt.Println("Успешно")
	size := getSizeFile(fileIn, *path)
	if size == -1 {
		os.Exit(1)
	}
	sizeFileConv, err := strconv.ParseInt(*sizeFile, 10, 32)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	generalSize, typeS, sizeFileConv := getGeneralSizeFile(size, sizeFileConv)
	fmt.Printf("Размер файла: %.2f %s\n\n", generalSize, typeS)

	reader := bufio.NewReader(fileIn)
	var position int64 = 0
	var seq int = 0

	var fileOut *os.File
	newFileName := getNewFileName(*path, seq)
	fileOut, err = os.Create(newFileName)
	if err != nil {
		fmt.Printf("Ошибка при создании файла %s: %s", newFileName, err.Error())
		os.Exit(1)
	}
	fmt.Println("Начинаю дробление. . .")
	writer := bufio.NewWriter(fileOut)
	for {
		data, err := reader.ReadByte()
		if err != nil {
			errFlush := writer.Flush()
			if errFlush != nil {
				fmt.Printf("Ошибка при сбросе буфера в файл: %s", errFlush.Error())
				os.Exit(1)
			}
			fileOut.Close()
			if err == io.EOF {
				break
			} else {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
		//fmt.Println(val)
		if position == sizeFileConv && data != 32 {
			err = writer.Flush()
			if err != nil {
				fmt.Printf("Ошибка при сбросе буфера в файл: %s", err.Error())
				os.Exit(1)
			}
			fileOut.Close()
			seq++
			position = 0
			newFileName = getNewFileName(*path, seq)
			fileOut, err = os.Create(newFileName)
			if err != nil {
				fmt.Printf("Ошибка при создании файла %s: %s", newFileName, err.Error())
				os.Exit(1)
			}
			writer = bufio.NewWriter(fileOut)
		}
		err = writer.WriteByte(data)
		if err != nil {
			fmt.Printf("Ошибка при записи: %s", err.Error())
			os.Exit(1)
		}
		position++
	}
	/* err = writer.Flush()
	if err != nil {
		fmt.Printf("Ошибка при сбросе буфера в файл: %s", err.Error())
		os.Exit(1)
	}
	fileOut.Close() */
	fmt.Println("Готово")
}

func getSizeFile(file *os.File, path string) int64 {
	stat, err := file.Stat()
	if err != nil {
		fmt.Printf("Произошла ошибка при получении метаданных файла %s: %s", path, err.Error())
		return -1
	}
	return stat.Size()
}

func getGeneralSizeFile(sizeFile int64, sizeIn int64) (float64, string, int64) {
	var resultType string
	var resultSize float64
	var resultSizeConv int64
	if sizeFile < KB {
		resultType = "byte"
		resultSize = float64(sizeFile)
		resultSizeConv = sizeIn
	} else if sizeFile < MB {
		resultType = "KB"
		resultSize = float64(sizeFile) / float64(KB)
		resultSizeConv = sizeIn * KB
	} else if sizeFile < GB {
		resultType = "MB"
		resultSize = float64(sizeFile) / float64(MB)
		resultSizeConv = sizeIn * MB
	} else {
		resultType = "GB"
		resultSize = float64(sizeFile) / float64(GB)
		resultSizeConv = sizeIn * GB
	}
	return resultSize, resultType, resultSizeConv
}

func getNewFileName(fileName string, seq int) string {
	return fileName + addedFile + strconv.Itoa(seq)
}
