package emc

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"os"
	"strings"

	"github.com/saracen/walker"
	classfileParser "github.com/wreulicke/go-java-class-parser/classfile"
)

func CountClassInStandardLibrary(version int) int64 {
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
			if findLambda && strings.HasPrefix(pathname, ".class") {
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
			if findLambda && strings.HasPrefix(f.Name, ".class") {
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
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return -1, err
	}
	var count int64 = 0
	classFile := classfileParser.Parse(bs)
	for _, a := range classFile.Attributes() {
		if v, ok := a.(*classfileParser.BootstrapMethodsAttribute); ok {
			for _, m := range v.BootstrapMethods {
				className := m.ClassName()
				methodName, _ := m.NameAndDescriptor()
				if className == "java/lang/invoke/LambdaMetafactory" && methodName == "metafactory" {
					count++
				}
			}
		}
	}
	return count, nil
}
