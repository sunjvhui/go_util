package go_util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// 字符substr在str中第n次出现的位置，不存在则返回-1。
func IndexN(str, substr string, n int) int {
	var num = 0
	return strings.IndexFunc(str, func(r rune) bool {
		if string(r) == substr {
			num++
		}
		if num == n {
			return true
		}
		return false
	})
}

// 返回字符串str第一次出现indexS 和最后一次出现lastIndexS 中间的字符, 任意一个不存在则返回空
/****
 str = "qweqe{123123123}qaweqwe",
 indexS="{" ,
 lastIndexS="}"
 return {123123123}
*****/
func InterceptString(str, indexS, lastIndexS string) string {
	var newString string
	index := strings.Index(str, indexS)
	lastIndex := strings.LastIndex(str, lastIndexS)
	if index != -1 && lastIndex != -1 {
		newString = str[index : lastIndex+1]
	}
	return newString
}

// search的字符串是否属于orgin中的字符串，字符串分隔使用sep
// 如origin=aa bb cc dd,用空格分隔，则seach=aa bb返回true,search=aa ee返回false
func ContainString(origin string, search string, sep string) bool {
	originSlice := strings.Split(origin, sep)
	searchSlice := strings.Split(search, sep)
	for _, vf := range searchSlice {
		isContain := false
		for _, vp := range originSlice {
			if vf == vp {
				isContain = true
				break
			}
		}
		if !isContain {
			return false
		}
	}
	return true
}

