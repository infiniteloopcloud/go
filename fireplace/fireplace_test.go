package fireplace

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"testing"
)

func TestHandler_Search(t *testing.T) {
	h := New()
	err := h.Add("file1", "./fixtures/file1")
	if err != nil {
		t.Error(err)
	}
	result, err := h.Search("A75AA268-4738-9C61-A132-4A3105E9440B")
	if err != nil {
		t.Error(err)
	}
	if len(result) != 1 {
		t.Fatal("result length must be 1")
	}
	t.Log(result[0].FileID)
	t.Log(result[0].LineNumber)
	t.Log(result[0].LineContent)
}

func TestHandler_SearchMultipleResult(t *testing.T) {
	h := New()
	err := h.Add("file3", "./fixtures/file3")
	if err != nil {
		t.Error(err)
	}
	results, err := h.Search("E3060FBC-EEB1-E0A8-EB94-511896939D2C")
	if err != nil {
		t.Error(err)
	}
	if len(results) < 1 {
		t.Fatal("results length must be greater than 0")
	}
	t.Log(len(results))
	for _, result := range results {
		t.Log(result.FileID)
		t.Log(result.LineNumber)
		t.Log(result.LineContent)
	}
}

func TestHandler_ReplaceMultiple(t *testing.T) {
	h := New()
	err := h.Add("file1", "./fixtures/file1")
	if err != nil {
		t.Error(err)
	}
	err = h.Add("file2", "./fixtures/file2")
	if err != nil {
		t.Error(err)
	}

	var lines = map[string]struct{}{
		"A75AA268-4738-9C61-A132-4A3105E9440B": {},
		"52378C39-4E9B-9400-058A-7B76C8F4C710": {},
		"6C662E3C-DFD6-94C6-2C13-72778A1535BE": {},
		"AB9802F1-283A-EBAA-52ED-8A096D946D1F": {},
	}
	var param = make(map[string][]*ReplaceOpts)
	for line := range lines {
		results, err := h.Search(line)
		if err != nil {
			t.Fatal(err)
		}
		if len(results) < 1 {
			t.Fatal("result length must not be zero")
		}
		for _, result := range results {
			param[result.FileID] = append(param[result.FileID], &ReplaceOpts{
				LineNumber: result.LineNumber,
				ReplaceTo:  result.LineContent + " replaced",
			})
		}
	}

	err = h.ReplaceMultiple(param)
	if err != nil {
		t.Error(err)
	}
}

func TestFileGen(t *testing.T) {
	f, err := os.OpenFile("./large_file", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	for i := 0; i < 4000000; i++ {
		id := pseudoUUID()
		if _, err = f.WriteString(id + "\n"); err != nil {
			t.Error(err)
		}
	}
}

func TestFileGen2(t *testing.T) {
	t.Skip()
	f, err := os.OpenFile("./large_file_2", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	repeatedID := pseudoUUID()

	for i := 0; i < 10000; i++ {
		id := pseudoUUID()
		if i%19 == 0 {
			id += "," + strings.Join(randomStringSlice(repeatedID, pseudoUUID), ",")
		}
		if _, err = f.WriteString(id + "\n"); err != nil {
			t.Error(err)
		}
	}
}

//nolint:unused
func randomStringSlice(id string, randomStr func() string) []string {
	randInt, err := rand.Int(rand.Reader, big.NewInt(10))
	if err != nil {
		randInt = big.NewInt(7)
	}
	randInt = randInt.Add(big.NewInt(1), big.NewInt(0))
	var res = make([]string, 0, int(randInt.Int64()))
	repeatTimes, err := rand.Int(rand.Reader, randInt)
	if err != nil {
		repeatTimes = big.NewInt(3)
	}

	for i := 0; i < int(randInt.Int64()+repeatTimes.Int64()); i++ {
		if i%2 == 0 {
			res = append(res, id)
		} else {
			res = append(res, randomStr())
		}
	}
	return res
}

func BenchmarkHandler_Search(b *testing.B) {
	h := New()
	err := h.Add("large_file", "./large_file")
	if err != nil {
		return
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := h.Search("1550139E-9942-C3DA-3F6D-69D74ACA53BA")
		if err != nil {
			b.Error(err)
		}
	}
}

func pseudoUUID() (uuid string) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}
