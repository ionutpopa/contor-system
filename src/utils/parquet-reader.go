package utils

import (
	"log"

	"github.com/xitongsys/parquet-go/reader"
)

func ReadParquets() {
	///read
	fr, err := NewLocalFileReader("logs/2024-11-24.parquet")
	if err != nil {
		log.Println("Can't open file")
		return
	}

	pr, err := reader.NewParquetReader(fr, new(LogEntry), 4)
	if err != nil {
		log.Println("Can't create parquet reader", err)
		return
	}
	num := int(pr.GetNumRows())
	for i := 0; i < num/10; i++ {
		// if i%2 == 0 {
		// 	pr.SkipRows(10) //skip 10 rows
		// 	continue
		// }
		stus := make([]LogEntry, 10) //read 10 rows
		if err = pr.Read(&stus); err != nil {
			log.Println("Read error", err)
		}
		log.Println(stus)
	}

	pr.ReadStop()
	fr.Close()
}
