package api

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

func getGoDoc(ident string) string {
	cmd := exec.Command("go", "doc", "-short", ident)
	buf := bytes.NewBuffer([]byte{})
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr

	logErr(cmd.Run())

	var (
		scanner = bufio.NewScanner(buf)
		sb      = strings.Builder{}
	)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			logErr(err)
			continue
		}

		txt := scanner.Text()
		if !strings.HasPrefix(txt, "    ") {
			continue
		}
		txt = strings.TrimPrefix(txt, "    ")
		sb.WriteString(txt)
		sb.WriteString(" ")
	}

	return strings.TrimSpace(sb.String())
}

func logErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}
