package cmd
import (
"github.com/urfave/cli/v2"
"github.com/logrusorgru/aurora"
log "github.com/sirupsen/logrus"
"io"
"os"
"github.com/adrg/xdg"
"fmt"
"time"
"github.com/a8m/envsubst"
"path/filepath"
"context"
"gopkg.in/yaml.v2"
"bjsh/installk8s/cluster"
)

type loghook struct {
	Writer    io.Writer
	Formatter log.Formatter

	levels []log.Level
}

var (
	configFlag = &cli.StringFlag{
		Name:      "config",
		Usage:     "Path to cluster config yaml. Use '-' to read from stdin.",
		Aliases:   []string{"c"},
		Value:     "k0sctl.yaml",
		TakesFile: true,
	}
	Colorize = aurora.NewAurora(false)
)

type ctxConfigKey struct{}

// initLogging initializes the logger   //这里引用了外部的log部件
func initLogging(ctx *cli.Context) error {
	log.SetLevel(log.TraceLevel)
	log.SetOutput(io.Discard)  
	return initFileLogger()   
}


const logPath = "k8sinstall/k8sinstall.log"

func LogFile() (io.Writer, error) {
	fn, err := xdg.SearchCacheFile(logPath)
	if err != nil {
		fn, err = xdg.CacheFile(logPath)
		if err != nil {
			return nil, err
		}
	}

	logFile, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open log %s: %s", fn, err.Error())
	}

	_, _ = fmt.Fprintf(logFile, "time=\"%s\" level=info msg=\"###### New session ######\"\n", time.Now().Format(time.RFC822))

	return logFile, nil
}

func initFileLogger() error {
	lf, err := LogFile()
	if err != nil {
		return err
	}
	log.AddHook(fileLoggerHook(lf))
	return nil
}

func fileLoggerHook(logFile io.Writer) *loghook {
	l := &loghook{
		Formatter: &log.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        time.RFC822,
			DisableLevelTruncation: true,
		},
		Writer: logFile,
	}

	l.SetLevel(log.DebugLevel)

	return l
}


func (h *loghook) SetLevel(level log.Level) {
	h.levels = []log.Level{}
	for _, l := range log.AllLevels {
		if level >= l {
			h.levels = append(h.levels, l)
		}
	}
}

func initConfig(ctx *cli.Context) error {
	f := ctx.String("config") //看看是否通过config设置了配置文件
	if f == "" {
		return nil
	}

	file, err := configReader(f)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	subst, err := envsubst.Bytes(content) //这里从读取内容生成环境变量
	if err != nil {
		return err
	}

	log.Debugf("Loaded configuration:\n%s", subst)   //这里又用到了外部库

	c := &cluster.Cluster{}
	if err := yaml.UnmarshalStrict(subst, c); err != nil {   //解析配置给c
		return err
	}

	m, err := yaml.Marshal(c) //序列化
	if err == nil {
		log.Debugf("unmarshaled configuration:\n%s", m)
	}
	ctx.Context = context.WithValue(ctx.Context, ctxConfigKey{}, c) //将解析到的配置给到context

	return nil
}

func configReader(f string) (io.ReadCloser, error) {
	//为-则从输入中读取配置
	if f == "-" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return nil, fmt.Errorf("can't stat stdin: %s", err.Error())
		}
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			return os.Stdin, nil
		}
		return nil, fmt.Errorf("can't read stdin")
	}

	variants := []string{f}
	// add .yml to default value lookup
	if f == "config.yaml" {
		variants = append(variants, "config.yml")
	}

	for _, fn := range variants {
		if _, err := os.Stat(fn); err != nil {
			continue
		}

		fp, err := filepath.Abs(fn)  //获取绝对路径
		if err != nil {
			return nil, err
		}
		file, err := os.Open(fp)  //打开配置文件
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	return nil, fmt.Errorf("failed to locate configuration")
}

func actions(funcs ...func(*cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		for _, f := range funcs {
			if err := f(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

func (h *loghook) Fire(entry *log.Entry) error {
	line, err := h.Formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to format log entry: %v", err)
		return err
	}
	_, err = h.Writer.Write(line)
	return err
}

func (h *loghook) Levels() []log.Level {
	return h.levels
}