package stub

import (
	"fmt"

	"bingo/pkg/encode"
	"bingo/pkg/log"
)


func check(d interface{}) error {
	if d == nil {
		return log.Error("stub is nil, make sure it had been opened[method On()] in Setup()")
	}
	return nil
}
func parseData(data []byte, obj interface{}) error {
	if obj != nil {
		if err := encode.SmartDecoder(data, obj); err != nil {
			return err
		}
		if obj == nil {
			return fmt.Errorf("target is nil")
		}
	}
	return nil
}
