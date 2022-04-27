package emc

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"os"
	"strings"

	"github.com/cli/safeexec"
	"github.com/paketo-buildpacks/libjvm/count"
	"github.com/saracen/walker"
	parser "github.com/wreulicke/classfile-parser"
)

func CountClassInStandardLibrary(version int) int64 {
	path, err := safeexec.LookPath("java")
	if err == nil {
		p, err := filepath.Abs(filepath.Join(filepath.Dir(path), "..", "lib", "modules"))
		if err == nil {
			v, err := count.ModuleClasses(p)
			if err == nil {
				return int64(v)
			}
		}
	}

	if version < 9 { // OpenJDK 64-Bit Server VM (AdoptOpenJDK)(build 25.232-b09, mixed mode) based
		return 30645
	}
	if version == 9 {
		return 33615 // OpenJDK 64-Bit Server VM (build 9.0.4+11, mixed mode)
	}
	return 30554 // OpenJDK 64-Bit Server VM AdoptOpenJDK (build 11.0.6+10, mixed mode)
}

func CountClassFileInDir(dir *os.File, findLambda bool) (int64, error) {
	var count int64 = 0
	walkfn := func(pathname string, fi os.FileInfo) error {
		if fi.IsDir() {
			return nil
		}
		if shouldCount(pathname) {
			count++
			if findLambda && strings.HasSuffix(pathname, ".class") {
				c, err := countLambdaClass(os.Open(pathname))
				if err == nil {
					count += c
				}
			}
			return nil
		}
		f, err := os.Open(pathname)
		if err != nil {
			return err
		}
		defer f.Close()
		r, err := CountClassFile(f, fi, findLambda)
		if err == nil {
			count += r
		}
		return nil
	}
	err := walker.Walk(dir.Name(), walkfn)
	return count, err
}

func CountClassFile(file *os.File, fi os.FileInfo, findLambda bool) (int64, error) {
	if fi.IsDir() {
		return CountClassFileInDir(file, findLambda)
	} else if strings.HasSuffix(file.Name(), ".jar") || strings.HasSuffix(file.Name(), ".jmods") {
		z, err := zip.NewReader(file, fi.Size())
		if err != nil {
			return -1, fmt.Errorf("cannot create zip reader. file=%s err=%v", file.Name(), err)
		}
		return CountClassFileInJar(z, file.Name()+"!", findLambda)
	}
	return 0, fmt.Errorf("file is not directory or jar. file=%s", file.Name())
}

func CountClassFileInJar(z *zip.Reader, pathPrefix string, findLambda bool) (int64, error) {
	var count int64 = 0
	for _, f := range z.File {
		if shouldCount(f.Name) {
			count++
			if findLambda && strings.HasSuffix(f.Name, ".class") {
				c, err := countLambdaClass(f.Open())
				if err == nil {
					count += c
				}
			}
		}
		if strings.HasSuffix(f.Name, ".jar") || strings.HasSuffix(f.Name, ".jmods") {
			r, err := f.Open()
			if err != nil {
				return -1, fmt.Errorf("cannot open %s err=%v", f.Name, err)
			}
			bs, err := ioutil.ReadAll(r)
			if err != nil {
				return -1, fmt.Errorf("cannot read %s err=%v", f.Name, err)
			}
			z, err := zip.NewReader(bytes.NewReader(bs), f.FileInfo().Size())
			if err != nil {
				return -1, fmt.Errorf("cannot create zip reader %s err=%v", f.Name, err)
			}
			c, err := CountClassFileInJar(z, pathPrefix+f.Name+"!", findLambda)
			if err != nil {
				return -1, err
			}
			count += c
		}
	}
	return count, nil
}

func shouldCount(fileName string) bool {
	return strings.HasSuffix(fileName, ".class") || strings.HasSuffix(fileName, ".groovy")
}

func countLambdaClass(r io.ReadCloser, err error) (int64, error) {
	if err != nil {
		return -1, err
	}
	defer r.Close()
	var count int64 = 0
	classFile, err := parser.New(r).Parse()
	if err != nil {
		return -1, err
	}
	constants := classFile.ConstantPool.Constants
	for _, a := range classFile.Attributes {
		if v, ok := a.(*parser.AttributeBootstrapMethods); ok {
			for _, m := range v.BootstrapMethods {
				methodRefIndex := m.BootstrapMethodRef
				methodHandle, ok := constants[methodRefIndex-1].(*parser.ConstantMethodHandle)
				if !ok {
					continue
				}
				methodRef, ok := constants[methodHandle.ReferenceIndex-1].(*parser.ConstantMethodref)
				if !ok {
					continue
				}
				class, ok := constants[methodRef.ClassIndex-1].(*parser.ConstantClass)
				if !ok {
					continue
				}
				nameAndType, ok := constants[methodRef.NameAndTypeIndex-1].(*parser.ConstantNameAndType)
				if !ok {
					continue
				}
				className := classFile.ConstantPool.LookupUtf8(class.NameIndex)
				methodName := classFile.ConstantPool.LookupUtf8(nameAndType.NameIndex)
				if className == nil || methodName == nil {
					continue
				}
				if className.String() == "java/lang/invoke/LambdaMetafactory" &&
					methodName.String() == "metafactory" {
					count++
				}
			}
		}
	}
	return count, nil
}
