package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pixb/go-server/internal/profile"
	"github.com/pixb/go-server/internal/version"
	"github.com/pixb/go-server/server"
	"github.com/pixb/go-server/store"
	"github.com/pixb/go-server/store/db/mysql"
	"github.com/pixb/go-server/store/db/postgresql"
	"github.com/pixb/go-server/store/db/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 定义 rootCmd 引用
// combra
// Use
// Short
// Long
// Run
var rootCmd = &cobra.Command{
	Use:   "go-server",
	Short: "go-server demo.",
	Long:  "go-server demo",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("===rootCmd Run...===")
		prof := &profile.Profile{
			Demo:   viper.GetBool("demo"),
			Addr:   viper.GetString("addr"),
			Port:   viper.GetInt("port"),
			Data:   viper.GetString("data"),
			Driver: viper.GetString("driver"),
			DSN:    viper.GetString("dsn"),
			Secret: viper.GetString("secret"),
		}
		prof.Version = version.GetCurrentVersion()

		return run(cmd.Context(), prof)
	},
}

// init() 方法
func init() {
	viper.SetDefault("demo", false)
	viper.SetDefault("driver", "sqlite")
	viper.SetDefault("port", 8081)

	rootCmd.PersistentFlags().Bool("demo", false, "enable demo")
	rootCmd.PersistentFlags().String("addr", "0.0.0.0", "bind address")
	rootCmd.PersistentFlags().Int("port", 8081, "port of server")
	rootCmd.PersistentFlags().String("data", "./data", "data directory")
	rootCmd.PersistentFlags().String("driver", "sqlite", "data driver")
	rootCmd.PersistentFlags().String("dsn", "", "database connection string")
	rootCmd.PersistentFlags().String("secret", "your-secret-key", "Secret key for authentication")

	if err := viper.BindPFlag("demo", rootCmd.PersistentFlags().Lookup("demo")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("data", rootCmd.PersistentFlags().Lookup("data")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("driver", rootCmd.PersistentFlags().Lookup("driver")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("dsn", rootCmd.PersistentFlags().Lookup("dsn")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("secret", rootCmd.PersistentFlags().Lookup("secret")); err != nil {
		panic(err)
	}

	viper.BindPFlags(rootCmd.Flags())
	viper.SetEnvPrefix("GO_SERVER")
	viper.AutomaticEnv()
}

// 运行服务器定义
// 1.检查配置是否正确
// 2.创建数据目录
// 3.创建数据驱动
// 4.创建存储实例并且迁移数据
// 5.创建服务实例，启动服务
// 6.处理优雅停机
func run(ctx context.Context, prof *profile.Profile) error {
	// 1.检查配置是否正确
	if err := prof.Validate(); err != nil {
		return err
	}

	fmt.Printf("== prof.Data: %s\n", prof.Data)

	// 2.创建数据目录
	if err := os.MkdirAll(prof.Data, 0755); err != nil {
		return err
	}

	// 3.创建数据驱动
	var dbDriver store.Driver
	var err error

	switch prof.Driver {
	case "sqlite":
		dbDriver, err = sqlite.NewDriver(prof)
	case "postgresql":
		dbDriver, err = postgresql.NewDriver(prof)
	case "mysql":
		dbDriver, err = mysql.NewDriver(prof)
	default:
		return fmt.Errorf("unsupported database driver: %s", prof.Driver)
	}

	if err != nil {
		return err
	}

	// 4.创建存储实例并且迁移数据
	storeInstance := store.New(dbDriver, prof)
	if err := storeInstance.Migrate(ctx); err != nil {
		return err
	}

	// 5.创建服务实例，启动服务
	s, err := server.NewServer(ctx, prof, storeInstance)
	if err != nil {
		return err
	}

	if err := s.Start(ctx); err != nil {
		return err
	}

	// 6.处理优雅停机
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return s.Shutdown(ctx)
}

// main方法执行 rootCmd
func main() {
	fmt.Println("==============main===================")
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
