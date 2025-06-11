package encryptmanage_test

import (
	"testing"

	encryptmanage "github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/managers/encrypt"
	"github.com/stretchr/testify/assert"
)

func TestEncodePassword_Unit(t *testing.T) {
	em := encryptmanage.NewEncrypter()

	cases := []struct {
		name string
		in   string
	}{
		{
			name: "test_1",
			in:   "rodion11",
		},
		{
			name: "test_2",
			in:   "",
		},
		{
			name: "test_3",
			in:   "121212121212==-=`1`1`1`",
		},
	}

	for _, c := range cases {
		hashedPassword, err := em.EncodePassword(c.in)
		t.Logf("%s-%s", c.in, hashedPassword)
		assert.NoError(t, err, c.name)
	}

}

func TestCheckPasswordValid_Unit(t *testing.T) {
	em := encryptmanage.NewEncrypter()

	cases := []struct {
		name string
		in   struct {
			password       string
			hashedPassword string
		}
	}{
		{
			name: "test_1",
			in: struct {
				password       string
				hashedPassword string
			}{
				password:       "rodion11",
				hashedPassword: "$2a$10$uUZt90THRxaIN8.qtIzCAuLdmidh/sFsHr8hNfZsKkcLcDP5BQ/gC",
			},
		},
		{
			name: "test_2",
			in: struct {
				password       string
				hashedPassword string
			}{
				password:       "",
				hashedPassword: "$2a$10$NhD/5UbxCv5UypRQaULbp.ldA2lt3JVWhaWbU.bQwOTuOE7a9yVVm",
			},
		},
		{
			name: "test_3",
			in: struct {
				password       string
				hashedPassword string
			}{
				password:       "121212121212==-=`1`1`1`",
				hashedPassword: "$2a$10$VqekMCH2mlwezj5qohOYu./OwtUret378vpTj567eHHBHW8lH.LDi",
			},
		},
	}

	for _, c := range cases {
		err := em.CheckPassword(c.in.password, c.in.hashedPassword)
		assert.NoError(t, err, c.name)
	}

}

func TestCheckPasswordInvalid_Unit(t *testing.T) {
	em := encryptmanage.NewEncrypter()

	cases := []struct {
		name string
		in   struct {
			password       string
			hashedPassword string
		}
	}{
		{
			name: "test_1",
			in: struct {
				password       string
				hashedPassword string
			}{
				password:       "rodion11123",
				hashedPassword: "$2a$10$uUZt90THRxaIN8.qtIzCAuLdmidh/sFsHr8hNfZsKkcLcDP5BQ/gC",
			},
		},
	}

	for _, c := range cases {
		err := em.CheckPassword(c.in.password, c.in.hashedPassword)
		assert.Error(t, err, c.name)
	}

}
