package main

import (
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client) //インスタンス生成  = var client client
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、AuthAvatar.GetAvatarURLはErrNoAvatarURLを返すべき")
	}
	//値をセット
	testUrl := "http://url-to-avatar/"
	client.userData = map[string]interface{}{
		"avatar_url": testUrl,
	}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("値が存在する場合、AuthAvaar.GetAvatarURLは" + "エラーを返すべきではい")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURLは正しいURLを返すべき")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userid": "abe9558e86fa0691266e7d79f6181112"}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではありません")
	}
	if url != "//www.gravatar.com/avatar/abe9558e86fa0691266e7d79f6181112" {
		t.Errorf("GravatarAvatar.GetAvatarURLが%sという誤った値を受け取った", url)
	}
}
