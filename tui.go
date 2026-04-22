package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Dùng mã màu ANSI để PowerShell/CMD không bị rác chữ
var styleTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true) // Màu Cam
var styleSus   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true) // Màu Đỏ

type model struct {
	cursor  int
	choices []string
}

func initialModel() model {
	return model{
		choices: []string{"Diệt tiến trình (Kill PID)", "Xóa tận gốc (Force Delete)", "Thoát"},
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	// Đổi tên thành TER-UNLOCKER-DEL theo ý bạn
	s := styleTitle.Render("--- TER-UNLOCKER-DEL (GO EDITION) ---") + "\n\n"
	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = styleSus.Render(" ") // Icon Nerd Font
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\n(Dùng phím mũi tên để di chuyển, Enter để chọn)\n"
	return s
}

// Hàm xử lý logic xóa có tích hợp Nerd Font icons
func handleForceDelete(path string) {
	// Giải mã biến môi trường %temp%, %appdata%...
	path = os.ExpandEnv(strings.ReplaceAll(path, "%", "$"))
	
	lockers, err := scanLocker(path)
	if err == nil && len(lockers) > 0 {
		fmt.Printf(styleSus.Render("\n  [!] Tìm thấy %d tiến trình đang khóa: %v\n"), len(lockers), lockers)
		fmt.Print("󰆴  Bạn có muốn diệt sạch bọn chúng không? (y/n): ")
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) == "y" {
			for _, pid := range lockers {
				killPID(pid)
			}
		}
	}
	
	err = forceDelete(path)
	if err != nil {
		fmt.Printf("󰅙  Lỗi: %v\n", err)
	} else {
		fmt.Println("󰗠  [+] Đã xóa sạch sẽ!")
	}
}

func runTUI() {
	m := initialModel()
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Lỗi TUI:", err)
		os.Exit(1)
	}

	sel := finalModel.(model)
	switch sel.cursor {
	case 0:
		var pid int32
		fmt.Print("󰍉  Nhập PID: ") // Icon kính lúp Nerd Font
		fmt.Scanln(&pid)
		killPID(pid)
	case 1:
		var path string
		fmt.Print("󱂵  Nhập đường dẫn: ") // Icon folder Nerd Font
		fmt.Scanln(&path)
		handleForceDelete(path)
	case 2:
		fmt.Println("Tạm biệt! ඞ")
	}
}