// orgin中的字符串不包含任何search中的字符串，search分隔使用sep
// 如origin=aabbccdd,seach=aa bb返回false,search=ee返回true
func StringNotContain(s string, subs string, subsSep string) bool {
	subSlice := strings.Split(subs, subsSep)
	for _, sub := range subSlice {
		if strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

// map int32 transform to string by separator
func Int32Map2String(dataMap map[int32]struct{}, separator string) string {
	if len(dataMap) == 0 {
		return ""
	}
	var str string
	for key, _ := range dataMap {
		str += fmt.Sprintf("%d%s", key, separator)
	}
	return strings.TrimRight(str, separator)
}

// map string transform to string by separator
func StringMap2String(dataMap map[string]struct{}, separator string) string {
	if len(dataMap) == 0 {
		return ""
	}
	var str string
	for key, _ := range dataMap {
		str += fmt.Sprintf("%s%s", key, separator)
	}
	return strings.TrimRight(str, separator)
}

// []string 去重
func ReDuplicateString(input []string) []string {
	result := make([]string, 0)
	tempMap := make(map[string]struct{}, len(input))
	for _, e := range input {
		l := len(tempMap)
		tempMap[e] = struct{}{}
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

// 字符串包含，不区分大小写
func StringFuzzyContains(s, substr string) bool {
	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)
	return strings.Contains(sLower, substrLower)
}

// 反转字符串
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// 反转号码以及替换特殊字符 (读取时，展示给页面)
// 先转为 ×#再反转
func ReadReverseNumber(str string) string {
	return ReverseString(EscapeToCharacter(str))
}

// 反转号码以及替换特殊字符 (写入时，对db操作)
// 先反转，在转译为_n
func WriteReverseNumber(str string) string {
	return EscapeToNum(ReverseString(str))
}

// 给字符串添加转译符号
func EscapeToSelect(str string) string {
	var escapeSrc string
	for _, v := range str {
		// 95=_ 39=' 92=\
		if v == 95 || v == 39 || v == 92 {
			escapeSrc += `\`
		}
		escapeSrc += string(v)
	}
	return escapeSrc
}

// _num转特殊字符
func EscapeToCharacter(str string) string {
	var escapeSrc string
	var character rune
	for _, v := range str {
		if v == 95 {
			character = v
			continue
		} else if character == 95 {
			character = 0
			switch v {
			case 49:
				escapeSrc += "#"
			case 50:
				escapeSrc += "("
			case 51:
				escapeSrc += ")"
			case 52:
				escapeSrc += "*"
			case 53:
				escapeSrc += "+"
			case 54:
				escapeSrc += "-"
			case 55:
				escapeSrc += "."
			default:
				escapeSrc += string(v)
			}
			continue
		}
		escapeSrc += string(v)
	}
	return escapeSrc
}

// 特殊字符转_num
func EscapeToNum(str string) string {
	// #==_1 (==_2 )==_3 *==_4 +==_5 -==_6 .==_7
	// #==35 (==40 )==41 *==42 +==43 -==45 .==46
	var escapeSrc string
	for _, v := range str {
		switch v {
		case 35:
			escapeSrc += "_1"
		case 40:
			escapeSrc += "_2"
		case 41:
			escapeSrc += "_3"
		case 42:
			escapeSrc += "_4"
		case 43:
			escapeSrc += "_5"
		case 45:
			escapeSrc += "_6"
		case 46:
			escapeSrc += "_7"
		default:
			escapeSrc += string(v)
		}
	}
	return escapeSrc
}

func GetDialPlate(str string) string {
	str = strings.ToLower(str)
	var dialPlateMap = map[string]string{
		"a": "2",
		"b": "2",
		"c": "2",
		"d": "3",
		"e": "3",
		"f": "3",
		"g": "4",
		"h": "4",
		"i": "4",
		"j": "5",
		"k": "5",
		"l": "5",
		"m": "6",
		"n": "6",
		"o": "6",
		"p": "7",
		"q": "7",
		"r": "7",
		"s": "7",
		"t": "8",
		"u": "8",
		"v": "8",
		"w": "9",
		"x": "9",
		"y": "9",
		"z": "9",
	}
	var DialPlateNum string
	for _, i2 := range []rune(str) {
		DialPlateNum += dialPlateMap[string(i2)]
	}
	return DialPlateNum
}

/**
["q","w","e","r"]
q,w,e,r
**/
func SliceToString(strSlice []string) string {
	var numString string
	var i int
	for _, v := range strSlice {
		if v == "" {
			continue
		}
		if i == 0 {
			numString += fmt.Sprintf("'%s'", v)
		} else {
			numString += fmt.Sprintf(",'%s'", v)
		}
		i++
	}
	return numString
}

/*
str: "a=1,b=2 c=3\n"
key: "b="
eof: "\n, "
return "2"
*/
func GetStringKeyValue(str string, key string, eof string) string {
	fn := func(c rune) bool {
		return strings.ContainsRune(eof, c)
	}
	strArr := strings.FieldsFunc(str, fn)
	prefix := len(key)
	for _, v := range strArr {
		if strings.HasPrefix(v, key) {
			return v[prefix:]
		}
	}
	return ""
}

// 大小写转换，空格转下划线
func ChangeContinent(data string) string {
	str := strings.Replace(data, " ", "_", -1)
	return strings.ToLower(str)
}

func IsStrContains(i interface{}, str string) bool {
	strSlice := make([]string, 0)
	switch i.(type) {
	case []string:
		strSlice = i.([]string)
	case string:
		s := i.(string)
		if len(s) <= 2 {
			return false
		}
		if err := json.Unmarshal([]byte(s), &strSlice); err != nil {
			//
			return false
		}
	}
	for _, v := range strSlice {
		if str == v {
			return true
		}
	}
	return false
}

func RemoveDuplicateStringSliceItem(stringSlice []string) []string {
	resultSlice := []string{}
	tmpMap := map[string]bool{}
	for _, v := range stringSlice {
		oldLen := len(tmpMap)
		tmpMap[v] = true
		if len(tmpMap) != oldLen {
			resultSlice = append(resultSlice, v)
		}
	}
	return resultSlice
}

// 校验 加号+数字+长度 true为合法  false为不合法oud
func CheckDialCode(data string, length int) bool {
	reg := regexp.MustCompile(`[^[:digit:]\+]`)
	result := reg.FindAllString(data, -1)
	if len(result) > 0 || len(data) < 1 {
		return false
	}

	if len(data) > length {
		return false
	}
	return true
}

// 判断a为空取b值，a不为空取a值
func GetExprIfEmpty(a, b string) string {
	if len(a) != 0 {
		return a
	} else {
		return b
	}
}

// map中是否包含slice中的值
func FindStrInMap(data map[string]interface{}, res []string) bool {
	for _, v := range res {
		if _, ok := data[v]; ok {
			return true
		}
	}
	return false
}

/*
	string转int32
*/
func StringToInt32(str string) (int32, error) {
	int10, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, err
	} else {
		return int32(int10), nil
	}
}

func InsertPrefix(str string, prefix string) string {
	return prefix + str
}

/*
	复制新map
*/
func CopyNewMap(vars map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{}, len(vars))
	for k, v := range vars {
		newMap[k] = v
	}
	return newMap
}

/*
	保持hold、静音mute、响铃状态转换
*/
func GetStatusByStr(status string) (holdStatus, muteStatus, ringStatus, askPop string) {
	// status : hold-mute-ring
	str := strings.Split(status, "-")
	holdStatus = "HoldOff"
	muteStatus = "MuteOff"
	ringStatus = "RingOff"
	askPop = "0"
	if len(str) != 4 {
		return
	}
	if str[0] == "1" {
		holdStatus = "HoldOn"
	}
	if str[1] == "1" {
		muteStatus = "MuteOn"
	}
	if str[2] == "1" {
		ringStatus = "RingOn"
	}
	askPop = str[3]
	return
}

// ContactUriNormalize 去除contactUrl中多余的字段，返回标准的url
func ContactUriNormalize(contactUri string) string {
	// 模拟话机接听时，发送过来的Channel带有板块 DAHDI/1-1 ，我们需要的是DAHDI/1, -1是随机生成的不需要
	dAHDIIndex := strings.Index(contactUri, "DAHDI")
	if dAHDIIndex != -1 {
		index2 := strings.Index(contactUri, "-")
		if index2 != -1 {
			contactUri = contactUri[:index2]
		}
	} else {
		//contactUri, webclient的url后面会加一个";ob". 需要兼容去掉
		astIndex := strings.LastIndex(contactUri, ";ob")
		if astIndex != -1 {
			contactUri = contactUri[:astIndex]
		}
		// contactUri 被叫时url前面会出现 "分机号/"，需要过滤掉
		index := strings.Index(contactUri, "/")
		if index != -1 {
			contactUri = contactUri[index+1:]
		}
	}
	return contactUri
}

/*
	字符串补0
	param:
		str 目标对象
		resultLen 补充位数
		reverse true为前置补0，false为后置补0
*/
func ZeroFillByStr(str string, resultLen int, reverse bool) string {
	if len(str) > resultLen || resultLen <= 0 {
		return str
	}
	if reverse {
		return fmt.Sprintf("%0*s", resultLen, str) //不足前置补零
	}
	result := str
	for i := 0; i < resultLen-len(str); i++ {
		result += "0"
	}
	return result
}

func SortNumString(data []string) {
	sort.Slice(data, func(i, j int) bool {
		numA, _ := strconv.Atoi(data[i])
		numB, _ := strconv.Atoi(data[j])
		return numA < numB
	})
}
