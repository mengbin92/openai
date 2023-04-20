package openai

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/pkg/errors"
)

type FormFactory interface {
	CreateFormFile(fieldname string, file *os.File) error
	WriteField(fieldname string, value string) error
	FormDataContentType() string
	Close() error
}

type defaultForm struct {
	writer *multipart.Writer
}

func newDefaultForm(body io.Writer) FormFactory {
	return &defaultForm{
		writer: multipart.NewWriter(body),
	}
}

func (f *defaultForm) CreateFormFile(fieldname string, file *os.File) error {
	fWriter, err := f.writer.CreateFormFile(fieldname, file.Name())
	if err != nil {
		return errors.Wrap(err, "CreateFormFile error")
	}

	_, err = io.Copy(fWriter, file)
	if err != nil {
		return errors.Wrap(err, "copy file from local error")
	}
	return nil
}

func (f *defaultForm) WriteField(fieldname string, value string) error {
	err := f.writer.WriteField(fieldname, value)
	if err != nil {
		return errors.Wrap(err, "WriteField error")
	}
	return nil
}

func (f *defaultForm) FormDataContentType() string {
	return f.writer.FormDataContentType()
}

func (f *defaultForm)Close()error{
	return f.writer.Close()
}