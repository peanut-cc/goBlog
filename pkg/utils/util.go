package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"reflect"
	"time"

	"github.com/google/uuid"

	"github.com/LyricTian/structs"
)

// StructMapToStruct 结构体映射
func StructMapToStruct(s, ts interface{}) error {
	if !structs.IsStruct(s) || !structs.IsStruct(ts) {
		return nil
	}

	ss, tss := structs.New(s), structs.New(ts)
	for _, tfield := range tss.Fields() {
		if !tfield.IsExported() {
			continue
		}

		if tfield.IsEmbedded() && tfield.Kind() == reflect.Struct {
			for _, tefield := range tfield.Fields() {
				if f, ok := ss.FieldOk(tefield.Name()); ok {
					tefield.Set2(f.Value())
				}
			}
		} else if f, ok := ss.FieldOk(tfield.Name()); ok {
			tfield.Set2(f.Value())
		}
	}

	return nil
}

// 读取目录
func ReadDir(dir string, filter func(name string) bool) (files []string) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, fi := range fis {
		if filter(fi.Name()) {
			continue
		}
		if fi.IsDir() {
			files = append(files, ReadDir(path.Join(dir, fi.Name()), filter)...)
			continue
		}
		files = append(files, path.Join(dir, fi.Name()))
	}
	return
}

// encrypt password
func EncryptPasswd(name, pass string) string {
	salt := "%$@w*)("
	h := sha256.New()
	io.WriteString(h, name)
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// UUID Define alias
type UUID = uuid.UUID

// NewUUID Create uuid
func NewUUID() (UUID, error) {
	return uuid.NewRandom()
}

// MustUUID Create uuid(Throw panic if something goes wrong)
func MustUUID() string {
	v, err := NewUUID()
	if err != nil {
		panic(err)
	}
	return v.String()
}

func CheckDate(date string) time.Time {
	if t, err := time.ParseInLocation("2006-01-02 15:04", date, time.Local);err == nil {
		return t
	}
	return time.Now()
}
