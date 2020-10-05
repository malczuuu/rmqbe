package auth

import "testing"

func TestShouldMatchNameByPrefix(t *testing.T) {
	res := matchNameByPattern("mqtt-subscription-*", "mqtt-subscription-mosq-eNQORgKjhUB2g2ij7Zqos1")

	if !res {
		t.Error("matchNameByPattern(...) didn't match name with tailing *")
	}
}

func TestShouldNotMatchNameBySuffix(t *testing.T) {
	res := matchNameByPattern("*-mqtt-subscription", "mqtt-subscription-mosq-eNQORgKjhUB2g2ij7Zqos1")

	if res {
		t.Error("matchNameByPattern(...) shouldn't match by suffix")
	}
}

func TestShouldMatchNameBySuffix(t *testing.T) {
	res := matchNameByPattern("*-mosq-eNQORgKjhUB2g2ij7Zqos1", "mqtt-subscription-mosq-eNQORgKjhUB2g2ij7Zqos1")

	if !res {
		t.Error("matchNameByPattern(...) didn't match name with heading *")
	}
}

func TestShouldNotMatchNameByPrefix(t *testing.T) {
	res := matchNameByPattern("mosq-eNQORgKjhUB2g2ij7Zqos1-*", "mqtt-subscription-mosq-eNQORgKjhUB2g2ij7Zqos1")

	if res {
		t.Error("matchNameByPattern(...) shouldn't match by prefix")
	}
}

func TestShouldMatchNameInside(t *testing.T) {
	res := matchNameByPattern("*-mosq-*", "mqtt-subscription-mosq-eNQORgKjhUB2g2ij7Zqos1")

	if !res {
		t.Error("matchNameByPattern(...) didn't match name with both heading and tailing *")
	}
}

func TestShouldNotMatchNameInside(t *testing.T) {
	res := matchNameByPattern("*-mosq-subscription-*", "mqtt-subscription-mosq-eNQORgKjhUB2g2ij7Zqos1")

	if res {
		t.Error("matchNameByPattern(...) shouldn't match both heading and tailing *")
	}
}
