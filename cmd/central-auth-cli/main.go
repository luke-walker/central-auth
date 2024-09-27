package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/joho/godotenv"

    "central-auth/internal/db"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("error loading .env file: %v", err)
    }

    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")
    dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

    db, err := database.NewDatabase(dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Println("=== Central Auth CLI ===")
        fmt.Println("1. List Servers")
        fmt.Println("2. Add Server")
        fmt.Println("3. Remove Server (WIP)")
        fmt.Println("4. Edit Server (WIP)")
        fmt.Println("5. List Users (WIP)")
        fmt.Println("6. Add User (WIP)")
        fmt.Println("7. Remove User (WIP)")
        fmt.Println("8. Promote/Demote User")
        fmt.Println("0. Exit")
        fmt.Print("Enter number: ")

        text := getInput(reader)
        fmt.Println()
        switch (text) {
        case "1":
            fmt.Println("=== Servers List ===")

            result, numRows, err := db.GetAllServerInfo()
            if err != nil {
                fmt.Printf("Failed to retrieve servers: %v\n", err)
                break
            }
            if numRows == 0 {
                fmt.Println("No servers found")
                break
            }

            for _, info := range result {
                fmt.Printf("%s:\n", info.Name)

                fmt.Println("\tAddresses...")
                if len(info.Addrs) == 0 {
                    fmt.Println("\n\n<not set>")
                }
                for _, addr := range info.Addrs {
                    fmt.Printf("\t\t%s\n", addr)
                }

                fmt.Print("\tProxy URL: ")
                if info.Proxy == "" {
                    fmt.Println("<not set>")
                } else {
                    fmt.Println(info.Proxy)
                }

                fmt.Print("\tRedirect URL: ")
                if info.Redirect == "" {
                    fmt.Println("<not set>")
                } else {
                    fmt.Println(info.Redirect)
                }
            }

            break
        case "2":
            fmt.Println("=== Add Server ===")

            name := ""
            for name == "" {
                fmt.Print("Name: ")
                name = getInput(reader)
            }

            fmt.Print("Addresses (comma-separated): ")
            addrsStr := getInput(reader)
            var addrs []string
            if addrsStr != "" {
                addrs = strings.Split(addrsStr, ",")
            }

            fmt.Print("Proxy URL: ")
            proxy := getInput(reader)

            redirect := ""
            for redirect == "" {
                fmt.Print("Redirect URL: ")
                redirect = getInput(reader)
            }

            if err = db.CreateServer(name, addrs, proxy, redirect); err != nil {
                fmt.Printf("Failed to create server: %v", err)
                break
            }
            fmt.Println("Successfully created server")

            break
        case "3":
            break
        case "4":
            break
        case "5":
            break
        case "6":
            break
        case "7":
            break
        case "8":
            fmt.Println("=== Promote/Demote User ===")
            
            username := ""
            for username == "" {
                fmt.Print("Username: ")
                username = getInput(reader)
            }

            userInfo, numRows, err := db.GetUserInfoByUsername(username)
            if err != nil {
                fmt.Printf("Error retrieving user: %v\n", err)
                break
            }
            if numRows == 0 {
                fmt.Printf("User '%s' not found\n", username)
                break
            }

            if userInfo.Admin {
                fmt.Printf("Would you like to demote '%s' from admin to normal privileges? ('yes' to confirm) ", username)
                if getInput(reader) != "yes" {
                    break
                }
                err = db.SetUserAdminStatusByToken(false, userInfo.Token)
            } else {
                fmt.Printf("Would you like to promote '%s' from normal to admin privileges? ('yes' to confirm) ", username)
                if getInput(reader) != "yes" {
                    break
                }
                err = db.SetUserAdminStatusByToken(true, userInfo.Token)
            }
            if err != nil {
                fmt.Printf("Error updating user admin status: %v\n", err)
                break
            }
            fmt.Println("User admin status updated")

            break
        case "0":
            os.Exit(0)
        default:
            fmt.Println("Invalid input.")
            break
        }
        fmt.Println()
    }
}

func getInput(reader *bufio.Reader) string {
    text, err := reader.ReadString('\n')
    if err != nil {
        return ""
    }
    return strings.TrimRight(text, "\n")
} 
