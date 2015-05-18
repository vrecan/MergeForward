package merge

// import (
// 	"errors"
// 	YAML "gopkg.in/yaml.v2"
// 	"io/ioutil"
// 	"reflect"
// )

// var (
// 	MaxDepth = 32
// )

// //Simple merge into a map for all values
// func Merge(src, dst map[interface{}]interface{}) (map[interface{}]interface{}, error) {
// 	return merge(src, dst, 0)
// }

// func merge(src, dst map[interface{}]interface{}, depth int) (map[interface{}]interface{}, error) {
// 	if depth > MaxDepth {
// 		return dst, errors.New("Recursed too deep.")
// 	}
// 	var err error
// 	for key, srcVal := range src {
// 		if dstVal, ok := dst[key]; ok {
// 			srcMap, srcMapOk := mapify(srcVal)
// 			dstMap, dstMapOk := mapify(dstVal)
// 			if srcMapOk && dstMapOk {
// 				srcVal, err = merge(srcMap, dstMap, depth+1)
// 				if nil != err {
// 					return dst, err
// 				}
// 			}
// 		}
// 		dst[key] = srcVal
// 	}
// 	return dst, err
// }

// func Mapify(i interface{}) (map[interface{}]interface{}, bool) {
// 	return mapify(i)
// }

// func mapify(i interface{}) (map[interface{}]interface{}, bool) {
// 	value := reflect.ValueOf(i)
// 	if value.Kind() == reflect.Map {
// 		m := map[interface{}]interface{}{}
// 		for _, k := range value.MapKeys() {
// 			if k.Kind() == reflect.String {
// 				m[k] = value.MapIndex(k).Interface()
// 			} else {
// 				kResult, _ := mapify(k)
// 				m[k.Interface()] = kResult

// 			}
// 		}
// 		return m, true
// 	}
// 	return map[interface{}]interface{}{}, false
// }

// //Write yaml from map
// func MarshalToYaml(m map[interface{}]interface{}, location string) (err error) {
// 	b, err := YAML.Marshal(&m)
// 	if nil != err {
// 		return err
// 	}
// 	err = ioutil.WriteFile(location, b, 0644)
// 	return err
// }
