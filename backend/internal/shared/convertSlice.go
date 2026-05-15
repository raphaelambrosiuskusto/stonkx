package shared

import (
	"fmt"
	"strconv"
)

func Convert(src []string) ([]int, error) {
	idlist := make([]int, 0, len(src))
	for _, id := range src {
		res, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("converting %v to int: %w", id, err)
		}
		idlist = append(idlist, res)
	}
	return idlist, nil
}
