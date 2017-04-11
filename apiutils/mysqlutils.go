package apiutils

import (
    
    "os"
    "io"
    "log"
    "fmt"
    "time"
    "strconv"
    "net/http"

    "github.com/Sirupsen/logrus"
    "github.com/ziutek/mymysql/mysql"
    _ "github.com/ziutek/mymysql/native"
    
)

// Channels for requesting and providing data connection
var (

    
    //Get MySQL and Mongo configuration from environment variables
    clientdb_server         = os.Getenv("MYSQL_SERVER")
    clientdb_port           = os.Getenv("MYSQL_PORT")
    clientdb_read_user      = os.Getenv("MYSQL_READ_USER")
    clientdb_read_pass      = os.Getenv("MYSQL_READ_PASS")
    clientdb_user           = os.Getenv("MYSQL_USER")
    clientdb_pass           = os.Getenv("MYSQL_PASS")
    clientdb_schema_db      = ""
    clientdb_dbtype         = "mysql"
    
    maxWorker, _ = strconv.Atoi(os.Getenv("MAX_WORKERS"))
    
)



type UpsertQuery struct {
    Query           string
    Parameters      []interface{}
}


// Function to be used by API functions to run MySQL query
func RunSelectQuery(schema, query string, parameters []interface{}, getTotalCount bool) ([]map[string]interface{}, int, string, int) {
    
    // Variables 
    var rowMapSlice []map[string]interface{}
    var totalCount int
    
    db, err := openMySqlDB(clientdb_server, clientdb_port, clientdb_read_user, clientdb_read_pass, "", clientdb_dbtype)
    if err != nil {
        return rowMapSlice, totalCount, err.Error(), http.StatusInternalServerError
    }
   
    // Change default schema to the one passed by function
    if err := db.Use(schema); err != nil {
        
        db.Close()
        return rowMapSlice, totalCount, fmt.Sprintf("Schema %s does not exist", schema), http.StatusUnprocessableEntity

    }
    
    
    // Prepare query. Using prepared queries because this forces MySQL to return the actual type
    if stmt, err := db.Prepare(query); err != nil { 
        
        logrus.WithFields(logrus.Fields{
            "err": err,
            "query": query,
        }).Warn("Error preparing query")
        
        // Send internal server error back
        db.Close()
        return rowMapSlice, totalCount, "", http.StatusInternalServerError

        
    } else {
        
        // Run query
        if res, err := stmt.Run(parameters...); err != nil {
        
            logrus.WithFields(logrus.Fields{
                "err": err,
                "query": query,
            }).Warn("Error running query")
            
            // Send internal server error back
            db.Close()
            return rowMapSlice, totalCount, "", http.StatusInternalServerError
            
            
        } else {
            
            // Process results
            rowMapSlice = rowsToMapSlice(res)
            
            // If function requests to get Total Count along with results
            if getTotalCount {
            
                // Get total count of detail rows disregarding the LIMIT in the query
                stmt2, _ := db.Prepare("select FOUND_ROWS()")
                res2, err2 := stmt2.Run()
                if err2 != nil {
                
                    logrus.WithFields(logrus.Fields{
                            "err": err2,
                    }).Warn("Error running FOUND_ROWS() detail query")
                    
                    // Send internal server error back
                    db.Close()
                    return rowMapSlice, totalCount, "", http.StatusInternalServerError
                }
                
                // Process resuts
                rowMapSlice2 := rowsToMapSlice(res2)
                
                 // Query only returns one row/one column, so only get that value
                for _, v := range rowMapSlice2[0] {
                    if b, ok := v.(int); ok {
                        totalCount = b
                    } else if b, ok := v.(int32); ok {
                        totalCount= int(b)
                    } else if b, ok := v.(int64); ok {
                        totalCount = int(b)
                    }else {

                        // Send internal server error back
                        db.Close()
                        return rowMapSlice, totalCount, "", http.StatusInternalServerError
                    }
                }
            }
        }

    }
    
    // Send results back
    db.Close()
    return rowMapSlice, totalCount, "", 0
    
}

            


