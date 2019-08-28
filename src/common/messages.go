package common

type Messages struct {
}

func (*Messages) LINK_FILE() int {
	return 10000
}

func (*Messages) LINK_FILE_TEXT() string {
	return "ln"
}

func (*Messages) JSTAT() int {
	return 10001
}

func (*Messages) JSTAT_TEXT() string {
	return "jstat"
}

func (*Messages) UPLOAD_FILE() int {
	return 10002
}

func (*Messages) UPLOAD_FILE_TEXT() string {
	return "upload"
}

func (*Messages) UPLOAD_PACKAGE() int {
	return 10003
}

func (*Messages) UPLOAD_PACKAGE_TEXT() string {
	return "uploadpackage"
}
