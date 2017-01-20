package setting

import (
    "gopkg.in/ini.v1"
    "github.com/gogits/gogs/modules/log"
    "fmt"
    "strings"
    "os"
    "path"
    "path/filepath"
    "os/exec"
)

var 	(
    // App settings.
    AppVer    string
    AppName   string
    ListeningPort string
    //AppUrl    string
    //AppSubUrl string
    AppRootPath string

    // Log setting
    LogRootPath string
    LogModes    []string
    LogConfigs  []string

    Cfg *ini.File
    MailService  *Mailer
)

// Mailer represents mail service.
type Mailer struct {
    Name              string
    Host              string
    From              string
    User, Passwd      string
    SkipVerify        bool
    UseCertificate    bool
    CertFile, KeyFile string
}


func ExecPath() (string, error) {
    file, err := exec.LookPath(os.Args[0])
    if err != nil {
        return "", err
    }
    p, err := filepath.Abs(file)
    if err != nil {
        return "", err
    }
    return p, nil
}

// WorkDir returns absolute path of work directory.
func WorkDir() (string, error) {
    execPath, err := ExecPath()
    return path.Dir(strings.Replace(execPath, "\\", "/", -1)), err
}

func NewConfigContext() {
    var err error
    if AppRootPath, err = WorkDir(); err != nil {
        log.Error(4, "%v", err)
        return
    }
    //sec := Cfg.Section("server")
    AppName = Cfg.Section("").Key("APP_NAME").MustString("Gogs: Go Git Service")
    //AppUrl = sec.Key("ROOT_URL").MustString("http://localhost:3000/")
    AppVer = Cfg.Section("").Key("APP_VERSION").MustString("1.0")
    ListeningPort = fmt.Sprintf(":%s", Cfg.Section("server").Key("HTTP_PORT").MustString("8080"))
    //log
    LogRootPath = Cfg.Section("log").Key("ROOT_PATH").MustString(path.Join(AppRootPath, "log"))
    LogModes = strings.Split(Cfg.Section("log").Key("MODE").MustString("console"), ",")

}

var logLevels = map[string]string{
    "Trace":    "0",
    "Debug":    "1",
    "Info":     "2",
    "Warn":     "3",
    "Error":    "4",
    "Critical": "5",
}

func NewLogService() {
    log.Info("%s %s", AppName, AppVer)

    LogConfigs = make([]string, len(LogModes))
    for i, mode := range LogModes {
        mode = strings.TrimSpace(mode)
        sec, err := Cfg.GetSection("log." + mode)
        if err != nil {
            log.Fatal(4, "Unknown log mode: %s", mode)
        }

        validLevels := []string{"Trace", "Debug", "Info", "Warn", "Error", "Critical"}
        // Log level.
        levelName := Cfg.Section("log."+mode).Key("LEVEL").In(
        Cfg.Section("log").Key("LEVEL").In("Trace", validLevels),
        validLevels)
        level, ok := logLevels[levelName]
        if !ok {
            log.Fatal(4, "Unknown log level: %s", levelName)
        }

        // Generate log configuration.
        switch mode {
            case "console":
            LogConfigs[i] = fmt.Sprintf(`{"level":%s}`, level)
            case "file":
            logPath := sec.Key("FILE_NAME").MustString(path.Join(LogRootPath, "gogs.log"))
            os.MkdirAll(path.Dir(logPath), os.ModePerm)
            LogConfigs[i] = fmt.Sprintf(
            `{"level":%s,"filename":"%s","rotate":%v,"maxlines":%d,"maxsize":%d,"daily":%v,"maxdays":%d}`, level,
            logPath,
            sec.Key("LOG_ROTATE").MustBool(true),
            sec.Key("MAX_LINES").MustInt(1000000),
            1<<uint(sec.Key("MAX_SIZE_SHIFT").MustInt(28)),
            sec.Key("DAILY_ROTATE").MustBool(true),
            sec.Key("MAX_DAYS").MustInt(7))
            case "conn":
            LogConfigs[i] = fmt.Sprintf(`{"level":%s,"reconnectOnMsg":%v,"reconnect":%v,"net":"%s","addr":"%s"}`, level,
            sec.Key("RECONNECT_ON_MSG").MustBool(),
            sec.Key("RECONNECT").MustBool(),
            sec.Key("PROTOCOL").In("tcp", []string{"tcp", "unix", "udp"}),
            sec.Key("ADDR").MustString(":7020"))
            case "smtp":
            LogConfigs[i] = fmt.Sprintf(`{"level":%s,"username":"%s","password":"%s","host":"%s","sendTos":"%s","subject":"%s"}`, level,
            sec.Key("USER").MustString("example@example.com"),
            sec.Key("PASSWD").MustString("******"),
            sec.Key("HOST").MustString("127.0.0.1:25"),
            sec.Key("RECEIVERS").MustString("[]"),
            sec.Key("SUBJECT").MustString("Diagnostic message from serve"))
            case "database":
            LogConfigs[i] = fmt.Sprintf(`{"level":%s,"driver":"%s","conn":"%s"}`, level,
            sec.Key("DRIVER").String(),
            sec.Key("CONN").String())
        }

        log.NewLogger(Cfg.Section("log").Key("BUFFER_LEN").MustInt64(10000), mode, LogConfigs[i])
        log.Info("Log Mode: %s(%s)", strings.Title(mode), levelName)
    }
}


func NewMailContext() {
    sec := Cfg.Section("mail")
    if enable := sec.Key("ENABLE").MustBool(); !enable {
        log.Info("Mail Service Disabled")
        return
    }
    MailService =  &Mailer{
        Name: sec.Key("NAME").String(),
        Host: sec.Key("HOST").String(),
        From: sec.Key("FROM").String(),
        User: sec.Key("USER").String(),
        Passwd: sec.Key("PASSWD").String(),
        SkipVerify: sec.Key("SKIPVERIFY").MustBool(),
    }
    log.Info("Mail Service Enabled")
}

func newServices() {
    NewConfigContext()
    NewLogService()
    NewMailContext()
}

func init() {
    log.NewLogger(0, "console", `{"level": 0}`)
    var err error
    Cfg, err = ini.Load("conf/app.ini")
    if err != nil {
        log.Error(0, "Load configuration file app.ini failed %v", err)
    }
    newServices()
}
