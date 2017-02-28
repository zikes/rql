package main

import (
	"fmt"

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

	var rootCmd = &cobra.Command{Use: "rql"}
	rootCmd.AddCommand(cmdSql)
	rootCmd.Execute()
}
