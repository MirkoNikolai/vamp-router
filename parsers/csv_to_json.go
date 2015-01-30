package parsers

import (
  "encoding/csv"
  "fmt"
  "io"
  "strings"
  "bytes"
)

// parses the raw stats CSV output to a json string
func CsvToJson(csvInput string) (string, error){

  csvReader := csv.NewReader(strings.NewReader(csvInput))
  lineCount := 0
  var headers []string
  var result bytes.Buffer
  var item bytes.Buffer
  result.WriteString("[")

  for {
    // read just one record, but we could ReadAll() as well
    record, err := csvReader.Read()

    if err == io.EOF {
      result.Truncate(int(len(result.String())-1))
      result.WriteString("]")
      break
    } else if err != nil {
      fmt.Println("Error:", err)
      return "", err
    }

    if lineCount == 0 {
      headers = record[:]
      lineCount += 1
    } else
    {
      item.WriteString("{")
      for i := 0; i < len(headers); i++ {
        item.WriteString("\"" + headers[i] + "\": \"" + record[i] + "\"")
        if i == (len(headers)-1) {
          item.WriteString("}")
        } else {
          item.WriteString(",")
        }
      }
      result.WriteString(item.String() + ",")
      item.Reset()
      lineCount += 1
    }
  }
  return result.String(), nil
}
