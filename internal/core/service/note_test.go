package service_test

import (
	"context"
	"testing"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateNoteInvalid_Unit(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	repo := mock.NewMockNoteRepository(cntrl)

	s := service.NewNoteService(repo)

	cases := []struct {
		name string
		in   entity.Note
		out  error
	}{
		{
			name: "invalid_status_format",
			in: entity.Note{
				NoteID:      1,
				Title:       "Test",
				Description: "Test",
				Status:      "invalid_status",
				AuthorID:    2,
			},
			out: service.ErrInvalidStatusFormat,
		},
		{
			name: "invalid_id",
			in: entity.Note{
				NoteID:      -5,
				Title:       "Test",
				Description: "Test",
				Status:      "SUCCESS",
				AuthorID:    2,
			},
			out: service.ErrInvalidIDFormat,
		},
		{
			name: "invalid_tittle_format",
			in: entity.Note{
				NoteID:      3,
				Title:       "",
				Description: "Test",
				Status:      "SUCCESS",
				AuthorID:    2,
			},
			out: service.ErrInvalidTitleFormat,
		},
	}

	for _, c := range cases {
		err := s.UpdateNote(context.Background(), c.in)
		assert.EqualError(t, err, c.out.Error(), c.name)
	}

}

func TestUpdateNoteValid_Unit(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	repo := mock.NewMockNoteRepository(cntrl)

	repo.EXPECT().UpdateNote(gomock.Any(), gomock.Any()).Return(nil).Times(3)

	s := service.NewNoteService(repo)

	cases := []struct {
		name string
		in   entity.Note
	}{
		{
			name: "valid_1",
			in: entity.Note{
				NoteID:      1,
				Title:       "Test",
				Description: "Test",
				Status:      "IN_PROGRESS",
				AuthorID:    2,
			},
		},
		{
			name: "valid_2",
			in: entity.Note{
				NoteID:      5,
				Title:       "Test",
				Description: "Test",
				Status:      "SUCCESS",
				AuthorID:    2,
			},
		},
		{
			name: "valid_3",
			in: entity.Note{
				NoteID:      3,
				Title:       "Labubu",
				Description: "Test",
				Status:      "NOT_START",
				AuthorID:    2,
			},
		},
	}

	for _, c := range cases {
		err := s.UpdateNote(context.Background(), c.in)
		assert.NoError(t, err, c.name)
	}

}

func TestDeleteNoteInvalid_Unit(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	repo := mock.NewMockNoteRepository(cntrl)

	s := service.NewNoteService(repo)

	cases := []struct {
		name string
		in   struct {
			noteID   int
			authorID int
		}
		out error
	}{
		{
			name: "invalid_id",
			in: struct {
				noteID   int
				authorID int
			}{
				noteID:   -5,
				authorID: 213,
			},
			out: service.ErrInvalidIDFormat,
		},
	}

	for _, c := range cases {
		err := s.DeleteNote(context.Background(), c.in.noteID, c.in.authorID)
		assert.EqualError(t, err, c.out.Error(), c.name)
	}

}

func TestDeleteValid_Unit(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	repo := mock.NewMockNoteRepository(cntrl)
	repo.EXPECT().DeleteNote(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	s := service.NewNoteService(repo)

	cases := []struct {
		name string
		in   struct {
			noteID   int
			authorID int
		}
		out error
	}{
		{
			name: "valid_1",
			in: struct {
				noteID   int
				authorID int
			}{
				noteID:   5,
				authorID: 213,
			},
			out: service.ErrInvalidIDFormat,
		},
	}

	for _, c := range cases {
		err := s.DeleteNote(context.Background(), c.in.noteID, c.in.authorID)
		assert.NoError(t, err, c.name)
	}

}

func TestCheckStatusInvalid_Unit(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "invalid_status",
			in:   "ivalid_status",
			out:  false,
		},
	}

	for _, c := range cases {
		b := service.CheckStatus(c.in)
		assert.Equal(t, b, false, c.name)
	}

}

func TestCheckStatusValid(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "valid_1",
			in:   "SUCCESS",
			out:  true,
		},
		{
			name: "valid_2",
			in:   "IN_PROGRESS",
			out:  true,
		},
		{
			name: "valid_3",
			in:   "NOT_START",
			out:  true,
		},
	}

	for _, c := range cases {
		b := service.CheckStatus(c.in)
		assert.Equal(t, b, true, c.name)
	}

}
