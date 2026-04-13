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

// Bereken snelheid in km/u met 1 decimaal
func paceToSpeedKmh(paceSeconds int) float64 {
	if paceSeconds == 0 {
		return 0.0
	}
	return 3600.0 / float64(paceSeconds) // seconden per km → km per uur
}

// Formateer snelheid naar 1 decimaal (bijv. 9.6)
func formatSpeedKmh(speed float64) string {
	return fmt.Sprintf("%.1f", speed)
}

// Totale tijd berekenen
func calculateTotalTime(paceSeconds int, km int) int {
	return paceSeconds * km
}

func formatTotalTime(totalSeconds int) string {
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	secs := totalSeconds % 60

	if hours > 0 {
		return fmt.Sprintf("%du %02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

func calculateTrainings(dagen int, dubbeldagen int, startMin, startSec, doelMin, doelSec int) []int {
	trainingsdagen := make([]int, dagen)

	startSeconds := 60*startMin + startSec
	doelSeconds := 60*doelMin + doelSec

	aantalStappen := dagen - dubbeldagen - 1
	if aantalStappen < 1 {
		aantalStappen = 1
	}

	teVerbeteren := startSeconds - doelSeconds
	split := teVerbeteren / aantalStappen
	rest := teVerbeteren % aantalStappen

	dubbeldagInterval := 0
	if dubbeldagen > 0 {
		dubbeldagInterval = dagen / dubbeldagen
	}

	trainingsdagen[0] = startSeconds
	current := startSeconds

	for i := 1; i < dagen; i++ {
		if dubbeldagInterval == 0 || (i+1)%dubbeldagInterval != 0 {
			current -= split
			if rest > 0 {
				current--
				rest--
			}
		}
		trainingsdagen[i] = current
	}

	trainingsdagen[dagen-1] = doelSeconds
	return trainingsdagen
}

// ==================== HTML GENERATIE ====================
func getHtmlFile(schema []int, km int) error {
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
        h1 { color: #333; }
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
        .pace { font-weight: bold; font-size: 1.1em; }
        .speed { font-weight: bold; color: #e67e22; }
        .total { font-weight: bold; color: #2c7a2c; }

        /* Print support */
        @media print {
            body { background-color: white; }
            th { background-color: #4CAF50 !important; color: white !important; }
            tr:nth-child(even) { background-color: #f2f2f2 !important; }
        }
    </style>
</head>
<body>
    <h1>🏃 Mijn Loopschema</h1>
    <p><strong>Aantal trainingsdagen:</strong> {{len .Schema}} &nbsp;&nbsp;&nbsp; <strong>Afstand:</strong> {{.Km}} km</p>
   
    <table>
        <thead>
            <tr>
                <th>Dag</th>
                <th>Tempo (min/km)</th>
                <th>Snelheid (km/u)</th>
                <th>Totale tijd ({{.Km}} km)</th>
            </tr>
        </thead>
        <tbody>
            {{range $index, $paceSeconds := .Schema}}
            <tr>
                <td><strong>{{add1 $index}}</strong></td>
                <td class="pace">{{convertPace $paceSeconds}}</td>
                <td class="speed">{{formatSpeedKmh (paceToSpeed $paceSeconds)}}</td>
                <td class="total">{{formatTotalTime (calcTotal $paceSeconds $.Km)}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>

    <p><em>Genereerd met Go - Succes met trainen! 💪</em></p>
</body>
</html>`

	type PageData struct {
		Schema []int
		Km     int
	}

	data := PageData{
		Schema: schema,
		Km:     km,
	}

	funcMap := template.FuncMap{
		"add1":            func(i int) int { return i + 1 },
		"convertPace":     convertSecondsToPace,
		"paceToSpeed":     paceToSpeedKmh,
		"formatSpeedKmh":  formatSpeedKmh,
		"calcTotal":       calculateTotalTime,
		"formatTotalTime": formatTotalTime,
	}

	t := template.Must(template.New("schema").Funcs(funcMap).Parse(tmpl))

	file, err := os.Create("loopschema.html")
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, data)
}

func main() {
	var dagen, dubbeldagen, startMin, startSec, doelMin, doelSec, km int

	fmt.Print("Aantal totale trainingsdagen: ")
	fmt.Scan(&dagen)

	fmt.Print("Aantal dubbeldagen (dagen met hetzelfde tempo): ")
	fmt.Scan(&dubbeldagen)

	fmt.Print("Aantal kilometers per training: ")
	fmt.Scan(&km)

	fmt.Print("Starttempo in min sec / km (bijv. 6 25): ")
	fmt.Scan(&startMin, &startSec)

	fmt.Print("Doeltempo in min sec / km (bijv. 4 30): ")
	fmt.Scan(&doelMin, &doelSec)

	schema := calculateTrainings(dagen, dubbeldagen, startMin, startSec, doelMin, doelSec)

	err := getHtmlFile(schema, km)
	if err != nil {
		fmt.Println("Fout bij genereren van HTML:", err)
	} else {
		fmt.Println("✅ HTML schema succesvol gegenereerd: loopschema.html")
	}

	fmt.Println("\nDruk op Enter om af te sluiten...")
	fmt.Scanln()
}
