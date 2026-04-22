package main

import (
	"fmt"
	"os"
)

// Khôi phục từ cli.go:13
func menu() {
	fmt.Fprintln(os.Stdout, "\n==============================")
	fmt.Fprintln(os.Stdout, "       AMONG SUS TOOL         ")
	fmt.Fprintln(os.Stdout, "==============================")
	fmt.Fprintln(os.Stdout, "1. Diệt tiến trình (Kill PID)")
	fmt.Fprintln(os.Stdout, "2. Xóa tận gốc (Force Delete)")
	fmt.Fprintln(os.Stdout, "0. Thoát")
	fmt.Fprint(os.Stdout, "Lựa chọn của bạn: ")
}