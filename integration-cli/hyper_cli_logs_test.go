package main

import (
	//"encoding/json"
	"fmt"
	//"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/pkg/integration/checker"
	"github.com/docker/docker/pkg/jsonlog"
	"github.com/go-check/check"
)

//TODO: get exited container log
// This used to work, it test a log of PageSize-1 (gh#4851)
/*func (s *DockerSuite) TestCliLogsContainerSmallerThanPage(c *check.C) {
	testRequires(c, DaemonIsLinux)
	testLen := 32767
	out, _ := dockerCmd(c, "run", "-d", "busybox", "sh", "-c", fmt.Sprintf("for i in $(seq 1 %d); do echo -n =; done; echo", testLen))

	id := strings.TrimSpace(out)
	dockerCmd(c, "stop", id)

	out, _ = dockerCmd(c, "logs", id)

	c.Assert(out, checker.HasLen, testLen+1)
}*/

//TODO: get exited container log
// Regression test: When going over the PageSize, it used to panic (gh#4851)
/*func (s *DockerSuite) TestCliLogsContainerBiggerThanPage(c *check.C) {
	testRequires(c, DaemonIsLinux)
	testLen := 32768
	out, _ := dockerCmd(c, "run", "-d", "busybox", "sh", "-c", fmt.Sprintf("for i in $(seq 1 %d); do echo -n =; done; echo", testLen))

	id := strings.TrimSpace(out)
	dockerCmd(c, "stop", id)

	out, _ = dockerCmd(c, "logs", id)

	c.Assert(out, checker.HasLen, testLen+1)
}*/

//TODO: get exited container log
/*// Regression test: When going much over the PageSize, it used to block (gh#4851)
func (s *DockerSuite) TestCliLogsContainerMuchBiggerThanPage(c *check.C) {
	testRequires(c, DaemonIsLinux)
	testLen := 33000
	out, _ := dockerCmd(c, "run", "-d", "busybox", "sh", "-c", fmt.Sprintf("for i in $(seq 1 %d); do echo -n =; done; echo", testLen))

	id := strings.TrimSpace(out)
	dockerCmd(c, "stop", id)

	out, _ = dockerCmd(c, "logs", id)

	c.Assert(out, checker.HasLen, testLen+1)
}*/

//TODO: get exited container log
func (s *DockerSuite) TestCliLogsTimestamps(c *check.C) {
	printTestCaseName()
	defer printTestDuration(time.Now())
	testRequires(c, DaemonIsLinux)

	testLen := 100
	pullImageIfNotExist("busybox")
	out, _ := dockerCmd(c, "run", "-d", "busybox", "sh", "-c", fmt.Sprintf("for i in $(seq 1 %d); do echo =; done;", testLen))
	time.Sleep(5 * time.Second)

	id := strings.TrimSpace(out)
	dockerCmd(c, "stop", id)
	time.Sleep(5 * time.Second)

	out, _ = dockerCmd(c, "logs", "-t", id)

	lines := strings.Split(out, "\n")

	c.Assert(lines, checker.HasLen, testLen+1)

	ts := regexp.MustCompile(`^.* `)

	for _, l := range lines {
		if l != "" {
			_, err := time.Parse(jsonlog.RFC3339NanoFixed+" ", ts.FindString(l))
			c.Assert(err, checker.IsNil, check.Commentf("Failed to parse timestamp from %v", l))
			// ensure we have padded 0's
			c.Assert(l[29], checker.Equals, uint8('Z'))
		}
	}
}

//TODO: get exited container log
func (s *DockerSuite) TestCliLogsSeparateStderrBasic(c *check.C) {
	printTestCaseName()
	defer printTestDuration(time.Now())
	testRequires(c, DaemonIsLinux)

	pullImageIfNotExist("busybox")
	msg := "stderr_log"
	out, _ := dockerCmd(c, "run", "-d", "busybox", "sh", "-c", fmt.Sprintf("echo %s 1>&2", msg))
	time.Sleep(5 * time.Second)

	id := strings.TrimSpace(out)
	dockerCmd(c, "stop", id)
	time.Sleep(5 * time.Second)

	stdout, stderr, _ := dockerCmdWithStdoutStderr(c, "logs", id)

	c.Assert(stdout, checker.Equals, "")

	stderr = strings.TrimSpace(stderr)

	c.Assert(stderr, checker.Equals, msg)
}

