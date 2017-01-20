package utils

import (
    "bytes"
    "encoding/gob"
    "golang.org/x/crypto/scrypt"
    "fmt"
    "encoding/base64"
    "prosnav.com/wxserver/modules/setting"
    "math/rand"
    "prosnav.com/wxserver/modules/results"
    "time"
    "os"
    "path"
    "path/filepath"
    "io"
    "errors"
    "image"
    "image/jpeg"
    "image/png"
    "image/gif"
    "prosnav.com/wxserver/modules/log"
    "net/url"
    "net/http"
    "encoding/json"
)

/**
   This function implements deepcopy, forever loop with a cycle reference
**/

func Clone(a, b interface{}) {
    buff := new(bytes.Buffer)
    enc := gob.NewEncoder(buff)
    dec := gob.NewDecoder(buff)
    enc.Encode(a)
    dec.Decode(b)
}

func Serialize(v interface{}) []byte {
    buff := new(bytes.Buffer)
    enc := gob.NewEncoder(buff)
    enc.Encode(v)
    return buff.Bytes()
}

func UnSerialize(data []byte, inst interface{}) {
    buf := new(bytes.Buffer)
    buf.Write(data)
    dec := gob.NewDecoder(buf)
    dec.Decode(inst)
    return
}

const (
    ENCRYPT_ALGORITHM = "pbkdf2"
    ENCRYPT_HASHER = "sha256"
    ENCRYPT_TIMES = 16384
)

var charset = []byte{
    'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
    'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
    'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
    'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
    '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

const (
    UPPER_CHARSET = iota
    LOWER_CHARSET
    INT_CHARSET
    UPPER_LOWER_CHARSET
    FULL_CHARSET
)

func RandomSpecStr(num int, set uint) string {
    var subcharset []byte

    switch(set) {
    case UPPER_CHARSET: subcharset = charset[:25]
    case LOWER_CHARSET: subcharset = charset[26:51]
    case INT_CHARSET: subcharset = charset[52:]
    case UPPER_LOWER_CHARSET: subcharset = charset[:51]
    case FULL_CHARSET: subcharset = charset
    default: panic(errors.New("Unknown charset identify:" + string(set)))
    }

    var buf = make([]byte, num)
    for i := 0; i < num; i++ {
        index := rand.Intn(len(subcharset))
        buf[i] = subcharset[index]
    }
    return string(buf)
}

func RandomStr(num int) string {
    return RandomSpecStr(num, FULL_CHARSET)
}

func EncryptPassword(password string, salts ...string) (encryptedPassword string, err error) {
    salt := RandomStr(6)
    if len(salts) > 0 {
        salt = salts[0]
    }

    dk, err := scrypt.Key([]byte(password), []byte(salt), ENCRYPT_TIMES, 8, 1, 32)

    if err != nil {
        return "", err
    }

    encryptedPassword = fmt.Sprintf("%s_%s_%s_%d_%s", ENCRYPT_ALGORITHM, ENCRYPT_HASHER, salt, ENCRYPT_TIMES, base64.StdEncoding.EncodeToString(dk))
    return encryptedPassword, nil
}

func In(e string, list []string) bool {
    for _, v := range list {
        if e == v {
            return true
        }
    }

    return false
}

type endpoint struct {
    AuthURL, TokenURL string
}

type oauth2Config struct {
    ClientID     string
    ClientSecret string
    Endpoint     endpoint
    RedirectURL  string
}

func GetAuthUrl(oathtype string) (authUrl string) {
    switch oathtype {
    case "wechat":
        authUrl = fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&#wechat_redirect",
            setting.Cfg.Section("wechat").Key("CORP_ID").String(),
            url.QueryEscape(setting.Cfg.Section("wechat").Key("REDIRECT_URI").String()))
    default: panic(new(results.NotAllowedError))
    }
    return
}

func FetchUrlWithJson(urlStr string, result interface{}) error {
    resp, err := http.Get(urlStr)
    if err != nil {
        log.Error(3, "%v", err)
        return err
    }
    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(result); err != nil {
        log.Error(3, "%v", err)
        return err
    }

    return nil
}

