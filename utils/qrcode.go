package utils

import (
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type QrcodeInfo struct {
	bgImg image.Image
	color *image.RGBA
	Qur   QrcodeUserInfo
}

type QrcodeUserInfo struct {
	url             string
	centerImagePath string
	qrcodeSavePath  string
	centerImageSize int
	imageHigh       int
	imageWidth      int
}

func (qr *QrcodeInfo) MakeNewQrcode() error {
	var err error
	var qrCode *qrcode.QRCode
	qrCode, err = qrcode.New(qr.Qur.url, qrcode.Highest)
	if err != nil {
		return err
	}
	qrCode.DisableBorder = true
	qr.bgImg = qrCode.Image(qr.Qur.centerImageSize)
	return nil
}

func (qr *QrcodeInfo) SetCenterImage() error {
	var err error
	avatarFile, err := os.Open(qr.Qur.centerImagePath)
	if err != nil {
		return err
	}
	avatarImg, err := png.Decode(avatarFile)
	if err != nil {
		return err
	}
	avatarImg = resize.Resize(uint(qr.Qur.imageWidth), uint(qr.Qur.imageHigh), avatarImg, resize.Lanczos3)
	b := qr.bgImg.Bounds()

	//set into the center
	offset := image.Pt((b.Max.X-avatarImg.Bounds().Max.X)/2, (b.Max.Y-avatarImg.Bounds().Max.Y)/2)
	m := image.NewRGBA(b)
	draw.Draw(m, b, qr.bgImg, image.Point{X: 0, Y: 0}, draw.Src)
	draw.Draw(m, avatarImg.Bounds().Add(offset), avatarImg, image.Point{X: 0, Y: 0}, draw.Over)
	qr.color = m
	return nil
}

func (qr *QrcodeInfo) SaveQrode() error {
	var err error
	file, err := os.Create(qr.Qur.qrcodeSavePath)
	if err != nil {
		return err
	}
	err = png.Encode(file, qr.color)
	if err != nil {
		return err
	}
	return nil
}

type QrcodeInfoOption func(qr *QrcodeInfo)

func NewQrCode(options ...QrcodeInfoOption) *QrcodeInfo {
	qr := &QrcodeInfo{
		bgImg: nil,
		color: nil,
		Qur:   QrcodeUserInfo{},
	}
	for _, option := range options {
		option(qr)
	}
	return qr
}

func WithQrcodeUrl(url string) QrcodeInfoOption {
	return func(qr *QrcodeInfo) {
		qr.Qur.url = url
	}
}

func WithQrcodeSize(w int, h int) QrcodeInfoOption {
	return func(qr *QrcodeInfo) {
		qr.Qur.imageWidth = w
		qr.Qur.imageHigh = h
	}
}

func WithCenterImage(imagePath string) QrcodeInfoOption {
	return func(qr *QrcodeInfo) {
		qr.Qur.centerImagePath = imagePath
	}
}

func WithSavePath(savePath string) QrcodeInfoOption {
	return func(qr *QrcodeInfo) {
		qr.Qur.qrcodeSavePath = savePath
	}
}

func WithCenterImageSize(size int) QrcodeInfoOption {
	return func(qr *QrcodeInfo) {
		qr.Qur.centerImageSize = size
	}
}
