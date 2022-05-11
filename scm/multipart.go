package scm

import "mime/multipart"

func (mw *MultipartWriter) Write(f, v string) {
	if mw.Error == nil {
		return
	}
	if v == "" {
		return
	}
	mw.Error = mw.Writer.WriteField(f, v)
}

func (mw *MultipartWriter) Close() {
	mw.Writer.Close()
}

func (mw *MultipartWriter) FormDataContentType() string {
	return mw.Writer.FormDataContentType()
}

type MultipartWriter struct {
	Writer *multipart.Writer
	Error  error
}
