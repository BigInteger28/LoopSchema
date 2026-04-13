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

	// Voorkom deling door 0 of negatief
	aantalNormaleDagen := dagen - dubbeldagen
	if aantalNormaleDagen < 2 {
		aantalNormaleDagen = 2
	}
	split := (startSeconds - doelSeconds) / (aantalNormaleDagen - 1)

	dubbeldagInterval := 0
	if dubbeldagen > 0 {
		dubbeldagInterval = dagen / dubbeldagen
	}

	trainingsdagen[0] = startSeconds

	for i := 1; i < dagen; i++ {
		if dubbeldagInterval == 0 || (i+1)%dubbeldagInterval != 0 {
			// Normale dag: tempo verbetert
			trainingsdagen[i] = trainingsdagen[i-1] - split
		} else {
			// Dubbeldag: hetzelfde tempo
			trainingsdagen[i] = trainingsdagen[i-1]
		}
	}

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