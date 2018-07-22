package util

import (
	"reflect"
)

func CanFlattenTo(obj map[string]interface{}, refType interface{}) (bRet bool) {
	t0, ok := refType.(reflect.Type)
	if !ok {
		t0 = reflect.TypeOf(refType)
	}

	t := t0
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	if t.NumField() < len(obj) {
		//DBG("reftype numfield is lesser than ReducerResult, reftype=", t.Name(), "; reducerresult=", obj)
		return
	}

	for k, v := range obj {
		if nil == v {
			continue
		}
		ft, has := t0.FieldByName(k)
		if !has {
			return
		}
		vt := reflect.TypeOf(v)
		if !ft.Type.AssignableTo(vt) && !canAssignFields(ft.Type, vt) {
			//DBG("canflatten was not assignable, ft.name=", ft.Name, "; ft typename=", ft.Type.Name(), "; v=", v)
			//DBG("ft.pkg=", ft.PkgPath)
			return
		}
	}

	return true
}

func FlattenToType(obj map[string]interface{}, refType interface{}) interface{} {
	t, ok := refType.(reflect.Type)
	if !ok {
		t = reflect.TypeOf(refType)
	}

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	out := reflect.New(t).Elem()

	for k, v := range obj {
		if nil == v {
			continue
		}
		fvout := out.FieldByName(k)
		fv := reflect.ValueOf(v)

		if fv.Type().AssignableTo(fvout.Type()) {
			fvout.Set(fv)
		} else {
			copyProps(fvout, fv)
		}
	}

	return out.Interface()
}

// copies assignable fields from b to a
func copyProps(a, b reflect.Value) {
	if b.Kind() == reflect.Ptr && b.IsNil() {
		return
	}
	for b.Type().Kind() == reflect.Ptr {
		b = b.Elem()
	}
	at, av := a.Type(), a
	bt, bv := b.Type(), b
	DBGf("copyProps a:%v, b:%v", at, bt)
	if bt.AssignableTo(at) {
		av.Set(bv)
		return
	}
	if (at.Kind() != reflect.Struct) || (bt.Kind() != reflect.Struct) {
		return
	}

	for i := 0; i < at.NumField(); i++ {
		fat := at.Field(i)
		_, has := bt.FieldByName(fat.Name)
		if !has {
			continue
		}

		bvfield := bv.FieldByName(fat.Name)

		for bvfield.Kind() == reflect.Ptr {
			bvfield = bvfield.Elem()
		}

		if bvfield.Type().AssignableTo(fat.Type) {
			av.Field(i).Set(bvfield)
		}
	}
}

func canAssignFields(a, b reflect.Type) (bRet bool) {
	for b.Kind() == reflect.Ptr {
		b = b.Elem()
	}

	if b.AssignableTo(a) {
		return true
	}

	DBGf("canAssignFields a:%v, b:%v", a, b)
	if a.NumField() < b.NumField() {
		return
	}
	for i := 0; i < b.NumField(); i++ {
		fb := b.Field(i)
		fa, has := a.FieldByName(fb.Name)
		if !has {
			return
		}
		if !fb.Type.AssignableTo(fa.Type) {
			return
		}
	}

	return true
}
