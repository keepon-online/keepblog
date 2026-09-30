package monitor

import "testing"

func TestMountPriority(t *testing.T) {
	cases := map[string]int{
		"/":            0,
		"/app":         1,
		"/data":        1,
		"/home":        1,
		"/var/www":     3,
		"/app/config":  5,
		"/app/logs":    5,
		"/boot":        10,
		"/srv/data":    10,
		"/etc/hosts":   10,
	}
	for mountpoint, want := range cases {
		if got := mountPriority(mountpoint); got != want {
			t.Errorf("mountPriority(%q) = %d, want %d", mountpoint, got, want)
		}
	}
}

func TestIsTempPartition(t *testing.T) {
	// Docker 注入容器的 bind mount 必须被过滤
	for _, mountpoint := range []string{"/etc/hosts", "/etc/hostname", "/etc/resolv.conf"} {
		if !(*Handler)(nil).isTempPartition(mountpoint) {
			t.Errorf("isTempPartition(%q) = false, want true", mountpoint)
		}
	}
	// 实际数据目录不能被过滤
	for _, mountpoint := range []string{"/", "/app", "/app/config", "/app/data", "/app/logs", "/home"} {
		if (*Handler)(nil).isTempPartition(mountpoint) {
			t.Errorf("isTempPartition(%q) = true, want false", mountpoint)
		}
	}
}

func TestIsVirtualFilesystem(t *testing.T) {
	// 容器根目录的 fstype 是 overlay，必须保留，交给挂载点/总量去重处理
	if (*Handler)(nil).isVirtualFilesystem("overlay") {
		t.Error("isVirtualFilesystem(overlay) = true, want false")
	}
	for _, fstype := range []string{"tmpfs", "proc", "cgroup2", "squashfs"} {
		if !(*Handler)(nil).isVirtualFilesystem(fstype) {
			t.Errorf("isVirtualFilesystem(%q) = false, want true", fstype)
		}
	}
}