func ImageSize(filePath string) (width, height int, err error) {
    f, err := os.Open(filePath)
    defer f.Close()
    if err != nil {
        return 0, 0, err
    }

    img, _, err := image.Decode(f)
    if err != nil {
        return 0, 0, err
    }
    p := img.Bounds().Size()

    return p.X, p.Y, nil
}

//func ScaleImage(filepath string, w, h int) (string, error) {
//    src, err := loadImage(filepath)
//    if err != nil {
//        log.Error(3, "%v", err)
//        return "", err
//    }
//
//    dst := image.NewRGBA(image.Rect(0, 0, w, h))
//    if err := graphics.Scale(dst, src); err != nil {
//        log.Error(3, "%v", err)
//        return "", err
//    }
//    dstpath := path.Join(path.Dir(filepath), fmt.Sprintf("scale%dx%d.%s", w, h, path.Base(filepath)))
//    if err := saveImage(dstpath, dst); err != nil {
//        log.Error(3, "%v", err)
//        return "", err
//    }
//    return dstpath, nil
//}

func loadImage(filepath string) (image.Image, error) {
    f, err := os.Open(filepath)
    defer f.Close()
    if err != nil {
        log.Error(3, "%v", err)
        return nil, err
    }
    img, _, err := image.Decode(f)
    if err != nil {
        log.Error(3, "%v", err)
        return nil, err
    }
    return img, err
}

func saveImage(filepath string, img image.Image) (err error) {
    imgfile, err := os.Create(filepath)
    defer imgfile.Close()
    switch path.Ext(filepath) {
    case ".jpg", ".jpeg":
        if err = jpeg.Encode(imgfile, img, nil); err != nil {
            log.Error(3, "%v", err)
            return
        }
    case ".png":
        if err = png.Encode(imgfile, img); err != nil {
            log.Error(3, "%v", err)
            return
        }
    case ".gif":
        if err = gif.Encode(imgfile, img, new(gif.Options)); err != nil {
            log.Error(3, "%v", err)
            return
        }
    }
    return
}

func SaveFile(src io.Reader, srcName string) string {
    suffix := filepath.Ext(srcName)
    fmt.Println(suffix)
    rootPath := setting.Cfg.Section("fileupload").Key("ROOT_PATH").MustString("/fileupload")
    fmt.Println(rootPath)
    now := time.Now()
    today := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
    dir := path.Join(rootPath, today)
    fmt.Println(dir)
    os.MkdirAll(dir, 0766)
    filename := RandomStr(12) + suffix
    filePath := path.Join(dir, filename)
    fmt.Println(filePath)
    dst, err := os.Create(filePath)
    if err != nil {
        panic(err)
    }
    io.Copy(dst, src)

    return filePath
}

func imageRegister() {
    image.RegisterFormat("jpeg", "FFD8FF?",
        func(reader io.Reader) (image.Image, error) {
            return jpeg.Decode(reader)
        }, func(reader io.Reader) (image.Config, error) {
            return jpeg.DecodeConfig(reader)
        })
    image.RegisterFormat("jpg", "FFD8FF?",
        func(reader io.Reader) (image.Image, error) {
            return jpeg.Decode(reader)
        }, func(reader io.Reader) (image.Config, error) {
            return jpeg.DecodeConfig(reader)
        })
    image.RegisterFormat("png", "89504E470D0A1A0A",
        func(reader io.Reader) (image.Image, error) {
            return png.Decode(reader)
        }, func(reader io.Reader) (image.Config, error) {
            return png.DecodeConfig(reader)
        })
    image.RegisterFormat("gif", "47494638?",
        func(reader io.Reader) (image.Image, error) {
            return gif.Decode(reader)
        }, func(reader io.Reader) (image.Config, error) {
            return gif.DecodeConfig(reader)
        })
}

func init() {
    imageRegister()
}


