package repositories_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGroupRepository_GetSaveGroupId(t *testing.T) {
	tcs := []struct {
		Value  int64
		Expect int64
	}{
		{
			Value:  5,
			Expect: 5,
		},
	}

	for _, tc := range tcs {
		err := GroupRepository.SaveGroupId(context.Background(), tc.Value)
		id, err := GroupRepository.GetGroupId(context.Background())

		assert.Nil(t, err)
		assert.Equal(t, id, tc.Expect)
	}
}