//TODO: get exited container log
func (s *DockerSuite) TestCliLogsStderrInStdout(c *check.C) {
	printTestCaseName()
	defer printTestDuration(time.Now())
	testRequires(c, DaemonIsLinux)

	pullImageIfNotExist("busybox")
	msg := "stderr_log"
	out, _ := dockerCmd(c, "run", "-d", "-t", "busybox", "sh", "-c", fmt.Sprintf("echo %s 1>&2", msg))
	time.Sleep(5 * time.Second)

	id := strings.TrimSpace(out)
	dockerCmd(c, "stop", id)
	time.Sleep(5 * time.Second)

	stdout, stderr, _ := dockerCmdWithStdoutStderr(c, "logs", id)
	c.Assert(stderr, checker.Equals, "")

	stdout = strings.TrimSpace(stdout)
	c.Assert(stdout, checker.Equals, msg)
}

//TODO: get exited container log
func (s *DockerSuite) TestCliLogsTailBasic(c *check.C) {
	printTestCaseName()
	defer printTestDuration(time.Now())
	testRequires(c, DaemonIsLinux)

	testLen := 100
	pullImageIfNotExist("busybox")
	out, _ := dockerCmd(c, "run", "-d", "busybox", "sh", "-c", fmt.Sprintf("for i in $(seq 1 %d); do echo =; done;", testLen))
	time.Sleep(5 * time.Second)

	id := strings.TrimSpace(out)
	dockerCmd(c, "stop", id)
	time.Sleep(5 * time.Second)

	out, _ = dockerCmd(c, "logs", "--tail", "5", id)

	lines := strings.Split(out, "\n")

	c.Assert(lines, checker.HasLen, 6)

	out, _ = dockerCmd(c, "logs", "--tail", "all", id)

	lines = strings.Split(out, "\n")

	c.Assert(lines, checker.HasLen, testLen+1)

	out, _, _ = dockerCmdWithStdoutStderr(c, "logs", "--tail", "random", id)

	lines = strings.Split(out, "\n")

	c.Assert(lines, checker.HasLen, testLen+1)
}

//TODO: fix #46
func (s *DockerSuite) TestCliLogsSince(c *check.C) {
	printTestCaseName()
	defer printTestDuration(time.Now())
	testRequires(c, DaemonIsLinux)

	pullImageIfNotExist("busybox")
	name := "testlogssince"
	dockerCmd(c, "run", "--name="+name, "-d", "busybox", "/bin/sh", "-c", "for i in $(seq 1 30); do sleep 2; echo log$i; done")
	time.Sleep(10 * time.Second)

	out, _ := dockerCmd(c, "logs", "-t", name)

	log2Line := strings.Split(strings.Split(out, "\n")[1], " ")
	t, err := time.Parse(time.RFC3339Nano, log2Line[0]) // the timestamp log2 is written
	c.Assert(err, checker.IsNil)
	since := t.Unix() + 1 // add 1s so log1 & log2 doesn't show up
	out, _ = dockerCmd(c, "logs", "-t", fmt.Sprintf("--since=%v", since), name)

	// Skip 2 seconds
	unexpected := []string{"log1", "log2"}
	for _, v := range unexpected {
		c.Assert(out, checker.Not(checker.Contains), v, check.Commentf("unexpected log message returned, since=%v", since))
	}

	// Test to make sure a bad since format is caught by the client
	out, _, _ = dockerCmdWithError("logs", "-t", "--since=2006-01-02T15:04:0Z", name)
	c.Assert(out, checker.Contains, "cannot parse \"0Z\" as \"05\"", check.Commentf("bad since format passed to server"))

	// Test with default value specified and parameter omitted
	expected := []string{"log1", "log2", "log3"}
	for _, cmd := range []*exec.Cmd{
		exec.Command(dockerBinary, "--region="+os.Getenv("DOCKER_HOST"), "logs", "-t", name),
		exec.Command(dockerBinary, "--region="+os.Getenv("DOCKER_HOST"), "logs", "-t", "--since=0", name),
	} {
		out, _, err = runCommandWithOutput(cmd)
		c.Assert(err, checker.IsNil, check.Commentf("failed to log container: %s", out))
		for _, v := range expected {
			c.Assert(out, checker.Contains, v)
		}
	}
}

