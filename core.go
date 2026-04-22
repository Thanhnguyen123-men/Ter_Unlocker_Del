package main

import (
	"context"
	"fmt"
	"os"
	"strings" // Đảm bảo có strings để xử lý đường dẫn
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// Khôi phục từ core.go:43
func killPID(pid int32) {
    ctx := context.Background()
    p, err := process.NewProcessWithContext(ctx, pid)
    if err != nil {
        fmt.Printf("Lỗi: Không tìm thấy tiến trình %d\n", pid)
        return
    }
    err = p.KillWithContext(ctx)
    if err != nil {
        // Sửa ở đây: Thêm biến , err vào sau pid
        fmt.Printf("Lỗi khi diệt %d: %v\n", pid, err) 
        return
    }
    fmt.Printf("Đã diệt thành công PID: %d\n", pid)
}

// Khôi phục từ core.go:54
func forceDelete(path string) error {
	// Giải mã các biến môi trường như %temp%, %appdata%
	// strings.ReplaceAll giúp chuyển % thành $ để os.ExpandEnv hiểu được
	cleanPath := os.ExpandEnv(strings.ReplaceAll(path, "%", "$"))

	for i := 0; i < 5; i++ {
		err := os.RemoveAll(cleanPath)
		if err == nil {
			// Check lại xem nó thực sự bay màu chưa
			if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
				return nil
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("không thể xóa: %s (vẫn còn tồn tại)", cleanPath)
}
// scanLocker tìm tất cả các PID đang chiếm giữ targetPath
func scanLocker(targetPath string) ([]int32, error) {
	var results []int32
	targetPath = strings.ToLower(targetPath)

	// Lấy tất cả process đang chạy
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	for _, p := range procs {
		// Bỏ qua các tiến trình hệ thống nhạy cảm (giống bản Python của bạn)
		name, _ := p.Name()
		if name == "System" || name == "Registry" || name == "Idle" {
			continue
		}

		// Lấy danh sách file đang mở bởi process này
		openFiles, err := p.OpenFiles()
		if err != nil {
			continue // Thường là lỗi Access Denied với process hệ thống
		}

		for _, f := range openFiles {
			if strings.Contains(strings.ToLower(f.Path), targetPath) {
				results = append(results, p.Pid)
				break // Tìm thấy rồi thì đổi sang process tiếp theo
			}
		}
	}
	return results, nil
}