package config

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/xxf098/lite-proxy/outbound"
	"github.com/xxf098/lite-proxy/utils"
)

var (
	NotSSRLink error = errors.New("not a shadowsocksR link")
)

func SSRLinkToSSROption(link string) (*outbound.ShadowSocksROption, error) {
	regex := regexp.MustCompile(`^ssr://([A-Za-z0-9+-=/_]+)`)
	res := regex.FindAllStringSubmatch(link, 1)
	b64 := ""
	if len(res) > 0 && len(res[0]) > 1 {
		b64 = res[0][1]
	}
	uri, err := utils.DecodeB64(b64)
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(uri, "/?", 2)
	links := strings.Split(parts[0], ":")
	if len(links) != 6 || len(parts) != 2 {
		return nil, NotSSRLink
	}
	port, err := strconv.Atoi(links[1])
	if err != nil {
		return nil, err
	}
	pass, err := utils.DecodeB64(links[5])
	if err != nil {
		return nil, err
	}
	ssrOption := &outbound.ShadowSocksROption{
		Name:     "ssr",
		Server:   links[0],
		Port:     port,
		Protocol: links[2],
		Cipher:   links[3],
		Obfs:     links[4],
		Password: pass,
		UDP:      false,
	}
	if rawQuery, err := url.ParseQuery(parts[1]); err == nil {
		obfsparam, err := utils.DecodeB64(rawQuery.Get("obfsparam"))
		if err != nil {
			return nil, err
		}
		ssrOption.ObfsParam = obfsparam
		protoparam, err := utils.DecodeB64(rawQuery.Get("protoparam"))
		if err != nil {
			return nil, err
		}
		ssrOption.ProtocolParam = protoparam
		remarks, err := utils.DecodeB64(rawQuery.Get("remarks"))
		if err == nil {
			ssrOption.Remarks = remarks
		}
	}
	return ssrOption, nil
}

func init() {
	outbound.RegisterDialerCreator("ssr", func(link string) (outbound.Dialer, error) {
		ssOption, err := SSRLinkToSSROption(link)
		if err != nil {
			return nil, err
		}
		return outbound.NewShadowSocksR(ssOption)
	})
}
