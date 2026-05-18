package buildinfo

import (
	"testing"
	"time"
)

func TestName(t *testing.T) {
	name := Name()
	if name == "" {
		t.Error("Name() should not be empty")
	}
}

func TestVersion(t *testing.T) {
	v := Version()
	if v.String() == "" {
		t.Error("Version() should not be empty")
	}
}

func TestBuildTime(t *testing.T) {
	bt := BuildTime()
	if bt.IsZero() {
		t.Error("BuildTime() should not be zero")
	}
	if bt.After(time.Now()) {
		t.Error("BuildTime() should not be in the future")
	}
}

func TestCommit(t *testing.T) {
	c := Commit()
	if c == "" {
		t.Error("Commit() should not be empty")
	}
}

func TestInfo(t *testing.T) {
	info := Info()
	if info == nil {
		t.Fatal("Info() should not be nil")
	}
	if info.Name == "" {
		t.Error("Info().Name should not be empty")
	}
	if info.OS == "" {
		t.Error("Info().OS should not be empty")
	}
	if info.Architecture == "" {
		t.Error("Info().Architecture should not be empty")
	}
}

func TestOS(t *testing.T) {
	if OS() == "" {
		t.Error("OS() should not be empty")
	}
}

func TestArch(t *testing.T) {
	if Arch() == "" {
		t.Error("Arch() should not be empty")
	}
}

func TestGithubRepo(t *testing.T) {
	if GithubRepo() != "haoli000/tttns" {
		t.Errorf("GithubRepo() = %q, want %q", GithubRepo(), "haoli000/tttns")
	}
}

func TestProjectName(t *testing.T) {
	if ProjectName() != "tttns" {
		t.Errorf("ProjectName() = %q, want %q", ProjectName(), "tttns")
	}
}
