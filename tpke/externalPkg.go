package tpke

import (
	"github.com/DE-labtory/cleisthenes/config"
	"github.com/DE-labtory/tpke"
)

type playGround struct {
	actors []*actor
	pkSet *tpke.PublicKeySet
}

type actor struct {
	id int
	skShare *tpke.SecretKeyShare
	pkShare *tpke.PublicKeyShare
	receivedMsg *tpke.CipherText
}


func setUp() *playGround {
	//获取指定文件的配置
	conf := config.Get()

	th := conf.Tpke.Threshold
	people := conf.HoneyBadger.NetworkSize

	//secretKeySet为私钥集合
	secretKeySet := tpke.RandomSecretKeySet(th)
	publicKeySet := secretKeySet.PublicKeySet()

	actors := make([]*actor, 0)
	i := 0
	for i < people {
		actors = append(actors, &actor{
			id: i,
			skShare: secretKeySet.KeyShare(i),
			pkShare: publicKeySet.KeyShare(i),
		})
		i++
	}
	return &playGround {
		actors: actors,
		pkSet: publicKeySet,
	}
}

func (pg *playGround) publishPubKey() *tpke.PublicKey {
	return pg.pkSet.PublicKey()
}
