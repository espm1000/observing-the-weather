$csv = Import-Csv .\weather-observability-github-project.csv


foreach ($row in $csv) {
    gh issue create `
        --title $row.Title `
        --body $row.Description
}
