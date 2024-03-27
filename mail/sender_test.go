package mail

import (
	"semantic_api/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "a test email"
	content := `
		<h1>hello biches</h1>
		<p>This is a test bich<a href="google.com">googel</a></p>
	`
	to := []string{"supahero999@gmail.com"}
	attachFiles := []string{"../instr"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
