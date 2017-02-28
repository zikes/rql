package main

import (
	"fmt"

	goquadapter "git.nwaonline.com/rune/rql/adapters/goqu"
	sqladapter "git.nwaonline.com/rune/rql/adapters/sql"
	rql "git.nwaonline.com/rune/rql/parse"
	"github.com/spf13/cobra"
)

func main() {
	var cmdSql = &cobra.Command{
		Use:   "sql [string to parse]",
		Short: "Converts RQL to SQL",
		Long:  `sql converts RQL input into SQL output`,
		Run: func(cmd *cobra.Command, args []string) {
			t, err := rql.New("root").Parse(args[0])
			if err != nil {
				fmt.Println("Error parsing RQL: %s", err)
				return
			}
			fmt.Println(sqladapter.ToSQL(t.Root))
		},
	}
	var cmdGoqu = &cobra.Command{
		Use:   "goqu [string to parse]",
		Short: "Converts RQL to SQL via goqu",
		Long:  `sql converts RQL input into SQL output via goqu`,
		Run: func(cmd *cobra.Command, args []string) {
			t, err := rql.New("root").Parse(args[0])
			if err != nil {
				fmt.Println("Error parsing RQL: %s", err)
				return
			}
			fmt.Println(goquadapter.ToSQL(t.Root))
		},
	}

	var rootCmd = &cobra.Command{Use: "rql"}
	rootCmd.AddCommand(cmdSql)
	rootCmd.AddCommand(cmdGoqu)
	rootCmd.Execute()
}
