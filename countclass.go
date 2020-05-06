package emc

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"log"
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
			log.Println("found", pathname)
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

func CountClassFileWithPath(path string, findLambda bool) (int64, error) {
	if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "https://") {
		resp, err := http.DefaultClient.Get(path)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		b := &bytes.Buffer{}
		_, err = io.Copy(b, resp.Body)
		if err != nil {
			return 0, err
		}
		r := bytes.NewReader(b.Bytes())
		zr, err := zip.NewReader(r, r.Size())
		if err != nil {
			return 0, err
		}
		return CountClassFileInJar(zr, "/", findLambda)
	}
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	fi, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return CountClassFile(file, fi, findLambda)
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
			log.Println("found", pathPrefix+f.Name)
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
	if count > 0 {
		log.Printf("found lambda %d", count)
	}
	return count, nil
}
