package main

import (
	"errors"
	"io/ioutil"
	"path"
)

// ErrNoAvatarURL returned
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL.\n")

// Avatar represent types representing user pictures
type Avatar interface {
	GetAvatarURL(u ChatUser) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c ChatUser) (string, error) {
	if len(c.AvatarURL()) > 0 {
		return c.AvatarURL(), nil
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "/avatars/d262423a4aa3dc61ea9a07306b053527.png", nil
	// return "", ErrNoAvatarURL
}
