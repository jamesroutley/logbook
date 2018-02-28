package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/jamesroutley/logbook/summary"
	"github.com/spf13/cobra"
)

var week bool

func init() {
	summaryCmd.Flags().BoolVarP(&week, "week", "w", false, "Print a summary of the week commencing on START. START must be a Monday.")
	rootCmd.AddCommand(summaryCmd)
}

var summaryCmd = &cobra.Command{
	Use:   "summary START [END]",
	Short: "Prints a summary of time period between START and END, inclusive",
	Long: `summary searches through all logfiles between START and END, inclusive, 
and prints out all summaries, organised by topic. You can add a new summary with the
syntax:

- SUMMARY:<Topic>: <summary>

Summaries with the same topic are grouped together. The topic is optional. If it isn't
included, the summary will be added to a Misc topic`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Number of args
		if week {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
		} else {
			if err := cobra.ExactArgs(2)(cmd, args); err != nil {
				return err
			}
		}
		// All args should be date-parseable
		for i, arg := range args {
			_, err := time.Parse("2006-01-02", arg)
			if err != nil {
				datePos := []string{"start", "end"}
				fmt.Printf("%s date should be in YYYY-MM-DD format\n", datePos[i])
				os.Exit(1)
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var start, end time.Time
		if week {
			start = forceStrToTime(args[0])
			if start.Weekday() != time.Monday {
				fmt.Println("start date must be a Monday")
				os.Exit(1)
			}
			end = start.AddDate(0, 0, 6)
		} else {
			start = forceStrToTime(args[0])
			end = forceStrToTime(args[1])
		}
		s := summary.Summarise(start, end)
		fmt.Printf(
			"James Routley Summary: %s to %s\n\n",
			start.Format("2006-01-02"), end.Format("2006-01-02"),
		)
		fmt.Println(s)
	},
}

func forceStrToTime(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}
