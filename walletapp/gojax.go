package main

 import (
 	"fmt"
 	"net/http"
 	"html/template"
 	"bytes"
    "strings"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "github.com/FactomProject/fctwallet/Wallet"
    "github.com/FactomProject/factoid"
    "log"
 )

 var chttp = http.NewServeMux()
 var myState IState
 
 type inputList struct {
    InputSize float64 `json:"inputSize"`
    InputAddress string `json:"inputAddress"`
    }
    
 type outputList struct {
    OutputSize float64 `json:"outputSize"`
    OutputAddress string `json:"outputAddress"`
    OutputType string `json:"outputType"`
    }
 
 func check(e error, shouldEnd bool) {
    if e != nil {
        if shouldEnd {
            log.Fatal("Produced error: ", e)
        } else {
     	    log.Print("Produced error: ", e)
     	}
    }
 }


 func Home(w http.ResponseWriter, r *http.Request) {
 
    if (strings.Contains(r.URL.Path, ".")) {
        chttp.ServeHTTP(w, r)
    } else {
        t, err := template.ParseFiles("fwallet.html")
        if err != nil {
            fmt.Println("err: ", err)
        }
        t.Execute(w, nil)
    }

 }
 
 func currRate(w http.ResponseWriter, r *http.Request) {
		v, err := GetRate(myState)
		if err != nil {
			fmt.Println(err)
            return
		}
		w.Write([]byte(factoid.ConvertDecimal(uint64(v))))
 }
 
 func reqFee(w http.ResponseWriter, r *http.Request) {
        txKey := r.FormValue("key")
        
        ib := myState.GetFS().GetDB().GetRaw([]byte(factoid.DB_BUILD_TRANS), []byte(txKey))
		trans, ok := ib.(factoid.ITransaction)
		if ib != nil && ok {
			
			
			v, err := GetRate(myState)
			if err != nil {
				fmt.Println(err)
				return
			}
			fee, err := trans.CalculateFee(uint64(v))
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write([]byte(strings.TrimSpace(factoid.ConvertDecimal(fee))))
	    } else {
	        w.Write([]byte("..."))
	    }
			
 }
 
 
 func craftTx(w http.ResponseWriter, r *http.Request) {
   		txKey := r.FormValue("key")

        execStrings := []string{"NewTransaction", txKey}
        newErr := myState.Execute(execStrings)
        if newErr != nil {
            myR := FactoidDeleteTx(txKey)
            if myR != nil {
                fmt.Println(myR)
            }
        } 
        
        
            myState.Execute(execStrings)

            var buffer bytes.Buffer
            buffer.WriteString("Transaction " + txKey + ":\n\n")
            
            inputStr := r.FormValue("inputs")
            var inRes []inputList
            err := json.Unmarshal([]byte(inputStr), &inRes)
	        if err != nil {
		        fmt.Println("error:", err)
	        }
	        
            
            outputStr := r.FormValue("outputs")
            var outRes []outputList
            json.Unmarshal([]byte(outputStr), &outRes)
            
            var testInputFeed []string
            var testOutputFeed []string
            totalInputs := 0.0
            totalOutputs := 0.0
            
            for _, inputElement := range(inRes) {
                testInputFeed = []string{"AddInput", string(txKey), string(inputElement.InputAddress), strconv.FormatFloat(inputElement.InputSize, 'f', -1, 64)}
                totalInputs += inputElement.InputSize
                doubleTestErr := myState.Execute(testInputFeed)
                if doubleTestErr != nil {
                    fmt.Println(doubleTestErr)
                }
                
                buffer.WriteString("\tInput: " + inputElement.InputAddress + " : " + strconv.FormatFloat(inputElement.InputSize, 'f', -1, 64) + "\n")
            }
            
            myTest := []string{"Print", string(txKey)}   
                    myTestErr := myState.Execute(myTest)
                    if myTestErr != nil {
                        fmt.Println(myTestErr)
                    }    
            
            
            for _, outputElement := range(outRes) {
                totalOutputs += outputElement.OutputSize
                if outputElement.OutputType == "fct" {
                    testOutputFeed = []string{"AddOutput", string(txKey), string(outputElement.OutputAddress), strconv.FormatFloat(outputElement.OutputSize, 'f', -1, 64)}
                } else {
                    testOutputFeed = []string{"AddECOutput", string(txKey), string(outputElement.OutputAddress), strconv.FormatFloat(outputElement.OutputSize, 'f', -1, 64)}
                }
                
                tripleTestErr := myState.Execute(testOutputFeed)
                if tripleTestErr != nil {
                    fmt.Println(tripleTestErr)
                }   
                
                buffer.WriteString("\tOutput: " + outputElement.OutputAddress + " : " + strconv.FormatFloat(outputElement.OutputSize, 'f', -1, 64) + " (" + outputElement.OutputType + 
                                   ") \n")
            }
      	    currFee := totalInputs - totalOutputs
      	    
      	    buffer.WriteString("\n\tFee: " + strconv.FormatFloat(currFee, 'f', -1, 64))

                                   
      	    if r.Method == "GET" {
                    w.Write(buffer.Bytes())
            } else { //r.Method == "POST"
            
                    testPrintTx := []string{"Print", string(txKey)}   

                    oneTestErr := myState.Execute(testPrintTx)
                    if oneTestErr != nil {
                        fmt.Println(oneTestErr)
                    }      
                    
                    testSignFeed := []string{"Sign", string(txKey)}    
                    quadrupleTestErr := myState.Execute(testSignFeed)
                    if quadrupleTestErr != nil {
                        fmt.Println(quadrupleTestErr)
                    }

                    testSubmitFeed := []string{"Submit", string(txKey)}    
                    fiveTestErr := myState.Execute(testSubmitFeed)
                    if fiveTestErr != nil {
                        fmt.Println(fiveTestErr)
                    }
                       
                    buffer.WriteString("\n\nTransaction ")
                    buffer.WriteString(txKey)
                    buffer.WriteString(" has been submitted to Factom.")
                    w.Write(buffer.Bytes())
            }
 }
 
 func FactoidDeleteTx(key string) error {
	// Make sure we have a key
	if len(key) == 0 {
		return fmt.Errorf("Missing transaction key")
	}
	// Wipe out the key
	myState.GetFS().GetDB().DeleteKey([]byte(factoid.DB_BUILD_TRANS), []byte(key))
	return nil
 }
 
 
 func receiveAjax(w http.ResponseWriter, r *http.Request) {
 	if r.Method == "POST" {
 		ajax_post_data := r.FormValue("ajax_post_data")
 		call_type := r.FormValue("call_type")
 		switch call_type {
 		    case "balance":
 		        printBal, err := Wallet.FactoidBalance(ajax_post_data)
 		        check(err, false)
 		        w.Write([]byte("Factoid Address " + ajax_post_data + " Balance: " + strings.Trim(factoid.ConvertDecimal(uint64(printBal)), " ") + " ⨎"))
 		    case "balances":
 		        printBal := GetBalances(myState)
 		        testErr := myState.Execute([]string{"balances"})
 		        if testErr != nil {
                    fmt.Println(testErr.Error())
                    w.Write([]byte(testErr.Error()))
                    return
                }
 		        w.Write(printBal)
  		    case "allTxs":
 		        txNames, _, err := Wallet.GetTransactions()
 		        if err != nil {
 		            fmt.Println(err.Error())
 		            w.Write([]byte(err.Error()))
 		            return
 		        }
 		        if len(txNames) == 0 {
                    w.Write([]byte("No transactions to display."))
 		         	return
 		        }
 		        sliceTxNames := []byte("")
 		        for i:= range txNames {
 		            sliceTxNames = append(sliceTxNames, txNames[i]...)
 		            if i < len(txNames) - 1 {
 		                sliceTxNames = append(sliceTxNames, byte('\n'))
 		            }
 		        }
 		        w.Write(sliceTxNames)
 		    case "addNewAddress":
 		        if ajax_post_data != "" {
     		        genErr := GenAddress(myState, "fct", ajax_post_data)
     		        if genErr != nil {
 		                w.Write([]byte(genErr.Error()))
     		            return
     		        }
     		        w.Write([]byte(ajax_post_data + " has been added to your wallet successfully."));
                }
 		    case "addNewEC":
 		        if ajax_post_data != "" {
     		        genErr := GenAddress(myState, "ec", ajax_post_data)
     		        if genErr != nil {
 		                w.Write([]byte(genErr.Error()))
     		            return
     		        }
     		        w.Write([]byte(ajax_post_data + " has been added to your wallet successfully."));
     		    }
 		    /*
 		    case "addNewTx":
 		        execStrings := []string{"NewTransaction", ajax_post_data}
                newErr := myState.Execute(execStrings)
     		 	if newErr != nil {
     		 	    if newErr.Error()[:13] == "Duplicate key" {
     		 	        return //w.Write([]byte("Already have TX: " + ajax_post_data))
     		 	    }
     		 	    return
     		 	}
     		 	//Wallet.FactoidNewTransaction(ajax_post_data)	

             	w.Write([]byte(ajax_post_data))*/
        }
 		//©
 	} else {
 	    helpText, err := ioutil.ReadFile("./extra/help.txt")
        check(err, false)
        w.Write([]byte(helpText))
 	}
 }

 func startServer(state IState) {
 	// http.Handler
 	myState = state
 	chttp.Handle("/", http.FileServer(http.Dir("./extra/")))
 	mux := http.NewServeMux()
 	mux.HandleFunc("/", Home)
 	mux.HandleFunc("/receive", receiveAjax)
 	mux.HandleFunc("/rate", currRate)
 	mux.HandleFunc("/tx", craftTx)
 	mux.HandleFunc("/fee", reqFee)
 	
 	http.ListenAndServe(":2337", mux)
 }