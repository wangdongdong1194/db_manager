package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// StrToIntArr 逗号分割字符串转 int 数组
func StrToIntArr(s string) ([]int, error) {
	// 1. 按逗号切割
	strArr := strings.Split(s, ",")

	// 2. 准备 int 数组
	res := make([]int, 0, len(strArr))

	// 3. 逐个转换
	for _, str := range strArr {
		// 去除空格（处理 "1, 2, 3" 这种带空格的情况）
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}

		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("转换失败: %w", err)
		}
		res = append(res, num)
	}

	return res, nil
}