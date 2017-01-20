package files
import (
    "github.com/gin-gonic/gin"
    "net/http"
    "prosnav.com/wxserver/modules/setting"
    "path/filepath"
    "prosnav.com/wxserver/modules/log"
    "os"
    "prosnav.com/wxserver/modules/results"
    "io"
)

var rootPath string

func FileServe(c *gin.Context) {
    productCode := c.Params.ByName("tagCode")
    newsDir := c.Params.ByName("newsdir")
    filename := c.Params.ByName("filename")
    targetFile := filepath.Join(rootPath, productCode, newsDir, filename)
    http.ServeFile(c.Writer, c.Request, targetFile)
}

type newsDirForm struct {
    NewsDir string  `form:"newsdir" json:"newsdir"`
}

func Mkdir(c *gin.Context) {
    productCode := c.Params.ByName("tagCode")
    var form newsDirForm
    c.Bind(&form)
    newsDir := form.NewsDir
    productDir := filepath.Join(rootPath, productCode)
    err := os.MkdirAll(productDir, 0777)
    if os.IsExist(err) {
        panic(results.NewBusinessError("5008"))
    }
    if newsDir == "" {
        return
    }
    targetFile := filepath.Join(productDir, newsDir)
    err = os.Mkdir(targetFile, 0777)
    if os.IsExist(err) {
        panic(results.NewBusinessError("5008"))
    }
    if err != nil {
        panic(err)
    }

    c.JSON(200, nil)
}

type fileForm struct {
    Name string     `json:"name"`
    Mode string     `json:"mode"`
    IsDir bool      `json:"isdir"`
}

func Ls(c *gin.Context) {
    productCode := c.Params.ByName("tagCode")
    newsDir := c.Params.ByName("newsdir")
    targetFile := filepath.Join(rootPath, productCode, newsDir)
    homeDir, err := os.Open(targetFile)
    if os.IsNotExist(err) {
        panic(results.NewBusinessError("5009"))
    }
    fl, err := homeDir.Readdir(-1)
    if err == io.EOF {
        log.Error(3, "Empty dir, %v", err)
        return
    }
    var fileList []fileForm
    for _, fi := range fl {
        fileList = append(fileList, fileForm{fi.Name(), fi.Mode().String(), fi.IsDir()})
    }

    c.JSON(200, fileList)
}

func Upload(c *gin.Context) {
    productCode := c.Params.ByName("tagCode")
    newsDir := c.Params.ByName("newsdir")
    err := c.Request.ParseMultipartForm(32 << 20)
    if err != nil {
        panic(err)
    }
    m := c.Request.MultipartForm
    files := m.File["images"]

    for _, f := range files {
        //for each fileheader, get a handle to the actual file
        src, err := f.Open()
        defer src.Close()
        if err != nil {
            panic(err)
        }
        //create destination file making sure the path is writeable.
        dst, err := os.Create(filepath.Join(rootPath, productCode, newsDir, f.Filename))
        defer dst.Close()
        if err != nil {
            panic(err)
        }
        //copy the uploaded file to the destination file
        if _, err := io.Copy(dst, src); err != nil {
            panic(err)
        }
    }

    c.JSON(200, nil)
}

func init() {
    rootPath = setting.Cfg.Section("fs").Key("ROOT_PATH").String()
}