// Function to be used by API functions to run MySQL query
func RunUpsertQueries(upsertQueries []UpsertQuery, getLastInsertId bool) (int, int, string, int) {
    
    lastInsertId := 0
    affectedRows := 0
    
    log.Printf("About to connect to mySqlDB %s. User %s\n", clientdb_server, clientdb_user)
    
    // Open connection to MySQL
    db, err := openMySqlDB(clientdb_server, clientdb_port, clientdb_user, clientdb_pass, "", clientdb_dbtype)
    if err != nil {
        return lastInsertId, affectedRows, err.Error(), http.StatusInternalServerError
    }
    
    log.Println("Connected to MySQLDB!")
    
    // Slice in which we will be placing all the prepared statements
    statements := make([]struct{
            stmt        mysql.Stmt
            parameters  []interface{}
    }, len(upsertQueries), len(upsertQueries))
    
    
    // Iterate through all upsertQueries and prepare them
    for i, q := range upsertQueries {
        
        // Prepare query
        if stmt, err := db.Prepare(q.Query); err != nil { 
            
            logrus.WithFields(logrus.Fields{
                "err": err,
                "query": q.Query,
            }).Warn("Error preparing query")
            
            // Send internal server error back
            db.Close()
            return lastInsertId, affectedRows, "", http.StatusInternalServerError
            
            
        } else {
            tmpStmt := struct{
                stmt        mysql.Stmt
                parameters  []interface{}
            }{stmt, q.Parameters}
            
            statements[i] = tmpStmt
        }
    }
    
    tx, err := db.Begin()
    if err != nil {
        
        db.Close()
        return lastInsertId, affectedRows, "", http.StatusInternalServerError
        
        

    }
    
    // Iterate through statements and run them
    for _, s := range statements {
        
        

        if res, err := tx.Do(s.stmt).Run(s.parameters...); err != nil {
            
            logrus.WithFields(logrus.Fields{
                "err": err,
            }).Warn("Error running query")
            
            // Send internal server error back
            tx.Rollback()
            db.Close()
            return lastInsertId, affectedRows, err.Error(), http.StatusInternalServerError

        } else {
            affectedRows = int(res.AffectedRows())
        }
    }
    
    // If we need to get the primary key of the last row that was inserted
    if getLastInsertId {
        
        // Run LAST_INSERT_ID() query
        if res, err := tx.Start("select LAST_INSERT_ID()"); err != nil {
            logrus.WithFields(logrus.Fields{
                "err": err,
            }).Warn("Error running LAST_INSERT_ID query")
            
            // Send internal server error back
            tx.Rollback()
            db.Close()
            return lastInsertId, affectedRows, "", http.StatusInternalServerError
            
        } else {
            
            // Iterate through returned rows (which should technically only be one)
            for {
                // GetRow
                if row, err := res.GetRow(); err != nil {
                    logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Warn("Error getting row from LAST_INSERT_ID query")
                    
                    tx.Rollback()
                    db.Close()
                    return lastInsertId, affectedRows, "", http.StatusInternalServerError
                    
                } else {
                    if row == nil {
                        // End of first result
                        break
                    }

                    // Assign primary key to lastInsertId
                    lastInsertId = row.Int(0)
                }
            }
        }
    }
    
    // Send results back
    tx.Commit()
    db.Close()
    return lastInsertId, affectedRows, "", 0
    
}



// Function that parses rows returned from a MySQL query and returns a map slice. Each map in the slice represents a single row in which the key is the column name and value is the actual value for that row.
func rowsToMapSlice(res mysql.Result) []map[string]interface{} {
    
    // Slice to return
    var retSlice []map[string]interface{}
    
    // Get []string of column names
    columns := res.Fields()
    
    row := res.MakeRow()
    for {
    
        // Scan rows one by one
        if err := res.ScanRow(row); err == io.EOF {
             // No more rows
             break
             
        // Error while scanning
        } else if err != nil {
            logrus.WithFields(logrus.Fields{
                "err": err,
            }).Warn("Error scanning row")
        }
        
        // Temporary map for current row. Will be appended to retSlice after
        tempMap := make(map[string]interface{})
        
        // Iterate through all values in row
        for i := 0; i < len(row); i++ {
                
            // v is of interface{} type because we don't know beforehand the type of variable we will be reading
            var v interface{}
            
            // Get actual value
            val := row[i]
            
            
            // If the data type of the actual value is []byte, convert to string
            if b, ok := val.([]byte); ok {
                v = string(b)
                
            // If the data type of the actual value is time.Time, convert to CustomTime, which unmarshals and Marshals into JSON in a desired format
            } else if t, ok := val.(time.Time); ok {
                var ct CustomTime
                ct.Time = t
                v = ct
            
            // Else, just return the value            
            } else {
                v = val
            }

            // Add new value to tempMap
            tempMap[columns[i].Name] = v

        }
        
        // Append tempMap to retSlice
        retSlice = append(retSlice, tempMap)
  
    }
    
    return retSlice
}






// openMySqlDB connects to MySQL database and launches go routine that provides db connection to processes. This function should only be called from main.go
func openMySqlDB(clientdb_server, clientdb_port, clientdb_user, clientdb_pass, clientdb_schema_db, clientdb_dbtype string) (mysql.Conn, error) {

    // Define DB connection
    db := mysql.New("tcp", "", fmt.Sprintf("%s:%s", clientdb_server, clientdb_port), clientdb_user, clientdb_pass)
    

    // Connect to MySQL DB
    if err := db.Connect(); err != nil {
        

        logrus.WithFields(logrus.Fields{
            "err": err,
            "server": clientdb_server,
        }).Warn("Error connecting to MySQL")
        
        return db, err
    }
    
    return db, nil

}


