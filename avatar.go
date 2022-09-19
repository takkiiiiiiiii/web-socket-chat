package main

import (
	"crypto/md5"
	"errors"
	"io"
	"strings"
)

//ErrNoAvatarはAvatarインスタンスがアバターのURLを返すことができない
//場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できない")

//Avatarはユーザーのプロフィール画像を表す型
type Avatar interface {
	//GetAvatarURLは指定されたクライアントのアバターのURlを取得
	//問題が発生した場合にはエラーを返す
	//特に、URLを取得できなかった場合にはErrNoAvatarURLを返す
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}
type GravatarAvatar struct{}

var UseAuthAvatar AuthAvatar
var UseGravatar GravatarAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(useridStr))
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}
