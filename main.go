package main

import (
	"fmt"
	"html/template"
	"os"
)

func convertSecondsToPace(seconds int) string {
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, secs)
}

func calculateTrainings(dagen int, dubbeldagen int, startMin, startSec, doelMin, doelSec int) []int {
	trainingsdagen := make([]int, dagen)

	startSeconds := 60*startMin + startSec
	doelSeconds := 60*doelMin + doelSec

	// Aantal dagen waarop het tempo écht moet dalen (normale dagen)
	aantalStappen := dagen - dubbeldagen - 1 // -1 omdat we al op dag 1 starten
	
	if aantalStappen < 1 {
		aantalStappen = 1
	}
	
	// Totale te verbeteren seconden
	teVerbeteren := startSeconds - doelSeconds

	// Seconden per normale stap (integer deling)
	split := teVerbeteren / aantalStappen

	// We houden de "rest" over om aan het eind exact uit te komen
	rest := teVerbeteren % aantalStappen

	dubbeldagInterval := 0
	if dubbeldagen > 0 {
		dubbeldagInterval = dagen / dubbeldagen
	}

	trainingsdagen[0] = startSeconds
	current := startSeconds

	for i := 1; i < dagen; i++ {
		if dubbeldagInterval == 0 || (i+1)%dubbeldagInterval != 0 {
			// Normale dag: tempo verbetert
			current -= split

			// Verdeel de rest over de eerste paar normale dagen zodat we exact eindigen
			if rest > 0 {
				current--
				rest--
			}
		} 
		// else: dubbeldag → current blijft hetzelfde

		trainingsdagen[i] = current
	}

	// Forceer de laatste dag altijd op exact het doeltempo (veiligheid)
	trainingsdagen[dagen-1] = doelSeconds

	return trainingsdagen
}

// getHtmlFile maakt een mooi HTML-tabelletje en slaat het op als loopschema.html
func getHtmlFile(schema []int) error {
	const tmpl = `
<!DOCTYPE html>
<html lang="nl">
<head>
    <meta charset="UTF-8">
    <title>Loopschema</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            text-align: center;
            background-color: #f9f9f9;
        }
        h1 {
            color: #333;
        }
        table {
            border-collapse: collapse;
            margin: 30px auto;
            box-shadow: 0 4px 8px rgba(0,0,0,0.1);
        }
        th, td {
            border: 1px solid #aaa;
            padding: 12px 20px;
            text-align: center;
        }
        th {
            background-color: #4CAF50;
            color: white;
        }
        tr:nth-child(even) {
            background-color: #f2f2f2;
        }
        tr:hover {
            background-color: #e0f7e0;
        }
        .pace {
            font-weight: bold;
            font-size: 1.1em;
        }
    </style>
</head>
<body>
    <h1>🏃 Mijn Loopschema</h1>
    <p><strong>Aantal trainingsdagen:</strong> {{len .}}</p>
    
    <table>
        <thead>
            <tr>
                <th>Dag</th>
                <th>Tempo (min/km)</th>
            </tr>
        </thead>
        <tbody>
            {{range $index, $seconds := .}}
            <tr>
                <td><strong>{{add1 $index}}</strong></td>
                <td class="pace">{{convertPace $seconds}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>

    <p><em>Genereerd met Go - Succes met trainen! 💪</em></p>
</body>
</html>`

	// Maak template met extra functies
	funcMap := template.FuncMap{
		"add1":       func(i int) int { return i + 1 },
		"convertPace": convertSecondsToPace,
	}

	t := template.Must(template.New("schema").Funcs(funcMap).Parse(tmpl))

	// Maak het HTML bestand
	file, err := os.Create("loopschema.html")
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, schema)
}

func main() {
	var dagen, dubbeldagen, startMin, startSec, doelMin, doelSec int

	fmt.Print("Aantal totale trainingsdagen: ")
	fmt.Scan(&dagen)

	fmt.Print("Aantal dubbeldagen (dagen met hetzelfde tempo): ")
	fmt.Scan(&dubbeldagen)

	fmt.Print("Starttempo in min sec / km (bijv. 6 25): ")
	fmt.Scan(&startMin, &startSec)

	fmt.Print("Doeltempo in min sec / km (bijv. 4 30): ")
	fmt.Scan(&doelMin, &doelSec)

	schema := calculateTrainings(dagen, dubbeldagen, startMin, startSec, doelMin, doelSec)

	err := getHtmlFile(schema)
	if err != nil {
		fmt.Println("Fout bij genereren van HTML:", err)
	} else {
		fmt.Println("✅ HTML schema succesvol gegenereerd: loopschema.html")
	}

	fmt.Println("\nDruk op Enter om af te sluiten...")
	fmt.Scanln()
}
