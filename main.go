package main

import (
	"fmt"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli"
)

var addIncidentCommand = cli.Command{
	Name: "add",
	Usage: `Add Incident to Store
		usage: add -d <DESCRIPTION>`,
	Action: cmdAddIncident,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "description, d",
			Usage: "Incident Description",
		},
	},
}

var solveIncidentCommand = cli.Command{
	Name: "solve",
	Usage: `Solve Incident from store
		usage: solve -id <INCIDENT_ID>`,
	Action: cmdSolveIncident,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "Incident Id",
		},
	},
}

var incidentStatusCommand = cli.Command{
	Name: "status",
	Usage: `Get incidents status from store
		usage: status -d1 <DATE1> -d2 <DATE2>
		where: DATE1 and DATE2 are dates in format "yyyy-mm-dd hh:mm"`,
	Action: cmdIncidentStatus,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "d1",
			Usage: "First date",
		},
		cli.StringFlag{
			Name:  "d2",
			Usage: "Second date",
		},
	},
}

var newStore = NewStore()

func main() {
	app := cli.NewApp()
	app.Name = "Simple Store"
	app.Usage = "Construct a simple Store that has a collection of Incidents and an incident_status method"
	app.Commands = []cli.Command{
		addIncidentCommand,
		solveIncidentCommand,
		incidentStatusCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func cmdAddIncident(c *cli.Context) error {
	description := c.String("d")
	result, err := newStore.AddIncident(description)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	fmt.Println(result)
	return nil
}

func cmdSolveIncident(c *cli.Context) error {
	id := c.String("id")
	uuidFromString, uuidErr := uuid.FromString(id)
	if uuidErr != nil {
		return cli.NewExitError(uuidErr, 1)
	}
	result, resultErr := newStore.SolveIncident(uuidFromString)
	if resultErr != nil {
		return cli.NewExitError(resultErr, 1)
	}
	fmt.Println(result)
	return nil
}

func cmdIncidentStatus(c *cli.Context) error {
	date1 := c.String("d1")
	date2 := c.String("d2")
	parsedDate1, err := time.Parse("2006-01-02 15:04", date1)
	parsedDate2, err := time.Parse("2006-01-02 15:04", date2)
	if err != nil {
		panic(err)
	}
	result, err := newStore.IncidentStatus(parsedDate1, parsedDate2)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	return nil
}
