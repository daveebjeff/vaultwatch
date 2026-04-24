package notify

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNormalizeNotifier_NilInner(t *testing.T) {
	_, err := NewNormalizeNotifier(nil)
	require.Error(t, err)
	assert.Equal(t, ErrNilInner, err)
}

func TestNewNormalizeNotifier_Valid(t *testing.T) {
	n, err := NewNormalizeNotifier(NewNoopNotifier())
	require.NoError(t, err)
	assert.NotNil(t, n)
}

func TestNormalizeNotifier_TrimsWhitespace(t *testing.T) {
	cap := &captureNotifier{}
	n, err := NewNormalizeNotifier(cap)
	require.NoError(t, err)

	msg := Message{Path: "secret/a", Body: "  hello world  "}
	require.NoError(t, n.Send(msg))
	assert.Equal(t, "hello world", cap.last.Body)
}

func TestNormalizeNotifier_CollapsesInternalSpaces(t *testing.T) {
	cap := &captureNotifier{}
	n, err := NewNormalizeNotifier(cap)
	require.NoError(t, err)

	msg := Message{Path: "secret/b", Body: "foo   bar\t\nbaz"}
	require.NoError(t, n.Send(msg))
	assert.Equal(t, "foo bar baz", cap.last.Body)
}

func TestNormalizeNotifier_LowerCaseOption(t *testing.T) {
	cap := &captureNotifier{}
	n, err := NewNormalizeNotifier(cap, WithLowerCase())
	require.NoError(t, err)

	msg := Message{Path: "secret/c", Body: "  VAULT Secret EXPIRING  "}
	require.NoError(t, n.Send(msg))
	assert.Equal(t, "vault secret expiring", cap.last.Body)
}

func TestNormalizeNotifier_EmptyBody(t *testing.T) {
	cap := &captureNotifier{}
	n, err := NewNormalizeNotifier(cap)
	require.NoError(t, err)

	msg := Message{Path: "secret/d", Body: "   "}
	require.NoError(t, n.Send(msg))
	assert.Equal(t, "", cap.last.Body)
}

func TestNormalizeNotifier_NonBodyFieldsUnchanged(t *testing.T) {
	cap := &captureNotifier{}
	n, err := NewNormalizeNotifier(cap)
	require.NoError(t, err)

	msg := Message{Path: "secret/e", Body: " ok ", Status: StatusExpired}
	require.NoError(t, n.Send(msg))
	assert.Equal(t, "secret/e", cap.last.Path)
	assert.Equal(t, StatusExpired, cap.last.Status)
}

func TestNormalizeNotifier_InnerErrorPropagated(t *testing.T) {
	fail := &failNotifier{err: assert.AnError}
	n, err := NewNormalizeNotifier(fail)
	require.NoError(t, err)

	err = n.Send(Message{Path: "secret/f", Body: "body"})
	assert.ErrorIs(t, err, assert.AnError)
}
