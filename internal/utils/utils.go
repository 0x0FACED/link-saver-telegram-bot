package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	MAX_DESC_SIZE = 32
)

var (
	// if len(parts) != 3 (/save link desc -> must be)
	ErrSaveMessageLen = errors.New("incorrect message format")
	// is parts[1] != uri (incorrect link)
	ErrNotUri = errors.New("incorrect url format")

	ErrDescSize = errors.New("desc size must be less than 32")

	ErrNotUTF8 = errors.New("there are characters not in utf8 encoding")

	ErrEmptyMessage = errors.New("empty message")
)

func GetCodeMsgFromError(err error) string {
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.OK:
			return fmt.Sprintf("Success: %v", st.Message())
		case codes.Canceled:
			return fmt.Sprintf("Request was canceled: %v", st.Message())
		case codes.Unknown:
			return fmt.Sprintf("Unknown error occurred: %v", st.Message())
		case codes.InvalidArgument:
			return fmt.Sprintf("Invalid argument provided: %v", st.Message())
		case codes.DeadlineExceeded:
			return fmt.Sprintf("Deadline exceeded before the operation could complete: %v", st.Message())
		case codes.NotFound:
			return fmt.Sprintf("Requested resource not found: %v", st.Message())
		case codes.AlreadyExists:
			return fmt.Sprintf("Resource already exists: %v", st.Message())
		case codes.PermissionDenied:
			return fmt.Sprintf("Permission denied for the requested operation: %v", st.Message())
		case codes.ResourceExhausted:
			return fmt.Sprintf("Resource exhausted, please try again later: %v", st.Message())
		case codes.FailedPrecondition:
			return fmt.Sprintf("Operation rejected due to failed precondition: %v", st.Message())
		case codes.Aborted:
			return fmt.Sprintf("Operation aborted: %v", st.Message())
		case codes.OutOfRange:
			return fmt.Sprintf("Operation out of range: %v", st.Message())
		case codes.Unimplemented:
			return fmt.Sprintf("Operation not implemented: %v", st.Message())
		case codes.Internal:
			return fmt.Sprintf("Internal server error: %v", st.Message())
		case codes.Unavailable:
			return fmt.Sprintf("Service unavailable: %v", st.Message())
		case codes.DataLoss:
			return fmt.Sprintf("Data loss occurred: %v", st.Message())
		case codes.Unauthenticated:
			return fmt.Sprintf("Unauthenticated request: %v", st.Message())
		default:
			return fmt.Sprintf("Unknown gRPC error: %v", st.Message())
		}
	}

	return "Unknown error"

}
func ValidateSaveMessage(msg string) error {
	if msg == "" {
		return ErrEmptyMessage
	}
	parts := strings.SplitN(msg, " ", 3)

	if len(parts) != 3 {
		return ErrSaveMessageLen
	}

	_, err := url.ParseRequestURI(parts[1])
	if err != nil {
		return ErrNotUri
	}

	if len(parts[2]) > MAX_DESC_SIZE {
		return ErrDescSize
	}

	if !utf8.ValidString(parts[2]) {
		return ErrNotUTF8
	}

	return nil
}