func (s *DockerSuite) TestCliLogsSinceFutureFollow(c *check.C) {
	printTestCaseName()
	defer printTestDuration(time.Now())
	testRequires(c, DaemonIsLinux)

	pullImageIfNotExist("busybox")
	out, _ := dockerCmd(c, "run", "-d", "busybox", "/bin/sh", "-c", `for i in $(seq 1 50); do date +%s; sleep 1; done`)
	time.Sleep(5 * time.Second)
	id := strings.TrimSpace(out)

	now := daemonTime(c).Unix()
	since := now - 5

	out, _ = dockerCmd(c, "logs", "-f", fmt.Sprintf("--since=%v", since), id)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	c.Assert(lines, checker.Not(checker.HasLen), 0)
	for _, v := range lines {
		ts, err := strconv.ParseInt(v, 10, 64)
		c.Assert(err, checker.IsNil, check.Commentf("cannot parse timestamp output from log: '%v'\nout=%s", v, out))
		c.Assert((ts+3) >= since, checker.Equals, true, check.Commentf("earlier log found. since=%v logdate=%v", since, ts))
	}
}

//TODO: fix Goroutine in multi-tenancy environment
/*func (s *DockerSuite) TestCliLogsFollowGoroutinesWithStdout(c *check.C) {
	printTestCaseName(); defer printTestDuration(time.Now())

	testRequires(c, DaemonIsLinux)
	out, _ := dockerCmd(c, "run", "-d", "busybox", "/bin/sh", "-c", "while true; do echo hello; sleep 2; done")
	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), checker.IsNil)

	type info struct {
		NGoroutines int
	}
	getNGoroutines := func() int {
		var i info
		status, b, err := sockRequest("GET", "/info", nil)
		c.Assert(err, checker.IsNil)
		c.Assert(status, checker.Equals, 200)
		c.Assert(json.Unmarshal(b, &i), checker.IsNil)
		return i.NGoroutines
	}

	nroutines := getNGoroutines()

	cmd := exec.Command(dockerBinary, "--region="+os.Getenv("DOCKER_HOST"), "logs", "-f", id)
	r, w := io.Pipe()
	cmd.Stdout = w
	c.Assert(cmd.Start(), checker.IsNil)

	// Make sure pipe is written to
	chErr := make(chan error)
	go func() {
		b := make([]byte, 1)
		_, err := r.Read(b)
		chErr <- err
	}()
	c.Assert(<-chErr, checker.IsNil)
	c.Assert(cmd.Process.Kill(), checker.IsNil)

	// NGoroutines is not updated right away, so we need to stop before failing
	t := time.After(30 * time.Second)
	for {
		select {
		case <-t:
			n := getNGoroutines()
			c.Assert(n <= nroutines, checker.Equals, true, check.Commentf("leaked goroutines: expected less than or equal to %d, got: %d", nroutines, n))

		default:
			if n := getNGoroutines(); n <= nroutines {
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (s *DockerSuite) TestCliLogsFollowGoroutinesNoOutput(c *check.C) {
	printTestCaseName(); defer printTestDuration(time.Now())

	testRequires(c, DaemonIsLinux)
	pullImageIfNotExist("busybox")
	out, _ := dockerCmd(c, "run", "-d", "busybox", "/bin/sh", "-c", "while true; do sleep 2; done")
	id := strings.TrimSpace(out)
	c.Assert(waitRun(id), checker.IsNil)

	type info struct {
		NGoroutines int
	}
	getNGoroutines := func() int {
		var i info
		status, b, err := sockRequest("GET", "/info", nil)
		c.Assert(err, checker.IsNil)
		c.Assert(status, checker.Equals, 200)
		c.Assert(json.Unmarshal(b, &i), checker.IsNil)
		return i.NGoroutines
	}

	nroutines := getNGoroutines()

	cmd := exec.Command(dockerBinary, "--region="+os.Getenv("DOCKER_HOST"), "logs", "-f", id)
	c.Assert(cmd.Start(), checker.IsNil)
	time.Sleep(200 * time.Millisecond)
	c.Assert(cmd.Process.Kill(), checker.IsNil)

	// NGoroutines is not updated right away, so we need to stop before failing
	t := time.After(30 * time.Second)
	for {
		select {
		case <-t:
			n := getNGoroutines()
			c.Assert(n <= nroutines, checker.Equals, true, check.Commentf("leaked goroutines: expected less than or equal to %d, got: %d", nroutines, n))

		default:
			if n := getNGoroutines(); n <= nroutines {
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
}*/

func (s *DockerSuite) TestCliLogsCLIContainerNotFound(c *check.C) {
	printTestCaseName()
	defer printTestDuration(time.Now())
	testRequires(c, DaemonIsLinux)

	name := "testlogsnocontainer"
	out, _, _ := dockerCmdWithError("logs", name)
	message := fmt.Sprintf("Error: No such container: %s\n", name)
	c.Assert(out, checker.Equals, message)
}
