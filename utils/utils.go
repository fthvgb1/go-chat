package utils

import "reflect"

func IsContain(elem interface{}, arr interface{}) int {
	r := reflect.ValueOf(arr)
	k := r.Kind()
	if k == reflect.Slice || k == reflect.Array {
		l := r.Len()
		for i := 0; i < l; i++ {
			if r.Index(i).Interface() == elem {
				return i
			}
		}
	}
	return -1
}
