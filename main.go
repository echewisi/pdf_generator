package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "runtime"

    "github.com/jung-kurt/gofpdf"
)

type Transaction struct {
    Date        string  `json:"date"`
    Description string  `json:"description"`
	MoneyOut    float64 `json:"money_out"`
	MoneyIn     float64 `json:"money_in"`
    Balance     float64 `json:"balance"`
}

type BalanceSummary struct {
	Product        string `json:"product"`
    OpeningBalance float64 `json:"opening_balance"`
    ClosingBalance float64 `json:"closing_balance"`
    MoneyIn        float64 `json:"money_in"`
    MoneyOut       float64 `json:"money_out"`
}

type AccountStatement struct {
    CompanyName     string        `json:"company_name"`
    CompanyAddress  string        `json:"company_address"`
    CustomerName    string        `json:"customer_name"`
    CustomerAddress string        `json:"customer_address"`
    AccountName     string        `json:"account_name"`
    AccountNumber   string        `json:"account_number"`
	ReportGenerationDate  string   `json:"report_generation_date"`
    BalanceSummary  []BalanceSummary `json:"balance_summary"`
    Transactions    []Transaction `json:"transactions"`
	// BalanceDetails []BalanceSummary `json:balance_summary`
}

func openPDF(filename string) {
    var cmd *exec.Cmd

    switch runtime.GOOS {
    case "linux":
        cmd = exec.Command("xdg-open", filename)
    case "windows":
        cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", filename)
    case "darwin":
        cmd = exec.Command("open", filename)
    default:
        log.Fatalf("Unsupported platform")
    }

    err := cmd.Start()
    if err != nil {
        log.Fatalf("Error opening PDF: %v", err)
    }
}

func formatCurrency(value interface{}) string {
    if value == "" {
        return ""
    }
    return fmt.Sprintf("$%.2f", value)
}

func main() {
    // Read the JSON file
    jsonFile, err := os.Open("account_statement.json")
    if err != nil {
        log.Fatalf("Error opening JSON file: %v", err)
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    var accountStatement AccountStatement
    json.Unmarshal(byteValue, &accountStatement)

    // Calculate Money In and Money Out
    // moneyIn, moneyOut := calculateMoneyInOut(accountStatement.Transactions)
    // accountStatement.BalanceSummary.MoneyIn = moneyIn
    // accountStatement.BalanceSummary.MoneyOut = moneyOut

    // Create a new PDF document
    pdf := gofpdf.New("P", "mm", "A4", "")

    // Add a page
    pdf.AddPage()

    // Add a logo
    pdf.Image("logo.png", 10, 10, 30, 0, false, "", 0, "")
	pdf.Ln(10)

    // Company name and address
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(0, 10, accountStatement.CompanyName)
    pdf.Ln(10)
    pdf.SetFont("Arial", "", 12)
    pdf.MultiCell(0, 10, accountStatement.CompanyAddress, "", "", false)
    pdf.Ln(10)

    // Customer name and address
    pdf.SetFont("Arial", "B", 12)
    pdf.Cell(0, 10, fmt.Sprintf("Customer: %s", accountStatement.CustomerName))
    pdf.Ln(10)
    pdf.SetFont("Arial", "", 12)
    pdf.MultiCell(0, 10, accountStatement.CustomerAddress, "", "", false)
    pdf.Ln(10)

    // Account details
    pdf.SetFont("Arial", "B", 14)
    pdf.Cell(3, 20, "USD Statement")
	pdf.Ln(10)
    pdf.SetFont("Arial", "I", 8)
    pdf.Cell(0, 10, fmt.Sprintf("Generated on: %s", accountStatement.ReportGenerationDate))
	pdf.Ln(10)
    pdf.SetFont("Arial", "I", 8)
    pdf.Cell(0, 10, fmt.Sprintf("issued by: %s", accountStatement.CompanyName))
    pdf.Ln(25)

    // Balance summary
    pdf.SetFont("Arial", "B", 12)
    pdf.Cell(0, 10, "Balance Summary")
    pdf.Ln(10)
    pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(40, 10, "Product", "1", 0, "C", false, 0, "" )
    pdf.CellFormat(40, 10, "Opening balance", "1", 0, "C", false, 0, "")
    pdf.CellFormat(40, 10, "Money Out", "1", 0, "C", false, 0, "")    
    pdf.CellFormat(40, 10, "Money In", "1", 0, "C", false, 0, "")   
	pdf.CellFormat(40, 10, "Closing balance", "1", 0, "C", false, 0, "")
	pdf.Ln(10)
    

	pdf.SetFont("Arial", "", 12)

	for _, details := range accountStatement.BalanceSummary {
		pdf.SetFont("Arial", "", 8)
		pdf.CellFormat(40, 10, details.Product, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, formatCurrency(details.OpeningBalance), "1", 0, "R", false, 0, " ")
		pdf.CellFormat(40, 10, formatCurrency(details.MoneyOut), "1", 0, "R", false, 0, " ")
		pdf.CellFormat(40, 10, formatCurrency(details.MoneyIn), "1", 0, "R", false, 0, " ")
		pdf.CellFormat(40, 10, formatCurrency(details.ClosingBalance), "1", 0, "R", false, 0, " ")
		pdf.Ln(10)
	}

    // Add table header
    pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Account Statement For Transactions In The Month So Far")
	pdf.Ln(10)
    pdf.SetFont("Arial", "", 8)
    pdf.CellFormat(40, 10, "Date", "1", 0, "C", false, 0, "")
    pdf.CellFormat(70, 10, "Description", "1", 0, "L", false, 0, "")
    pdf.CellFormat(30, 10, "Money out", "1", 0, "L", false, 0, "")
    pdf.CellFormat(30, 10, "Money in", "1", 0, "L", false, 0, "")
    pdf.CellFormat(30, 10, "Balance", "1", 0, "L", false, 0, "")
    pdf.Ln(10)

    // Set font for table content
    pdf.SetFont("Arial", "", 8)

    // Add table content
    for _, t := range accountStatement.Transactions {
		pdf.SetFont("Arial", "", 8)
        pdf.CellFormat(40, 10, t.Date, "1", 0, "", false, 0, "")
        pdf.CellFormat(70, 10, t.Description, "1", 0, "", false, 0, "")
        pdf.CellFormat(30, 10, formatCurrency(t.MoneyOut), "1", 0, "R", false, 0, "")
        pdf.CellFormat(30, 10, formatCurrency(t.MoneyIn), "1", 0, "R", false, 0, "")
        pdf.CellFormat(30, 10, formatCurrency(t.Balance), "1", 0, "R", false, 0, "")
        pdf.Ln(10)
    }

    // Save the PDF to a file
    outputFilename := "account_statement.pdf"
    err = pdf.OutputFileAndClose(outputFilename)
    if err != nil {
        log.Fatalf("Error creating PDF file: %v", err)
    }

    fmt.Println("PDF generated successfully.")

    // Open the PDF
    openPDF(outputFilename)
}
