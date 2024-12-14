package utils

import (
	"fmt"
	"strconv"
	"testing"
)

func TestGenerate_shortener(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"length 5", 5},
		{"length 10", 10},
		{"length 0", 0},
		{"length 20", 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Generate_shortener(tt.n)
			if len(got) != tt.n {
				t.Errorf("Generate_shortener(%d) = %v, want length %d", tt.n, got, tt.n)
			}
		})
	}
}

func TestGenerate_shortenerRepeat(t *testing.T) {
	t.Run("% of repeat a shortener", func(t *testing.T) {
		epoch := 100_000
		recorder := make(map[string]int8)
		count := 0
		for range epoch {
			got := Generate_shortener(8)
			if _, ok := recorder[got]; ok {
				count++
			} else {
				recorder[got] = 1
			}
		}
		t.Logf("Repetiu %d do total de %d números", count, epoch)
	})
}

func TestList_Add(t *testing.T) {
	tests := []struct {
		name string
		data string
		key  string
	}{
		{"Add first element", "data1", "key1"},
		{"Add second element", "data2", "key2"},
		{"Add third element", "data3", "key3"},
		{"Add fourth element", "data4", "key4"},
		{"Add fifth element (should remove oldest)", "data5", "key5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := make(map[string]*Node)
			list := &List{}
			for i := 1; i <= 4; i++ {
				list.Add("data"+strconv.Itoa(i), cache, "key"+strconv.Itoa(i))
			}
			list.Add(tt.data, cache, tt.key)
			fmt.Println(cache)

			if len(cache) > 4 {
				t.Errorf("cache size = %d, want <= 4", len(cache))
			}

			if list.Size > 4 {
				t.Errorf("list size = %d, want <= 4", list.Size)
			}

			if list.Head.Data != tt.data {
				t.Errorf("head data = %s, want %s", list.Head.Data, tt.data)
			}
		})
	}
}

func TestList_Get(t *testing.T) {
	tests := []struct {
		name      string
		setupData []string
		key       string
		want      string
		wantErr   bool
	}{
		{"Get existing element", []string{"data1", "data2", "data3", "data4"}, "data2", "data2", false},
		{"Get non-existing element", []string{"data1", "data2", "data3", "data4"}, "data5", "", true},
		{"Get head element", []string{"data1", "data2", "data3", "data4"}, "data4", "data4", false},
		{"Get tail element", []string{"data1", "data2", "data3", "data4"}, "data1", "data1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := make(map[string]*Node)
			list := &List{}
			for _, data := range tt.setupData {
				list.Add(data, cache, data)
			}

			got, err := list.Get(tt.key, cache)
			if (err != nil) != tt.wantErr {
				t.Errorf("List.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && *got != tt.want {
				t.Errorf("List.Get() = %v, want %v", *got, tt.want)
			}
		})
	}
}
