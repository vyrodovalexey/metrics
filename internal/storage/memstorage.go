package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type MemStorage struct {
	GaugeMap   map[string]Gauge
	CounterMap map[string]Counter
}

func (m *MemStorage) New() {
	m.GaugeMap = make(map[string]Gauge)
	m.CounterMap = make(map[string]Counter)
}

func (m *MemStorage) AddCounterAsString(name string, item string) error {
	counter, err := strconv.ParseInt(item, 10, 64)
	if err == nil {
		m.CounterMap[name] = m.CounterMap[name] + counter
	}
	return err

}

func (m *MemStorage) AddCounter(name string, item Counter) {
	m.CounterMap[name] = m.CounterMap[name] + item

}

func (m *MemStorage) AddGaugeAsString(name string, item string) error {
	gauge, err := strconv.ParseFloat(item, 64)
	if err == nil {
		m.GaugeMap[name] = gauge
	}
	return err

}

func (m *MemStorage) AddGauge(name string, item Gauge) {

	m.GaugeMap[name] = item

}

func (m *MemStorage) GetAllMetricNames() (map[string]string, map[string]string) {
	//to debug
	//names := make([]string, 0, len(storage.GaugeMap)+len(storage.CounterMap))
	gvalues := make(map[string]string, len(m.GaugeMap))

	cvalues := make(map[string]string, len(m.CounterMap))
	// Iterate over the map and collect the keys
	for name := range m.GaugeMap {
		//to debug
		//names = append(names, name)
		gv, _ := m.GetGauge(name)
		gvalues[name] = fmt.Sprintf("%v", gv)
	}

	for name := range m.CounterMap {
		//to debug
		//names = append(names, name)
		cv, _ := m.GetCounter(name)
		cvalues[name] = fmt.Sprintf("%v", cv)
	}

	return gvalues, cvalues
}

func (m *MemStorage) GetGauge(name string) (Gauge, bool) {
	res, e := m.GaugeMap[name]
	if e {
		return res, e
	}
	return 0, false

}

func (m *MemStorage) GetCounter(name string) (Counter, bool) {
	res, e := m.CounterMap[name]

	if e {
		return res, e
	}
	return 0, false
}

func (m *MemStorage) Load(f *os.File) error {
	m.New()
	// Read the file's content
	byteValue, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	if len(byteValue) > 0 {
		// Unmarshal the JSON data into the struct
		err = json.Unmarshal(byteValue, m)
		fmt.Println(m)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
		}
	}
	return err
}

func (m *MemStorage) SaveAsync(f *os.File, interval int) {

	for {
		mst, err := json.Marshal(m)

		if err != nil {
			fmt.Println("Error move to json:", err)

		}
		err = f.Truncate(0)
		if err != nil {
			fmt.Println("Can't truncate file error:", err)
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			fmt.Println("Can't seek on start error:", err)
		}
		_, err = f.Write(mst)
		if err != nil {
			fmt.Println("Error writing to file:", err)

		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func (m *MemStorage) Save(f *os.File) {

	mst, err := json.Marshal(m)

	if err != nil {
		fmt.Println("Error move to json:", err)

	}
	err = f.Truncate(0)
	if err != nil {
		fmt.Println("Can't truncate file error:", err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		fmt.Println("Can't seek on start error:", err)
	}
	_, err = f.Write(mst)
	if err != nil {
		fmt.Println("Error writing to file:", err)

	}

}
