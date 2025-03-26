package schemas

import (
    "database/sql"
    "log"
    "os"
    "path/filepath"
    "sort"

    _ "github.com/lib/pq" 
)

func InitializeDatabase(db *sql.DB) error {
    
    log.Println("Executing update schema files...")
    err := executeSQLDir(db, filepath.Join("schemas", "update-schema"))
    if err != nil {
        return err
    }
    
    log.Println("Executing patch schema files...")
    err = executeSQLDirInOrder(db, filepath.Join("schemas", "patch-schema"))
    if err != nil {
        return err
    }

	log.Println("Executing bundled schema files...")
    err = executeSQLDir(db, filepath.Join("schemas", "bundled-schema"))
    if err != nil {
        return err
    }
    
    log.Println("Database initialization complete")
    return nil
}

func executeSQL(db *sql.DB, filePath string) error {
    log.Printf("Executing SQL file: %s", filePath)
    
    content, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }
    
    _, err = db.Exec(string(content))
    return err
}

func executeSQLDir(db *sql.DB, dirPath string) error {
    files, err := os.ReadDir(dirPath)
    if err != nil {
        return err
    }
    
    for _, file := range files {
        if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
            err = executeSQL(db, filepath.Join(dirPath, file.Name()))
            if err != nil {
                return err
            }
        }
    }
    
    return nil
}

func executeSQLDirInOrder(db *sql.DB, dirPath string) error {
    files, err := os.ReadDir(dirPath)
    if err != nil {
        return err
    }
    
    var sqlFiles []string
    for _, file := range files {
        if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" && file.Name() != ".gitkeep" {
            sqlFiles = append(sqlFiles, file.Name())
        }
    }
    
    sort.Strings(sqlFiles)
    
    for _, fileName := range sqlFiles {
        err = executeSQL(db, filepath.Join(dirPath, fileName))
        if err != nil {
            return err
        }
    }
    
    return nil
}